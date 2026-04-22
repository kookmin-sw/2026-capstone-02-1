package traceinspector

import (
	"traceinspector/domain"
	"traceinspector/imp"
)

type ArrayDomain[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] interface {
	domain.AbstractDomain[ArrayDomainImpl]
	From_ArrayLitExpr(imp.ArrayLitExpr) ArrayDomainImpl
	Index(IntDomainImpl) AbstractValue[IntDomainImpl, ArrayDomainImpl]
	Len() IntDomainImpl
	Make_array(IntDomainImpl, AbstractValue[IntDomainImpl, ArrayDomainImpl]) ArrayDomainImpl
}
