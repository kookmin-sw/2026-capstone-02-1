package traceinspector

import (
	"fmt"
	"strings"
	"traceinspector/algebra"
	"traceinspector/domain"
	"traceinspector/imp"
)

// An AbstractState is the pair (l, M^#) ↪ (l', M^#') used in the abstract transition relation
// node_location: node location to be interpreted
// abstract_mem: the input abstract memory state wrt the node should be interpreted
type AbstractState[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl ArrayDomain[IntDomainImpl, ArrayDomainImpl]] struct {
	node_location CFGNodeLocation
	abstract_mem  AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]
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
	return new_st
}

// Step: Given an input state (l, m^#), execute the abstract step relation for l under memory state m^#, and
// Return the subsequent states {(l', m^#')} ∈ P(L * M^#)
type AbstractSemantics[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl ArrayDomain[IntDomainImpl, ArrayDomainImpl]] interface {
	Step(AbstractState[IntDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, ArrayDomainImpl]
}

// Abstract transition semantics for Imp wrt to arbitrary abstract domain impelmentations

// ImpFunctionInterpreter performs abstract interpretation of a function body from a given initial state. The
// interpreter performs interpretation until it collects the fixpoint semantics for the function body, and hence the
// return value. The interpreter will spawn another ImpFunctionInterpreter in the case a function call is invoked.
type ImpFunctionInterpreter[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl ArrayDomain[IntDomainImpl, ArrayDomainImpl]] struct {
	func_cfg_map          FunctionCFGMap
	func_name             imp.ImpFunctionName
	func_info_map         imp.ImpFunctionMap
	abstract_function_mem *AbstractFunctionMem[IntDomainImpl, ArrayDomainImpl] // joined global state
	intdomain_default     IntDomainImpl                                        // an instantiation of the integer domain impl
	booldomain_default    domain.BoolDomain                                    // an instantiation of the boolean domain
	arraydomain_default   ArrayDomainImpl                                      // an instantiation of the array domain impl
	settings              AnalysisSettings
}

// TODO: fix this part and integrate with get_abstract_value_from_expr
// Compute the abstract value of an expression expr under an abstract memory state abs_mem
func (interpreter *ImpFunctionInterpreter[IntDomainImpl, ArrayDomainImpl]) Eval_expr(node_location CFGNodeLocation, expr imp.Expr, abs_mem AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]) AbstractValue[IntDomainImpl, ArrayDomainImpl] {
	switch expr_ty := expr.(type) {
	case *imp.VarExpr:
		return abs_mem[expr_ty.Name]
	case *imp.IntLitExpr:
		intdom_result := interpreter.intdomain_default.From_IntLitExpr(*expr_ty)
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: intdom_result}
	case *imp.BoolLitExpr:
		booldom_result := interpreter.booldomain_default.From_BoolLitExpr(*expr_ty)
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: booldom_result}
	case *imp.StringLitExpr:
		// TODO do something about strings
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{}
	case *imp.NegExpr:
		subexpr_val := interpreter.Eval_expr(node_location, expr_ty.Subexpr, abs_mem)
		if subexpr_val.domain_kind != IntDomainKind {
			write_error(node_location, fmt.Sprintf("Result of arithmetic negation returned %s instead if IntDomain", subexpr_val.domain_kind))
		}
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: subexpr_val.Get_int().Neg()}
	case *imp.AddExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			write_error(node_location, "Add expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Add(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.SubExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			write_error(node_location, "Sub expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Sub(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.MulExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			write_error(node_location, "Mul expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Mul(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.DivExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			write_error(node_location, "Div expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Div(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.ModExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			write_error(node_location, "Add expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Mod(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.ParenExpr:
		return interpreter.Eval_expr(node_location, expr_ty.Subexpr, abs_mem)
	case *imp.ArrayIndexExpr:
		arr_val := interpreter.Eval_expr(node_location, expr_ty.Base, abs_mem)
		index_val := interpreter.Eval_expr(node_location, expr_ty.Index, abs_mem)
		if arr_val.domain_kind != ArrayDomainKind {
			write_error(node_location, fmt.Sprintf("'%s' : expected arr to have arr domain type", expr_ty))
		}
		if !index_val.Get_int().Leq(arr_val.Get_array().Len().Sub(index_val.int_domain.From_IntLitExpr(imp.IntLitExpr{Value: 1}))).IsTrue() {
			write_warning(node_location, fmt.Sprintf("Potentially unsafe array indexing: index has value %s, but %s.Len has value %s.", index_val.Get_int(), get_varname_from_lvalue(expr_ty.Base), arr_val.Get_array().Len()))
		}
		result_val := arr_val.Get_array().GetIndex(index_val.Get_int())
		return result_val
	case *imp.LenExpr:
		arr_val := interpreter.Eval_expr(node_location, expr_ty.Subexpr, abs_mem)
		if arr_val.domain_kind != ArrayDomainKind {
			write_error(node_location, fmt.Sprintf("'%s' : expected arr to have arr domain type", expr_ty))
		}
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: arr_val.Get_array().Len()}
	case *imp.MakeArrayExpr:
		default_val := interpreter.Eval_expr(node_location, expr_ty.Value, abs_mem)
		size_val := interpreter.Eval_expr(node_location, expr_ty.Size, abs_mem)
		if size_val.domain_kind != IntDomainKind {
			write_error(node_location, fmt.Sprintf("Size argument '%s' in make_array expected to have int domain type, but has type %s", expr_ty.Size, size_val.domain_kind))
		}
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: ArrayDomainKind, array_domain: interpreter.arraydomain_default.Make_array(size_val.Get_int(), default_val)}
	case *imp.EqExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if lhs_val.domain_kind != rhs_val.domain_kind {
			write_error(node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different", expr_ty))
		}
		result_val := lhs_val.Get_int().Eq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.NeqExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if lhs_val.domain_kind != rhs_val.domain_kind {
			write_error(node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different", expr_ty))
		}
		result_val := lhs_val.Get_int().Neq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.LeqExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && lhs_val.domain_kind == rhs_val.domain_kind) {
			write_error(node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different (%s vs %s)", expr_ty, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		result_val := lhs_val.Get_int().Leq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.GeqExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && lhs_val.domain_kind == rhs_val.domain_kind) {
			write_error(node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different (%s vs %s)", expr_ty, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		result_val := lhs_val.Get_int().Geq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.LessthanExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && lhs_val.domain_kind == rhs_val.domain_kind) {
			write_error(node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different (%s vs %s)", expr_ty, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		result_val := lhs_val.Get_int().Lessthan(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.GreaterthanExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && lhs_val.domain_kind == rhs_val.domain_kind) {
			write_error(node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different (%s vs %s)", expr_ty, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		result_val := lhs_val.Get_int().Greaterthan(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
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
func (interpreter *ImpFunctionInterpreter[IntDomainImpl, ArrayDomainImpl]) get_abstract_value_from_lvalue_expr(expr imp.Expr, state AbstractState[IntDomainImpl, ArrayDomainImpl]) AbstractValue[IntDomainImpl, ArrayDomainImpl] {
	switch expr_ty := expr.(type) {
	case *imp.VarExpr:
		val, var_in_mem := state.abstract_mem[expr_ty.Name]
		if !var_in_mem {
			write_error(state.node_location, fmt.Sprintf("Variable '%s' not defined", expr_ty.Name))
		}
		return val
	case *imp.LenExpr:
		arr_val := interpreter.get_abstract_value_from_lvalue_expr(expr_ty.Subexpr, state)
		if arr_val.domain_kind != ArrayDomainKind {
			write_error(state.node_location, fmt.Sprintf("Called len() on a non-array value '%s'", expr_ty.Subexpr))
		}
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: arr_val.Get_array().Len()}
	case *imp.ArrayIndexExpr:
		arr_val := interpreter.get_abstract_value_from_lvalue_expr(expr_ty.Base, state)
		if arr_val.domain_kind != ArrayDomainKind {
			write_error(state.node_location, fmt.Sprintf("Tried to index a non-array value '%s'", expr_ty.Base))
		}
		index_val := interpreter.Eval_expr(state.node_location, expr_ty.Index, state.abstract_mem)
		if index_val.domain_kind != ArrayDomainKind {
			write_error(state.node_location, fmt.Sprintf("Index expression '%s' is not an integerdomain", expr_ty.Index))
		}
		return arr_val.Get_array().GetIndex(index_val.Get_int())
	}
	write_error(state.node_location, fmt.Sprintf("get_abstract_value_from_expr: Unsupported expression '%s'", expr))
	panic("unreachable")
}

// TODO: fix this part
// Given an lvalue expression lhs, the rhs AbstractValue, and the state, update state.abstract_mem accordingly *in-place*
func (interpreter *ImpFunctionInterpreter[IntDomainImpl, ArrayDomainImpl]) set_abstract_value_from_expr(lhs imp.Expr, rhs_val AbstractValue[IntDomainImpl, ArrayDomainImpl], state *AbstractState[IntDomainImpl, ArrayDomainImpl]) {
	switch lhs_node := lhs.(type) {
	case *imp.VarExpr:
		lhs_val, lhs_exists := state.abstract_mem[lhs_node.Name]
		if lhs_exists && lhs_val.domain_kind != rhs_val.domain_kind {
			write_error(state.node_location, fmt.Sprintf("LHS of variable %s has type %s, but RHS has type %s.", lhs_node.Name, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		state.abstract_mem[lhs_node.Name] = rhs_val
	case *imp.ArrayIndexExpr:
		index_val := interpreter.Eval_expr(state.node_location, lhs_node.Index, state.abstract_mem)
		if index_val.domain_kind != IntDomainKind {
			write_error(state.node_location, "Only 1-dimensional arrays implemented")
		}
		if index_val.Get_int().IsBot() {
			// if index is bot, don't do assignment
			return
		}
		arr_varname := get_varname_from_lvalue(lhs_node.Base)
		lhs_val, lhs_exists := state.abstract_mem[arr_varname]
		if !lhs_exists {
			write_error(state.node_location, fmt.Sprintf("Attempting to index nonexisting variable '%s'", arr_varname))
		}
		if lhs_val.domain_kind != ArrayDomainKind {
			write_error(state.node_location, fmt.Sprintf("Attempting to index non-array variable '%s'", arr_varname))
		}
		// fmt.Println(index_val.Get_int(), "<=", lhs_val.Get_array().Len().Sub(index_val.int_domain.From_IntLitExpr(imp.IntLitExpr{Value: 1})), "=", index_val.Get_int().Leq(lhs_val.Get_array().Len().Sub(index_val.int_domain.From_IntLitExpr(imp.IntLitExpr{Value: 1}))))
		if !index_val.Get_int().Leq(lhs_val.Get_array().Len().Sub(index_val.int_domain.From_IntLitExpr(imp.IntLitExpr{Value: 1}))).IsTrue() {
			write_warning(state.node_location, fmt.Sprintf("Potentially unsafe array indexing: index has value %s, but %s.Len has value %s.", index_val.Get_int(), arr_varname, lhs_val.Get_array().Len()))
		}
		state.abstract_mem[arr_varname] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: ArrayDomainKind, array_domain: lhs_val.Get_array().SetVal(index_val.Get_int(), rhs_val)}
	default:
		write_error(state.node_location, fmt.Sprintf("Unsupported LHS expr '%s' with type %T", lhs, lhs))
	}
}

func (interpreter *ImpFunctionInterpreter[IntDomainImpl, ArrayDomainImpl]) Step(in_state AbstractState[IntDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, ArrayDomainImpl] {
	cfg_node, cfg_node_exists := interpreter.func_cfg_map[interpreter.func_name].Node_map[in_state.node_location.Id]
	if !cfg_node_exists {
		write_error(create_empty_node_location(), fmt.Sprintf("The designated CFG Node %s doesn't exist", in_state.node_location))
	}

	global_state, _ := interpreter.abstract_function_mem.pre_mem_node_map[in_state.node_location.Id]

	cond_node, is_cond_node := cfg_node.(*CFGCondNode)
	loop_widened := false
	if is_cond_node && cond_node.Is_loop_head && interpreter.abstract_function_mem.n_visits[in_state.node_location.Id] > interpreter.settings.Loop_iters_before_Widening {
		// Apply widening if visit count is greater than threshold
		interpreter.abstract_function_mem.n_visits[in_state.node_location.Id] = 1
		global_state.Widen_inplace(in_state.abstract_mem)
		write_update_node_state(in_state.node_location, global_state.String(), "Widen global memory state")
		loop_widened = true
	} else {
		// When we receive a new pair, update the global state with its join
		state_changed := global_state.Join_inplace(in_state.abstract_mem)
		if !state_changed && interpreter.abstract_function_mem.n_visits[in_state.node_location.Id] > 0 {
			// no updates to the state
			write_info(in_state.node_location, "No updates to node state")
			return nil
		}
		write_update_node_state(in_state.node_location, global_state.String(), "Join global memory state")
	}

	interpreter.abstract_function_mem.n_visits[in_state.node_location.Id]++
	// Executed on the joined node state
	in_state.abstract_mem = global_state.Clone()

	var return_states []AbstractState[IntDomainImpl, ArrayDomainImpl]
	switch cfg_node := cfg_node.(type) {
	case *CFGNode:
		switch stmt := cfg_node.Ast.(type) {
		case *imp.AssignStmt:
			// assignment should overwrite the value, instead of join
			rhs_val := interpreter.Eval_expr(in_state.node_location, stmt.Rhs, in_state.abstract_mem)
			interpreter.set_abstract_value_from_expr(stmt.Lhs, rhs_val, &in_state)

		case *imp.SkipStmt:
			// do nothing
		case *imp.BreakStmt:
			// do nothing; follows CFG edge
		case *imp.ContinueStmt:
			// do nothing; follows CFG edge
		case *imp.PrintStmt:
			// do nothing to state; follows CFG edge
		// case *imp.ScanfStmt:
		// 	// by default scanf initializes variables to top
		// 	for index, fmt_str := range strings.Split(stmt.Format_string, " ") {
		// 		switch fmt_str {
		// 		case "%d":
		// 			// integer type
		// 		}
		// 	}

		default:
			panic(fmt.Sprintf("unimplemented stmt %T", stmt))
		}

	case *CFGCondNode:
		cond_edge, is_cond_edge := interpreter.func_cfg_map[interpreter.func_name].Edge_map_from[in_state.node_location.Id].(*CFGCondEdge)
		if !is_cond_edge {
			write_error(in_state.node_location, "Condition stmt does not have outgoing edge of CondEdge type.")
		}
		// If at in_state the prop evaluates to either true or false,
		// We can just execute only the corresponding branch.
		// Otherwise filter for each branch and join the result
		cond_val := interpreter.Eval_expr(in_state.node_location, cfg_node.Cond_expr, in_state.abstract_mem)
		// Try to represent it as SimpleProp
		cond_simpleprop, simpleprop_success := algebra.Imp_expr_to_simple_prop(cfg_node.Cond_expr)
		if !simpleprop_success {
			write_warning(in_state.node_location, fmt.Sprintf("Could not represent '%s' as SimpleProp. Analysis precision may severely deterioriate.", cfg_node.Cond_expr))
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
					rhs_dom_val := interpreter.Eval_expr(in_state.node_location, filter.Rhs_expr, in_state.abstract_mem)
					if rhs_dom_val.domain_kind == IntDomainKind {
						updated_intdom := interpreter.get_abstract_value_from_lvalue_expr(filter.Term_expr, new_state).Get_int().Filter(filter.Query_type, rhs_dom_val.Get_int())
						updated_val := AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: updated_intdom}
						interpreter.set_abstract_value_from_expr(filter.Term_expr, updated_val, &new_state)
					}
				}
			}
			return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: cond_edge.To_true_node_loc, abstract_mem: new_state.abstract_mem})
		}

		if (cond_val.Get_bool().IsFalse() || cond_val.Get_bool().IsTop() || loop_widened) && cond_edge.To_false_node_loc.NodeExists() {
			// run just the false branch on filter_false(in_state)
			new_state := in_state.Clone()
			if simpleprop_success {
				false_filters := domain.Filter_false_query_simpleprop(cond_simpleprop)
				// TODO: This doesn't refine array lengths
				for _, filter := range false_filters {
					rhs_dom_val := interpreter.Eval_expr(in_state.node_location, filter.Rhs_expr, in_state.abstract_mem)
					if rhs_dom_val.domain_kind == IntDomainKind {
						updated_intdom := interpreter.get_abstract_value_from_lvalue_expr(filter.Term_expr, new_state).Get_int().Filter(filter.Query_type, rhs_dom_val.Get_int())
						updated_val := AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: updated_intdom}
						interpreter.set_abstract_value_from_expr(filter.Term_expr, updated_val, &new_state)
					}
				}
			}
			return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: cond_edge.To_false_node_loc, abstract_mem: new_state.abstract_mem})
		}
	}

	switch outgoing_edge := interpreter.func_cfg_map[interpreter.func_name].Edge_map_from[in_state.node_location.Id].(type) {
	case *CFGEdge:
		new_state := in_state.Clone()
		return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: outgoing_edge.To_node_loc, abstract_mem: new_state.abstract_mem})
	case *CFGCondEdge:
		// handle condition edges within their respective stmt handling
	}
	return return_states
}
