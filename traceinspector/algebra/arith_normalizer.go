package algebra

import (
	"fmt"
	"traceinspector/imp"
)

// Given an integer equality/inequality, rewrite to canonical form
// 1. e1 <= e2 -> e1 - e2 <= 0  I named this as zero-rhs form
// 2. ax + by + c <= 0 where a and b are integer constants and x y are identifiers

// Given an integer (in)equality expression of the form `e1 ☉ e2“, convert to `e1 - (e2) ☉ 0`,
// where ☉ is a comparison operator int -> int -> bool.
func zero_rhs(expr imp.Expr) (imp.Expr, error) {
	switch expr_ty := expr.(type) {
	case *imp.EqExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.EqExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	case *imp.NeqExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.NeqExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	case *imp.LessthanExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.LessthanExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	case *imp.GreaterthanExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.GreaterthanExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	case *imp.LeqExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.LeqExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	case *imp.GeqExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.GeqExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	default:
		return nil, fmt.Errorf("zero_rhs: Unsupported boolean expression %s", expr)
	}
}

// --e -> e
// e - e -> e + -e
// Given an intege expression, convert all subtraction into addition, and simplify any double negations
func convert_subtraction_to_neg(expr imp.Expr, negate bool) imp.Expr {
	switch expr_ty := expr.(type) {
	case *imp.VarExpr, *imp.ArrayIndexExpr, *imp.LenExpr:
		if negate {
			return &imp.NegExpr{Node: imp.Node{Line_num: expr.GetLineNum()}, Subexpr: expr_ty}
		} else {
			return expr
		}
	case *imp.IntLitExpr:
		if negate {
			return &imp.IntLitExpr{Node: expr_ty.Node, Value: -expr_ty.Value}
		} else {
			return expr
		}
	case *imp.NegExpr:
		if negate {
			return expr_ty.Subexpr
		} else {
			return &imp.NegExpr{Node: expr_ty.Node, Subexpr: convert_subtraction_to_neg(expr_ty.Subexpr, true)}
		}
	case *imp.AddExpr:
		return &imp.AddExpr{Node: expr_ty.Node, Lhs: convert_subtraction_to_neg(expr_ty.Lhs, negate), Rhs: convert_subtraction_to_neg(expr_ty.Rhs, negate)}
	case *imp.SubExpr:
		// - (lhs - rhs) -> -lhs + rhs
		// (lhs - rhs) -> lhs + -rhs
		if negate {
			return &imp.AddExpr{Node: expr_ty.Node, Lhs: convert_subtraction_to_neg(expr_ty.Lhs, true), Rhs: convert_subtraction_to_neg(expr_ty.Rhs, false)}
		} else {
			return &imp.AddExpr{Node: expr_ty.Node, Lhs: convert_subtraction_to_neg(expr_ty.Lhs, false), Rhs: convert_subtraction_to_neg(expr_ty.Rhs, true)}
		}
	case *imp.MulExpr:
		// multiplication, division should propogate sign to one of its arguments
		return &imp.MulExpr{Node: expr_ty.Node, Lhs: convert_subtraction_to_neg(expr_ty.Lhs, negate), Rhs: convert_subtraction_to_neg(expr_ty.Rhs, false)}
	case *imp.DivExpr:
		return &imp.DivExpr{Node: expr_ty.Node, Lhs: convert_subtraction_to_neg(expr_ty.Lhs, negate), Rhs: convert_subtraction_to_neg(expr_ty.Rhs, false)}
	case *imp.ModExpr:
		return &imp.ModExpr{Node: expr_ty.Node, Lhs: convert_subtraction_to_neg(expr_ty.Lhs, negate), Rhs: convert_subtraction_to_neg(expr_ty.Rhs, false)}
	case *imp.ParenExpr:
		// parenexpr does though
		return &imp.ParenExpr{Node: expr_ty.Node, Subexpr: convert_subtraction_to_neg(expr_ty.Subexpr, negate)}
	default:
		return expr_ty
	}
}

// Represents a linear arithmetic Polynomial ax ⊙ by ⊙ ... ⊙ cz + C, where a, b, c are coefficient exprs
// and x, y, z are addressible.
//
// variable_expr: ax ⊙ by ⊙ ... ⊙ cz
// constant: C
type Polynomial struct {
	variable_expr imp.Expr
	constant      int
}

// Build the normalized polynomial representation of the integer expression
func build_polynomial(expr imp.Expr) (Polynomial, error) {
	accumulated := Polynomial{}
	switch expr_ty := expr.(type) {
	case *imp.VarExpr:
		accumulated.variable_expr = expr_ty
	case *imp.IntLitExpr:
		accumulated.constant = expr_ty.Value
	case *imp.ArrayLitExpr:
		accumulated.variable_expr = expr_ty
	case *imp.ArrayIndexExpr:
		accumulated.variable_expr = expr_ty
	case *imp.MakeArrayExpr:
		accumulated.variable_expr = expr_ty
	case *imp.LenExpr:
		accumulated.variable_expr = expr_ty
	case *imp.CallExpr:
		accumulated.variable_expr = expr_ty
	case *imp.AddExpr:
		lhs_poly, err := build_polynomial(expr_ty.Lhs)
		if err != nil {
			return Polynomial{}, err
		}
		rhs_poly, err := build_polynomial(expr_ty.Rhs)
		if err != nil {
			return Polynomial{}, err
		}
		if lhs_poly.variable_expr == nil {
			rhs_poly.constant += lhs_poly.constant
			return rhs_poly, nil
		} else if rhs_poly.variable_expr == nil {
			lhs_poly.constant += rhs_poly.constant
			return lhs_poly, nil
		} else {
			accumulated.variable_expr = &imp.AddExpr{Node: expr_ty.Node, Lhs: lhs_poly.variable_expr, Rhs: rhs_poly.variable_expr}
			accumulated.constant = lhs_poly.constant + rhs_poly.constant
		}
	case *imp.SubExpr:
		lhs_poly, err := build_polynomial(expr_ty.Lhs)
		if err != nil {
			return Polynomial{}, err
		}
		rhs_poly, err := build_polynomial(expr_ty.Rhs)
		if err != nil {
			return Polynomial{}, err
		}
		if lhs_poly.variable_expr == nil {
			rhs_poly.constant += lhs_poly.constant
			return rhs_poly, nil
		} else if rhs_poly.variable_expr == nil {
			lhs_poly.constant -= rhs_poly.constant
			return lhs_poly, nil
		} else {
			accumulated.variable_expr = &imp.SubExpr{Node: expr_ty.Node, Lhs: lhs_poly.variable_expr, Rhs: rhs_poly.variable_expr}
			accumulated.constant = lhs_poly.constant - rhs_poly.constant
		}
	case *imp.MulExpr:
		// For the case of multiplication
		lhs_poly, err := build_polynomial(expr_ty.Lhs)
		if err != nil {
			return Polynomial{}, err
		}
		rhs_poly, err := build_polynomial(expr_ty.Rhs)
		if err != nil {
			return Polynomial{}, err
		}
		if lhs_poly.variable_expr == nil && rhs_poly.variable_expr == nil {
			// both subexprs are constants
			accumulated.constant = lhs_poly.constant * rhs_poly.constant

			// Don't do "constant folding" for subexprs yet
			// } else if lhs_poly.variable_expr == nil {
			// 	// LHS is constant, but RHS isn't so LHS should be used as coefficient
			// 	accumulated.variable_expr = &imp.MulExpr{Node: expr_ty.Node, Lhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: lhs_poly.constant}, Rhs: expr_ty.Rhs}
			// } else if rhs_poly.variable_expr == nil {
			// 	// same goes for RHS
			// 	accumulated.variable_expr = &imp.MulExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: rhs_poly.constant}}
			// } else {
		} else {
			// if both are not constants, return the original
			accumulated.variable_expr = expr_ty
		}
	case *imp.DivExpr:
		lhs_poly, err := build_polynomial(expr_ty.Lhs)
		if err != nil {
			return Polynomial{}, err
		}
		rhs_poly, err := build_polynomial(expr_ty.Rhs)
		if err != nil {
			return Polynomial{}, err
		}
		if lhs_poly.variable_expr == nil && rhs_poly.variable_expr == nil {
			accumulated.constant = lhs_poly.constant / rhs_poly.constant
		} else {
			accumulated.variable_expr = expr_ty
		}
	case *imp.ModExpr:
		lhs_poly, err := build_polynomial(expr_ty.Lhs)
		if err != nil {
			return Polynomial{}, err
		}
		rhs_poly, err := build_polynomial(expr_ty.Rhs)
		if err != nil {
			return Polynomial{}, err
		}
		if lhs_poly.variable_expr == nil && rhs_poly.variable_expr == nil {
			accumulated.constant = lhs_poly.constant % rhs_poly.constant
		} else {
			accumulated.variable_expr = expr_ty
		}
	case *imp.NegExpr:
		sub_poly, err := build_polynomial(expr_ty.Subexpr)
		if err != nil {
			return Polynomial{}, err
		}
		if sub_poly.variable_expr != nil {
			accumulated.variable_expr = &imp.NegExpr{Node: expr_ty.Node, Subexpr: sub_poly.variable_expr}

		}
		accumulated.constant = -sub_poly.constant
	case *imp.ParenExpr:
		sub_poly, err := build_polynomial(expr_ty.Subexpr)
		if err != nil {
			return Polynomial{}, err
		}
		if sub_poly.variable_expr != nil {
			accumulated.variable_expr = &imp.ParenExpr{Node: expr_ty.Node, Subexpr: sub_poly.variable_expr}

		}
		accumulated.constant = sub_poly.constant
	default:
		return Polynomial{}, fmt.Errorf("build_polynomial: unsupported expressions %s", expr_ty)
	}
	return accumulated, nil
}

// Given an arbitrary integer expression, normalize to the form
// ax ☉ by ☉ ... ☉ cz + C, where C is an integer constant
func normalize_integer_expr(expr imp.Expr) (imp.Expr, error) {
	poly, err := build_polynomial(convert_subtraction_to_neg(expr, false))
	if err != nil {
		return nil, err
	} else {
		return &imp.AddExpr{Node: imp.Node{Line_num: expr.GetLineNum()}, Lhs: poly.variable_expr, Rhs: &imp.IntLitExpr{Node: imp.Node{Line_num: expr.GetLineNum()}, Value: poly.constant}}, nil
	}
}
