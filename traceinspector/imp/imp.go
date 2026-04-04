package imp

import (
	"fmt"
	"strings"
)

////////////
// Type declarations

type ImpTypes interface {
	isType()
}

type IntType ImpTypes

type BoolType ImpTypes
type StringType ImpTypes
type NonType ImpTypes

type ImpArray struct {
	Element_type ImpTypes
	Len          int
}

func (ty ImpArray) String() string {
	return fmt.Sprintf("ArrayType<%s>", ty.Element_type)
}

func (ty ImpArray) isType()

// Given a Value and a specified Type, check that the Value is of the specified type
func check_val_type_match(val ImpValues, _type ImpTypes) bool {
	switch val_ty := val.(type) {
	case *IntVal:
		_, type_is_int := _type.(IntType)
		return type_is_int
	case *BoolVal:
		_, type_is_bool := _type.(BoolType)
		return type_is_bool
	case *ArrayVal:
		array_ty, type_is_array := _type.(ImpArray)
		return type_is_array && array_ty.Element_type == val_ty.element_type
	case *StringVal:
		_, type_is_string := _type.(StringType)
		return type_is_string
	}
	return false
}

////////////////////////

type ImpValues interface {
	isValue()
}

type IntVal struct {
	val int
}

func (val IntVal) String() string {
	return val.String()
}

func (*IntVal) isValue() {}

type BoolVal struct {
	val bool
}

func (val BoolVal) String() string {
	return val.String()
}

func (*BoolVal) isValue() {}

type ArrayVal struct {
	element_type ImpTypes
	val          []ImpValues
}

func (val ArrayVal) String() string {
	return val.String()
}

func (*ArrayVal) isValue() {}

type StringVal struct {
	val string
}

func (val StringVal) String() string {
	return val.val
}

func (*StringVal) isValue() {}

type NoneVal struct{}

func (val NoneVal) String() string {
	return "NoneVal"
}

func (*NoneVal) isValue() {}

//////////////////////
// Function definitions

type ArgPair struct {
	name     string
	arg_type ImpTypes
}

type ImpFunction struct {
	Arg_names   []ArgPair
	Body        []Stmt
	Return_type ImpTypes
}

// ////////////////////
type Node struct {
	Line_num int // return the line number corresponding to the original source
}

type Expr interface {
	isExpr()
}

type Stmt interface {
	isStmt()
}

// expressions

type VarExpr struct {
	name     string
	var_type ImpTypes
	Node
}

func (*VarExpr) isExpr() {}

func (expr VarExpr) String() string {
	return expr.name
}

type IntLitExpr struct {
	value int
	Node
}

func (*IntLitExpr) isExpr() {}

func (expr IntLitExpr) String() string {
	return fmt.Sprintf("%d", expr.value)
}

type BoolLitExpr struct {
	value bool
}

func (*BoolLitExpr) isExpr() {}

func (expr BoolLitExpr) String() string {
	return fmt.Sprintf("%t", expr.value)
}

type AddExpr struct {
	lhs, rhs Expr
}

func (*AddExpr) isExpr() {}

func (expr AddExpr) String() string {
	return fmt.Sprintf("%s + %s", expr.lhs, expr.rhs)
}

type SubExpr struct {
	lhs, rhs Expr
}

func (*SubExpr) isExpr() {}

func (expr SubExpr) String() string {
	return fmt.Sprintf("%s - %s", expr.lhs, expr.rhs)
}

type MulExpr struct {
	lhs, rhs Expr
}

func (*MulExpr) isExpr() {}

func (expr MulExpr) String() string {
	return fmt.Sprintf("%s * %s", expr.lhs, expr.rhs)
}

type DivExpr struct {
	lhs, rhs Expr
}

func (*DivExpr) isExpr() {}

func (expr DivExpr) String() string {
	return fmt.Sprintf("%s / %s", expr.lhs, expr.rhs)
}

type ParenExpr struct {
	subexpr Expr
}

func (*ParenExpr) isExpr() {}

func (expr ParenExpr) String() string {
	return fmt.Sprintf("(%s)", expr.subexpr)
}

type ArrayIndexExpr struct {
	base  Expr
	index Expr
}

func (*ArrayIndexExpr) isExpr() {}

func (expr ArrayIndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", expr.base, expr.index)
}

type EqExpr struct {
	lhs, rhs Expr
}

func (*EqExpr) isExpr() {}

func (expr EqExpr) String() string {
	return fmt.Sprintf("%s == %s", expr.lhs, expr.rhs)
}

type NeqExpr struct {
	lhs, rhs Expr
}

func (*NeqExpr) isExpr() {}

func (expr NeqExpr) String() string {
	return fmt.Sprintf("%s != %s", expr.lhs, expr.rhs)
}

type NotExpr struct {
	subexpr Expr
}

func (*NotExpr) isExpr() {}

func (expr NotExpr) String() string {
	return fmt.Sprintf("!%s", expr.subexpr)
}

type AndExpr struct {
	lhs, rhs Expr
}

func (*AndExpr) isExpr() {}

func (expr AndExpr) String() string {
	return fmt.Sprintf("%s && %s", expr.lhs, expr.rhs)
}

type OrExpr struct {
	lhs, rhs Expr
}

func (*OrExpr) isExpr() {}

func (expr OrExpr) String() string {
	return fmt.Sprintf("%s || %s", expr.lhs, expr.rhs)
}

type CallExpr struct {
	func_name string
	args      []Expr
}

func (*CallExpr) isExpr() {}

func (expr CallExpr) String() string {
	var func_args []string
	for _, arg := range expr.args {
		func_args = append(func_args, fmt.Sprintf("%s", arg))
	}
	return fmt.Sprintf("%s(%s)", expr.func_name, strings.Join(func_args, ", "))
}

// statements

type SkipStmt struct {
}

func (*SkipStmt) isStmt() {}

type AssignStmt struct {
	lhs Expr
	rhs Expr
}

func (*AssignStmt) isStmt() {}

func (stmt AssignStmt) String() string {
	return fmt.Sprintf("%s = %s", stmt.lhs, stmt.rhs)
}

type IfElseStmt struct {
	cond       Expr
	true_stmt  []Stmt
	false_stmt []Stmt
}

func (*IfElseStmt) isStmt() {}

func (stmt IfElseStmt) String() string {
	var true_stmts []string
	var false_stmts []string
	for _, stmt := range stmt.true_stmt {
		true_stmts = append(true_stmts, fmt.Sprintf("    %s", stmt))
	}
	for _, stmt := range stmt.false_stmt {
		false_stmts = append(false_stmts, fmt.Sprintf("    %s", stmt))
	}
	return fmt.Sprintf("if %s {\n%s} else {\n%s}\n", stmt.cond, strings.Join(true_stmts, "\n"), strings.Join(false_stmts, "\n"))
}

type WhileStmt struct {
	cond      Expr
	body_stmt []Stmt
}

func (*WhileStmt) isStmt() {}

func (stmt WhileStmt) String() string {
	var substmts []string
	for _, stmt := range stmt.body_stmt {
		substmts = append(substmts, fmt.Sprintf("    %s", stmt))
	}
	return fmt.Sprintf("while %s {\n%s}\n", stmt.cond, strings.Join(substmts, "\n"))
}

type CallStmt struct {
	func_name string
	args      []Expr
}

func (*CallStmt) isStmt() {}

func (stmt CallStmt) String() string {
	var args []string
	for _, arg := range stmt.args {
		args = append(args, fmt.Sprintf("%s", arg))
	}
	return fmt.Sprintf("%s(%s)\n", stmt.func_name, strings.Join(args, ", "))
}

type PrintStmt struct {
	args []Expr
}

func (*PrintStmt) isStmt() {}

func (stmt PrintStmt) String() string {
	var args []string
	for _, arg := range stmt.args {
		args = append(args, fmt.Sprintf("%s", arg))
	}
	return fmt.Sprintf("print(%s)\n", strings.Join(args, ", "))
}

type ScanfStmt struct {
	format_string    string
	assign_locations []Expr
}

func (*ScanfStmt) isStmt() {}

func (stmt ScanfStmt) String() string {
	var args []string
	for _, arg := range stmt.assign_locations {
		args = append(args, fmt.Sprintf("%s", arg))
	}
	return fmt.Sprintf("scanf(%s, %s)\n", stmt.format_string, strings.Join(args, ", "))
}

type ReturnStmt struct {
	arg Expr
}

func (*ReturnStmt) isStmt() {}

func (stmt ReturnStmt) String() string {
	return fmt.Sprintf("return %s\n", stmt.arg)
}
