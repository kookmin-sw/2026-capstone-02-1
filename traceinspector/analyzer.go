package traceinspector

import (
	"fmt"
	"strings"
	"traceinspector/algebra"
	"traceinspector/domain"
	"traceinspector/imp"
)

type AnalysisSettings struct {
	Loop_iters_before_widening int // number of loop interations before widening
	Max_call_stack_depth       int // number of nested function calls before returning top
}

// An AbstractState is the pair (l, M^#) ↪ (l', M^#') used in the abstract transition relation
// node_location: node location to be interpreted
// abstract_mem: the input abstract memory state wrt the node should be interpreted
// remaining_call_depth: the number of nested function calls allowed before widening to top. Negative number means no limit.
type AbstractState[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl ArrayDomain[IntDomainImpl, ArrayDomainImpl]] struct {
	node_location        CFGNodeLocation
	abstract_mem         AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]
	remaining_call_depth int
}

func (state AbstractState[IntDomainImpl, ArrayDomainImpl]) String() string {
	var ret []string
	for key, val := range state.abstract_mem {
		ret = append(ret, fmt.Sprintf("%s : %s", key, val))
	}
	return fmt.Sprintf("%s - {%s}", state.node_location, strings.Join(ret, ", "))
}

func (state AbstractState[IntDomainImpl, ArrayDomainImpl]) Clone() AbstractState[IntDomainImpl, ArrayDomainImpl] {
	new_st := AbstractState[IntDomainImpl, ArrayDomainImpl]{}
	new_st.node_location = state.node_location
	new_st.abstract_mem = state.abstract_mem.Clone()
	new_st.remaining_call_depth = state.remaining_call_depth
	return new_st
}

/////////////////

// TODO: fix this part and integrate with get_abstract_value_from_expr
// Compute the abstract value of an expression expr under an abstract memory state abs_mem
func (interpreter *AbstractAnalyzer[IntDomainImpl, ArrayDomainImpl]) Eval_expr(abs_state AbstractState[IntDomainImpl, ArrayDomainImpl], expr imp.Expr) AbstractValue[IntDomainImpl, ArrayDomainImpl] {
	switch expr_ty := expr.(type) {
	case *imp.VarExpr:
		return abs_state.abstract_mem[expr_ty.Name]
	case *imp.IntLitExpr:
		intdom_result := interpreter.Intdomain_default.From_IntLitExpr(*expr_ty)
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: intdom_result}
	case *imp.BoolLitExpr:
		booldom_result := domain.BoolDomain{}.From_BoolLitExpr(*expr_ty)
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: booldom_result}
	case *imp.StringLitExpr:
		// TODO do something about strings
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{}
	case *imp.NegExpr:
		subexpr_val := interpreter.Eval_expr(abs_state, expr_ty.Subexpr)
		if subexpr_val.domain_kind != IntDomainKind {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("Result of arithmetic negation returned %s instead if IntDomain", subexpr_val.domain_kind))
		}
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: subexpr_val.Get_int().Neg()}
	case *imp.AddExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			interpreter.output_handler.write_error(abs_state.node_location, "Add expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Add(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.SubExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			interpreter.output_handler.write_error(abs_state.node_location, "Sub expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Sub(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.MulExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			interpreter.output_handler.write_error(abs_state.node_location, "Mul expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Mul(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.DivExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			interpreter.output_handler.write_error(abs_state.node_location, "Div expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Div(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.ModExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			interpreter.output_handler.write_error(abs_state.node_location, "Add expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Mod(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.ParenExpr:
		return interpreter.Eval_expr(abs_state, expr_ty.Subexpr)
	case *imp.ArrayIndexExpr:
		arr_val := interpreter.Eval_expr(abs_state, expr_ty.Base)
		index_val := interpreter.Eval_expr(abs_state, expr_ty.Index)
		if arr_val.domain_kind != ArrayDomainKind {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("'%s' : expected arr to have arr domain type", expr_ty))
		}
		if !index_val.Get_int().Leq(arr_val.Get_array().Len().Sub(index_val.int_domain.From_IntLitExpr(imp.IntLitExpr{Value: 1}))).IsTrue() {
			interpreter.output_handler.write_warning(abs_state.node_location, fmt.Sprintf("Potentially unsafe array indexing: index has value %s, but %s.Len has value %s.", index_val.Get_int(), get_varname_from_lvalue(expr_ty.Base), arr_val.Get_array().Len()))
		}
		result_val := arr_val.Get_array().GetIndex(index_val.Get_int())
		return result_val
	case *imp.LenExpr:
		arr_val := interpreter.Eval_expr(abs_state, expr_ty.Subexpr)
		if arr_val.domain_kind != ArrayDomainKind {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("'%s' : expected arr to have arr domain type", expr_ty))
		}
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: arr_val.Get_array().Len()}
	case *imp.MakeArrayExpr:
		default_val := interpreter.Eval_expr(abs_state, expr_ty.Value)
		size_val := interpreter.Eval_expr(abs_state, expr_ty.Size)
		if size_val.domain_kind != IntDomainKind {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("Size argument '%s' in make_array expected to have int domain type, but has type %s", expr_ty.Size, size_val.domain_kind))
		}
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: ArrayDomainKind, array_domain: interpreter.Arraydomain_default.Make_array(size_val.Get_int(), default_val)}
	case *imp.EqExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if lhs_val.domain_kind != rhs_val.domain_kind {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different", expr_ty))
		}
		result_val := lhs_val.Get_int().Eq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.NeqExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if lhs_val.domain_kind != rhs_val.domain_kind {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different", expr_ty))
		}
		result_val := lhs_val.Get_int().Neq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.LeqExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if !(lhs_val.domain_kind == IntDomainKind && lhs_val.domain_kind == rhs_val.domain_kind) {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different (%s vs %s)", expr_ty, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		result_val := lhs_val.Get_int().Leq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.GeqExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if !(lhs_val.domain_kind == IntDomainKind && lhs_val.domain_kind == rhs_val.domain_kind) {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different (%s vs %s)", expr_ty, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		result_val := lhs_val.Get_int().Geq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.LessthanExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if !(lhs_val.domain_kind == IntDomainKind && lhs_val.domain_kind == rhs_val.domain_kind) {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different (%s vs %s)", expr_ty, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		result_val := lhs_val.Get_int().Lessthan(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.GreaterthanExpr:
		lhs_val := interpreter.Eval_expr(abs_state, expr_ty.Lhs)
		rhs_val := interpreter.Eval_expr(abs_state, expr_ty.Rhs)
		if !(lhs_val.domain_kind == IntDomainKind && lhs_val.domain_kind == rhs_val.domain_kind) {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different (%s vs %s)", expr_ty, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		result_val := lhs_val.Get_int().Greaterthan(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.CallExpr:
		func_info, func_exists := interpreter.Function_defs[expr_ty.Func_name]
		if !(func_exists) {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("Function %s does not exist", expr_ty.Func_name))
			panic("unreachable")
		}
		if len(func_info.Arg_pairs) != len(expr_ty.Args) {
			interpreter.output_handler.write_error(abs_state.node_location, fmt.Sprintf("Function %s expectes %d arguments, but got %d", expr_ty.Func_name, len(func_info.Arg_pairs), len(expr_ty.Args)))
			panic("unreachable")
		}
		function_entry_varmem := make(AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl])
		for arg_index, arg_info := range func_info.Arg_pairs {
			// make sure nested function calls are checked against the limit
			subexpr_state := abs_state.Clone()
			subexpr_state.remaining_call_depth = abs_state.remaining_call_depth - 1
			arg_val := interpreter.Eval_expr(subexpr_state, expr_ty.Args[arg_index])
			function_entry_varmem[arg_info.Name] = arg_val
		}
		if abs_state.remaining_call_depth == 0 {
			interpreter.output_handler.write_warning(abs_state.node_location, fmt.Sprintf("Maximum function call depth(%d) reached. The call to %s will not be interpreted, but assumed to be ⊤.", interpreter.Settings.Max_call_stack_depth, expr_ty.Func_name))
			switch func_info.Return_type.(type) {
			case imp.IntType:
				return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind}.Make_top()
			case imp.BoolType:
				return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind}.Make_top()
			case imp.ArrayType:
				return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: ArrayDomainKind}.Make_top()
			default:
				return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: InvalidKind}
			}
		}
		// return the call value
		return interpreter.Interpret_function(func_info, function_entry_varmem, abs_state.remaining_call_depth-1)

	}
	return AbstractValue[IntDomainImpl, ArrayDomainImpl]{}
}

func get_varname_from_lvalue(expr imp.Expr) string {
	switch expr_ty := expr.(type) {
	case *imp.VarExpr:
		return expr_ty.Name
	case *imp.ArrayIndexExpr:
		return get_varname_from_lvalue(expr_ty.Base)
	}
	panic(fmt.Sprintf("get_varname_from_lvalue: unimplemented expr type %T", expr))
}

// TODO: fix this part
// Given an lvalue expression, get the AbstractValue from the state
func (interpreter *AbstractAnalyzer[IntDomainImpl, ArrayDomainImpl]) get_abstract_value_from_lvalue_expr(expr imp.Expr, state AbstractState[IntDomainImpl, ArrayDomainImpl]) AbstractValue[IntDomainImpl, ArrayDomainImpl] {
	switch expr_ty := expr.(type) {
	case *imp.VarExpr:
		val, var_in_mem := state.abstract_mem[expr_ty.Name]
		if !var_in_mem {
			interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("Variable '%s' not defined", expr_ty.Name))
		}
		return val
	case *imp.LenExpr:
		arr_val := interpreter.get_abstract_value_from_lvalue_expr(expr_ty.Subexpr, state)
		if arr_val.domain_kind != ArrayDomainKind {
			interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("Called len() on a non-array value '%s'", expr_ty.Subexpr))
		}
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: arr_val.Get_array().Len()}
	case *imp.ArrayIndexExpr:
		arr_val := interpreter.get_abstract_value_from_lvalue_expr(expr_ty.Base, state)
		if arr_val.domain_kind != ArrayDomainKind {
			interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("Tried to index a non-array value '%s'", expr_ty.Base))
		}
		index_val := interpreter.Eval_expr(state, expr_ty.Index)
		if index_val.domain_kind != ArrayDomainKind {
			interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("Index expression '%s' is not an integerdomain", expr_ty.Index))
		}
		return arr_val.Get_array().GetIndex(index_val.Get_int())
	}
	interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("get_abstract_value_from_expr: Unsupported expression '%s'", expr))
	panic("unreachable")
}

// TODO: fix this part
// Given an lvalue expression lhs, the rhs AbstractValue, and the state, update state.abstract_mem accordingly *in-place*
func (interpreter *AbstractAnalyzer[IntDomainImpl, ArrayDomainImpl]) set_abstract_value_from_expr(lhs imp.Expr, rhs_val AbstractValue[IntDomainImpl, ArrayDomainImpl], state *AbstractState[IntDomainImpl, ArrayDomainImpl]) {
	switch lhs_node := lhs.(type) {
	case *imp.VarExpr:
		lhs_val, lhs_exists := state.abstract_mem[lhs_node.Name]
		if lhs_exists && lhs_val.domain_kind != rhs_val.domain_kind {
			interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("LHS of variable %s has type %s, but RHS has type %s.", lhs_node.Name, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		state.abstract_mem[lhs_node.Name] = rhs_val
	case *imp.ArrayIndexExpr:
		index_val := interpreter.Eval_expr(*state, lhs_node.Index)
		if index_val.domain_kind != IntDomainKind {
			interpreter.output_handler.write_error(state.node_location, "Only 1-dimensional arrays implemented")
		}
		if index_val.Get_int().IsBot() {
			// if index is bot, don't do assignment
			return
		}
		arr_varname := get_varname_from_lvalue(lhs_node.Base)
		lhs_val, lhs_exists := state.abstract_mem[arr_varname]
		if !lhs_exists {
			interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("Attempting to index nonexisting variable '%s'", arr_varname))
		}
		if lhs_val.domain_kind != ArrayDomainKind {
			interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("Attempting to index non-array variable '%s'", arr_varname))
		}
		// fmt.Println(index_val.Get_int(), "<=", lhs_val.Get_array().Len().Sub(index_val.int_domain.From_IntLitExpr(imp.IntLitExpr{Value: 1})), "=", index_val.Get_int().Leq(lhs_val.Get_array().Len().Sub(index_val.int_domain.From_IntLitExpr(imp.IntLitExpr{Value: 1}))))
		if !index_val.Get_int().Leq(lhs_val.Get_array().Len().Sub(index_val.int_domain.From_IntLitExpr(imp.IntLitExpr{Value: 1}))).IsTrue() {
			interpreter.output_handler.write_warning(state.node_location, fmt.Sprintf("Potentially unsafe array indexing: index has value %s, but %s.Len has value %s.", index_val.Get_int(), arr_varname, lhs_val.Get_array().Len()))
		}
		state.abstract_mem[arr_varname] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: ArrayDomainKind, array_domain: lhs_val.Get_array().SetVal(index_val.Get_int(), rhs_val)}
	case *imp.LenExpr:
		arr_varname := get_varname_from_lvalue(lhs_node.Subexpr)
		lhs_val, lhs_exists := state.abstract_mem[arr_varname]
		if !lhs_exists {
			interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("Attempting to get length of nonexisting variable '%s'", arr_varname))
		}
		if lhs_val.domain_kind != ArrayDomainKind {
			interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("Attempting to get length of non-array variable '%s'", arr_varname))
		}
		if rhs_val.domain_kind != IntDomainKind {
			interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("subexpr of len is not an IntDomainKind, but %s", rhs_val.domain_kind))
		}
		state.abstract_mem[arr_varname] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: ArrayDomainKind, array_domain: lhs_val.Get_array().SetLen(rhs_val.Get_int())}
	default:
		interpreter.output_handler.write_error(state.node_location, fmt.Sprintf("set_abstract_value_from_expr: Unsupported LHS expr '%s' with type %T", lhs, lhs))
	}
}

func (interpreter *AbstractAnalyzer[IntDomainImpl, ArrayDomainImpl]) Step(in_state AbstractState[IntDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, ArrayDomainImpl] {
	cfg_node, cfg_node_exists := interpreter.Function_cfgs[in_state.node_location.Function_name].Node_map[in_state.node_location.Id]
	if !cfg_node_exists {
		interpreter.output_handler.write_error(create_empty_node_location(), fmt.Sprintf("The designated CFG Node %s doesn't exist", in_state.node_location))
	}

	func_name := in_state.node_location.Function_name
	global_state, _ := interpreter.function_pre_mem_map[func_name].pre_mem_node_map[in_state.node_location.Id] // the global state held by the interpreter

	cond_node, is_cond_node := cfg_node.(*CFGCondNode)
	loop_widened := false
	if is_cond_node && cond_node.Is_loop_head && interpreter.function_pre_mem_map[func_name].n_visits[in_state.node_location.Id] > interpreter.Settings.Loop_iters_before_widening {
		// Apply widening if visit count is greater than threshold
		interpreter.function_pre_mem_map[func_name].n_visits[in_state.node_location.Id] = 1
		global_state.Widen_inplace(in_state.abstract_mem)
		interpreter.output_handler.write_update_node_state(in_state.node_location, global_state.String(), "Widen global memory state")
		loop_widened = true
	} else {
		// When we receive a new pair, update the global state with its join
		state_changed := global_state.Join_inplace(in_state.abstract_mem, in_state.node_location)
		if !state_changed && interpreter.function_pre_mem_map[func_name].n_visits[in_state.node_location.Id] > 0 {
			// no updates to the state
			interpreter.output_handler.write_info(in_state.node_location, "No updates to node state")
			return nil
		}
		interpreter.output_handler.write_update_node_state(in_state.node_location, global_state.String(), "Join global memory state")
	}

	interpreter.function_pre_mem_map[func_name].n_visits[in_state.node_location.Id]++
	// Executed on the joined node state
	in_state.abstract_mem = global_state.Clone()

	var return_states []AbstractState[IntDomainImpl, ArrayDomainImpl]
	switch cfg_node := cfg_node.(type) {
	case *CFGNode:
		switch stmt := cfg_node.Ast.(type) {
		case *imp.AssignStmt:
			// assignment should overwrite the value, instead of join
			rhs_val := interpreter.Eval_expr(in_state, stmt.Rhs)
			interpreter.set_abstract_value_from_expr(stmt.Lhs, rhs_val, &in_state)

		case *imp.SkipStmt:
			// do nothing
		case *imp.BreakStmt:
			// do nothing; follows CFG edge
		case *imp.ContinueStmt:
			// do nothing; follows CFG edge
		case *imp.PrintStmt:
			// eval subexprs and follow CFG edge
			for _, val := range stmt.Args {
				interpreter.Eval_expr(in_state, val)
			}
		case *imp.ScanfStmt:
			// by default scanf initializes variables to top
			if len(strings.Split(stmt.Format_string, " ")) != len(stmt.Assign_locations) {
				interpreter.output_handler.write_error(in_state.node_location, "Scanf must have matching amount of format specifiers and assignment locations")
			}
			for index, fmt_str := range strings.Split(stmt.Format_string, " ") {
				switch fmt_str {
				case "%d":
					// integer type
					default_val := AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind}.Make_top()
					interpreter.set_abstract_value_from_expr(stmt.Assign_locations[index], default_val, &in_state)
				case "%t":
					default_val := AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind}.Make_top()
					interpreter.set_abstract_value_from_expr(stmt.Assign_locations[index], default_val, &in_state)
				default:
					interpreter.output_handler.write_error(in_state.node_location, fmt.Sprintf("Unknown scanf format specifier '%s'. only %%d for integers and %%t for bools allowed", fmt_str))
				}
			}
		case *imp.CallStmt:
			// the processing routine is the same as Eval_expr for CallExpr. But we just ignore the return result.
			_ = interpreter.Eval_expr(in_state, &imp.CallExpr{Node: stmt.Node, Func_name: stmt.Func_name, Args: stmt.Args})
		case *imp.ReturnStmt:
			val := interpreter.Eval_expr(in_state, stmt.Arg)
			joined_val, _ := interpreter.function_pre_mem_map[func_name].return_value.Join(val, in_state.node_location)
			interpreter.function_pre_mem_map[func_name].return_value = joined_val
		default:
			panic(fmt.Sprintf("unimplemented stmt %T", stmt))
		}

	case *CFGCondNode:
		cond_edge, is_cond_edge := interpreter.Function_cfgs[func_name].Edge_map_from[in_state.node_location.Id].(*CFGCondEdge)
		if !is_cond_edge {
			interpreter.output_handler.write_error(in_state.node_location, "Condition stmt does not have outgoing edge of CondEdge type.")
		}
		// If at in_state the prop evaluates to either true or false,
		// We can just execute only the corresponding branch.
		// Otherwise filter for each branch and join the result
		cond_val := interpreter.Eval_expr(in_state, cfg_node.Cond_expr)
		// Try to represent it as SimpleProp
		cond_simpleprop, simpleprop_success := algebra.Imp_expr_to_simple_prop(cfg_node.Cond_expr)
		if !simpleprop_success {
			interpreter.output_handler.write_warning(in_state.node_location, fmt.Sprintf("Could not represent '%s' as SimpleProp. Analysis precision may severely deterioriate.", cfg_node.Cond_expr))
		}
		if cond_val.Get_bool().IsBot() { // dead branch
			return nil
		}

		if (cond_val.Get_bool().IsTrue() || cond_val.Get_bool().IsTop()) && cond_edge.To_true_node_loc.NodeExists() {
			// run just the true branch on filter_true(in_state)
			new_state := in_state.Clone()
			if simpleprop_success {
				true_filters := domain.Filter_true_query_simpleprop(cond_simpleprop)
				// TODO: This doesn't refine array lengths
				for _, filter := range true_filters {
					rhs_dom_val := interpreter.Eval_expr(in_state, filter.Rhs_expr)
					if rhs_dom_val.domain_kind == IntDomainKind {
						updated_intdom := interpreter.get_abstract_value_from_lvalue_expr(filter.Term_expr, new_state).Get_int().Filter(filter.Query_type, rhs_dom_val.Get_int())
						updated_val := AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: updated_intdom}
						interpreter.set_abstract_value_from_expr(filter.Term_expr, updated_val, &new_state)
					}
				}
			}
			return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: cond_edge.To_true_node_loc, abstract_mem: new_state.abstract_mem, remaining_call_depth: in_state.remaining_call_depth})
		}

		if (cond_val.Get_bool().IsFalse() || cond_val.Get_bool().IsTop() || loop_widened) && cond_edge.To_false_node_loc.NodeExists() {
			// run just the false branch on filter_false(in_state)
			new_state := in_state.Clone()
			if simpleprop_success {
				false_filters := domain.Filter_false_query_simpleprop(cond_simpleprop)
				// TODO: This doesn't refine array lengths
				for _, filter := range false_filters {
					rhs_dom_val := interpreter.Eval_expr(in_state, filter.Rhs_expr)
					if rhs_dom_val.domain_kind == IntDomainKind {
						updated_intdom := interpreter.get_abstract_value_from_lvalue_expr(filter.Term_expr, new_state).Get_int().Filter(filter.Query_type, rhs_dom_val.Get_int())
						updated_val := AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: updated_intdom}
						interpreter.set_abstract_value_from_expr(filter.Term_expr, updated_val, &new_state)
					}
				}
			}
			return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: cond_edge.To_false_node_loc, abstract_mem: new_state.abstract_mem, remaining_call_depth: in_state.remaining_call_depth})
		}
	}

	switch outgoing_edge := interpreter.Function_cfgs[func_name].Edge_map_from[in_state.node_location.Id].(type) {
	case *CFGEdge:
		new_state := in_state.Clone()
		return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: outgoing_edge.To_node_loc, abstract_mem: new_state.abstract_mem, remaining_call_depth: in_state.remaining_call_depth})
	case *CFGCondEdge:
		// handle condition edges within their respective stmt handling
	}
	return return_states
}

// This is the main struct for driving abstract interpretation.
type AbstractAnalyzer[IntDom domain.IntegerDomain[IntDom], ArrDom ArrayDomain[IntDom, ArrDom]] struct {
	Function_cfgs        FunctionCFGMap
	Function_defs        imp.ImpFunctionMap
	Settings             AnalysisSettings
	Intdomain_default    IntDom
	Arraydomain_default  ArrDom
	function_pre_mem_map map[imp.ImpFunctionName]*AbstractFunctionMem[IntDom, ArrDom] // map from function name to pre-states
	output_handler       *AnalyzerOutputHandler
}

func (analyzer *AbstractAnalyzer[IntDomainImpl, ArrayDomainImpl]) Initialize() {
	analyzer.function_pre_mem_map = make(map[imp.ImpFunctionName]*AbstractFunctionMem[IntDomainImpl, ArrayDomainImpl])
}

func (analyzer *AbstractAnalyzer[IntDom, ArrDom]) Run_analysis() {
	analyzer.Interpret_function(analyzer.Function_defs["main"], nil, analyzer.Settings.Max_call_stack_depth)
}

// Perform abstract interpretation/analysis on the given function, setting the pre-node state of the entry node as initial_node_mem
// Returns the abstract return value
func (analyzer *AbstractAnalyzer[IntDomainImpl, ArrayDomainImpl]) Interpret_function(function_def imp.ImpFunction, initial_node_mem AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl], remaining_call_depth int) AbstractValue[IntDomainImpl, ArrayDomainImpl] {
	function_name := function_def.Name
	analyzer.function_pre_mem_map[function_name] = &AbstractFunctionMem[IntDomainImpl, ArrayDomainImpl]{}
	analyzer.function_pre_mem_map[function_name].Initialize(function_def, analyzer.Function_cfgs[function_name], initial_node_mem)

	initial_state := AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: analyzer.Function_cfgs[function_name].Entry_node, abstract_mem: make(AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]), remaining_call_depth: remaining_call_depth}
	worklist := []AbstractState[IntDomainImpl, ArrayDomainImpl]{initial_state}
	for len(worklist) > 0 {
		front_val := worklist[0]
		worklist = worklist[1:]
		// fmt.Println("Process state", front_val)
		for _, val := range analyzer.Step(front_val) {
			worklist = append(worklist, val)
		}
	}
	return analyzer.function_pre_mem_map[function_name].return_value
	// fmt.Println("Final mem", analyzer.function_pre_mem_map[function_name])
}

func Test(func_cfg_map FunctionCFGMap, func_name imp.ImpFunctionName, func_info_map imp.ImpFunctionMap, debug bool) {
	g := AbstractAnalyzer[domain.IntervalDomain, ArraySummaryDomain[domain.IntervalDomain]]{
		Function_cfgs: func_cfg_map,
		Function_defs: func_info_map,
		Settings: AnalysisSettings{
			Loop_iters_before_widening: 5,
			Max_call_stack_depth:       5,
		},
		Intdomain_default:   domain.IntervalDomain{},
		Arraydomain_default: ArraySummaryDomain[domain.IntervalDomain]{},
		output_handler:      &AnalyzerOutputHandler{},
	}
	g.Initialize()
	g.Run_analysis()
	g.output_handler.Print()
	if !debug {
		return
	}
	fmt.Println("Finished. Final state:")
	for key, val := range g.function_pre_mem_map {
		fmt.Println("---------")
		fmt.Println(key)
		for nid, nval := range val.pre_mem_node_map {
			fmt.Println(nid, ":", nval)
		}
	}

	// modify the cfg so we print mermaid with global state
	for fun_name, val := range g.function_pre_mem_map {
		fmt.Println("---------")
		fmt.Println(fun_name)
		updated_nodes := make(map[NodeID]CFGNodeClass)
		for node_id, node_state := range val.pre_mem_node_map {
			switch node := func_cfg_map[fun_name].Node_map[node_id].(type) {
			case *CFGNode:
				updated_nodes[node_id] = &CFGNode{Id: node.Id, Code: node_state.String() + "\n" + node.Code, Node_type: node.Node_type}
			case *CFGCondNode:
				updated_nodes[node_id] = &CFGCondNode{Id: node.Id, Code: node_state.String() + "\n" + node.Code, Node_type: node.Node_type}
			}
		}
		func_cfg_map[fun_name].Node_map = updated_nodes
		fmt.Println("```")
		fmt.Println(func_cfg_map[fun_name].To_mermaid())
		fmt.Println("```")
	}
}
