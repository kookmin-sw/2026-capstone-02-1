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
	prop_type SimplePropType  // type of the prop
	x_expr    imp.Expr        // variable or arrayindexop or len()
	x_coeff   SimplePropCoeff // whether coefficient of x is positive
	y_expr    imp.Expr        // same as x_expr
	y_coeff   SimplePropCoeff // whether coefficient of y is positive
	constant  int             // constant value
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

func (ieq SimpleProp) String() string {
	var x_sign, y_sign, prop_op string

	switch ieq.prop_type {
	case SimplePropType_Invalid:
		return "INVALID_SIMPLEPROP"
	case SimplePropType_Eq:
		prop_op = "="
	case SimplePropType_Neq:
		prop_op = "!="
	case SimplePropType_Leq:
		prop_op = "<="
	}
	switch ieq.x_coeff {
	case SimplePropCoeff_negative:
		x_sign = "-"
	}
	switch ieq.y_coeff {
	case SimplePropCoeff_zero:
		return fmt.Sprintf("%s%s %s %d", x_sign, ieq.x_expr, prop_op, ieq.constant)
	case SimplePropCoeff_negative:
		y_sign = "-"
	}

	return fmt.Sprintf("%s%s + %s%s %s %d", x_sign, ieq.x_expr, y_sign, ieq.y_expr, prop_op, ieq.constant)

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
func convert_lhs_to_simple_prop(prop_type SimplePropType, lhs imp.Expr) (SimpleProp, bool) {
	// pull constants out of LHS by representing LHS as Polynomial struct
	lhs_poly, err := build_polynomial(convert_subtraction_to_neg(lhs, false))
	if err != nil {
		return SimpleProp{}, false
	}
	created_prop := SimpleProp{prop_type: prop_type}
	created_prop.constant = -lhs_poly.constant // send constant to other side of leq

	// check if the polynomial is the form `±x + C`
	single_expr, single_coeff := _check_if_var(lhs_poly.variable_expr)
	if single_expr != nil {
		created_prop.x_expr = single_expr
		created_prop.y_coeff = SimplePropCoeff_zero
		created_prop.x_coeff = single_coeff
		return created_prop, true
	}

	// check if the polynomial is the form `±x ±y`
	x_expr, x_coeff, y_expr, y_coeff := _check_binary_expr(lhs_poly.variable_expr)

	if x_expr != nil && y_expr != nil {
		created_prop.x_expr = x_expr
		created_prop.x_coeff = x_coeff
		created_prop.y_expr = y_expr
		created_prop.y_coeff = y_coeff
		return created_prop, true
	} else {
		return SimpleProp{}, false
	}
}

// Given an imp bool expression, try and convert the expression into a SimpleProp.
// Returns SimpleProp, and a boolean indicating whether the conversion was possible.
// Very naive and lazy implementation btw
func imp_expr_to_simple_prop(expr imp.Expr) (SimpleProp, bool) {
	switch expr_ty := expr.(type) {
	case *imp.LessthanExpr:
		// convert to leq
		// lhs < rhs -> lhs <= rhs - 1
		return imp_expr_to_simple_prop(&imp.LeqExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: &imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 1}}})
	case *imp.GreaterthanExpr:
		// lhs > rhs -> rhs < lhs
		return imp_expr_to_simple_prop(&imp.LessthanExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: expr_ty.Lhs})
	case *imp.GeqExpr:
		// lhs >= rhs -> rhs <= lhs
		return imp_expr_to_simple_prop(&imp.LeqExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: expr_ty.Lhs})
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
		return convert_lhs_to_simple_prop(SimplePropType_Leq, zero_expr_leq.Lhs)

	case *imp.EqExpr:
		zero_expr, err := zero_rhs(expr)
		if err != nil {
			return SimpleProp{}, false
		}
		zero_expr_eq, is_eq_expr := zero_expr.(*imp.EqExpr)
		if !is_eq_expr {
			return SimpleProp{}, false
		}

		return convert_lhs_to_simple_prop(SimplePropType_Eq, zero_expr_eq.Lhs)
	case *imp.NeqExpr:
		zero_expr, err := zero_rhs(expr)
		if err != nil {
			return SimpleProp{}, false
		}
		zero_expr_neq, is_neq_expr := zero_expr.(*imp.NeqExpr)
		if !is_neq_expr {
			return SimpleProp{}, false
		}

		return convert_lhs_to_simple_prop(SimplePropType_Neq, zero_expr_neq.Lhs)
	}
	return SimpleProp{}, false
}
