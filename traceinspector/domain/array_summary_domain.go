package domain

import "fmt"

type ArraySummaryDomain[ElemDomain AbstractDomain[ElemDomain]] struct {
	val               ElemDomain
	is_bottom, is_top bool
}

func (domain ArraySummaryDomain[ElemDomain]) String() string {
	if domain.is_bottom {
		return "⊥_bool"
	} else if domain.is_top {
		return "⊤_bool"
	} else {
		return fmt.Sprintf("%s", domain.val)
	}
}

func (domain ArraySummaryDomain[ElemDomain]) IsBot() bool {
	return domain.is_bottom
}

func (domain ArraySummaryDomain[ElemDomain]) IsTop() bool {
	return domain.is_top
}

func (lhs ArraySummaryDomain[ElemDomain]) Join(rhs ArraySummaryDomain[ElemDomain]) ArraySummaryDomain[ElemDomain] {
	if lhs.is_bottom {
		return rhs
	} else if rhs.is_bottom {
		return lhs
	} else if lhs.is_top || rhs.is_top {
		return ArraySummaryDomain[ElemDomain]{is_top: true}
	} else {
		return ArraySummaryDomain[ElemDomain]{val: lhs.val.Join(rhs.val)}
	}
}

func (lhs ArraySummaryDomain[ElemDomain]) Incl(rhs ArraySummaryDomain[ElemDomain]) bool {
	return lhs.val.Incl(rhs.val)
}

func (lhs ArraySummaryDomain[ElemDomain]) Widen(rhs ArraySummaryDomain[ElemDomain]) ArraySummaryDomain[ElemDomain] {
	if lhs.is_bottom {
		return rhs
	}
	if rhs.is_bottom {
		return lhs
	}
	if lhs.is_top || rhs.is_top {
		return ArraySummaryDomain[ElemDomain]{is_top: true}
	}
	if lhs.val.Incl(rhs.val) {
		return lhs
	} else {
		return ArraySummaryDomain[ElemDomain]{val: lhs.val.Widen(rhs.val)}
	}
}

// expression evaluation
