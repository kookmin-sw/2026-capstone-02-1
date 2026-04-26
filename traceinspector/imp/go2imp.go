package imp

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"slices"
	"strconv"
	"strings"
)

// Return an AST node's code string
func nodeString(n ast.Node) string {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	format.Node(&buf, fset, n)
	return buf.String()
}

type ImpFunctionMap map[ImpFunctionName]ImpFunction

type Go2ImpTranslator struct {
	Fset           *token.FileSet
	for_post_stmts [][]Stmt // stack of post stmts, in case of nested loop
}

func (nh *Go2ImpTranslator) get_top_post_stmt() []Stmt {
	return nh.for_post_stmts[len(nh.for_post_stmts)-1]
}

func (nh *Go2ImpTranslator) push_post_stmt(stmts []Stmt) {
	nh.for_post_stmts = append(nh.for_post_stmts, stmts)
}

func (nh *Go2ImpTranslator) pop_post_stmt() {
	nh.for_post_stmts = nh.for_post_stmts[:len(nh.for_post_stmts)-1]
}

func (nh *Go2ImpTranslator) translate_ast_node_as_ImpType(ast_node ast.Node) ImpTypes {
	switch node_ty := ast_node.(type) {
	case *ast.Ident:
		switch node_ty.Name {
		case "int":
			return IntType{}
		case "bool":
			return BoolType{}
		default:
			panic(fmt.Sprintf("go2imp: Unsupported type '%s'\n", node_ty.Name))
		}
	case *ast.ArrayType:
		return ArrayType{Element_type: nh.translate_ast_node_as_ImpType(node_ty.Elt)}
	default:
		panic(fmt.Sprintf("go2imp: Unsupported node as type: '%s'\n", nodeString(ast_node)))
	}
}

func (nh *Go2ImpTranslator) get_ast_linenum(ast_node ast.Node) int {
	return nh.Fset.Position(ast_node.Pos()).Line
}

func (nh *Go2ImpTranslator) create_node_struct_from_ast(ast_node ast.Node) Node {
	return Node{Line_num: nh.get_ast_linenum(ast_node)}
}

func (nh *Go2ImpTranslator) translate_BasicLit(expr *ast.BasicLit) Expr {
	switch expr.Kind {
	case token.INT:
		i, err := strconv.Atoi(expr.Value)
		if err != nil {
			panic(fmt.Sprintf("go2imp translate_BasicLit: Error while translating %s: %s", nodeString(expr), err))
		}
		return &IntLitExpr{Node: Node{nh.get_ast_linenum(expr)}, Value: i}
	case token.STRING:
		// str, _ := strconv.Unquote(expr.Value)
		// return &StringLitExpr{Node: nh.create_node_struct_from_ast(expr), Value: str}
		return &StringLitExpr{Node: nh.create_node_struct_from_ast(expr), Value: strings.Trim(expr.Value, "\"")}

	default:
		panic(fmt.Sprintf("go2imp translate_BasicLit: Unsupported literal '%s'", nodeString(expr)))
	}
}

func (nh *Go2ImpTranslator) translate_BinaryExpr(expr *ast.BinaryExpr) Expr {
	lhs_expr := nh.Translate_Expr(expr.X)
	rhs_expr := nh.Translate_Expr(expr.Y)
	node := nh.create_node_struct_from_ast(expr)
	switch expr.Op {
	case token.ADD:
		return &AddExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.SUB:
		return &SubExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.MUL:
		return &MulExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.QUO:
		return &DivExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.REM:
		return &ModExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.LAND:
		return &AndExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.LOR:
		return &OrExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.EQL:
		return &EqExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.NEQ:
		return &NeqExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.LSS:
		return &LessthanExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.GTR:
		return &GreaterthanExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.LEQ:
		return &LeqExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	case token.GEQ:
		return &GeqExpr{Node: node, Lhs: lhs_expr, Rhs: rhs_expr}
	default:
		panic(fmt.Sprintf("go2imp: Unsupported token.Token %s", expr.Op))
	}
}

func (nh *Go2ImpTranslator) translate_CallExpr(expr *ast.CallExpr) Expr {
	func_ident, func_is_ident := expr.Fun.(*ast.Ident)
	if !func_is_ident {
		panic("go2imp: Only idents are allowed as functions in ast.CallExpr")
	}
	var translated_args []Expr
	for _, expr_node := range expr.Args {
		translated_args = append(translated_args, nh.Translate_Expr(expr_node))
	}
	switch func_ident.Name {
	case "make_array":
		if len(translated_args) != 2 {
			panic(fmt.Sprintf("go2imp: make_array() expects 2 arguments, but got %d", len(translated_args)))
		}
		return &MakeArrayExpr{Node: nh.create_node_struct_from_ast(expr), Size: translated_args[0], Value: translated_args[1]}
	case "len":
		if len(translated_args) != 1 {
			panic(fmt.Sprintf("go2imp: len() expects 2 arguments, but got %d", len(translated_args)))
		}
		return &LenExpr{Node: nh.create_node_struct_from_ast(expr), Subexpr: translated_args[0]}
	}
	return &CallExpr{Node: nh.create_node_struct_from_ast(expr), Func_name: ImpFunctionName(func_ident.Name), Args: translated_args}
}

func (nh *Go2ImpTranslator) translate_CompositeLit(expr *ast.CompositeLit) Expr {
	var translated_elements []Expr
	for _, elem_node := range expr.Elts {
		translated_elements = append(translated_elements, nh.Translate_Expr(elem_node))
	}
	return &ArrayLitExpr{Node: nh.create_node_struct_from_ast(expr), Elements: translated_elements}
}

func (nh *Go2ImpTranslator) translate_Ident(expr *ast.Ident) Expr {
	switch expr.Name {
	case "true":
		return &BoolLitExpr{Node: nh.create_node_struct_from_ast(expr), Value: true}
	case "false":
		return &BoolLitExpr{Node: nh.create_node_struct_from_ast(expr), Value: false}
	}
	return &VarExpr{Node: nh.create_node_struct_from_ast(expr), Name: expr.Name}
}

func (nh *Go2ImpTranslator) translate_IndexExpr(expr *ast.IndexExpr) Expr {
	base_expr := nh.Translate_Expr(expr.X)
	index_expr := nh.Translate_Expr(expr.Index)
	return &ArrayIndexExpr{Node: nh.create_node_struct_from_ast(expr), Base: base_expr, Index: index_expr}
}

func (nh *Go2ImpTranslator) translate_ParenExpr(expr *ast.ParenExpr) Expr {
	subexpr := nh.Translate_Expr(expr.X)
	return &ParenExpr{Node: nh.create_node_struct_from_ast(expr), Subexpr: subexpr}
}

func (nh *Go2ImpTranslator) translate_UnaryExpr(expr *ast.UnaryExpr) Expr {
	switch expr.Op {
	case token.NOT:
		return &NotExpr{Node: nh.create_node_struct_from_ast(expr), Subexpr: nh.Translate_Expr(expr.X)}
	case token.SUB:
		return &NegExpr{Node: nh.create_node_struct_from_ast(expr), Subexpr: nh.Translate_Expr(expr.X)}
	default:
		panic(fmt.Sprintf("go2imp: Unsupported unary operator type '%s'\n", expr.Op))
	}
}

func (nh *Go2ImpTranslator) Translate_Expr(expr ast.Expr) Expr {
	switch expr := (expr).(type) {
	case *ast.BasicLit:
		return nh.translate_BasicLit(expr)
	case *ast.BinaryExpr:
		return nh.translate_BinaryExpr(expr)
	case *ast.CallExpr:
		return nh.translate_CallExpr(expr)
	case *ast.CompositeLit:
		return nh.translate_CompositeLit(expr)
	case *ast.Ident:
		return nh.translate_Ident(expr)
	case *ast.IndexExpr:
		return nh.translate_IndexExpr(expr)
	case *ast.ParenExpr:
		return nh.translate_ParenExpr(expr)
	case *ast.UnaryExpr:
		return nh.translate_UnaryExpr(expr)
	default:
		panic(fmt.Sprintf("go2imp translate_Expr: unsupported ast.Expr node type: %T", expr))
	}
}

func (nh *Go2ImpTranslator) translate_AssignStmt(stmt *ast.AssignStmt) []Stmt {
	if !(len(stmt.Lhs) == len(stmt.Rhs) && len(stmt.Lhs) == 1) {
		panic("go2imp: Number of LHS and RHS expressions in an assignment must be 1\n")
	}
	return []Stmt{&AssignStmt{Node: nh.create_node_struct_from_ast(stmt), Lhs: nh.Translate_Expr(stmt.Lhs[0]), Rhs: nh.Translate_Expr(stmt.Rhs[0])}}
}

func (nh *Go2ImpTranslator) translate_BlockStmt(stmt *ast.BlockStmt) []Stmt {
	var return_stmts []Stmt
	for _, ast_stmt := range stmt.List {
		return_stmts = slices.Concat(return_stmts, nh.Translate_Stmt(ast_stmt))
	}
	return return_stmts
}

func (nh *Go2ImpTranslator) translate_EmptyStmt(stmt *ast.EmptyStmt) []Stmt {
	return []Stmt{&SkipStmt{Node: nh.create_node_struct_from_ast(stmt)}}
}

func (nh *Go2ImpTranslator) translate_ExprStmt(stmt *ast.ExprStmt) []Stmt {
	call_subexpr, is_callexpr := stmt.X.(*ast.CallExpr)
	if !is_callexpr {
		panic(fmt.Sprintf("go2imp: Subexpression of ExprStmt must be a CallExpr, but got '%s'\n", nodeString(stmt.X)))
	}

	var translated_args []Expr
	for _, expr_node := range call_subexpr.Args {
		translated_args = append(translated_args, nh.Translate_Expr(expr_node))
	}

	create_scanf := func() Stmt {
		format_string_impexpr, format_string_is_string := translated_args[0].(*StringLitExpr)
		if !format_string_is_string {
			panic("go2imp: First argument of Scanf/fmt.Scanf must be a string literal")
		}
		return &ScanfStmt{Node: nh.create_node_struct_from_ast(stmt), Format_string: format_string_impexpr.Value, Assign_locations: translated_args[1:]}
	}

	func_ident, func_is_ident := call_subexpr.Fun.(*ast.Ident)
	if !func_is_ident {
		func_selectorexpr, func_is_selector := call_subexpr.Fun.(*ast.SelectorExpr)
		if func_is_selector && func_selectorexpr.Sel.Name == "Scanf" {
			return []Stmt{create_scanf()}
		} else if func_is_selector && func_selectorexpr.Sel.Name == "Print" {
			return []Stmt{&PrintStmt{Node: nh.create_node_struct_from_ast(stmt), Args: translated_args}}
		}
		panic("go2imp: Only idents or fmt.Scanf/fmt.Print are allowed as functions in ast.CallExpr")
	}
	switch func_ident.Name {
	case "Scanf":
		return []Stmt{create_scanf()}
	case "Print":
		return []Stmt{&PrintStmt{Node: nh.create_node_struct_from_ast(stmt), Args: translated_args}}
	default:
		return []Stmt{&CallStmt{Node: nh.create_node_struct_from_ast(stmt), Func_name: ImpFunctionName(func_ident.Name), Args: translated_args}}
	}
}

func (nh *Go2ImpTranslator) translate_ForStmt(stmt *ast.ForStmt) []Stmt {
	var init_stmts []Stmt
	if stmt.Init != nil {
		init_stmts = nh.Translate_Stmt(stmt.Init)
	}
	cond_expr := nh.Translate_Expr(stmt.Cond)
	if stmt.Post != nil {
		nh.push_post_stmt(nh.Translate_Stmt(stmt.Post))
	} else {
		nh.push_post_stmt(nil)
	}

	body_stmt := nh.Translate_Stmt(stmt.Body)
	if stmt.Post != nil {
		body_stmt = slices.Concat(body_stmt, nh.get_top_post_stmt())
	}
	nh.pop_post_stmt()
	return append(init_stmts, &WhileStmt{Node: nh.create_node_struct_from_ast(stmt), Cond: cond_expr, Body_stmt: body_stmt})
}

func (nh *Go2ImpTranslator) translate_IfStmt(stmt *ast.IfStmt) []Stmt {
	var else_stmts []Stmt
	if stmt.Else != nil {
		else_stmts = nh.Translate_Stmt(stmt.Else)
	}
	true_stmts := nh.Translate_Stmt(stmt.Body)
	return []Stmt{&IfElseStmt{Node: nh.create_node_struct_from_ast(stmt), Cond: nh.Translate_Expr(stmt.Cond), True_stmt: true_stmts, False_stmt: else_stmts}}
}

func (nh *Go2ImpTranslator) translate_IncDecStmt(stmt *ast.IncDecStmt) []Stmt {
	switch stmt.Tok {
	case token.INC:
		// return []Stmt{&IncStmt{Node: nh.create_node_struct_from_ast(stmt), Subexpr: nh.Translate_Expr(stmt.X)}}
		return []Stmt{&AssignStmt{Node: nh.create_node_struct_from_ast(stmt), Lhs: nh.Translate_Expr(stmt.X), Rhs: &AddExpr{Node: nh.create_node_struct_from_ast(stmt.X), Lhs: nh.Translate_Expr(stmt.X), Rhs: &IntLitExpr{Node: nh.create_node_struct_from_ast(stmt.X), Value: 1}}}}
	case token.DEC:
		// return []Stmt{&DecStmt{Node: nh.create_node_struct_from_ast(stmt), Subexpr: nh.Translate_Expr(stmt.X)}}
		return []Stmt{&AssignStmt{Node: nh.create_node_struct_from_ast(stmt), Lhs: nh.Translate_Expr(stmt.X), Rhs: &SubExpr{Node: nh.create_node_struct_from_ast(stmt.X), Lhs: nh.Translate_Expr(stmt.X), Rhs: &IntLitExpr{Node: nh.create_node_struct_from_ast(stmt.X), Value: 1}}}}
	default:
		panic(fmt.Sprintf("go2imp: Unsupported IncDecStmt token '%s'\n", stmt.Tok))
	}
}

func (nh *Go2ImpTranslator) translate_ReturnStmt(stmt *ast.ReturnStmt) []Stmt {
	if len(stmt.Results) != 1 {
		panic(fmt.Sprintf("go2imp: only a single return expression allowed, but return returns %d exprs at line %d\n", len(stmt.Results), nh.get_ast_linenum(stmt)))
	}
	return []Stmt{&ReturnStmt{Node: nh.create_node_struct_from_ast(stmt), Arg: nh.Translate_Expr(stmt.Results[0])}}
}

func (nh *Go2ImpTranslator) translate_BranchStmt(stmt *ast.BranchStmt) []Stmt {
	switch stmt.Tok {
	case token.BREAK:
		return []Stmt{&BreakStmt{Node: nh.create_node_struct_from_ast(stmt)}}
	case token.CONTINUE:
		if len(nh.get_top_post_stmt()) > 0 {
			// make sure post stmt is added
			return append(nh.get_top_post_stmt(), &ContinueStmt{Node: nh.create_node_struct_from_ast(stmt)})
		} else {
			return []Stmt{&ContinueStmt{Node: nh.create_node_struct_from_ast(stmt)}}
		}
	default:
		panic(fmt.Sprintf("go2imp: Unsupported BranchStmt token %s\n", stmt.Tok))
	}
}

func (nh *Go2ImpTranslator) Translate_Stmt(stmt ast.Stmt) []Stmt {
	switch stmt := (stmt).(type) {
	case *ast.AssignStmt:
		return nh.translate_AssignStmt(stmt)
	case *ast.BlockStmt:
		return nh.translate_BlockStmt(stmt)
	// case *ast.DeclStmt:
	// 	return nh.translate_DeclStmt(stmt)
	case *ast.EmptyStmt:
		return nh.translate_EmptyStmt(stmt)
	case *ast.ExprStmt:
		return nh.translate_ExprStmt(stmt)
	case *ast.ForStmt:
		return nh.translate_ForStmt(stmt)
	case *ast.IfStmt:
		return nh.translate_IfStmt(stmt)
	case *ast.IncDecStmt:
		return nh.translate_IncDecStmt(stmt)
	case *ast.ReturnStmt:
		return nh.translate_ReturnStmt(stmt)
	case *ast.BranchStmt:
		return nh.translate_BranchStmt(stmt)
	default:
		panic(fmt.Sprintf("go2imp: translate_Stmt: unsupported ast.Stmt node type: %T", stmt))
	}
}

func Translate_ast_file_to_imp(go_input_file *ast.File, fset *token.FileSet) ImpFunctionMap {
	translator := Go2ImpTranslator{Fset: fset}
	output := make(ImpFunctionMap)
	for _, decl := range go_input_file.Decls {
		func_decl, is_func_decl := decl.(*ast.FuncDecl)
		if is_func_decl {
			var func_argpairs []ArgPair
			for _, field := range func_decl.Type.Params.List {
				if len(field.Names) != 1 {
					panic("go2imp: Function arguments must be individually defined")
				}
				func_argpairs = append(func_argpairs, ArgPair{Name: field.Names[0].Name, Arg_type: translator.translate_ast_node_as_ImpType(field.Type)})
			}
			var return_type ImpTypes = NoneType{}
			if func_decl.Type.Results != nil {
				if len(func_decl.Type.Results.List) > 1 {
					panic(fmt.Sprintf("go2imp: Function must return at most 1 argument, but function '%s' returns %d values\n", func_decl.Name, len(func_decl.Type.Results.List)))
				}
				return_type = translator.translate_ast_node_as_ImpType(func_decl.Type.Results.List[0].Type)
			}
			output[ImpFunctionName(func_decl.Name.Name)] = ImpFunction{Name: ImpFunctionName(func_decl.Name.Name), Arg_pairs: func_argpairs, Body: translator.Translate_Stmt(func_decl.Body), Return_type: return_type}
		}
	}
	return output
}
