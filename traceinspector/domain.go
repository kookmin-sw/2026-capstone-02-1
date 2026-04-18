package traceinspector

type AbstractDomain interface {
	IsBot() bool
	IsTop() bool
	Incl(rhs AbstractDomain) bool            // inclusion operator `lhs ⊑ rhs`
	Join(rhs AbstractDomain) AbstractDomain  // abstract join operator `lhs ⊔ rhs`
	Widen(rhs AbstractDomain) AbstractDomain // widening operator `lhs ▽ rhs`
	String() string                          // return string representation of the domain value
}
