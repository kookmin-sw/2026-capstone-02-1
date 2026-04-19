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

// AbstractNodeState maps from variables to AbstractValue
type AbstractNodeState[IntDomainImpl AbstractDomain[IntDomainImpl], BoolDomainImpl AbstractDomain[BoolDomainImpl], ArrayDomainImpl AbstractDomain[ArrayDomainImpl]] map[string]AbstractValue[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]

// An abstract state for a function holds a map from nodes to AbstractImpState
type FunctionAbstractState[IntDomainImpl AbstractDomain[IntDomainImpl], BoolDomainImpl AbstractDomain[BoolDomainImpl], ArrayDomainImpl AbstractDomain[ArrayDomainImpl]] struct {
	states                map[int]AbstractNodeState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
	current_function_name string
	return_value          AbstractValue[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
}
