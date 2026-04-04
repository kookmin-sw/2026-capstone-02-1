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

type IntType struct{}

func (ty IntType) isType()

type BoolType struct{}

func (ty BoolType) isType()

type StringType struct{}

func (ty StringType) isType()

type NoneType struct{}

func (ty NoneType) isType()

type ArrayType struct {
	Element_type ImpTypes
}

func (ty ArrayType) String() string {
	return fmt.Sprintf("ArrayType<%s>", ty.Element_type)
}

func (ty ArrayType) isType()

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
		array_ty, type_is_array := _type.(ArrayType)
		return type_is_array && array_ty.Element_type == val_ty.element_type
	case *StringVal:
		_, type_is_string := _type.(StringType)
		return type_is_string
	}
	return false
}

func get_type(val ImpValues) ImpTypes {
	switch val_ty := val.(type) {
	case *IntVal:
		return IntType{}
	case *BoolVal:
		return BoolType{}
	case *StringVal:
		return StringType{}
	case *ArrayVal:
		return ArrayType{Element_type: val_ty.element_type}
	default:
		panic(fmt.Sprintf("get_type: got unknown type for value %s\n", val))
	}
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
	Node
	name     string
	var_type ImpTypes
}

func (*VarExpr) isExpr() {}

func (expr VarExpr) String() string {
	return expr.name
}

type IntLitExpr struct {
	Node
	value int
}

func (*IntLitExpr) isExpr() {}

func (expr IntLitExpr) String() string {
	return fmt.Sprintf("%d", expr.value)
}

type BoolLitExpr struct {
	Node
	value bool
}

func (*BoolLitExpr) isExpr() {}

func (expr BoolLitExpr) String() string {
	return fmt.Sprintf("%t", expr.value)
}

type StringLitExpr struct {
	Node
	value string
}

func (*StringLitExpr) isExpr() {}

func (expr StringLitExpr) String() string {
	return expr.value
}

type ArrayLitExpr struct {
	Node
	element_type ImpTypes
	elements     []Expr
}

func (*ArrayLitExpr) isExpr() {}

func (expr ArrayLitExpr) String() string {
	var elems []string
	for _, elem := range expr.elements {
		elems = append(elems, fmt.Sprintf("%s", elem))
	}
	return fmt.Sprintf("{%s}", strings.Join(elems, ", "))
}

type AddExpr struct {
	Node
	lhs, rhs Expr
}

func (*AddExpr) isExpr() {}

func (expr AddExpr) String() string {
	return fmt.Sprintf("%s + %s", expr.lhs, expr.rhs)
}

type SubExpr struct {
	Node
	lhs, rhs Expr
}

func (*SubExpr) isExpr() {}

func (expr SubExpr) String() string {
	return fmt.Sprintf("%s - %s", expr.lhs, expr.rhs)
}

type MulExpr struct {
	Node
	lhs, rhs Expr
}

func (*MulExpr) isExpr() {}

func (expr MulExpr) String() string {
	return fmt.Sprintf("%s * %s", expr.lhs, expr.rhs)
}

type DivExpr struct {
	Node
	lhs, rhs Expr
}

func (*DivExpr) isExpr() {}

func (expr DivExpr) String() string {
	return fmt.Sprintf("%s / %s", expr.lhs, expr.rhs)
}

type ParenExpr struct {
	Node
	subexpr Expr
}

func (*ParenExpr) isExpr() {}

func (expr ParenExpr) String() string {
	return fmt.Sprintf("(%s)", expr.subexpr)
}

type ArrayIndexExpr struct {
	Node
	base  Expr
	index Expr
}

func (*ArrayIndexExpr) isExpr() {}

func (expr ArrayIndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", expr.base, expr.index)
}

type EqExpr struct {
	Node
	lhs, rhs Expr
}

func (*EqExpr) isExpr() {}

func (expr EqExpr) String() string {
	return fmt.Sprintf("%s == %s", expr.lhs, expr.rhs)
}

type NeqExpr struct {
	Node
	lhs, rhs Expr
}

func (*NeqExpr) isExpr() {}

func (expr NeqExpr) String() string {
	return fmt.Sprintf("%s != %s", expr.lhs, expr.rhs)
}

type NotExpr struct {
	Node
	subexpr Expr
}

func (*NotExpr) isExpr() {}

func (expr NotExpr) String() string {
	return fmt.Sprintf("!%s", expr.subexpr)
}

type AndExpr struct {
	Node
	lhs, rhs Expr
}

func (*AndExpr) isExpr() {}

func (expr AndExpr) String() string {
	return fmt.Sprintf("%s && %s", expr.lhs, expr.rhs)
}

type OrExpr struct {
	Node
	lhs, rhs Expr
}

func (*OrExpr) isExpr() {}

func (expr OrExpr) String() string {
	return fmt.Sprintf("%s || %s", expr.lhs, expr.rhs)
}

type CallExpr struct {
	Node
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

type MakeArrayExpr struct {
	Node
	size  Expr
	value Expr
}

func (*MakeArrayExpr) isExpr() {}

func (expr MakeArrayExpr) String() string {
	return fmt.Sprintf("make_array(%s, %s)", expr.size, expr.value)
}

type LenExpr struct {
	Node
	subexpr Expr
}

func (*LenExpr) isExpr() {}

func (expr LenExpr) String() string {
	return fmt.Sprintf("len(%s)", expr.subexpr)
}

// statements

type SkipStmt struct {
}

func (*SkipStmt) isStmt() {}

type AssignStmt struct {
	Node
	lhs Expr
	rhs Expr
}

func (*AssignStmt) isStmt() {}

func (stmt AssignStmt) String() string {
	return fmt.Sprintf("%s = %s", stmt.lhs, stmt.rhs)
}

type IfElseStmt struct {
	Node
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
	Node
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
	Node
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
	Node
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
	Node
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
	Node
	arg Expr
}

func (*ReturnStmt) isStmt() {}

func (stmt ReturnStmt) String() string {
	return fmt.Sprintf("return %s\n", stmt.arg)
}
