package algebra

import (
	"fmt"
	"traceinspector/imp"
)

// Represent a prop of the form ±x ±y ⊙ C and ±x ⊙ C, where x and y are addressible, and
// ⊙ is either "=", "!=", or "<="
// x and y must be one of imp.VarExpr, imp.ArrayIndexExpr, imp.LenExpr.
// If y_coeff is SimpleInequalityCoeff_zero, it means y doesn't exist
type SimpleProp struct {
	Prop_type SimplePropType  // type of the prop
	X_expr    imp.Expr        // variable or arrayindexop or len()
	X_coeff   SimplePropCoeff // whether coefficient of x is positive
	Y_expr    imp.Expr        // same as x_expr
	Y_coeff   SimplePropCoeff // whether coefficient of y is positive
	Constant  int             // constant value
}

type SimplePropType int

const (
	SimplePropType_Invalid SimplePropType = iota
	SimplePropType_Eq
	SimplePropType_Neq
	SimplePropType_Leq
)

type SimplePropCoeff int

const (
	SimplePropCoeff_zero SimplePropCoeff = iota
	SimplePropCoeff_negative
	SimplePropCoeff_positive
)

func (coeff SimplePropCoeff) Negate() SimplePropCoeff {
	switch coeff {
	case SimplePropCoeff_zero:
		return SimplePropCoeff_zero
	case SimplePropCoeff_negative:
		return SimplePropCoeff_positive
	case SimplePropCoeff_positive:
		return SimplePropCoeff_negative
	}
	panic("Negate(): this should never happen(cosmic ray bitflip)")
}

func (ieq SimpleProp) String() string {
	var x_sign, y_sign, prop_op string

	switch ieq.Prop_type {
	case SimplePropType_Invalid:
		return "INVALID_SIMPLEPROP"
	case SimplePropType_Eq:
		prop_op = "="
	case SimplePropType_Neq:
		prop_op = "!="
	case SimplePropType_Leq:
		prop_op = "<="
	}
	switch ieq.X_coeff {
	case SimplePropCoeff_negative:
		x_sign = "-"
	}
	switch ieq.Y_coeff {
	case SimplePropCoeff_zero:
		return fmt.Sprintf("%s%s %s %d", x_sign, ieq.X_expr, prop_op, ieq.Constant)
	case SimplePropCoeff_negative:
		y_sign = "-"
	}

	return fmt.Sprintf("%s%s + %s%s %s %d", x_sign, ieq.X_expr, y_sign, ieq.Y_expr, prop_op, ieq.Constant)

}

// Given an imp.Expr, check if the expr is of the form +-expr. Only used within imp_expr_to_simp_inequality
func _check_if_var(expr imp.Expr) (imp.Expr, SimplePropCoeff) {
	switch expr_ty := expr.(type) {
	case *imp.NegExpr:
		subexpr, sub_coeff := _check_if_var(expr_ty.Subexpr)
		switch sub_coeff {
		case SimplePropCoeff_zero:
			return nil, SimplePropCoeff_zero
		case SimplePropCoeff_negative:
			return subexpr, SimplePropCoeff_positive
		case SimplePropCoeff_positive:
			return subexpr, SimplePropCoeff_negative
		}
	case *imp.VarExpr:
		return expr_ty, SimplePropCoeff_positive
	case *imp.LenExpr:
		return expr_ty, SimplePropCoeff_positive
	case *imp.ArrayIndexExpr:
		return expr_ty, SimplePropCoeff_positive
	}
	return nil, SimplePropCoeff_zero
}

// Also verify that an expression is either a valid simpleprop expr or a negation of it
func _check_binary_expr(expr imp.Expr) (imp.Expr, SimplePropCoeff, imp.Expr, SimplePropCoeff) {
	switch expr_ty := expr.(type) {
	case *imp.AddExpr:
		lhs_expr, lhs_coeff := _check_if_var(expr_ty.Lhs)
		rhs_expr, rhs_coeff := _check_if_var(expr_ty.Rhs)
		return lhs_expr, lhs_coeff, rhs_expr, rhs_coeff
	case *imp.ParenExpr:
		return _check_binary_expr(expr_ty.Subexpr)
	}
	return nil, SimplePropCoeff_zero, nil, SimplePropCoeff_zero
}

// given a LHS expr of a integer prop, try to convert the expr into the SimpleProp of the prop_type
func _convert_lhs_to_simple_prop(prop_type SimplePropType, lhs imp.Expr) (SimpleProp, bool) {
	// pull constants out of LHS by representing LHS as Polynomial struct
	lhs_poly, err := build_polynomial(Convert_subtraction_to_neg(lhs, false))
	if err != nil {
		return SimpleProp{}, false
	}
	created_prop := SimpleProp{Prop_type: prop_type}
	created_prop.Constant = -lhs_poly.constant // send constant to other side of leq

	// check if the polynomial is the form `±x + C`
	single_expr, single_coeff := _check_if_var(lhs_poly.variable_expr)
	if single_expr != nil {
		created_prop.X_expr = single_expr
		created_prop.Y_coeff = SimplePropCoeff_zero
		created_prop.X_coeff = single_coeff
		return created_prop, true
	}

	// check if the polynomial is the form `±x ±y`
	x_expr, x_coeff, y_expr, y_coeff := _check_binary_expr(lhs_poly.variable_expr)

	if x_expr != nil && y_expr != nil {
		created_prop.X_expr = x_expr
		created_prop.X_coeff = x_coeff
		created_prop.Y_expr = y_expr
		created_prop.Y_coeff = y_coeff
		return created_prop, true
	} else {
		return SimpleProp{}, false
	}
}

// Given an imp bool expression, try and convert the expression into a SimpleProp.
// Returns SimpleProp, and a boolean indicating whether the conversion was possible.
// Very naive and lazy implementation btw
func Imp_expr_to_simple_prop(expr imp.Expr) (SimpleProp, bool) {
	switch expr_ty := expr.(type) {
	case *imp.LessthanExpr:
		// convert to leq
		// lhs < rhs -> lhs <= rhs - 1
		return Imp_expr_to_simple_prop(&imp.LeqExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: &imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 1}}})
	case *imp.GreaterthanExpr:
		// lhs > rhs -> rhs < lhs
		return Imp_expr_to_simple_prop(&imp.LessthanExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: expr_ty.Lhs})
	case *imp.GeqExpr:
		// lhs >= rhs -> rhs <= lhs
		return Imp_expr_to_simple_prop(&imp.LeqExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: expr_ty.Lhs})
	case *imp.LeqExpr:
		// move all terms to lhs
		zero_expr, err := zero_rhs(expr)
		if err != nil {
			return SimpleProp{}, false
		}
		zero_expr_leq, is_leq_expr := zero_expr.(*imp.LeqExpr)
		if !is_leq_expr {
			return SimpleProp{}, false
		}
		return _convert_lhs_to_simple_prop(SimplePropType_Leq, zero_expr_leq.Lhs)

	case *imp.EqExpr:
		zero_expr, err := zero_rhs(expr)
		if err != nil {
			return SimpleProp{}, false
		}
		zero_expr_eq, is_eq_expr := zero_expr.(*imp.EqExpr)
		if !is_eq_expr {
			return SimpleProp{}, false
		}

		return _convert_lhs_to_simple_prop(SimplePropType_Eq, zero_expr_eq.Lhs)
	case *imp.NeqExpr:
		zero_expr, err := zero_rhs(expr)
		if err != nil {
			return SimpleProp{}, false
		}
		zero_expr_neq, is_neq_expr := zero_expr.(*imp.NeqExpr)
		if !is_neq_expr {
			return SimpleProp{}, false
		}

		return _convert_lhs_to_simple_prop(SimplePropType_Neq, zero_expr_neq.Lhs)
	}
	return SimpleProp{}, false
}

// Compute the negated form of the SimpleProp
func (sp SimpleProp) Negate() SimpleProp {
	switch sp.Prop_type {
	case SimplePropType_Eq:
		sp.Prop_type = SimplePropType_Neq
	case SimplePropType_Neq:
		sp.Prop_type = SimplePropType_Eq
	case SimplePropType_Leq:
		// !(±x ±y <= C) = ±x ±y > C = ±x ±y >= C + 1 = ∓x ∓y <= -C - 1
		return SimpleProp{
			Prop_type: SimplePropType_Leq,
			X_expr:    sp.X_expr,
			X_coeff:   sp.X_coeff.Negate(),
			Y_expr:    sp.Y_expr,
			Y_coeff:   sp.Y_coeff.Negate(),
			Constant:  -sp.Constant - 1,
		}
	}
	return sp
}

// Force the X_coeff to be non-negative, adjusting the signs of other terms and the prop type
// e.g) -x + y <= 3 -> x - y >= -3
func (sp SimpleProp) Set_xcoeff_positive() SimpleProp {
	switch sp.X_coeff {
	case SimplePropCoeff_zero, SimplePropCoeff_positive:
		return sp
	}
	//x_coeff is negative
	switch sp.Prop_type {
	case SimplePropType_Invalid:
		return sp
	case SimplePropType_Eq:
		return SimpleProp{
			Prop_type: SimplePropType_Eq,
			X_expr:    sp.X_expr,
			X_coeff:   sp.X_coeff.Negate(),
			Y_expr:    sp.Y_expr,
			Y_coeff:   sp.Y_coeff.Negate(),
			Constant:  -sp.Constant,
		}
	}
}
