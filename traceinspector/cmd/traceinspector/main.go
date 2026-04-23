package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"traceinspector"
	"traceinspector/imp"
)

func main() {
	// https://pkg.go.dev/flag#String
	input_path := flag.String("gofile", "", "")
	cfg_json_argname := "print-cfg-json"
	cfg_mermaid_argname := "print-cfg-mermaid"
	print_imp_argname := "print-imp"
	interpret_imp_argname := "interpret-imp"

	_ = flag.Bool(cfg_json_argname, false, "whether to just print cfg and exit")
	_ = flag.Bool(cfg_mermaid_argname, false, "whether to just print the mermaid graph and exit")
	_ = flag.Bool(print_imp_argname, false, "whether to just print the translated Imp code and exit")
	_ = flag.Bool(interpret_imp_argname, false, "whether to just interpret the translated Imp code and exit")
	flag.Parse()
	if *input_path == "" {
		panic("need to pass input go file path with --gofile")
	}

	just_print_cfg := false
	just_print_imp := false
	just_print_mermaid := false
	just_interpret := false
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case cfg_json_argname:
			just_print_cfg = true
		case cfg_mermaid_argname:
			just_print_mermaid = true
		case print_imp_argname:
			just_print_imp = true
		case interpret_imp_argname:
			just_interpret = true
		}
	})

	fset := token.NewFileSet()
	// https://pkg.go.dev/go/parser#Mode
	file, err := parser.ParseFile(fset, *input_path, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error while parsing", input_path, "-", err)
		return
	}

	imp_functions := imp.Translate_ast_file_to_imp(file, fset)

	cfg_map := traceinspector.Create_cfg(imp_functions)

	if just_print_cfg {
		traceinspector.Print_cfg_map_json(cfg_map)
		return
	}

	if just_print_mermaid {
		for fun_name, fun_cfg := range cfg_map {
			fmt.Println(fun_name)
			fmt.Println("----------------")
			traceinspector.Print_mermaid(fun_cfg)
		}
		return
	}

	if just_print_imp {
		for _, fun := range imp_functions {
			fmt.Println(fun)
		}
		return
	}

	if just_interpret {
		interpreter := imp.ImpInterpreter{Functions: imp_functions}
		interpreter.Interpret_main()
		return
	}
	traceinspector.Test(cfg_map, "main", imp_functions)
}
