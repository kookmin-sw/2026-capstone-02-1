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
