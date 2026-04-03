package imp

////////////
// Type declarations

type ImpTypes interface {
	isType()
}

type IntType ImpTypes
type BoolType ImpTypes

type ImpArray struct {
	element_type ImpTypes
	len          int
}

func (ty ImpArray) isType()

//////////////////////

type Node interface {
	return_code_linenum() int // return the line number corresponding to the original source
	repr() string             // return the string representation of the code
}

// expressions

type Expr interface {
	Node
}

type VarExpr struct {
	name     string
	var_type ImpTypes
}

type IntValueExpr struct {
	value int
}

type BoolValueExpr struct {
	value bool
}

type AddExpr struct {
	lhs, rhs Expr
}

type SubExpr struct {
	lhs, rhs Expr
}

type MulExpr struct {
	lhs, rhs Expr
}

type DivExpr struct {
	lhs, rhs Expr
}

type ArrayIndexExpr struct {
	base  Expr
	index Expr
}

type EqExpr struct {
	lhs, rhs Expr
}

type NeqExpr struct {
	lhs, rhs Expr
}

type NotExpr struct {
	subexpr Expr
}

type AndExpr struct {
	lhs, rhs Expr
}

type OrExpr struct {
	lhs, rhs Expr
}

// statements

type SeqStmt struct {
	lhs, rhs Stmt
}

type SkipStmt struct {
}

type AssignStmt struct {
	lhs Expr
	rhs Expr
}

type InputStmt struct {
	assign_var VarExpr
}

type IfElseStmt struct {
	cond       Expr
	true_stmt  Stmt
	false_stmt Stmt
}

type WhileStmt struct {
	cond      Expr
	body_stmt Stmt
}

type Stmt interface {
	Node
}
