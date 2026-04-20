package algebra

import "strconv"

// This file defines an "extended" integer type, the integer set Z augmented with positive and negative infinity

type ExtInt struct {
	value              int
	is_inf, is_neg_inf bool
}

func (lhs ExtInt) String() string {
	if lhs.is_inf {
		return "∞"
	} else if lhs.is_neg_inf {
		return "-∞"
	} else {
		return strconv.Itoa(lhs.value)
	}
}

func (eint ExtInt) Value() int {
	return eint.value
}

func ExtInt_Finite(val int) ExtInt {
	return ExtInt{value: val}
}

func ExtInt_Zero() ExtInt {
	return ExtInt{value: 0}
}

func ExtInt_Infty() ExtInt {
	return ExtInt{is_inf: true}
}

func ExtInt_NegInfty() ExtInt {
	return ExtInt{is_neg_inf: true}
}

func (lhs ExtInt) IsFinite() bool {
	return !(lhs.IsInfty() || lhs.IsNegInfty())
}

func (lhs ExtInt) IsInfty() bool {
	return lhs.is_inf
}

func (lhs ExtInt) IsNegInfty() bool {
	return lhs.is_neg_inf
}

func (lhs ExtInt) IsPositive() bool {
	return (lhs.IsFinite() && lhs.value > 0) || lhs.is_inf
}

func (lhs ExtInt) IsNegative() bool {
	return (lhs.IsFinite() && lhs.value < 0) || lhs.is_neg_inf
}

func (lhs ExtInt) Eq(rhs ExtInt) bool {
	return (lhs.is_inf && rhs.is_inf) || (lhs.is_neg_inf && rhs.is_neg_inf) || (lhs.IsFinite() && rhs.IsFinite() && lhs.value == rhs.value)
}

func (lhs ExtInt) Leq(rhs ExtInt) bool {
	// trivial case
	if lhs.is_neg_inf || rhs.is_inf {
		return true
	}
	if lhs.is_inf {
		return rhs.is_inf
	}
	if rhs.is_neg_inf {
		return lhs.is_neg_inf
	}

	// remaining case is lhs, rhs = Z
	return lhs.value <= rhs.value
}

func (lhs ExtInt) Neg() ExtInt {
	if lhs.is_inf {
		return ExtInt_NegInfty()
	} else if lhs.is_neg_inf {
		return ExtInt_Infty()
	} else {
		return ExtInt{value: -lhs.value}
	}
}

func (lhs ExtInt) Add(rhs ExtInt) ExtInt {
	if lhs.is_inf || rhs.is_inf {
		return ExtInt_Infty()
	}
	if lhs.is_neg_inf || rhs.is_neg_inf {
		return ExtInt_NegInfty()
	}
	return ExtInt{value: lhs.value + rhs.value}
}

func (lhs ExtInt) Sub(rhs ExtInt) ExtInt {
	// note that I model infty ± infty as infty, and -infty ± infty as -infty. This is not mathematically correct
	if lhs.is_inf || rhs.is_inf {
		return ExtInt_Infty()
	}
	if lhs.is_neg_inf || rhs.is_neg_inf {
		return ExtInt_NegInfty()
	}
	return ExtInt{value: lhs.value - rhs.value}
}

func (lhs ExtInt) Mul(rhs ExtInt) ExtInt {
	if lhs.Eq(ExtInt_Zero()) || rhs.Eq(ExtInt_Zero()) {
		// Note that 0 * infty is undefined, but we return 0. Again this is not mathematically correct
		return ExtInt_Zero()
	}
	if lhs.IsFinite() && rhs.IsFinite() {
		return ExtInt{value: lhs.value * rhs.value}
	}
	// at least one value is +- inf, so the value is inf; just have to define the sign
	switch lhs.IsPositive() {
	case true:
		switch rhs.IsPositive() {
		case true:
			return ExtInt_Infty()
		case false:
			return ExtInt_NegInfty()
		}
	case false:
		// lhs = -
		switch rhs.IsNegative() {
		case true:
			return ExtInt_Infty()
		case false:
			return ExtInt_NegInfty()
		}
	}
	panic("This should never ever happen")
}

// Compute the minimum of two ExtInts
func (lhs ExtInt) Min(rhs ...ExtInt) ExtInt {
	return_val := lhs
	for _, val := range rhs {
		if val.Leq(return_val) {
			return_val = val
		}
	}
	return return_val
}

// Compute the maximum of two ExtInts
func (lhs ExtInt) Max(rhs ...ExtInt) ExtInt {
	return_val := lhs
	for _, val := range rhs {
		if return_val.Leq(val) {
			return_val = val
		}
	}
	return return_val
}
