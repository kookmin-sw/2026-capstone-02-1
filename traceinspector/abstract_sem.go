package traceinspector

type AnalyzerSettings struct {
	loop_iters_before_widening int
}

// An AbstractState is the pair (l, M^#) ↪ (l', M^#') used in the abstract transition relation
// node_id: node ID to be interpreted
// abstract_mem: the input abstract memory state wrt the node should be interpreted
type AbstractState[IntDomainImpl AbstractDomain[IntDomainImpl], BoolDomainImpl AbstractDomain[BoolDomainImpl], ArrayDomainImpl AbstractDomain[ArrayDomainImpl]] struct {
	node_id      NodeID
	abstract_mem AbstractNodeMem[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
}

// Step: Given an input state (l, m^#), execute the abstract step relation for l under memory state m^#, and
// Return the subsequent states {(l', m^#')} ∈ P(L * M^#)
type AbstractSemantics[IntDomainImpl AbstractDomain[IntDomainImpl], BoolDomainImpl AbstractDomain[BoolDomainImpl], ArrayDomainImpl AbstractDomain[ArrayDomainImpl]] interface {
	Step(AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
}

// type AnalyzerState struct {
// 	cfg_graph        *CFGGraph
// 	imp_function_map imp.ImpFunctionMap
// 	node_state_map   *NodeStateMap // map from node ID to abstract program state
// 	worklist []
// }

// func (a_state *AnalyzerState) step_function(function_name string, function_args ...any) {
// 	imp_function, imp_function_exists := a_state.imp_function_map[function_name]
// }
