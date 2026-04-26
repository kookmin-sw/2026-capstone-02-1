package traceinspector

import (
	"fmt"
	"strings"
	"traceinspector/domain"
	"traceinspector/imp"
)

// Abstract states wrt arbitrary pointwise abstract domains

// AbstractValue maps from types to abstract domain values
type AbstractDomainKind int

const (
	InvalidKind AbstractDomainKind = iota
	IntDomainKind
	BoolDomainKind
	ArrayDomainKind
)

func (dom_kind AbstractDomainKind) String() string {
	switch dom_kind {
	case InvalidKind:
		return "InvalidDomainKind"
	case IntDomainKind:
		return "IntDomainKind"
	case BoolDomainKind:
		return "BoolDomainKind"
	case ArrayDomainKind:
		return "ArrayDomainKind"
	}
	panic("This should be unreachable")
}

// AbstractValue holds the abstract domain value for a variable
type AbstractValue[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] struct {
	domain_kind  AbstractDomainKind
	int_domain   IntDomainImpl
	bool_domain  domain.BoolDomain
	array_domain ArrayDomainImpl
}

func (val AbstractValue[IntDomainImpl, ArrayDomainImpl]) Get_int() IntDomainImpl {
	if val.domain_kind != IntDomainKind {
		panic(fmt.Sprintf("Attempted to get non-intkind(%s) abstractvalue as int", val.domain_kind))
	}
	return val.int_domain
}

func (val AbstractValue[IntDomainImpl, ArrayDomainImpl]) Get_bool() domain.BoolDomain {
	if val.domain_kind != BoolDomainKind {
		panic("Attempted to get non-boolkind abstractvalue as bool")
	}
	return val.bool_domain
}

func (val AbstractValue[IntDomainImpl, ArrayDomainImpl]) Get_array() ArrayDomainImpl {
	if val.domain_kind != ArrayDomainKind {
		panic("Attempted to get non-arraykind abstractvalue as array")
	}
	return val.array_domain
}

func (val AbstractValue[IntDomainImpl, ArrayDomainImpl]) String() string {
	switch val.domain_kind {
	case InvalidKind:
		return "INVALID"
	case IntDomainKind:
		return val.Get_int().String()
	case BoolDomainKind:
		return val.Get_bool().String()
	case ArrayDomainKind:
		return val.Get_array().String()
	}
	return ""
}

// AbstractVarMemMap maps from variables to AbstractValue
type AbstractVarMemMap[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] map[string]AbstractValue[IntDomainImpl, ArrayDomainImpl]

func (node_mem AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]) String() string {
	var ret []string
	for key, val := range node_mem {
		ret = append(ret, fmt.Sprintf("%s : %s", key, val))
	}
	return "{" + strings.Join(ret, ", ") + "}"
}

func (node_mem AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]) Clone() AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl] {
	new_mem := make(AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl])
	for key, val := range node_mem {
		switch val.domain_kind {
		case InvalidKind:
			new_mem[key] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{}
		case IntDomainKind:
			new_mem[key] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: val.Get_int().Clone()}
		case BoolDomainKind:
			new_mem[key] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: val.Get_bool().Clone()}
		case ArrayDomainKind:
			new_mem[key] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: ArrayDomainKind, array_domain: val.Get_array().Clone()}
		}
	}
	return new_mem
}

// Modify the values in node_mem inplace so that they are the result of joining with values in other_mem
// Returns bool indicating whether the mem was changed
func (node_mem AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]) Join_inplace(other_mem AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]) bool {
	changed := false
	for key, val := range other_mem {
		original_val, original_exists := node_mem[key]
		var joined AbstractValue[IntDomainImpl, ArrayDomainImpl]
		if original_exists {
			joined.domain_kind = original_val.domain_kind
			indiv_changed := false
			switch joined.domain_kind {
			case IntDomainKind:
				joined.int_domain, indiv_changed = original_val.Get_int().Clone().Join(val.Get_int())
			case BoolDomainKind:
				joined.bool_domain, indiv_changed = original_val.Get_bool().Clone().Join(val.Get_bool())
			case ArrayDomainKind:
				joined.array_domain, indiv_changed = original_val.Get_array().Clone().Join(val.Get_array())
			}
			changed = changed || indiv_changed
		} else {
			changed = true
			joined.domain_kind = val.domain_kind
			switch joined.domain_kind {
			case IntDomainKind:
				joined.int_domain = val.Get_int()
			case BoolDomainKind:
				joined.bool_domain = val.Get_bool()
			case ArrayDomainKind:
				joined.array_domain = val.Get_array()
			}
		}
		node_mem[key] = joined
	}
	return changed
}

func (node_mem AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]) Widen_inplace(other_mem AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]) {
	for key, val := range other_mem {
		original_val, original_exists := node_mem[key]
		var widened AbstractValue[IntDomainImpl, ArrayDomainImpl]
		if original_exists {
			widened.domain_kind = original_val.domain_kind
			switch widened.domain_kind {
			case IntDomainKind:
				widened.int_domain = original_val.Get_int().Clone().Widen(val.Get_int())
			case BoolDomainKind:
				widened.bool_domain = original_val.Get_bool().Clone().Widen(val.Get_bool())
			case ArrayDomainKind:
				widened.array_domain = original_val.Get_array().Clone().Widen(val.Get_array())
			}
		} else {
			widened.domain_kind = val.domain_kind
			switch widened.domain_kind {
			case IntDomainKind:
				widened.int_domain = val.Get_int()
			case BoolDomainKind:
				widened.bool_domain = val.Get_bool()
			case ArrayDomainKind:
				widened.array_domain = val.Get_array()
			}
		}
		node_mem[key] = widened
	}
}

type AbstractNodeMemMap[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] map[NodeID]AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]

func (mem_map AbstractNodeMemMap[IntDomainImpl, ArrayDomainImpl]) String() string {
	var ret []string
	for key, val := range mem_map {
		ret = append(ret, fmt.Sprintf("%d : %s", key, val))
	}
	return "{" + strings.Join(ret, ", ") + "}"
}

// An abstract Memory state for a function holds a map from nodes to AbstractNodeMem
// pre_mem represents the memory state at the **entry of a node - before executing the node**.
// the return value is also an abstraction of the possible return values
type AbstractFunctionMem[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] struct {
	pre_mem_node_map AbstractNodeMemMap[IntDomainImpl, ArrayDomainImpl]
	function_name    imp.ImpFunctionName
	n_visits         map[NodeID]int
	return_value     AbstractValue[IntDomainImpl, ArrayDomainImpl]
}

func (func_mem AbstractFunctionMem[IntDomainImpl, ArrayDomainImpl]) String() string {
	return fmt.Sprintf("%s : %s, returns %s", func_mem.function_name, func_mem.pre_mem_node_map, func_mem.return_value)
}

func (func_mem *AbstractFunctionMem[IntDomainImpl, ArrayDomainImpl]) Initialize(function_name imp.ImpFunctionName, function_cfg *CFGGraph, initial_node_mem AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl]) {
	func_mem.pre_mem_node_map = make(AbstractNodeMemMap[IntDomainImpl, ArrayDomainImpl])
	func_mem.n_visits = make(map[NodeID]int)
	func_mem.function_name = function_name
	for node_id := range function_cfg.Node_map {
		if node_id == function_cfg.Entry_node.Id {
			func_mem.pre_mem_node_map[node_id] = initial_node_mem
		} else {
			func_mem.pre_mem_node_map[node_id] = make(AbstractVarMemMap[IntDomainImpl, ArrayDomainImpl])
		}
	}
	func_mem.return_value = AbstractValue[IntDomainImpl, ArrayDomainImpl]{}
}
