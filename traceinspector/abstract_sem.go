package traceinspector

import "traceinspector/domain"

type AnalyzerSettings struct {
	loop_iters_before_widening int
}

// An AbstractState is the pair (l, M^#) ↪ (l', M^#') used in the abstract transition relation
// node_id: node ID to be interpreted
// abstract_mem: the input abstract memory state wrt the node should be interpreted
type AbstractState[IntDomainImpl domain.AbstractDomain[IntDomainImpl], BoolDomainImpl domain.AbstractDomain[BoolDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] struct {
	node_id      NodeID
	abstract_mem AbstractNodeMem[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
}

// Step: Given an input state (l, m^#), execute the abstract step relation for l under memory state m^#, and
// Return the subsequent states {(l', m^#')} ∈ P(L * M^#)
type AbstractSemantics[IntDomainImpl domain.AbstractDomain[IntDomainImpl], BoolDomainImpl domain.AbstractDomain[BoolDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] interface {
	Step(AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
}
