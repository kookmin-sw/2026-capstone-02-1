package algebra

import (
	"fmt"
	"traceinspector/imp"
)

// Represent an inequality of the form ±x ±y <= C and ±x <= C, where x and y are addressible.
// x and y must be one of imp.VarExpr, imp.ArrayIndexExpr, imp.LenExpr.
// If y_coeff is SimpleInequalityCoeff_zero, it means y doesn't exist
type SimpleInequality struct {
	x_expr   imp.Expr              // variable or arrayindexop or len()
	x_coeff  SimpleInequalityCoeff // whether coefficient of x is positive
	y_expr   imp.Expr              // same as x_expr
	y_coeff  SimpleInequalityCoeff // whether coefficient of y is positive
	constant int                   // constant value
}

type SimpleInequalityCoeff int

const (
	SimpleInequalityCoeff_zero SimpleInequalityCoeff = iota
	SimpleInequalityCoeff_negative
	SimpleInequalityCoeff_positive
)

func (ieq SimpleInequality) String() string {
	var x_sign, y_sign string
	switch ieq.x_coeff {
	case SimpleInequalityCoeff_negative:
		x_sign = "-"
	}
	switch ieq.y_coeff {
	case SimpleInequalityCoeff_zero:
		return fmt.Sprintf("%s%s <= %d", x_sign, ieq.x_expr, ieq.constant)
	case SimpleInequalityCoeff_negative:
		y_sign = "-"
	}
	return fmt.Sprintf("%s%s + %s%s <= %d", x_sign, ieq.x_expr, y_sign, ieq.y_expr, ieq.constant)

}

// Given an imp.Expr, check if the expr is of the form +-expr. Only used within imp_expr_to_simp_inequality
func _check_if_var(expr imp.Expr) (imp.Expr, SimpleInequalityCoeff) {
	switch expr_ty := expr.(type) {
	case *imp.NegExpr:
		subexpr, sub_coeff := _check_if_var(expr_ty.Subexpr)
		switch sub_coeff {
		case SimpleInequalityCoeff_zero:
			return nil, SimpleInequalityCoeff_zero
		case SimpleInequalityCoeff_negative:
			return subexpr, SimpleInequalityCoeff_positive
		case SimpleInequalityCoeff_positive:
			return subexpr, SimpleInequalityCoeff_negative
		}
	case *imp.VarExpr:
		return expr_ty, SimpleInequalityCoeff_positive
	case *imp.LenExpr:
		return expr_ty, SimpleInequalityCoeff_positive
	case *imp.ArrayIndexExpr:
		return expr_ty, SimpleInequalityCoeff_positive
	}
	return nil, SimpleInequalityCoeff_zero
}

// Also verify that an expression is either a valid simpleineq expr or a negation of it
func _check_binary_expr(expr imp.Expr) (imp.Expr, SimpleInequalityCoeff, imp.Expr, SimpleInequalityCoeff) {
	switch expr_ty := expr.(type) {
	case *imp.AddExpr:
		lhs_expr, lhs_coeff := _check_if_var(expr_ty.Lhs)
		rhs_expr, rhs_coeff := _check_if_var(expr_ty.Rhs)
		return lhs_expr, lhs_coeff, rhs_expr, rhs_coeff
	case *imp.ParenExpr:
		return _check_binary_expr(expr_ty.Subexpr)
	}
	return nil, SimpleInequalityCoeff_zero, nil, SimpleInequalityCoeff_zero
}

// Given an imp leq expression, try and convert the expression into a SimpleInequality.
// Returns SimpleInequality, and a boolean indicating whether the conversion was possible.
// Very naive and lazy implementation btw
func imp_expr_to_simp_inequality(expr imp.Expr) (SimpleInequality, bool) {
	switch expr_ty := expr.(type) {
	case *imp.LessthanExpr:
		// convert to leq
		// lhs < rhs -> lhs <= rhs - 1
		return imp_expr_to_simp_inequality(&imp.LeqExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: &imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 1}}})
	case *imp.GreaterthanExpr:
		// lhs > rhs -> rhs < lhs
		return imp_expr_to_simp_inequality(&imp.LessthanExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: expr_ty.Lhs})
	case *imp.GeqExpr:
		// lhs >= rhs -> rhs <= lhs
		return imp_expr_to_simp_inequality(&imp.LeqExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: expr_ty.Lhs})
	case *imp.LeqExpr:
		// move all terms to lhs
		zero_expr, err := zero_rhs(expr)
		if err != nil {
			return SimpleInequality{}, false
		}
		zero_expr_leq, is_leq_expr := zero_expr.(*imp.LeqExpr)
		if !is_leq_expr {
			return SimpleInequality{}, false
		}

		// pull constants out of LHS by representing LHS as Polynomial struct
		lhs_poly, err := build_polynomial(convert_subtraction_to_neg(zero_expr_leq.Lhs, false))
		// fmt.Println(expr, "->", zero_expr_leq.Lhs, "->", convert_subtraction_to_neg(zero_expr_leq.Lhs, false), "||", lhs_poly.variable_expr, lhs_poly.constant)
		created_ineq := SimpleInequality{}
		created_ineq.constant = -lhs_poly.constant // send constant to other side of leq

		// check if the polynomial is the form `±x + C`
		single_expr, single_coeff := _check_if_var(lhs_poly.variable_expr)
		if single_expr != nil {
			created_ineq.x_expr = single_expr
			created_ineq.y_coeff = SimpleInequalityCoeff_zero
			created_ineq.x_coeff = single_coeff
			return created_ineq, true
		}

		// check if the polynomial is the form `±x ±y`
		x_expr, x_coeff, y_expr, y_coeff := _check_binary_expr(lhs_poly.variable_expr)

		if x_expr != nil && y_expr != nil {
			created_ineq.x_expr = x_expr
			created_ineq.x_coeff = x_coeff
			created_ineq.y_expr = y_expr
			created_ineq.y_coeff = y_coeff
			return created_ineq, true
		}
	}
	return SimpleInequality{}, false
}
