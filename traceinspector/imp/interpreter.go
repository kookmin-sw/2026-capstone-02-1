package imp

import (
	"fmt"
	"strings"
)

type ImpState struct {
	vars                  map[string]ImpValues // all exprs are reduced to go values
	current_function_name ImpFunctionName
	return_value          ImpValues // The return value of the current local scope, if exists
}

// This designates the control flow result of executing statements
type ControlflowResult int

// ControlNormal: continue execution of statements as normal
// ControlBreak: break from the current loop
// ControlContinue: return to the loop head(continue) of the current loop
// ControlReturn: return from the current function, stopping further execution of the function body
const (
	ControlNormal ControlflowResult = iota
	ControlBreak
	ControlContinue
	ControlReturn
)

type ImpInterpreter struct {
	States    []*ImpState // a stack of program states to represent scopes
	Functions ImpFunctionMap
}

// Interpret Imp code starting from main()
func (interpreter *ImpInterpreter) Interpret_main() {
	interpreter.eval_function_call("main", nil, 0)
}

func (interpreter *ImpInterpreter) get_top_state() *ImpState {
	return interpreter.States[len(interpreter.States)-1]
}

// Return the topmost variable name and bool indicating if the variable exists
func (interpreter *ImpInterpreter) get_variable(name string) (ImpValues, bool) {
	toplevel_func_name := interpreter.States[len(interpreter.States)-1].current_function_name
	for stack_index := len(interpreter.States) - 1; stack_index >= 0; stack_index-- {
		// Check that the scope is within the function stack
		if toplevel_func_name != interpreter.States[stack_index].current_function_name {
			break
		}

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

func (ImpInterpreter *ImpInterpreter) deepcopy_impvalue(val ImpValues) ImpValues {
	switch val_ty := val.(type) {
	case *IntVal:
		return &IntVal{Val: val_ty.Val}
	case *BoolVal:
		return &BoolVal{Val: val_ty.Val}
	case *StringVal:
		return &StringVal{Val: val_ty.Val}
	case *ArrayVal:
		arrayval := ArrayVal{Element_type: val_ty.Element_type}
		var copied_slice []ImpValues
		for _, elem := range val_ty.Val {
			copied_slice = append(copied_slice, ImpInterpreter.deepcopy_impvalue(elem))
		}
		arrayval.Val = copied_slice
		return &arrayval
	default:
		panic(fmt.Sprintf("deepcopy_impvalue: Unknown value type %T\n", val))
	}
}

////////////////////////

func (interpreter *ImpInterpreter) eval_VarExpr(node VarExpr) ImpValues {
	var_value, var_exists := interpreter.get_variable(node.Name)
	if !var_exists {
		panic(fmt.Sprintf("Line %d: Unknown variable '%s'", node.Node.Line_num, node.Name))
	}
	return var_value
}

func (interpreter *ImpInterpreter) eval_Expr_lvalue(lhs Expr, rhs_ty ImpTypes) ImpValues {
	lhs_var, lhs_is_var := lhs.(*VarExpr)
	if lhs_is_var {
		_, lhs_exists := interpreter.get_variable(lhs_var.Name)
		if !lhs_exists {
			switch ty := rhs_ty.(type) {
			case IntType:
				interpreter.get_top_state().vars[lhs_var.Name] = &IntVal{}
			case BoolType:
				interpreter.get_top_state().vars[lhs_var.Name] = &BoolVal{}
			case ArrayType:
				interpreter.get_top_state().vars[lhs_var.Name] = &ArrayVal{Element_type: ty.Element_type}
			default:
				panic(fmt.Sprintf("Line %d: Unknown rhs type %T", lhs_var.Line_num, ty))
			}
		}
		var_val, _ := interpreter.get_variable(lhs_var.Name)
		return var_val
	} else {
		return interpreter.eval_Expr(lhs)
	}
}

func (interpreter *ImpInterpreter) eval_IntLitExpr(node IntLitExpr) ImpValues {
	return &IntVal{Val: node.Value}
}

func (interpreter *ImpInterpreter) eval_BoolLitExpr(node BoolLitExpr) ImpValues {
	return &BoolVal{Val: node.Value}
}

func (interpreter *ImpInterpreter) eval_StringLitExpr(node StringLitExpr) ImpValues {
	return &StringVal{Val: node.Value}
}

func (interpreter *ImpInterpreter) eval_ArrayLitExpr(node ArrayLitExpr) ImpValues {
	var elem_vals []ImpValues
	var elem_type ImpTypes
	for _, elem := range node.Elements {
		elem_val := interpreter.eval_Expr(elem)
		if elem_type == nil {
			elem_type = get_type(elem_val)
		}
		if !check_val_type_match(elem_val, elem_type) {
			panic(fmt.Sprintf("Line %d: Array element type is identified as %s, but got expr '%s' with value %s\n", node.Line_num, elem_type, elem, elem_val))
		}
		elem_vals = append(elem_vals, elem_val)
	}
	return &ArrayVal{Element_type: elem_type, Val: elem_vals}
}

func (interpreter *ImpInterpreter) eval_AddExpr(node AddExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.Lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.Rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: LHS of addition should be an int value, but got '%s'", node.Line_num, node.Lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("Line %d: RHS of addition should be an int value, but got '%s'", node.Line_num, node.Rhs))
	}
	return &IntVal{Val: lhs_val.Val + rhs_val.Val}
}

func (interpreter *ImpInterpreter) eval_SubExpr(node SubExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.Lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.Rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: LHS of subtraction should be an int value, but got '%s'", node.Line_num, node.Lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of subtraction should be an int value, but got '%s'", node.Rhs))
	}
	return &IntVal{Val: lhs_val.Val - rhs_val.Val}
}

func (interpreter *ImpInterpreter) eval_MulExpr(node MulExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.Lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.Rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: LHS of multiplication should be an int value, but got '%s'", node.Line_num, node.Lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("Line %d: RHS of multiplication should be an int value, but got '%s'", node.Line_num, node.Rhs))
	}
	return &IntVal{Val: lhs_val.Val * rhs_val.Val}
}

func (interpreter *ImpInterpreter) eval_DivExpr(node DivExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.Lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.Rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: LHS of division should be an int value, but got '%s'", node.Line_num, node.Lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("Line %d: RHS of division should be an int value, but got '%s'", node.Line_num, node.Rhs))
	}
	return &IntVal{Val: lhs_val.Val / rhs_val.Val}
}

func (interpreter *ImpInterpreter) eval_ModExpr(node ModExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.Lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.Rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: LHS of modulus should be an int value, but got '%s'", node.Line_num, node.Lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("Line %d: RHS of modulus should be an int value, but got '%s'", node.Line_num, node.Rhs))
	}
	return &IntVal{Val: lhs_val.Val % rhs_val.Val}
}

func (interpreter *ImpInterpreter) eval_ParenExpr(node ParenExpr) ImpValues {
	return interpreter.eval_Expr(node.Subexpr)
}

func (interpreter *ImpInterpreter) eval_ArrayIndexExpr(node ArrayIndexExpr) ImpValues {
	index_val, index_is_int := interpreter.eval_Expr(node.Index).(*IntVal)
	if !index_is_int {
		panic(fmt.Sprintf("Line %d: Index of array indexing should be an int value, but got '%s'", node.Line_num, node.Index))
	}
	base_val, base_is_arrayval := interpreter.eval_Expr(node.Base).(*ArrayVal)
	if !base_is_arrayval {
		panic(fmt.Sprintf("Line %d: Expr '%s' is not an array", node.Line_num, node.Base))
	}
	return base_val.Val[index_val.Val]
}

func (interpreter *ImpInterpreter) eval_EqExpr(node EqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.Lhs)
	rhs_val := interpreter.eval_Expr(node.Rhs)
	if !check_vals_type_equal(lhs_val, rhs_val) {
		panic(fmt.Sprintf("Line %d: Unsupported '==' between '%s' and '%s'", node.Line_num, lhs_val, rhs_val))
	}
	switch lhs_val := lhs_val.(type) {
	case *IntVal:
		rhs_val, _ := rhs_val.(*IntVal)
		return &BoolVal{Val: lhs_val.Val == rhs_val.Val}
	case *BoolVal:
		rhs_val, _ := rhs_val.(*BoolVal)
		return &BoolVal{Val: lhs_val.Val == rhs_val.Val}
	case *StringVal:
		rhs_val, _ := rhs_val.(*StringVal)
		return &BoolVal{Val: lhs_val.Val == rhs_val.Val}
	case *NoneVal:
		return &BoolVal{Val: true}
	default:
		panic(fmt.Sprintf("Line %d: Unsupported '==' between %s and %s", node.Line_num, lhs_val, rhs_val))
	}
}

func (interpreter *ImpInterpreter) eval_NeqExpr(node NeqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.Lhs)
	rhs_val := interpreter.eval_Expr(node.Rhs)
	if !check_vals_type_equal(lhs_val, rhs_val) {
		panic(fmt.Sprintf("Line %d: Unsupported '!=' between %s and %s", node.Line_num, lhs_val, rhs_val))
	}
	switch lhs_val := lhs_val.(type) {
	case *IntVal:
		rhs_val, _ := rhs_val.(*IntVal)
		return &BoolVal{Val: lhs_val.Val != rhs_val.Val}
	case *BoolVal:
		rhs_val, _ := rhs_val.(*BoolVal)
		return &BoolVal{Val: lhs_val.Val != rhs_val.Val}
	case *StringVal:
		rhs_val, _ := rhs_val.(*StringVal)
		return &BoolVal{Val: lhs_val.Val != rhs_val.Val}
	case *NoneVal:
		return &BoolVal{Val: false}
	default:
		panic(fmt.Sprintf("Line %d: Unsupported '!=' between %s and %s", node.Line_num, lhs_val, rhs_val))
	}
}

func (interpreter *ImpInterpreter) eval_LessthanExpr(node LessthanExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.Lhs)
	rhs_val := interpreter.eval_Expr(node.Rhs)
	lhs_intvar, lhs_is_int := lhs_val.(*IntVal)
	rhs_intvar, rhs_is_int := rhs_val.(*IntVal)
	if !(lhs_is_int && rhs_is_int) {
		panic(fmt.Sprintf("Line %d: Lessthan operator must be applied between two integer values", node.Line_num))
	}
	return &BoolVal{Val: lhs_intvar.Val < rhs_intvar.Val}
}

func (interpreter *ImpInterpreter) eval_GreaterthanExpr(node GreaterthanExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.Lhs)
	rhs_val := interpreter.eval_Expr(node.Rhs)
	lhs_intvar, lhs_is_int := lhs_val.(*IntVal)
	rhs_intvar, rhs_is_int := rhs_val.(*IntVal)
	if !(lhs_is_int && rhs_is_int) {
		panic(fmt.Sprintf("Line %d: Greaterthan operator must be applied between two integer values", node.Line_num))
	}
	return &BoolVal{Val: lhs_intvar.Val > rhs_intvar.Val}
}

func (interpreter *ImpInterpreter) eval_LeqExpr(node LeqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.Lhs)
	rhs_val := interpreter.eval_Expr(node.Rhs)
	lhs_intvar, lhs_is_int := lhs_val.(*IntVal)
	rhs_intvar, rhs_is_int := rhs_val.(*IntVal)
	if !(lhs_is_int && rhs_is_int) {
		panic(fmt.Sprintf("Line %d: Leq operator must be applied between two integer values", node.Line_num))
	}
	return &BoolVal{Val: lhs_intvar.Val <= rhs_intvar.Val}
}

func (interpreter *ImpInterpreter) eval_GeqExpr(node GeqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.Lhs)
	rhs_val := interpreter.eval_Expr(node.Rhs)
	lhs_intvar, lhs_is_int := lhs_val.(*IntVal)
	rhs_intvar, rhs_is_int := rhs_val.(*IntVal)
	if !(lhs_is_int && rhs_is_int) {
		panic(fmt.Sprintf("Line %d: Geq operator must be applied between two integer values", node.Line_num))
	}
	return &BoolVal{Val: lhs_intvar.Val >= rhs_intvar.Val}
}

func (interpreter *ImpInterpreter) eval_NegExpr(node NegExpr) ImpValues {
	subexpr_val, subexpr_is_int := interpreter.eval_Expr(node.Subexpr).(*IntVal)
	if !subexpr_is_int {
		panic(fmt.Sprintf("Line %d: Subexpr %s of Unary neg operator should be of type int", node.Line_num, node.Subexpr))
	}
	return &IntVal{Val: -subexpr_val.Val}
}

func (interpreter *ImpInterpreter) eval_NotExpr(node NotExpr) ImpValues {
	subexpr_val, subexpr_is_bool := interpreter.eval_Expr(node.Subexpr).(*BoolVal)
	if !subexpr_is_bool {
		panic(fmt.Sprintf("Line %d: Subexpr %s of NOT operator should be of type bool", node.Line_num, node.Subexpr))
	}
	return &BoolVal{Val: !subexpr_val.Val}
}

func (interpreter *ImpInterpreter) eval_AndExpr(node AndExpr) ImpValues {
	lhs_val, lhs_is_bool := interpreter.eval_Expr(node.Lhs).(*BoolVal)
	rhs_val, rhs_is_bool := interpreter.eval_Expr(node.Rhs).(*BoolVal)

	if !lhs_is_bool {
		panic(fmt.Sprintf("Line %d: LHS of AND should be a bool value, but got '%s'", node.Line_num, node.Lhs))
	}

	if !rhs_is_bool {
		panic(fmt.Sprintf("Line %d: RHS of AND should be a bool value, but got '%s'", node.Line_num, node.Rhs))
	}
	return &BoolVal{Val: lhs_val.Val && rhs_val.Val}
}

func (interpreter *ImpInterpreter) eval_OrExpr(node OrExpr) ImpValues {
	lhs_val, lhs_is_bool := interpreter.eval_Expr(node.Lhs).(*BoolVal)
	rhs_val, rhs_is_bool := interpreter.eval_Expr(node.Rhs).(*BoolVal)

	if !lhs_is_bool {
		panic(fmt.Sprintf("Line %d: LHS of OR should be a bool value, but got '%s'", node.Line_num, node.Lhs))
	}

	if !rhs_is_bool {
		panic(fmt.Sprintf("Line %d: RHS of OR should be a bool value, but got '%s'", node.Line_num, node.Rhs))
	}
	return &BoolVal{Val: lhs_val.Val || rhs_val.Val}
}

// Imp is pass-by-value for int/bool, but arrays are passed references
func (interpreter *ImpInterpreter) eval_function_call(func_name ImpFunctionName, args []Expr, line_num int) ImpValues {
	// copy values if primitive
	prepare_args := func(arg ImpValues) ImpValues {
		switch arg_ty := arg.(type) {
		case *IntVal:
			return &IntVal{Val: arg_ty.Val}
		case *BoolVal:
			return &BoolVal{Val: arg_ty.Val}
		case *ArrayVal:
			return arg
		}
		panic(fmt.Sprintf("Line %d: Unknown arg type '%s' for function %s\n", line_num, get_type(arg), func_name))
	}
	func_local_state := ImpState{vars: make(map[string]ImpValues), current_function_name: func_name, return_value: &NoneVal{}}
	imp_function, function_exists := interpreter.Functions[func_name]
	if !function_exists {
		panic(fmt.Sprintf("Line %d: Unknown function '%s'\n", line_num, func_name))
	}
	for index, arg_expr := range args {
		arg_info := imp_function.Arg_pairs[index]
		arg_val := prepare_args(interpreter.eval_Expr(arg_expr))
		if !check_val_type_match(arg_val, arg_info.Arg_type) {
			panic(fmt.Sprintf("Line %d: Argument '%s' of function '%s' is defined as type %s, but passed expr '%s' of type %s", line_num, arg_info.Name, func_name, arg_info.Arg_type, arg_expr, get_type(arg_val)))
		}
		func_local_state.vars[arg_info.Name] = arg_val
	}
	interpreter.push_state(func_local_state)
	interpreter.eval_Stmt(imp_function.Body)
	return_value := interpreter.get_top_state().return_value
	interpreter.pop_state()

	return return_value
}

func (interpreter *ImpInterpreter) eval_CallExpr(node CallExpr) ImpValues {
	return interpreter.eval_function_call(node.Func_name, node.Args, node.Line_num)
}

func (interpreter *ImpInterpreter) eval_MakeArrayExpr(node MakeArrayExpr) ImpValues {
	len_node := interpreter.eval_Expr(node.Size)
	len_intval, len_is_int := len_node.(*IntVal)
	if !len_is_int {
		panic(fmt.Sprintf("Line %d: %s - length expression %s is not an integer value", node.Line_num, node, node.Size))
	}
	default_val := interpreter.eval_Expr(node.Value)
	generated := make([]ImpValues, len_intval.Val)
	for i := 0; i < len_intval.Val; i++ {
		generated[i] = interpreter.deepcopy_impvalue(default_val)
	}
	return &ArrayVal{get_type(default_val), generated}
}

func (interpreter *ImpInterpreter) eval_LenExpr(node LenExpr) ImpValues {
	array_node := interpreter.eval_Expr(node.Subexpr)
	array_val, is_array := array_node.(*ArrayVal)
	if !is_array {
		panic(fmt.Sprintf("Line %d: len() - Non-array value %s passed to len()", node.Line_num, node.Subexpr))
	}
	return &IntVal{Val: len(array_val.Val)}
}

func (interpreter *ImpInterpreter) eval_Expr(node Expr) ImpValues {
	switch node_ty := node.(type) {
	case *VarExpr:
		return interpreter.eval_VarExpr(*node_ty)
	case *IntLitExpr:
		return interpreter.eval_IntLitExpr(*node_ty)
	case *BoolLitExpr:
		return interpreter.eval_BoolLitExpr(*node_ty)
	case *StringLitExpr:
		return interpreter.eval_StringLitExpr(*node_ty)
	case *ArrayLitExpr:
		return interpreter.eval_ArrayLitExpr(*node_ty)
	case *AddExpr:
		return interpreter.eval_AddExpr(*node_ty)
	case *SubExpr:
		return interpreter.eval_SubExpr(*node_ty)
	case *MulExpr:
		return interpreter.eval_MulExpr(*node_ty)
	case *DivExpr:
		return interpreter.eval_DivExpr(*node_ty)
	case *ModExpr:
		return interpreter.eval_ModExpr(*node_ty)
	case *ParenExpr:
		return interpreter.eval_ParenExpr(*node_ty)
	case *ArrayIndexExpr:
		return interpreter.eval_ArrayIndexExpr(*node_ty)
	case *EqExpr:
		return interpreter.eval_EqExpr(*node_ty)
	case *NeqExpr:
		return interpreter.eval_NeqExpr(*node_ty)
	case *LessthanExpr:
		return interpreter.eval_LessthanExpr(*node_ty)
	case *GreaterthanExpr:
		return interpreter.eval_GreaterthanExpr(*node_ty)
	case *LeqExpr:
		return interpreter.eval_LeqExpr(*node_ty)
	case *GeqExpr:
		return interpreter.eval_GeqExpr(*node_ty)
	case *NegExpr:
		return interpreter.eval_NegExpr(*node_ty)
	case *NotExpr:
		return interpreter.eval_NotExpr(*node_ty)
	case *AndExpr:
		return interpreter.eval_AndExpr(*node_ty)
	case *OrExpr:
		return interpreter.eval_OrExpr(*node_ty)
	case *CallExpr:
		return interpreter.eval_CallExpr(*node_ty)
	case *MakeArrayExpr:
		return interpreter.eval_MakeArrayExpr(*node_ty)
	case *LenExpr:
		return interpreter.eval_LenExpr(*node_ty)
	default:
		panic(fmt.Sprintf(" Unimplemented expr type %s", node))
	}
}

/////////////////////////////////
// statements

func (interpreter *ImpInterpreter) eval_SkipStmt(SkipStmt) ControlflowResult {
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_AssignStmt(node AssignStmt) ControlflowResult {
	rhs_val := interpreter.eval_Expr(node.Rhs)
	switch lhs_loc := interpreter.eval_Expr_lvalue(node.Lhs, get_type(rhs_val)).(type) {
	case *IntVal:
		rhs_intval, rhs_is_intval := rhs_val.(*IntVal)
		if !rhs_is_intval {
			panic(fmt.Sprintf("Line %d: Attempted to assign RHS '%s' of type %s to LHS '%s' of type %s", node.Line_num, node.Rhs, get_type(rhs_val), node.Lhs, get_type(lhs_loc)))
		}
		lhs_loc.Val = rhs_intval.Val
	case *BoolVal:
		rhs_boolval, rhs_is_boolval := rhs_val.(*BoolVal)
		if !rhs_is_boolval {
			panic(fmt.Sprintf("Line %d: Attempted to assign RHS '%s' of type %s to LHS '%s' of type %s", node.Line_num, node.Rhs, get_type(rhs_val), node.Lhs, get_type(lhs_loc)))
		}
		lhs_loc.Val = rhs_boolval.Val
	case *ArrayVal:
		rhs_arrval, rhs_is_arrayval := rhs_val.(*ArrayVal)
		if !rhs_is_arrayval {
			panic(fmt.Sprintf("Line %d: Attempted to assign RHS '%s' of type %s to LHS '%s' of type %s", node.Line_num, node.Rhs, get_type(rhs_val), node.Lhs, get_type(lhs_loc)))
		}
		lhs_loc.Val = rhs_arrval.Val
	default:
		panic(fmt.Sprintf("Line %d: LHS expr '%s' has unresolved value type %T\n", node.Line_num, node.Lhs, lhs_loc))
	}
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_IfElseStmt(node IfElseStmt) ControlflowResult {
	cond_val := interpreter.eval_Expr(node.Cond)
	cond_boolval, cond_is_bool := cond_val.(*BoolVal)
	if !cond_is_bool {
		panic(fmt.Sprintf("Line %d: If statement got non-boolean condition '%s'\n", node.Line_num, node.Cond))
	}
	if cond_boolval.Val {
		return interpreter.eval_Stmt(node.True_stmt)
	} else {
		return interpreter.eval_Stmt(node.False_stmt)
	}
}

func (interpreter *ImpInterpreter) eval_WhileStmt(node WhileStmt) ControlflowResult {
	for true {
		cond_val := interpreter.eval_Expr(node.Cond)
		cond_boolval, cond_is_bool := cond_val.(*BoolVal)
		if !cond_is_bool {
			panic(fmt.Sprintf("Line %d: While statement got non-boolean condition '%s'\n", node.Line_num, node.Cond))
		}
		if cond_boolval.Val == false {
			break
		}
		stmt_result := interpreter.eval_Stmt(node.Body_stmt)
		switch stmt_result {
		case ControlBreak:
			return ControlNormal
		case ControlContinue:
			continue
		case ControlReturn:
			return ControlReturn
		}
	}
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_BreakStmt(_ BreakStmt) ControlflowResult {
	return ControlBreak
}

func (interpreter *ImpInterpreter) eval_ContinueStmt(_ ContinueStmt) ControlflowResult {
	return ControlContinue
}

func (interpreter *ImpInterpreter) eval_IncStmt(node IncStmt) ControlflowResult {
	lhs_val_int, lhs_is_int := interpreter.eval_Expr(node.Subexpr).(*IntVal)
	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: Attempted to increment non-integer value '%s'\n", node.Line_num, node))
	}
	lhs_val_int.Val++
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_DecStmt(node DecStmt) ControlflowResult {
	lhs_val_int, lhs_is_int := interpreter.eval_Expr(node.Subexpr).(*IntVal)
	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: Attempted to decrement non-integer value '%s'\n", node.Line_num, node))
	}
	lhs_val_int.Val--
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_CallStmt(node CallStmt) ControlflowResult {
	interpreter.eval_function_call(node.Func_name, node.Args, node.Line_num)
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_PrintStmt(node PrintStmt) ControlflowResult {
	var outputs []string
	for _, arg := range node.Args {
		arg_val := interpreter.eval_Expr(arg)
		outputs = append(outputs, fmt.Sprintf("%s", arg_val))
	}
	fmt.Print(strings.ReplaceAll(strings.Join(outputs, " "), "\\n", "\n"))
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_ScanfStmt(node ScanfStmt) ControlflowResult {
	for index, fmt_str := range strings.Split(node.Format_string, " ") {
		var imp_val ImpValues
		switch fmt_str {
		case "%d":
			var input int
			fmt.Scanf(fmt_str, &input)
			imp_val = &IntVal{Val: input}
		case "%t":
			var input bool
			fmt.Scanf(fmt_str, &input)
			imp_val = &BoolVal{Val: input}
		default:
			panic(fmt.Sprintf("Line %d: scanf - Unsupported formatting specifier %s\n", node.Line_num, fmt_str))
		}

		lhs_val := node.Assign_locations[index]
		switch lhs_loc := interpreter.eval_Expr_lvalue(lhs_val, get_type(imp_val)).(type) {
		case *IntVal:
			rhs_intval, rhs_is_intval := imp_val.(*IntVal)
			if !rhs_is_intval {
				panic(fmt.Sprintf("Line %d: scanf - Attempted to assign input '%s' of type %T to variable '%s' of type %T\n", node.Line_num, imp_val, imp_val, lhs_val, lhs_loc))
			}
			lhs_loc.Val = rhs_intval.Val
		case *BoolVal:
			rhs_intval, rhs_is_boolval := imp_val.(*BoolVal)
			if !rhs_is_boolval {
				panic(fmt.Sprintf("Line %d: scanf - Attempted to assign input '%s' of type %T to variable '%s' of type %T\n", node.Line_num, imp_val, imp_val, lhs_val, lhs_loc))
			}
			lhs_loc.Val = rhs_intval.Val
		}
	}
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_ReturnStmt(node ReturnStmt) ControlflowResult {
	top_state := interpreter.get_top_state()
	top_state.return_value = interpreter.eval_Expr(node.Arg)
	if top_state.current_function_name != "" {
		expected_return_type := interpreter.Functions[top_state.current_function_name].Return_type
		if !check_val_type_match(top_state.return_value, expected_return_type) {
			panic(fmt.Sprintf("Line %d: Function %s should return value of type %s, but actually returned '%s' of type %s\n", node.Line_num, top_state.current_function_name, expected_return_type, node.Arg, top_state.return_value))
		}
	}
	return ControlReturn
}

// eval_Stmt evaluates a sequence of statements
// The ControlflowResult return type designates whether execution of statements should be affected
func (interpreter *ImpInterpreter) eval_Stmt(nodes []Stmt) ControlflowResult {
	var control_status ControlflowResult = ControlNormal
	for _, stmt := range nodes {
		switch stmt := stmt.(type) {
		case *SkipStmt:
			control_status = interpreter.eval_SkipStmt(*stmt)
		case *AssignStmt:
			control_status = interpreter.eval_AssignStmt(*stmt)
		case *IfElseStmt:
			control_status = interpreter.eval_IfElseStmt(*stmt)
		case *WhileStmt:
			control_status = interpreter.eval_WhileStmt(*stmt)
		case *BreakStmt:
			control_status = interpreter.eval_BreakStmt(*stmt)
		case *ContinueStmt:
			control_status = interpreter.eval_ContinueStmt(*stmt)
		case *IncStmt:
			control_status = interpreter.eval_IncStmt(*stmt)
		case *DecStmt:
			control_status = interpreter.eval_DecStmt(*stmt)
		case *CallStmt:
			control_status = interpreter.eval_CallStmt(*stmt)
		case *PrintStmt:
			control_status = interpreter.eval_PrintStmt(*stmt)
		case *ScanfStmt:
			control_status = interpreter.eval_ScanfStmt(*stmt)
		case *ReturnStmt:
			control_status = interpreter.eval_ReturnStmt(*stmt)
		}
		if control_status != ControlNormal {
			return control_status
		}
	}
	return control_status
}
