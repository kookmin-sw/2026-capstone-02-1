package traceinspector

// Abstract states wrt arbitrary pointwise abstract domains

// AbstractValue maps from types to abstract domain values
type AbstractDomainKind int

const (
	IntDomainKind AbstractDomainKind = iota
	BoolDomainKind
	ArrayDomainKind
)

// AbstractValue holds the abstract domain value for a variable
type AbstractValue[IntDomainImpl AbstractDomain[IntDomainImpl], BoolDomainImpl AbstractDomain[BoolDomainImpl], ArrayDomainImpl AbstractDomain[ArrayDomainImpl]] struct {
	domain_kind  AbstractDomainKind
	int_domain   IntDomainImpl
	bool_domain  BoolDomainImpl
	array_domain ArrayDomainImpl
}

func (val AbstractValue[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) get_int() IntDomainImpl {
	if val.domain_kind != IntDomainKind {
		panic("Attempted to get non-intkind abstractvalue as int")
	}
	return val.int_domain
}

func (val AbstractValue[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) get_bool() BoolDomainImpl {
	if val.domain_kind != IntDomainKind {
		panic("Attempted to get non-intkind abstractvalue as int")
	}
	return val.bool_domain
}

func (val AbstractValue[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) get_array() ArrayDomainImpl {
	if val.domain_kind != IntDomainKind {
		panic("Attempted to get non-intkind abstractvalue as int")
	}
	return val.array_domain
}

// AbstractNodeMem maps from variables to AbstractValue
type AbstractNodeMem[IntDomainImpl AbstractDomain[IntDomainImpl], BoolDomainImpl AbstractDomain[BoolDomainImpl], ArrayDomainImpl AbstractDomain[ArrayDomainImpl]] map[string]AbstractValue[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]

// An abstract Memory state for a function holds a map from nodes to AbstractNodeMem
// mem represents the memory state at the **entry, before executing the node**.
type FunctionAbstractMem[IntDomainImpl AbstractDomain[IntDomainImpl], BoolDomainImpl AbstractDomain[BoolDomainImpl], ArrayDomainImpl AbstractDomain[ArrayDomainImpl]] struct {
	mem                   map[NodeID]AbstractNodeMem[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
	current_function_name string
	return_value          AbstractValue[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
}
