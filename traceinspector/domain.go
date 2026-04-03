package traceinspector

type AbstractDomain interface {
	incl(lhs AbstractDomain, rhs AbstractDomain) bool          // inclusion operator `lhs ⊑ rhs`
	join(a1 AbstractDomain, a2 AbstractDomain) AbstractDomain  // abstract join operator `a1 ⊔ a2`
	widen(a1 AbstractDomain, a2 AbstractDomain) AbstractDomain // widening operator `a1 ▽ a2`
	repr() string                                              // return string representation of the domain value
}

type Bottom AbstractDomain
