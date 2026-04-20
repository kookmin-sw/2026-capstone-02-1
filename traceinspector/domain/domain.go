package domain

type AbstractDomain[DomainImpl any] interface {
	IsBot() bool
	IsTop() bool
	Incl(DomainImpl) bool        // inclusion operator `lhs ⊑ rhs`
	Join(DomainImpl) DomainImpl  // abstract join operator `lhs ⊔ rhs`
	Widen(DomainImpl) DomainImpl // widening operator `lhs ▽ rhs`
	String() string              // return string representation of the domain value
}
