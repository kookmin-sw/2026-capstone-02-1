package domain

import "traceinspector/algebra"

// The SimplePropDomain domation is a relation domain representing
// ±x ±y ⊙ C or ±x ⊙ C
type SimplePropDomain struct {
	prop           algebra.SimpleProp
	is_bot, is_top bool
}

func (dom SimplePropDomain) IsBot() bool {
	return dom.is_bot
}

func (dom SimplePropDomain) IsTop() bool {
	return dom.is_top
}

func (dom SimplePropDomain) CreateBot() SimplePropDomain {
	return SimplePropDomain{is_bot: true}
}

func (dom SimplePropDomain) CreateTop() SimplePropDomain {
	return SimplePropDomain{is_top: true}
}

// For an LHS proposition ±x ±y ⊙ C and RHS prop ±x' ±y' ⊙ C', Check if
// ±x ±y ⊙_1 C -> ±x' ±y' ⊙_2 C'.
//
// If ⊙_1 != ⊙_2 return false
// For equality and inequality, the prop must be the same, except in the case the constant is zero.
// Then for all signs of x, y, x', y' the prop is the same.
// For inequalities, case analysis the signs and determine whether C <= C' or C >= C'
func (lhs SimplePropDomain) Incl(rhs SimplePropDomain) bool {
	if lhs.prop.Prop_type != rhs.prop.Prop_type {
		return false
	}
	switch lhs.prop.Prop_type {
	case algebra.SimplePropType_Eq:
		if rhs.prop.Prop_type != algebra.SimplePropType_Eq {
			return false
		}
		if lhs.prop.Constant == rhs.prop.Constant {
			if lhs.prop.Constant == 0 {
				return true
			} else {
				return lhs.prop.X_coeff == rhs.prop.Y_coeff && lhs.prop.X_coeff == rhs.prop.Y_coeff
			}
		}
	case algebra.SimplePropType_Neq:
		if rhs.prop.Prop_type != algebra.SimplePropType_Neq {
			return false
		}
		if lhs.prop.Constant == rhs.prop.Constant {
			if lhs.prop.Constant == 0 {
				return true
			} else {
				return lhs.prop.X_coeff == rhs.prop.Y_coeff && lhs.prop.X_coeff == rhs.prop.Y_coeff
			}
		}
	case algebra.SimplePropType_Leq:
		// fill in here
	}
}
