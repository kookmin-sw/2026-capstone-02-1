package domain

import (
	"fmt"
	"traceinspector/algebra"
)

//////////////////////////////////

type IntervalDomain struct {
	lower, upper algebra.ExtInt
	is_bottom    bool
}

func (domain IntervalDomain) String() string {
	if domain.is_bottom {
		return "⊥"
	} else {
		return fmt.Sprintf("[%s, %s]", domain.lower, domain.upper)
	}
}

func (domain IntervalDomain) Clone() IntervalDomain {
	return IntervalDomain{lower: domain.lower, upper: domain.upper, is_bottom: domain.is_bottom}
}

func IntervalBot() IntervalDomain {
	return IntervalDomain{is_bottom: true}
}

func (domain IntervalDomain) IsBot() bool {
	return domain.is_bottom
}

func (domain IntervalDomain) IsTop() bool {
	return domain.lower.IsNegInfty() && domain.upper.IsInfty()
}

func (domain IntervalDomain) Is_bounded() bool {
	return domain.lower.IsFinite() && domain.upper.IsFinite()
}

func (lhs IntervalDomain) Join(rhs IntervalDomain) IntervalDomain {
	return IntervalDomain{lower: lhs.lower.Min(rhs.lower), upper: lhs.upper.Max(rhs.upper)}
}

// `lhs ⊑ rhs` = lhs.lower >= rhs.lower && lhs.upper <= rhs.upper
func (lhs IntervalDomain) Incl(rhs IntervalDomain) bool {
	return rhs.lower.Leq(lhs.lower) && lhs.upper.Leq(rhs.upper)
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

	var lower_val algebra.ExtInt
	var upper_val algebra.ExtInt
	if rhs.upper.Leq(lhs.upper) {
		upper_val = rhs.upper
	} else {
		upper_val = algebra.ExtInt_Infty()
	}

	if lhs.lower.Leq(rhs.lower) {
		lower_val = lhs.lower
	} else {
		lower_val = algebra.ExtInt_NegInfty()
	}
	return IntervalDomain{lower: lower_val, upper: upper_val}
}

func (lhs IntervalDomain) Add(rhs IntervalDomain) IntervalDomain {
	// [x1, x2] + [y1, y2] = [x1 + y1, x2 + y2]
	if lhs.is_bottom || rhs.is_bottom {
		return IntervalBot()
	}
	return IntervalDomain{lower: lhs.lower.Add(rhs.lower), upper: lhs.upper.Add(rhs.upper)}
}

func (lhs IntervalDomain) Sub(rhs IntervalDomain) IntervalDomain {
	// [x1, x2] - [y1, y2] = [x1 - y1, x2 - y2]
	if lhs.is_bottom || rhs.is_bottom {
		return IntervalBot()
	}
	return IntervalDomain{lower: lhs.lower.Sub(rhs.lower), upper: lhs.upper.Sub(rhs.upper)}
}

func (lhs IntervalDomain) Mul(rhs IntervalDomain) IntervalDomain {
	// [x1, x2] * [y1, y2] = [min(x1y1, x1y2, x2y1, x2y2) , max(x1y1, x1y2, x2y1, x2y2)]
	x1y1 := lhs.lower.Mul(rhs.lower)
	x1y2 := lhs.lower.Mul(rhs.upper)
	x2y1 := lhs.upper.Mul(rhs.lower)
	x2y2 := lhs.upper.Mul(rhs.upper)
	return IntervalDomain{lower: x1y1.Min(x1y2, x2y1, x2y2), upper: x1y1.Max(x1y2, x2y1, x2y2)}
}
