package domain

import "traceinspector/imp"

type AbstractDomain[DomainImpl any] interface {
	IsBot() bool
	IsTop() bool
	Incl(DomainImpl) bool        // inclusion operator `lhs ⊑ rhs`
	Join(DomainImpl) DomainImpl  // abstract join operator `lhs ⊔ rhs`
	Widen(DomainImpl) DomainImpl // widening operator `lhs ▽ rhs`
	String() string              // return string representation of the domain value
}

type IntegerDomain[DomainImpl any] interface {
	AbstractDomain[DomainImpl]
	From_IntLitExpr(imp.IntLitExpr) DomainImpl
	Add(DomainImpl) DomainImpl
	Sub(DomainImpl) DomainImpl
	Mul(DomainImpl) DomainImpl
	Div(DomainImpl) DomainImpl
	Mod(DomainImpl) DomainImpl
	Eq(DomainImpl) BoolDomain
	NEq(DomainImpl) BoolDomain
	Lessthan(DomainImpl) BoolDomain
	Greaterthan(DomainImpl) BoolDomain
	Leq(DomainImpl) BoolDomain
	Geq(DomainImpl) BoolDomain
	Neg(DomainImpl) DomainImpl
}
