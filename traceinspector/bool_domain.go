package traceinspector

import "fmt"

type BoolDomain struct {
	val               bool
	is_bottom, is_top bool
}

func (domain BoolDomain) String() string {
	if domain.is_bottom {
		return "⊥_bool"
	} else if domain.is_top {
		return "⊤_bool"
	} else {
		return fmt.Sprintf("%t", domain.val)
	}
}

func (domain BoolDomain) IsBot() bool {
	return domain.is_bottom
}

func (domain BoolDomain) IsTop() bool {
	return domain.is_top
}

func (lhs BoolDomain) Join(rhs BoolDomain) BoolDomain {
	if lhs.is_top || rhs.is_top {
		return BoolDomain{is_top: true}
	} else if lhs.is_bottom {
		return BoolDomain{val: rhs.val, is_bottom: rhs.is_bottom}
	} else if rhs.is_bottom {
		return BoolDomain{val: lhs.val, is_bottom: lhs.is_bottom}
	} else if lhs.val == rhs.val {
		return BoolDomain{val: lhs.val}
	} else {
		return BoolDomain{is_top: true}
	}
}

func (lhs BoolDomain) Incl(rhs BoolDomain) bool {
	if rhs.is_top {
		return true
	} else if lhs.is_bottom {
		return true
	} else if lhs.is_top {
		return rhs.is_top
	} else if rhs.is_bottom { // lhs concrete
		return false
	} else { // lhs = concrete, rhs = concrete
		return lhs.val == rhs.val
	}
}

func (lhs BoolDomain) Widen(rhs BoolDomain) BoolDomain {
	if lhs.is_bottom {
		return rhs
	}
	if rhs.is_bottom {
		return lhs
	}
	if lhs.is_top || rhs.is_top {
		return BoolDomain{is_top: true}
	}
	if lhs.val == rhs.val {
		return lhs
	} else {
		return BoolDomain{is_top: true}
	}
}

// expression evaluation
