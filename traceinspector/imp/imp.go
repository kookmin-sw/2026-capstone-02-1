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

type IntValueExpr struct {
	value int
	Node
}

func (*IntValueExpr) isExpr() {}

type BoolValueExpr struct {
	value bool
}

func (*BoolValueExpr) isExpr() {}

type AddExpr struct {
	lhs, rhs Expr
}

func (*AddExpr) isExpr() {}

type SubExpr struct {
	lhs, rhs Expr
}

func (*SubExpr) isExpr() {}

type MulExpr struct {
	lhs, rhs Expr
}

func (*MulExpr) isExpr() {}

type DivExpr struct {
	lhs, rhs Expr
}

func (*DivExpr) isExpr() {}

type ArrayIndexExpr struct {
	base  Expr
	index Expr
}

func (*ArrayIndexExpr) isExpr() {}

type EqExpr struct {
	lhs, rhs Expr
}

func (*EqExpr) isExpr() {}

type NeqExpr struct {
	lhs, rhs Expr
}

func (*NeqExpr) isExpr() {}

type NotExpr struct {
	subexpr Expr
}

func (*NotExpr) isExpr() {}

type AndExpr struct {
	lhs, rhs Expr
}

func (*AndExpr) isExpr() {}

type OrExpr struct {
	lhs, rhs Expr
}

func (*OrExpr) isExpr() {}

type CallExpr struct {
	func_name string
	args      []Expr
}

func (*CallExpr) isExpr() {}

// statements

type SeqStmt struct {
	lhs, rhs Stmt
}

func (*SeqStmt) isStmt() {}

type SkipStmt struct {
}

func (*SkipStmt) isStmt() {}

type AssignStmt struct {
	lhs Expr
	rhs Expr
}

func (*AssignStmt) isStmt() {}

type InputStmt struct {
	assign_var VarExpr
}

func (*InputStmt) isStmt() {}

type IfElseStmt struct {
	cond       Expr
	true_stmt  Stmt
	false_stmt Stmt
}

func (*IfElseStmt) isStmt() {}

type WhileStmt struct {
	cond      Expr
	body_stmt Stmt
}

func (*WhileStmt) isStmt() {}

type CallStmt struct {
	func_name string
	args      []Expr
}

func (*CallStmt) isStmt() {}

type ReturnStmt struct {
	arg Expr
}

func (*ReturnStmt) isStmt() {}
