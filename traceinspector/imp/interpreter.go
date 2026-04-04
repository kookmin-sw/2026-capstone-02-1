package imp

import (
	"fmt"
	"strings"
)

type ImpState struct {
	vars                  map[string]ImpValues // all exprs are reduced to go values
	current_function_name string
	return_value          ImpValues // The return value of the current local scope, if exists
}

type ImpInterpreter struct {
	States    []*ImpState // a stack of program states to represent scopes
	Functions map[string]ImpFunction
}

func (interpreter *ImpInterpreter) get_top_state() *ImpState {
	return interpreter.States[len(interpreter.States)-1]
}

// Return the topmost variable name and bool indicating if the variable exists
func (interpreter *ImpInterpreter) get_variable(name string) (ImpValues, bool) {
	for stack_index := len(interpreter.States) - 1; stack_index >= 0; stack_index-- {
		var_value, var_exists := interpreter.States[stack_index].vars[name]
		if var_exists {
			return var_value, true
		}
	}
	return nil, false
}

func (interpreter *ImpInterpreter) push_state(state ImpState) {
	interpreter.States = append(interpreter.States, &state)
}

func (interpreter *ImpInterpreter) pop_state() {
	interpreter.States = interpreter.States[:len(interpreter.States)-1]
}

func (interpreter *ImpInterpreter) eval_VarExpr(node VarExpr) ImpValues {
	var_value, var_exists := interpreter.get_variable(node.name)
	if !var_exists {
		panic("Unknown variable " + node.name)
	}
	return var_value
}

func (interpreter *ImpInterpreter) eval_Expr(node Expr) ImpValues {
	switch node_ty := node.(type) {
	case *VarExpr:
		return interpreter.eval_VarExpr(*node_ty)
	case *IntLitExpr:
		return interpreter.eval_IntLitExpr(*node_ty)
	case *BoolLitExpr:
		return interpreter.eval_BoolLitExpr(*node_ty)
	case *AddExpr:
		return interpreter.eval_AddExpr(*node_ty)
	case *SubExpr:
		return interpreter.eval_SubExpr(*node_ty)
	case *MulExpr:
		return interpreter.eval_MulExpr(*node_ty)
	case *DivExpr:
		return interpreter.eval_DivExpr(*node_ty)
	case *ParenExpr:
		return interpreter.eval_ParenExpr(*node_ty)
	case *ArrayIndexExpr:
		return interpreter.eval_ArrayIndexExpr(*node_ty)
	case *EqExpr:
		return interpreter.eval_EqExpr(*node_ty)
	case *NeqExpr:
		return interpreter.eval_NeqExpr(*node_ty)
	case *NotExpr:
		return interpreter.eval_NotExpr(*node_ty)
	case *AndExpr:
		return interpreter.eval_AndExpr(*node_ty)
	case *OrExpr:
		return interpreter.eval_OrExpr(*node_ty)
	case *CallExpr:
		return interpreter.eval_CallExpr(*node_ty)
	default:
		panic(fmt.Sprintf("Unimplemented expr type %s", node))
	}
}

func (interpreter *ImpInterpreter) eval_Expr_lvalue(lhs Expr, rhs ImpValues) ImpValues {
	lhs_var, lhs_is_var := lhs.(*VarExpr)
	if lhs_is_var {
		_, lhs_exists := interpreter.get_variable(lhs_var.name)
		if !lhs_exists {
			switch rhs.(type) {
			case *IntVal:
				interpreter.get_top_state().vars[lhs_var.name] = &IntVal{}
			case *BoolVal:
				interpreter.get_top_state().vars[lhs_var.name] = &BoolVal{}
			case *ArrayVal:
				interpreter.get_top_state().vars[lhs_var.name] = &ArrayVal{}
			}
		}
		var_val, _ := interpreter.get_variable(lhs_var.name)
		return var_val
	} else {
		return interpreter.eval_Expr(lhs)
	}
}

func (interpreter *ImpInterpreter) eval_IntLitExpr(node IntLitExpr) ImpValues {
	return &IntVal{val: node.value}
}

func (interpreter *ImpInterpreter) eval_BoolLitExpr(node BoolLitExpr) ImpValues {
	return &BoolVal{val: node.value}
}

func (interpreter *ImpInterpreter) eval_AddExpr(node AddExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("LHS of addition should be an int value, but got '%s'", node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of addition should be an int value, but got '%s'", node.rhs))
	}
	return &IntVal{val: lhs_val.val + rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_SubExpr(node SubExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("LHS of subtraction should be an int value, but got '%s'", node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of subtraction should be an int value, but got '%s'", node.rhs))
	}
	return &IntVal{val: lhs_val.val - rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_MulExpr(node MulExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("LHS of multiplication should be an int value, but got '%s'", node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of multiplication should be an int value, but got '%s'", node.rhs))
	}
	return &IntVal{val: lhs_val.val * rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_DivExpr(node DivExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("LHS of division should be an int value, but got '%s'", node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of division should be an int value, but got '%s'", node.rhs))
	}
	return &IntVal{val: lhs_val.val / rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_ParenExpr(node ParenExpr) ImpValues {
	return interpreter.eval_Expr(node.subexpr)
}

func (interpreter *ImpInterpreter) eval_ArrayIndexExpr(node ArrayIndexExpr) ImpValues {
	index_val, index_is_int := interpreter.eval_Expr(node.index).(*IntVal)
	if !index_is_int {
		panic(fmt.Sprintf("Index of array indexing should be an int value, but got '%s'", node.index))
	}
	base_val, base_is_arrayval := interpreter.eval_Expr(node.base).(*ArrayVal)
	if !base_is_arrayval {
		panic(fmt.Sprintf("Expr %s is not an array", node.base))
	}
	return base_val.val[index_val.val]
}

func (interpreter *ImpInterpreter) eval_EqExpr(node EqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.lhs)
	rhs_val := interpreter.eval_Expr(node.rhs)
	return &BoolVal{val: lhs_val == rhs_val}
}

func (interpreter *ImpInterpreter) eval_NeqExpr(node NeqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.lhs)
	rhs_val := interpreter.eval_Expr(node.rhs)
	return &BoolVal{val: lhs_val != rhs_val}
}

func (interpreter *ImpInterpreter) eval_NotExpr(node NotExpr) ImpValues {
	subexpr_val, subexpr_is_bool := interpreter.eval_Expr(node.subexpr).(*BoolVal)
	if !subexpr_is_bool {
		panic(fmt.Sprintf("Subexpr %s of NOT operator should be of type bool", node.subexpr))
	}
	return &BoolVal{val: !subexpr_val.val}
}

func (interpreter *ImpInterpreter) eval_AndExpr(node AndExpr) ImpValues {
	lhs_val, lhs_is_bool := interpreter.eval_Expr(node.lhs).(*BoolVal)
	rhs_val, rhs_is_bool := interpreter.eval_Expr(node.rhs).(*BoolVal)

	if !lhs_is_bool {
		panic(fmt.Sprintf("LHS of AND should be a bool value, but got '%s'", node.lhs))
	}

	if !rhs_is_bool {
		panic(fmt.Sprintf("RHS of AND should be a bool value, but got '%s'", node.rhs))
	}
	return &BoolVal{val: lhs_val.val && rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_OrExpr(node OrExpr) ImpValues {
	lhs_val, lhs_is_bool := interpreter.eval_Expr(node.lhs).(*BoolVal)
	rhs_val, rhs_is_bool := interpreter.eval_Expr(node.rhs).(*BoolVal)

	if !lhs_is_bool {
		panic(fmt.Sprintf("LHS of OR should be a bool value, but got '%s'", node.lhs))
	}

	if !rhs_is_bool {
		panic(fmt.Sprintf("RHS of OR should be a bool value, but got '%s'", node.rhs))
	}
	return &BoolVal{val: lhs_val.val || rhs_val.val}
}

// Imp is pass-by-value for int/bool, but arrays are passed references
func (interpreter *ImpInterpreter) eval_function_call(func_name string, args []Expr) ImpValues {
	// copy values if primitive
	prepare_args := func(arg ImpValues) ImpValues {
		switch arg_ty := arg.(type) {
		case *IntVal:
			return &IntVal{val: arg_ty.val}
		case *BoolVal:
			return &BoolVal{val: arg_ty.val}
		case *ArrayVal:
			return arg_ty
		}
		panic(fmt.Sprintf("Function call: Unknown arg type %s", arg))
	}
	func_local_state := ImpState{vars: make(map[string]ImpValues)}
	for index, arg_expr := range args {
		arg_info := interpreter.Functions[func_name].Arg_names[index]
		arg_val := prepare_args(interpreter.eval_Expr(arg_expr))
		if !check_val_type_match(arg_val, arg_info.arg_type) {
			panic(fmt.Sprintf("Argument '%s' of function '%s' is defined as type %s, but passed expr '%s' of type %s", arg_info.name, func_name, arg_info.arg_type, arg_expr, arg_val))
		}
		func_local_state.vars[arg_info.name] = arg_val
	}
	interpreter.push_state(func_local_state)
	interpreter.eval_Stmt(interpreter.Functions[func_name].Body)
	return_value := interpreter.get_top_state().return_value
	interpreter.pop_state()

	if return_value == nil {
		return &NoneVal{}
	} else {
		return return_value
	}
}

func (interpreter *ImpInterpreter) eval_CallExpr(node CallExpr) ImpValues {
	return interpreter.eval_function_call(node.func_name, node.args)
}

func (interpreter *ImpInterpreter) eval_SkipStmt(SkipStmt) {}

func (interpreter *ImpInterpreter) eval_AssignStmt(node AssignStmt) {
	rhs_val := interpreter.eval_Expr(node.rhs)
	switch lhs_loc := interpreter.eval_Expr_lvalue(node.lhs, rhs_val).(type) {
	case *IntVal:
		rhs_intval, rhs_is_intval := rhs_val.(*IntVal)
		if !rhs_is_intval {
			panic(fmt.Sprintf("Attempted to assign RHS '%s' of type %s to LHS '%s' of type %s", node.rhs, rhs_val, node.lhs, lhs_loc))
		}
		lhs_loc.val = rhs_intval.val
	case *BoolVal:
		rhs_intval, rhs_is_boolval := rhs_val.(*BoolVal)
		if !rhs_is_boolval {
			panic(fmt.Sprintf("Attempted to assign RHS '%s' of type %s to LHS '%s' of type %s", node.rhs, rhs_val, node.lhs, lhs_loc))
		}
		lhs_loc.val = rhs_intval.val
	case *ArrayVal:
		rhs_intval, rhs_is_arrayval := rhs_val.(*ArrayVal)
		if !rhs_is_arrayval {
			panic(fmt.Sprintf("Attempted to assign RHS '%s' of type %s to LHS '%s' of type %s", node.rhs, rhs_val, node.lhs, lhs_loc))
		}
		lhs_loc.val = rhs_intval.val
	}
}

func (interpreter *ImpInterpreter) eval_IfElseStmt(node IfElseStmt) {
	cond_val := interpreter.eval_Expr(node.cond)
	cond_boolval, cond_is_bool := cond_val.(*BoolVal)
	if !cond_is_bool {
		panic(fmt.Sprintf("If statement got non-boolean condition '%s'\n", node.cond))
	}
	if cond_boolval.val {
		interpreter.eval_Stmt(node.true_stmt)
	} else {
		interpreter.eval_Stmt(node.false_stmt)
	}
}

func (interpreter *ImpInterpreter) eval_WhileStmt(node WhileStmt) {
	for true {
		cond_val := interpreter.eval_Expr(node.cond)
		cond_boolval, cond_is_bool := cond_val.(*BoolVal)
		if !cond_is_bool {
			panic(fmt.Sprintf("While statement got non-boolean condition '%s'\n", node.cond))
		}
		if cond_boolval.val == false {
			break
		}
		interpreter.eval_Stmt(node.body_stmt)
	}
}

func (interpreter *ImpInterpreter) eval_CallStmt(node CallStmt) {
	interpreter.eval_function_call(node.func_name, node.args)
}

func (interpreter *ImpInterpreter) eval_PrintStmt(node PrintStmt) {
	var outputs []string
	for _, arg := range node.args {
		arg_val := interpreter.eval_Expr(arg)
		outputs = append(outputs, fmt.Sprintf("%s", arg_val))
	}
	fmt.Print(strings.Join(outputs, " "))
}

func (interpreter *ImpInterpreter) eval_ScanfStmt(node ScanfStmt) {
	for index, fmt_str := range strings.Split(node.format_string, " ") {
		var imp_val ImpValues
		switch fmt_str {
		case "%d":
			var input int
			fmt.Scanf(fmt_str, &input)
			imp_val = &IntVal{val: input}
		case "%t":
			var input bool
			fmt.Scanf(fmt_str, &input)
			imp_val = &BoolVal{val: input}
		default:
			panic(fmt.Sprintf("scanf: Unsupported formatting specifier %s\n", fmt_str))
		}

		lhs_val := node.assign_locations[index]
		switch lhs_loc := interpreter.eval_Expr_lvalue(lhs_val, imp_val).(type) {
		case *IntVal:
			rhs_intval, rhs_is_intval := imp_val.(*IntVal)
			if !rhs_is_intval {
				panic(fmt.Sprintf("scanf: Attempted to assign input '%s' of type %T to variable '%s' of type %T\n", imp_val, imp_val, lhs_loc, lhs_loc))
			}
			lhs_loc.val = rhs_intval.val
		case *BoolVal:
			rhs_intval, rhs_is_boolval := imp_val.(*BoolVal)
			if !rhs_is_boolval {
				panic(fmt.Sprintf("scanf: Attempted to assign input '%s' of type %T to variable '%s' of type %T\n", imp_val, imp_val, lhs_loc, lhs_loc))
			}
			lhs_loc.val = rhs_intval.val
		}
	}
}

func (interpreter *ImpInterpreter) eval_ReturnStmt(node ReturnStmt) {
	top_state := interpreter.get_top_state()
	top_state.return_value = interpreter.eval_Expr(node.arg)
	if top_state.current_function_name != "" {
		expected_return_type := interpreter.Functions[top_state.current_function_name].Return_type
		if !check_val_type_match(top_state.return_value, expected_return_type) {
			panic(fmt.Sprintf("Function %s should return value of type %s, but actually returned '%s' of type %s\n", top_state.current_function_name, expected_return_type, node.arg, top_state.return_value))
		}
	}
}

func (interpreter *ImpInterpreter) eval_Stmt(nodes []Stmt) {
	for _, stmt := range nodes {
		switch stmt := stmt.(type) {
		case *SkipStmt:
			interpreter.eval_SkipStmt(*stmt)
		case *AssignStmt:
			interpreter.eval_AssignStmt(*stmt)
		case *IfElseStmt:
			interpreter.eval_IfElseStmt(*stmt)
		case *WhileStmt:
			interpreter.eval_WhileStmt(*stmt)
		case *CallStmt:
			interpreter.eval_CallStmt(*stmt)
		case *PrintStmt:
			interpreter.eval_PrintStmt(*stmt)
		case *ScanfStmt:
			interpreter.eval_ScanfStmt(*stmt)
		case *ReturnStmt:
			interpreter.eval_ReturnStmt(*stmt)
		}
	}
}
