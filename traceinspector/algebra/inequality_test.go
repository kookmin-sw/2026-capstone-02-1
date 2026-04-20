package algebra

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
	"traceinspector/imp"
)

func Test_inequality(t *testing.T) {
	test_go_files, err := filepath.Glob("test_files/*.go")
	if err != nil {
		panic(err)
	}
	for _, gofile := range test_go_files {
		fset := token.NewFileSet()
		// https://pkg.go.dev/go/parser#Mode
		translator := imp.Go2ImpTranslator{Fset: fset}
		runner := func(n ast.Node) bool {
			switch expr := n.(type) {
			case *ast.BinaryExpr:
				switch expr.Op {
				case token.EQL, token.NEQ, token.GEQ, token.LEQ, token.GTR, token.LSS:
					original := translator.Translate_Expr(expr)
					original_zero, _ := zero_rhs(original)
					new_expr, safe := imp_expr_to_simp_inequality(original_zero)
					if safe {
						t.Errorf("ineq: %s -> %s\n", original, new_expr)
					}
				}
			}
			return true
		}

		file, err := parser.ParseFile(fset, gofile, nil, parser.ParseComments)
		if err != nil {
			t.Error("Error while parsing", gofile, "-", err)
			return
		}
		ast.Inspect(file, runner)
	}
}
