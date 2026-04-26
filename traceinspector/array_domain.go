package traceinspector

import (
	"traceinspector/domain"
)

type ArrayDomain[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] interface {
	domain.AbstractDomain[ArrayDomainImpl]
	From_AbstractValues([]AbstractValue[IntDomainImpl, ArrayDomainImpl]) ArrayDomainImpl
	GetIndex(IntDomainImpl) AbstractValue[IntDomainImpl, ArrayDomainImpl]
	SetVal(IntDomainImpl, AbstractValue[IntDomainImpl, ArrayDomainImpl]) ArrayDomainImpl
	Len() IntDomainImpl
	SetLen(IntDomainImpl) ArrayDomainImpl
	Make_array(IntDomainImpl, AbstractValue[IntDomainImpl, ArrayDomainImpl]) ArrayDomainImpl
}
