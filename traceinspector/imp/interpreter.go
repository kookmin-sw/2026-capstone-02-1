package imp

import "fmt"

type ImpState struct {
	vars map[string]any // all exprs are reduced to go values
}

type ImpInterpreter struct {
	states    []*ImpState // a stack of program states
	functions map[string]*Stmt
}

func (interpreter *ImpInterpreter) get_top_state() *ImpState {
	return interpreter.states[len(interpreter.states)-1]
}

func (interpreter *ImpInterpreter) eval_VarExpr(node VarExpr) any {
	var_value, var_exists := interpreter.get_top_state().vars[node.name]
	if !var_exists {
		panic("Unknown variable " + node.name)
	}
	return var_value
}

func (interpreter *ImpInterpreter) eval_Expr(node Expr) any {

}

func (interpreter *ImpInterpreter) eval_IntValueExpr(node IntValueExpr) int {
	return node.value
}

func (interpreter *ImpInterpreter) eval_BoolValueExpr(node BoolValueExpr) bool {
	return node.value
}

func (interpreter *ImpInterpreter) eval_AddExpr(node AddExpr) int {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(int)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(int)

	if !lhs_is_int {
		panic(fmt.Sprintf("LHS of addition should be an int value, but got '%s'", node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of addition should be an int value, but got '%s'", node.rhs))
	}
	return lhs_val + rhs_val
}

func (interpreter *ImpInterpreter) eval_SubExpr(node SubExpr) int {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(int)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(int)

	if !lhs_is_int {
		panic(fmt.Sprintf("LHS of subtraction should be an int value, but got '%s'", node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of subtraction should be an int value, but got '%s'", node.rhs))
	}
	return lhs_val - rhs_val
}
