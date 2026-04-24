package domain

import (
	"fmt"
	"traceinspector/algebra"
	"traceinspector/imp"
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

// Check if the domain is valid, and if upper < lower, change to bot
func (domain IntervalDomain) CheckValid() IntervalDomain {
	// upper <= lower - 1
	if !domain.is_bottom && domain.upper.Leq(domain.lower.Sub(algebra.ExtInt_Finite(1))) {
		return IntervalBot()
	}
	return domain
}

func (domain IntervalDomain) Clone() IntervalDomain {
	return IntervalDomain{lower: domain.lower, upper: domain.upper, is_bottom: domain.is_bottom}
}

func IntervalBot() IntervalDomain {
	return IntervalDomain{is_bottom: true}
}

func IntervalTop() IntervalDomain {
	return IntervalDomain{lower: algebra.ExtInt_NegInfty(), upper: algebra.ExtInt_Infty()}
}

func (domain IntervalDomain) IsBot() bool {
	return domain.is_bottom
}

func (domain IntervalDomain) IsTop() bool {
	return domain.lower.IsNegInfty() && domain.upper.IsInfty()
}

func (domain IntervalDomain) Is_bounded() bool {
	return !domain.IsBot() && domain.lower.IsFinite() && domain.upper.IsFinite()
}

func (lhs IntervalDomain) Join(rhs IntervalDomain) (IntervalDomain, bool) {
	if lhs.is_bottom {
		return rhs, rhs.IsBot()
	}
	if rhs.is_bottom {
		return lhs, false
	}
	return IntervalDomain{lower: lhs.lower.Min(rhs.lower), upper: lhs.upper.Max(rhs.upper)}, !(lhs.lower.Eq(lhs.lower.Min(rhs.lower)) && lhs.upper.Eq(lhs.upper.Max(rhs.upper)))
}

func (lhs IntervalDomain) Intersection(rhs IntervalDomain) IntervalDomain {
	if lhs.Disjoint(rhs) {
		return IntervalBot()
	}
	return IntervalDomain{lower: lhs.lower.Max(rhs.lower), upper: lhs.upper.Min(rhs.upper)}.CheckValid()
}

// `lhs ⊑ rhs` = lhs.lower >= rhs.lower && lhs.upper <= rhs.upper
func (lhs IntervalDomain) Incl(rhs IntervalDomain) bool {
	if lhs.IsBot() {
		return true
	}
	if rhs.IsBot() {
		return lhs.IsBot()
	}
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

func (lhs IntervalDomain) From_IntLitExpr(expr imp.IntLitExpr) IntervalDomain {
	return IntervalDomain{lower: algebra.ExtInt_Finite(expr.Value), upper: algebra.ExtInt_Finite(expr.Value)}
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
	if lhs.is_bottom || rhs.is_bottom {
		return IntervalBot()
	}
	x1y1 := lhs.lower.Mul(rhs.lower)
	x1y2 := lhs.lower.Mul(rhs.upper)
	x2y1 := lhs.upper.Mul(rhs.lower)
	x2y2 := lhs.upper.Mul(rhs.upper)
	return IntervalDomain{lower: x1y1.Min(x1y2, x2y1, x2y2), upper: x1y1.Max(x1y2, x2y1, x2y2)}
}

func (lhs IntervalDomain) Div(rhs IntervalDomain) IntervalDomain {
	if lhs.is_bottom || rhs.is_bottom {
		return IntervalBot()
	}
	return IntervalDomain{lower: algebra.ExtInt_NegInfty(), upper: algebra.ExtInt_Infty()}
}

func (lhs IntervalDomain) Mod(rhs IntervalDomain) IntervalDomain {
	if lhs.is_bottom || rhs.is_bottom {
		return IntervalBot()
	}
	return IntervalDomain{lower: algebra.ExtInt_NegInfty(), upper: algebra.ExtInt_Infty()}
}

func (lhs IntervalDomain) Neg() IntervalDomain {
	if lhs.is_bottom {
		return IntervalBot()
	}
	return IntervalDomain{lower: lhs.upper.Neg(), upper: lhs.lower.Neg()}
}

// Returns whether the two domains are disjoint(do they not overlap?)
func (lhs IntervalDomain) Disjoint(rhs IntervalDomain) bool {
	// [--------]           [-------]
	//            [------]
	if lhs.is_bottom || rhs.is_bottom {
		return true
	}
	return lhs.upper.Leq(rhs.lower.Sub(algebra.ExtInt_Finite(1))) || rhs.upper.Leq(lhs.lower.Add(algebra.ExtInt_Finite(1)))
}

func (lhs IntervalDomain) Eq(rhs IntervalDomain) BoolDomain {
	if lhs.IsBot() || rhs.IsBot() {
		return BoolDomain{is_bottom: true}
	}
	if lhs.Disjoint(rhs) {
		return BoolDomain{val: false}
	}
	// [x, x] = [x, x]
	if lhs.lower.IsFinite() && rhs.lower.IsFinite() && lhs.upper.IsFinite() && rhs.upper.IsFinite() && lhs.lower.Eq(lhs.upper) && lhs.lower.Eq(rhs.lower) && lhs.upper.Eq(rhs.upper) {
		return BoolDomain{val: true}
	}
	// sound
	return BoolDomain{is_top: true}
}

func (lhs IntervalDomain) Neq(rhs IntervalDomain) BoolDomain {
	if lhs.IsBot() || rhs.IsBot() {
		return BoolDomain{is_bottom: true}
	}
	if lhs.Disjoint(rhs) {
		return BoolDomain{val: true}
	}
	// [x, x] == [x, x] <-> !([x, x] == [x, x]) = false
	if lhs.lower.IsFinite() && rhs.lower.IsFinite() && lhs.upper.IsFinite() && rhs.upper.IsFinite() && lhs.lower.Eq(lhs.upper) && lhs.lower.Eq(rhs.lower) && lhs.upper.Eq(rhs.upper) {
		return BoolDomain{val: false}
	}
	return BoolDomain{is_top: true}
}

func (lhs IntervalDomain) Geq(rhs IntervalDomain) BoolDomain {
	if lhs.IsBot() || rhs.IsBot() {
		return BoolDomain{is_bottom: true}
	}
	if !lhs.Disjoint(rhs) {
		return BoolDomain{is_top: true}
	}
	// [x, y] >= [a, b] <-> b <= x /\ b, x are finite
	if lhs.lower.IsFinite() && rhs.upper.IsFinite() && rhs.upper.Leq(lhs.lower) {
		return BoolDomain{val: true}
	}
	return BoolDomain{is_top: true}
}

func (lhs IntervalDomain) Greaterthan(rhs IntervalDomain) BoolDomain {
	if lhs.IsBot() || rhs.IsBot() {
		return BoolDomain{is_bottom: true}
	}
	if !lhs.Disjoint(rhs) {
		return BoolDomain{is_top: true}
	}
	// [x, y] > [a, b] <-> b <= x - 1 /\ b,x are finite
	if lhs.lower.IsFinite() && rhs.upper.IsFinite() && rhs.upper.Leq(lhs.lower.Sub(algebra.ExtInt_Finite(1))) {
		return BoolDomain{val: true}
	}
	return BoolDomain{is_top: true}
}

func (lhs IntervalDomain) Leq(rhs IntervalDomain) BoolDomain {
	if lhs.IsBot() || rhs.IsBot() {
		return BoolDomain{is_bottom: true}
	}
	if !lhs.Disjoint(rhs) {
		return BoolDomain{is_top: true}
	}
	// [x, y] <= [a, b] <-> y <= a /\ y, a are finite
	if lhs.upper.IsFinite() && rhs.lower.IsFinite() && lhs.upper.Leq(rhs.lower) {
		return BoolDomain{val: true}
	}
	return BoolDomain{is_top: true}
}

func (lhs IntervalDomain) Lessthan(rhs IntervalDomain) BoolDomain {
	if lhs.IsBot() || rhs.IsBot() {
		return BoolDomain{is_bottom: true}
	}
	if !lhs.Disjoint(rhs) {
		return BoolDomain{is_top: true}
	}
	// [x, y] < [a, b] <-> y <= a - 1 /\ y, a are finite
	if lhs.upper.IsFinite() && rhs.lower.IsFinite() && lhs.upper.Leq(rhs.lower.Sub(algebra.ExtInt_Finite(1))) {
		return BoolDomain{val: true}
	}
	return BoolDomain{is_top: true}
}

func (lhs IntervalDomain) Filter(filter_type FilterQueryType, rhs IntervalDomain) IntervalDomain {
	switch filter_type {
	case FilterQueryType_Eq:
		return lhs.Intersection(rhs)
	case FilterQueryType_Neq:
		// imprecise?
		return lhs
	case FilterQueryType_Leq:
		// fmt.Println("Leq filter interval lhs:", lhs, "rhs:", rhs)
		lhs.upper = lhs.upper.Min(rhs.upper)
	case FilterQueryType_Geq:
		// fmt.Println("Geq filter interval lhs:", lhs, "rhs:", rhs)
		lhs.lower = lhs.lower.Max(rhs.lower)
	}
	return lhs.CheckValid()
}
