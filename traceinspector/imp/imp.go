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

func (ty IntType) isType() {}

func (ty IntType) String() string {
	return "int"
}

type BoolType struct{}

func (ty BoolType) isType() {}

func (ty BoolType) String() string {
	return "bool"
}

type StringType struct{}

func (ty StringType) isType() {}

func (ty StringType) String() string {
	return "string"
}

type NoneType struct{}

func (ty NoneType) isType() {}

func (ty NoneType) String() string {
	return "none"
}

type ArrayType struct {
	Element_type ImpTypes
}

func (ty ArrayType) String() string {
	return fmt.Sprintf("ArrayType<%s>", ty.Element_type)
}

func (ty ArrayType) isType() {}

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
		return type_is_array && array_ty.Element_type == val_ty.Element_type
	case *StringVal:
		_, type_is_string := _type.(StringType)
		return type_is_string
	}
	return false
}

// Get the ImpType of an ImpVal
func get_type(val ImpValues) ImpTypes {
	switch val_ty := val.(type) {
	case *IntVal:
		return IntType{}
	case *BoolVal:
		return BoolType{}
	case *StringVal:
		return StringType{}
	case *ArrayVal:
		return ArrayType{Element_type: val_ty.Element_type}
	default:
		panic(fmt.Sprintf("get_type: got unknown type for value %s\n", val))
	}
}

// Given two ImpVals, check if they have the same ImpType
func check_vals_type_equal(val1 ImpValues, val2 ImpValues) bool {
	return check_val_type_match(val1, get_type(val2))
}

////////////////////////

type ImpValues interface {
	isValue()
}

type IntVal struct {
	Val int
}

func (val IntVal) String() string {
	return fmt.Sprintf("%d", val.Val)
}

func (*IntVal) isValue() {}

type BoolVal struct {
	Val bool
}

func (val BoolVal) String() string {
	return fmt.Sprintf("%t", val.Val)
}

func (*BoolVal) isValue() {}

type ArrayVal struct {
	Element_type ImpTypes
	Val          []ImpValues
}

func (val ArrayVal) String() string {
	return fmt.Sprintf("%s", val.Val)
}

func (*ArrayVal) isValue() {}

type StringVal struct {
	Val string
}

func (val StringVal) String() string {
	return val.Val
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
	Name     string
	Arg_type ImpTypes
}

func (arg ArgPair) String() string {
	return fmt.Sprintf("%s %s", arg.Name, arg.Arg_type)
}

type ImpFunctionName string

type ImpFunction struct {
	Name        ImpFunctionName
	Arg_pairs   []ArgPair
	Body        []Stmt
	Return_type ImpTypes
}

func (fun ImpFunction) String() string {
	var args []string
	for _, arg := range fun.Arg_pairs {
		args = append(args, fmt.Sprintf("%s", arg))
	}
	var stmts []string
	for index, stmt := range fun.Body {
		stmts = append(stmts, fmt.Sprintf("%s", stmt))
		stmts[index] = strings.ReplaceAll(stmts[index], "\n", "\n\t")
	}
	return fmt.Sprintf("fun %s(%s) %s {\n\t%s\n}\n", fun.Name, strings.Join(args, ", "), fun.Return_type, strings.Join(stmts, "\n\t"))
}

// ////////////////////
type Node struct {
	Line_num int // return the line number corresponding to the original source
}

func (node *Node) GetLineNum() int {
	return node.Line_num
}

type Expr interface {
	isExpr()
	GetLineNum() int
}

type Stmt interface {
	isStmt()
	String() string
	GetLineNum() int
}

// expressions

type VarExpr struct {
	Node
	Name string
}

func (*VarExpr) isExpr() {}

func (expr VarExpr) String() string {
	return expr.Name
}

type IntLitExpr struct {
	Node
	Value int
}

func (*IntLitExpr) isExpr() {}

func (expr IntLitExpr) String() string {
	return fmt.Sprintf("%d", expr.Value)
}

type BoolLitExpr struct {
	Node
	Value bool
}

func (*BoolLitExpr) isExpr() {}

func (expr BoolLitExpr) String() string {
	return fmt.Sprintf("%t", expr.Value)
}

type StringLitExpr struct {
	Node
	Value string
}

func (*StringLitExpr) isExpr() {}

func (expr StringLitExpr) String() string {
	return "\"" + expr.Value + "\""
}

type ArrayLitExpr struct {
	Node
	Elements []Expr
}

func (*ArrayLitExpr) isExpr() {}

func (expr ArrayLitExpr) String() string {
	var elems []string
	for _, elem := range expr.Elements {
		elems = append(elems, fmt.Sprintf("%s", elem))
	}
	return fmt.Sprintf("{%s}", strings.Join(elems, ", "))
}

type AddExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*AddExpr) isExpr() {}

func (expr AddExpr) String() string {
	return fmt.Sprintf("%s + %s", expr.Lhs, expr.Rhs)
}

type SubExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*SubExpr) isExpr() {}

func (expr SubExpr) String() string {
	return fmt.Sprintf("%s - %s", expr.Lhs, expr.Rhs)
}

type MulExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*MulExpr) isExpr() {}

func (expr MulExpr) String() string {
	return fmt.Sprintf("%s * %s", expr.Lhs, expr.Rhs)
}

type DivExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*DivExpr) isExpr() {}

func (expr DivExpr) String() string {
	return fmt.Sprintf("%s / %s", expr.Lhs, expr.Rhs)
}

type ModExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*ModExpr) isExpr() {}

func (expr ModExpr) String() string {
	return fmt.Sprintf("%s %% %s", expr.Lhs, expr.Rhs)
}

type ParenExpr struct {
	Node
	Subexpr Expr
}

func (*ParenExpr) isExpr() {}

func (expr ParenExpr) String() string {
	return fmt.Sprintf("(%s)", expr.Subexpr)
}

type ArrayIndexExpr struct {
	Node
	Base  Expr
	Index Expr
}

func (*ArrayIndexExpr) isExpr() {}

func (expr ArrayIndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", expr.Base, expr.Index)
}

type EqExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*EqExpr) isExpr() {}

func (expr EqExpr) String() string {
	return fmt.Sprintf("%s == %s", expr.Lhs, expr.Rhs)
}

type NeqExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*NeqExpr) isExpr() {}

func (expr NeqExpr) String() string {
	return fmt.Sprintf("%s != %s", expr.Lhs, expr.Rhs)
}

type LessthanExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*LessthanExpr) isExpr() {}

func (expr LessthanExpr) String() string {
	return fmt.Sprintf("%s < %s", expr.Lhs, expr.Rhs)
}

type GreaterthanExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*GreaterthanExpr) isExpr() {}

func (expr GreaterthanExpr) String() string {
	return fmt.Sprintf("%s > %s", expr.Lhs, expr.Rhs)
}

type LeqExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*LeqExpr) isExpr() {}

func (expr LeqExpr) String() string {
	return fmt.Sprintf("%s <= %s", expr.Lhs, expr.Rhs)
}

type GeqExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*GeqExpr) isExpr() {}

func (expr GeqExpr) String() string {
	return fmt.Sprintf("%s >= %s", expr.Lhs, expr.Rhs)
}

// arithmetic negation
type NegExpr struct {
	Node
	Subexpr Expr
}

func (*NegExpr) isExpr() {}

func (expr NegExpr) String() string {
	return fmt.Sprintf("-%s", expr.Subexpr)
}

type NotExpr struct {
	Node
	Subexpr Expr
}

func (*NotExpr) isExpr() {}

func (expr NotExpr) String() string {
	return fmt.Sprintf("!%s", expr.Subexpr)
}

type AndExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*AndExpr) isExpr() {}

func (expr AndExpr) String() string {
	return fmt.Sprintf("%s && %s", expr.Lhs, expr.Rhs)
}

type OrExpr struct {
	Node
	Lhs, Rhs Expr
}

func (*OrExpr) isExpr() {}

func (expr OrExpr) String() string {
	return fmt.Sprintf("%s || %s", expr.Lhs, expr.Rhs)
}

type CallExpr struct {
	Node
	Func_name ImpFunctionName
	Args      []Expr
}

func (*CallExpr) isExpr() {}

func (expr CallExpr) String() string {
	var func_args []string
	for _, arg := range expr.Args {
		func_args = append(func_args, fmt.Sprintf("%s", arg))
	}
	return fmt.Sprintf("%s(%s)", expr.Func_name, strings.Join(func_args, ", "))
}

type MakeArrayExpr struct {
	Node
	Size  Expr
	Value Expr
}

func (*MakeArrayExpr) isExpr() {}

func (expr MakeArrayExpr) String() string {
	return fmt.Sprintf("make_array(%s, %s)", expr.Size, expr.Value)
}

type LenExpr struct {
	Node
	Subexpr Expr
}

func (*LenExpr) isExpr() {}

func (expr LenExpr) String() string {
	return fmt.Sprintf("len(%s)", expr.Subexpr)
}

// statements

type SkipStmt struct {
	Node
}

func (*SkipStmt) isStmt() {}

func (*SkipStmt) String() string {
	return "Skip"
}

type AssignStmt struct {
	Node
	Lhs Expr
	Rhs Expr
}

func (*AssignStmt) isStmt() {}

func (stmt AssignStmt) String() string {
	return fmt.Sprintf("%s = %s", stmt.Lhs, stmt.Rhs)
}

type IfElseStmt struct {
	Node
	Cond       Expr
	True_stmt  []Stmt
	False_stmt []Stmt
}

func (*IfElseStmt) isStmt() {}

func (stmt IfElseStmt) String() string {
	var true_stmts []string
	var false_stmts []string
	for index, stmt := range stmt.True_stmt {
		true_stmts = append(true_stmts, fmt.Sprintf("\t%s", stmt))
		true_stmts[index] = strings.ReplaceAll(true_stmts[index], "\n", "\n\t")
	}
	for index, stmt := range stmt.False_stmt {
		false_stmts = append(false_stmts, fmt.Sprintf("\t%s", stmt))
		false_stmts[index] = strings.ReplaceAll(false_stmts[index], "\n", "\n\t")
	}
	return fmt.Sprintf("if %s {\n%s\n} else {\n%s\n}", stmt.Cond, strings.Join(true_stmts, "\n"), strings.Join(false_stmts, "\n"))
}

type WhileStmt struct {
	Node
	Cond      Expr
	Body_stmt []Stmt
}

func (*WhileStmt) isStmt() {}

func (stmt WhileStmt) String() string {
	var substmts []string
	for index, stmt := range stmt.Body_stmt {
		substmts = append(substmts, fmt.Sprintf("\t%s", stmt))
		substmts[index] = strings.ReplaceAll(substmts[index], "\n", "\n\t")
	}
	return fmt.Sprintf("while %s {\n%s\n}", stmt.Cond, strings.Join(substmts, "\n"))
}

type BreakStmt struct {
	Node
}

func (*BreakStmt) isStmt() {}

func (stmt BreakStmt) String() string {
	return "break"
}

type ContinueStmt struct {
	Node
}

func (*ContinueStmt) isStmt() {}

func (stmt ContinueStmt) String() string {
	return "continue"
}

type IncStmt struct {
	Node
	Subexpr Expr
}

func (*IncStmt) isStmt() {}

func (stmt IncStmt) String() string {
	return fmt.Sprintf("%s++", stmt.Subexpr)
}

type DecStmt struct {
	Node
	Subexpr Expr
}

func (*DecStmt) isStmt() {}

func (stmt DecStmt) String() string {
	return fmt.Sprintf("%s--", stmt.Subexpr)
}

type CallStmt struct {
	Node
	Func_name ImpFunctionName
	Args      []Expr
}

func (*CallStmt) isStmt() {}

func (stmt CallStmt) String() string {
	var args []string
	for _, arg := range stmt.Args {
		args = append(args, fmt.Sprintf("%s", arg))
	}
	return fmt.Sprintf("%s(%s)", stmt.Func_name, strings.Join(args, ", "))
}

type PrintStmt struct {
	Node
	Args []Expr
}

func (*PrintStmt) isStmt() {}

func (stmt PrintStmt) String() string {
	var args []string
	for _, arg := range stmt.Args {
		args = append(args, fmt.Sprintf("%s", arg))
	}
	return fmt.Sprintf("print(%s)", strings.Join(args, ", "))
}

type ScanfStmt struct {
	Node
	Format_string    string
	Assign_locations []Expr
}

func (*ScanfStmt) isStmt() {}

func (stmt ScanfStmt) String() string {
	var args []string
	for _, arg := range stmt.Assign_locations {
		args = append(args, fmt.Sprintf("%s", arg))
	}
	return fmt.Sprintf("Scanf(%s, %s)", stmt.Format_string, strings.Join(args, ", "))
}

type ReturnStmt struct {
	Node
	Arg Expr
}

func (*ReturnStmt) isStmt() {}

func (stmt ReturnStmt) String() string {
	return fmt.Sprintf("return %s", stmt.Arg)
}
