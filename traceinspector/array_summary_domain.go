package traceinspector

import (
	"fmt"
	"traceinspector/domain"
	"traceinspector/imp"
)

type ArraySummaryDomain[IntDomainImpl domain.IntegerDomain[IntDomainImpl]] struct {
	val            IntDomainImpl
	length         IntDomainImpl
	is_bot, is_top bool
}

func (domain ArraySummaryDomain[IntDomainImpl]) String() string {
	if domain.is_bot {
		return "⊥_bool"
	} else if domain.is_top {
		return "⊤_bool"
	} else {
		return fmt.Sprintf("ArraySummary{len: %s, val: %s}", domain.length, domain.val)
	}
}

func (domain ArraySummaryDomain[ElemDomain]) Clone() ArraySummaryDomain[ElemDomain] {
	return ArraySummaryDomain[ElemDomain]{length: domain.length.Clone(), val: domain.val.Clone(), is_bot: domain.is_bot, is_top: domain.is_top}
}

func (domain ArraySummaryDomain[ElemDomain]) IsBot() bool {
	return domain.is_bot
}

func (domain ArraySummaryDomain[ElemDomain]) IsTop() bool {
	return domain.is_top
}

func (lhs ArraySummaryDomain[ElemDomain]) Join(rhs ArraySummaryDomain[ElemDomain]) (ArraySummaryDomain[ElemDomain], bool) {
	if lhs.is_bot {
		return rhs, !rhs.is_bot
	} else if rhs.is_bot {
		return lhs, !lhs.is_bot
	} else if lhs.is_top || rhs.is_top {
		return ArraySummaryDomain[ElemDomain]{is_top: true}, !lhs.is_top
	} else {
		elem_joined, elem_changed := lhs.val.Join(rhs.val)
		len_joined, len_changed := lhs.length.Join(rhs.length)
		return ArraySummaryDomain[ElemDomain]{length: len_joined, val: elem_joined}, elem_changed || len_changed
	}
}

func (lhs ArraySummaryDomain[ElemDomain]) Incl(rhs ArraySummaryDomain[ElemDomain]) bool {
	return lhs.val.Incl(rhs.val)
}

func (lhs ArraySummaryDomain[ElemDomain]) Widen(rhs ArraySummaryDomain[ElemDomain]) ArraySummaryDomain[ElemDomain] {
	if lhs.is_bot {
		return rhs
	}
	if rhs.is_bot {
		return lhs
	}
	if lhs.is_top || rhs.is_top {
		return ArraySummaryDomain[ElemDomain]{is_top: true}
	}
	if lhs.val.Incl(rhs.val) {
		return lhs
	} else {
		return ArraySummaryDomain[ElemDomain]{length: lhs.length.Widen(rhs.length), val: lhs.val.Widen(rhs.val)}
	}
}

// expression evaluation

func (arr ArraySummaryDomain[IntDomainImpl]) GetIndex(val IntDomainImpl) AbstractValue[IntDomainImpl, ArraySummaryDomain[IntDomainImpl]] {
	return AbstractValue[IntDomainImpl, ArraySummaryDomain[IntDomainImpl]]{domain_kind: IntDomainKind, int_domain: arr.val}
}

func (arr ArraySummaryDomain[IntDomainImpl]) Len() IntDomainImpl {
	return arr.length
}

func (arr ArraySummaryDomain[IntDomainImpl]) SetLen(val IntDomainImpl) ArraySummaryDomain[IntDomainImpl] {
	arr.length = val
	return arr
}

func (arr ArraySummaryDomain[IntDomainImpl]) CreateTop() ArraySummaryDomain[IntDomainImpl] {
	return ArraySummaryDomain[IntDomainImpl]{val: arr.length.CreateTop(), length: arr.length.CreateTop(), is_top: true}
}

func (arr ArraySummaryDomain[IntDomainImpl]) CreateBot() ArraySummaryDomain[IntDomainImpl] {
	return ArraySummaryDomain[IntDomainImpl]{val: arr.length.CreateBot(), length: arr.length.CreateBot(), is_bot: true}
}

func (arr ArraySummaryDomain[IntDomainImpl]) From_AbstractValues(vals []AbstractValue[IntDomainImpl, ArraySummaryDomain[IntDomainImpl]]) ArraySummaryDomain[IntDomainImpl] {
	if len(vals) == 0 {
		// use len as a hack for accessing intdomain constructor
		return ArraySummaryDomain[IntDomainImpl]{length: arr.Len().From_IntLitExpr(imp.IntLitExpr{Node: imp.Node{}, Value: 0}), val: arr.length.CreateBot()}
	}
	base := vals[0]
	for index, val := range vals {
		if index == 0 {
			continue
		}
		switch base.domain_kind {
		case IntDomainKind:
			base.int_domain, _ = base.Get_int().Join(val.Get_int())

		case BoolDomainKind:
			panic("arraysummarydomain From_AbstractValues unimplemented")
		case ArrayDomainKind:
			panic("arraysummarydomain From_AbstractValues unimplemented")
		}
	}
	return ArraySummaryDomain[IntDomainImpl]{val: base.Get_int()}
}

func (arr ArraySummaryDomain[IntDomainImpl]) Make_array(len_dom IntDomainImpl, default_value AbstractValue[IntDomainImpl, ArraySummaryDomain[IntDomainImpl]]) ArraySummaryDomain[IntDomainImpl] {
	return ArraySummaryDomain[IntDomainImpl]{length: len_dom, val: default_value.Get_int()}
}

func (arr ArraySummaryDomain[IntDomainImpl]) SetVal(index IntDomainImpl, val AbstractValue[IntDomainImpl, ArraySummaryDomain[IntDomainImpl]]) ArraySummaryDomain[IntDomainImpl] {

	new_val, _ := arr.val.Join(val.Get_int())
	if !index.Leq(arr.Len().Sub(index.From_IntLitExpr(imp.IntLitExpr{Value: 1}))).IsTrue() {
		// potentially unsafe indexing
		return ArraySummaryDomain[IntDomainImpl]{val: new_val.CreateBot(), length: arr.length}
	} else {
		return ArraySummaryDomain[IntDomainImpl]{length: arr.length, val: new_val}
	}
}
