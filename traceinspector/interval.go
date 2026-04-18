package traceinspector

import (
	"fmt"
	"strconv"
)

// We don't use math.Inf since it results in a float
type IntervalDomainValue struct {
	value                  int
	is_infty, is_neg_infty bool // disregard value field if any of these values is true
}

func (v IntervalDomainValue) String() string {
	if v.is_infty {
		return "∞"
	} else if v.is_neg_infty {
		return "∞"
	} else {
		return strconv.Itoa(v.value)
	}
}

// returns whether v is a finite value
func (v IntervalDomainValue) is_finite() bool {
	return !(v.is_infty || v.is_neg_infty)
}

// Compute the minimum of two IntervalDomainValues
func (l1 IntervalDomainValue) min(l2 IntervalDomainValue) IntervalDomainValue {
	if l1.is_neg_infty || l2.is_neg_infty {
		return IntervalDomainValue{is_neg_infty: true} // zero values are 0 and false
	} else if l1.is_infty && l2.is_infty {
		return IntervalDomainValue{is_infty: true}
	} else if l1.is_infty {
		return IntervalDomainValue{value: l2.value}
	} else if l2.is_infty {
		return IntervalDomainValue{value: l1.value}
	} else {
		return IntervalDomainValue{value: min(l1.value, l2.value)}
	}
}

// Compute the maximum of two IntervalDomainValues
func (l1 IntervalDomainValue) max(l2 IntervalDomainValue) IntervalDomainValue {
	if l1.is_infty || l2.is_infty {
		return IntervalDomainValue{is_infty: true}
	} else if l1.is_neg_infty && l2.is_neg_infty {
		return IntervalDomainValue{is_neg_infty: true}
	} else if l1.is_neg_infty {
		return IntervalDomainValue{value: l2.value}
	} else if l2.is_neg_infty {
		return IntervalDomainValue{value: l1.value}
	} else {
		return IntervalDomainValue{value: max(l1.value, l2.value)}
	}
}

func (lhs IntervalDomainValue) eq(rhs IntervalDomainValue) bool {
	if lhs.is_infty {
		return rhs.is_infty
	} else if lhs.is_neg_infty {
		return rhs.is_neg_infty
	} else if rhs.is_infty || rhs.is_neg_infty {
		return false
	} else {
		return lhs.value == rhs.value
	}
}

func (lhs IntervalDomainValue) leq(rhs IntervalDomainValue) bool {
	if lhs.is_neg_infty {
		return true
	} else if rhs.is_infty { // lhs > -infty && rhs = infty
		return true
	} else if rhs.is_neg_infty { // rhs = -infty && lhs > -infty
		return false
	} else if lhs.is_infty { // lhs = infty && -infty < rhs < infty
		return false
	} else {
		return lhs.value <= rhs.value
	}
}

//////////////////////////////////

type IntervalDomain struct {
	lower, upper IntervalDomainValue
	is_bottom    bool
}

func (domain IntervalDomain) String() string {
	if domain.is_bottom {
		return "⊥"
	} else {
		return fmt.Sprintf("[%s, %s]", domain.lower, domain.upper)
	}
}

func (domain IntervalDomain) IsBot() bool {
	return domain.is_bottom
}

func (domain IntervalDomain) IsTop() bool {
	return domain.lower.is_neg_infty && domain.upper.is_infty
}

func (domain IntervalDomain) is_bounded() bool {
	return domain.lower.is_finite() && domain.upper.is_finite()
}

func (lhs IntervalDomain) Join(rhs IntervalDomain) IntervalDomain {
	return IntervalDomain{lower: lhs.lower.min(rhs.lower), upper: lhs.upper.max(rhs.upper)}
}

// `lhs ⊑ rhs` = lhs.lower >= rhs.lower && lhs.upper <= rhs.upper
func (lhs IntervalDomain) Incl(rhs IntervalDomain) bool {
	return rhs.lower.leq(lhs.lower) && lhs.upper.leq(rhs.upper)
}

// replace increasing chains with infty/-infty
// [n, u1] ▽ [n, u2] = if u2 <= u1 then [n, u1] else [n, infty]
// [l1, n] ▽ [l2, n] = if l1 <= l2 then [l1, n] else [-infty, n]
func (lhs IntervalDomain) Widen(rhs IntervalDomain) IntervalDomain {
	if lhs.is_bottom {
		return rhs
	}
	if rhs.is_bottom {
		return lhs
	}

	var lower_val IntervalDomainValue
	var upper_val IntervalDomainValue
	if rhs.upper.leq(lhs.upper) {
		upper_val = rhs.upper
	} else {
		upper_val = IntervalDomainValue{is_infty: true}
	}

	if lhs.lower.leq(rhs.lower) {
		lower_val = lhs.lower
	} else {
		lower_val = IntervalDomainValue{is_neg_infty: true}
	}
	return IntervalDomain{lower: lower_val, upper: upper_val}
}
