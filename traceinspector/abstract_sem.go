package traceinspector

import (
	"fmt"
	"traceinspector/domain"
	"traceinspector/imp"
)

type AnalyzerSettings struct {
	loop_iters_before_widening int
}

// An AbstractState is the pair (l, M^#) ↪ (l', M^#') used in the abstract transition relation
// node_location: node location to be interpreted
// abstract_mem: the input abstract memory state wrt the node should be interpreted
type AbstractState[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl ArrayDomain[IntDomainImpl, ArrayDomainImpl]] struct {
	node_location CFGNodeLocation
	abstract_mem  AbstractNodeMem[IntDomainImpl, ArrayDomainImpl]
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
	func_cfg_map        FunctionCFGMap
	func_name           imp.ImpFunctionName
	func_info_map       imp.ImpFunctionMap
	abstract_mem        *FunctionAbstractMem[IntDomainImpl, ArrayDomainImpl] // joined global state
	intdomain_default   IntDomainImpl                                        // an instantiation of the integer domain impl
	booldomain_default  domain.BoolDomain                                    // an instantiation of the boolean domain
	arraydomain_default ArrayDomainImpl                                      // an instantiation of the array domain impl
}

// Compute the abstract value of an expression expr under an abstract memory state abs_mem
func (interpreter *ImpFunctionInterpreter[IntDomainImpl, ArrayDomainImpl]) Eval_expr(node_location CFGNodeLocation, expr imp.Expr, abs_mem AbstractNodeMem[IntDomainImpl, ArrayDomainImpl]) AbstractValue[IntDomainImpl, ArrayDomainImpl] {
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
		result_val := arr_val.Get_array().Index(index_val.Get_int())
		return result_val
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
	}
	return AbstractValue[IntDomainImpl, ArrayDomainImpl]{}
}

func (interpreter *ImpFunctionInterpreter[IntDomainImpl, ArrayDomainImpl]) Step(in_state AbstractState[IntDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, ArrayDomainImpl] {
	cfg_node, cfg_node_exists := interpreter.func_cfg_map[interpreter.func_name].Node_map[in_state.node_location.Id]
	if !cfg_node_exists {
		write_error(create_empty_node_location(), fmt.Sprintf("The designated CFG Node %d doesn't exist", in_state.node_location))
	}
	var return_states []AbstractState[IntDomainImpl, ArrayDomainImpl]
	switch cfg_node := cfg_node.(type) {
	case *CFGNode:
		switch stmt := cfg_node.Ast.(type) {
		case *imp.AssignStmt:
			rhs_val := interpreter.Eval_expr(in_state.node_location, stmt.Rhs, in_state.abstract_mem)
			switch lhs_ty := stmt.Lhs.(type) {
			case *imp.VarExpr:
				_, var_exists := in_state.abstract_mem[lhs_ty.Name]
				if var_exists {
					// join here
				} else {
					in_state.abstract_mem[lhs_ty.Name] = rhs_val
				}
			}
		}
	case *CFGCondNode:
	}
}
