package domain

import "fmt"

type ArrayTopDomain struct {
	length IntervalDomain
	is_top bool
}

func (domain ArrayTopDomain) String() string {
	return fmt.Sprintf("[⊤, len : %s]", domain.length.String())
}

func (domain ArrayTopDomain) IsBot() bool {
	return false
}

func (domain ArrayTopDomain) IsTop() bool {
	return true
}

func (lhs ArrayTopDomain) Join(rhs ArrayTopDomain) ArrayTopDomain {
	length_joined := lhs.length.Join(rhs.length)
	return ArrayTopDomain{length: length_joined, is_top: true}
}

func (lhs ArrayTopDomain) Incl(rhs ArrayTopDomain) bool {
	return lhs.length.Incl(rhs.length)
}

func (lhs ArrayTopDomain) Widen(rhs ArrayTopDomain) ArrayTopDomain {
	return ArrayTopDomain{length: lhs.length.Widen(rhs.length), is_top: true}
}
