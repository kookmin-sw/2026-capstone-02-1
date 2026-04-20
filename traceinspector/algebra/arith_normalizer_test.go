package algebra

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
	"traceinspector/imp"
)

func Test_zero_rhs(t *testing.T) {
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
				original := translator.Translate_Expr(expr)
				switch expr.Op {
				case token.EQL, token.NEQ, token.GEQ, token.LEQ, token.GTR, token.LSS:
					new_stmt, _ := zero_rhs(original)
					t.Logf("%s -> %s\n", original, new_stmt)
				}
			}
			return true
		}

		file, err := parser.ParseFile(fset, gofile, nil, parser.ParseComments)
		if err != nil {
			fmt.Println("Error while parsing", gofile, "-", err)
			return
		}
		ast.Inspect(file, runner)
	}
}

func Test_normalize_integer_expr(t *testing.T) {
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
					original, _ = zero_rhs(original)
					switch original_ty := original.(type) {
					case *imp.EqExpr:
						original = original_ty.Lhs
					case *imp.NeqExpr:
						original = original_ty.Lhs
					case *imp.GeqExpr:
						original = original_ty.Lhs
					case *imp.LeqExpr:
						original = original_ty.Lhs
					case *imp.LessthanExpr:
						original = original_ty.Lhs
					case *imp.GreaterthanExpr:
						original = original_ty.Lhs
					}
					new_expr, err := normalize_integer_expr(original)
					if err == nil {
						t.Logf("polynomial: %s -> %s\n", original, new_expr)
					} else {
						t.Error(err)
					}
				}
			}
			return true
		}

		file, err := parser.ParseFile(fset, gofile, nil, parser.ParseComments)
		if err != nil {
			fmt.Println("Error while parsing", gofile, "-", err)
			return
		}
		ast.Inspect(file, runner)
	}
}

func Test_convert_subtraction(t *testing.T) {
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
					original, _ = zero_rhs(original)
					switch original_ty := original.(type) {
					case *imp.EqExpr:
						original = original_ty.Lhs
					case *imp.NeqExpr:
						original = original_ty.Lhs
					case *imp.GeqExpr:
						original = original_ty.Lhs
					case *imp.LeqExpr:
						original = original_ty.Lhs
					case *imp.LessthanExpr:
						original = original_ty.Lhs
					case *imp.GreaterthanExpr:
						original = original_ty.Lhs
					}
					new_expr := convert_subtraction_to_neg(original, false)
					t.Logf("convert subtraction: %s -> %s\n", original, new_expr)
				}
			}
			return true
		}

		file, err := parser.ParseFile(fset, gofile, nil, parser.ParseComments)
		if err != nil {
			fmt.Println("Error while parsing", gofile, "-", err)
			return
		}
		ast.Inspect(file, runner)
	}
}
