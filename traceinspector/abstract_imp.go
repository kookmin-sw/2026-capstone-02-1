package traceinspector

import "traceinspector/domain"

// Abstract transition semantics for Imp wrt to arbitrary abstract domain impelmentations

type ImpAbstractSemantics[IntDomainImpl domain.AbstractDomain[IntDomainImpl], BoolDomainImpl domain.AbstractDomain[BoolDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] struct{}

// func (interpreter *ImpAbstractInterpreter[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) Step(in_state AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl] {
// 	cfg_node, cfg_node_exists := interpreter.cfg.Node_map[in_state.node_id]
// 	if !cfg_node_exists {
// 		write_error(create_empty_node_location(), fmt.Sprintf("The designated CFG Node %d doesn't exist", in_state.node_id))
// 	}
// 	switch cfg_node := cfg_node.(type) {
// 	case *CFGCondNode:
// 	}
// }
