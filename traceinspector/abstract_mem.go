package traceinspector

import (
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

// AbstractValue holds the abstract domain value for a variable
type AbstractValue[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] struct {
	domain_kind  AbstractDomainKind
	int_domain   IntDomainImpl
	bool_domain  domain.BoolDomain
	array_domain ArrayDomainImpl
}

func (val AbstractValue[IntDomainImpl, ArrayDomainImpl]) Get_int() IntDomainImpl {
	if val.domain_kind != IntDomainKind {
		panic("Attempted to get non-intkind abstractvalue as int")
	}
	return val.int_domain
}

func (val AbstractValue[IntDomainImpl, ArrayDomainImpl]) Get_bool() domain.BoolDomain {
	if val.domain_kind != IntDomainKind {
		panic("Attempted to get non-boolkind abstractvalue as int")
	}
	return val.bool_domain
}

func (val AbstractValue[IntDomainImpl, ArrayDomainImpl]) Get_array() ArrayDomainImpl {
	if val.domain_kind != IntDomainKind {
		panic("Attempted to get non-arraykind abstractvalue as int")
	}
	return val.array_domain
}

// AbstractNodeMem maps from variables to AbstractValue
type AbstractNodeMem[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] map[string]AbstractValue[IntDomainImpl, ArrayDomainImpl]

func (node_mem AbstractNodeMem[IntDomainImpl, ArrayDomainImpl]) Clone() AbstractNodeMem[IntDomainImpl, ArrayDomainImpl] {
	new_mem := make(AbstractNodeMem[IntDomainImpl, ArrayDomainImpl])
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
func (node_mem AbstractNodeMem[IntDomainImpl, ArrayDomainImpl]) Join_inplace(other_mem AbstractNodeMem[IntDomainImpl, ArrayDomainImpl]) {
	for key, val := range other_mem {
		original_val, original_exists := node_mem[key]
		var joined AbstractValue[IntDomainImpl, ArrayDomainImpl]
		if original_exists {
			joined.domain_kind = original_val.domain_kind
			switch joined.domain_kind {
			case IntDomainKind:
				joined.int_domain = original_val.Get_int().Clone()
			case BoolDomainKind:
				joined.bool_domain = original_val.Get_bool().Clone()
			case ArrayDomainKind:
				joined.array_domain = original_val.Get_array().Clone()
			}
		}
		switch val.domain_kind {
		case IntDomainKind:
			joined.int_domain = joined.Get_int().Join(val.Get_int())
		case BoolDomainKind:
			joined.bool_domain = joined.Get_bool().Join(val.Get_bool())
		case ArrayDomainKind:
			joined.array_domain = joined.Get_array().Join(val.Get_array())
		}
		node_mem[key] = joined
	}
}

// An abstract Memory state for a function holds a map from nodes to AbstractNodeMem
// mem represents the memory state at the **entry of a node - before executing the node**.
// the return value is also an abstraction of the possible return values
type FunctionAbstractMem[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] struct {
	mem           map[NodeID]AbstractNodeMem[IntDomainImpl, ArrayDomainImpl]
	function_name imp.ImpFunctionName
	return_value  AbstractValue[IntDomainImpl, ArrayDomainImpl]
}

func (func_mem *FunctionAbstractMem[IntDomainImpl, ArrayDomainImpl]) Initialize(function_name imp.ImpFunctionName) {
	func_mem.mem = make(map[NodeID]AbstractNodeMem[IntDomainImpl, ArrayDomainImpl])
	func_mem.function_name = function_name
}
