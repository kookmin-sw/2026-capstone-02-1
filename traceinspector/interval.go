package traceinspector


type IntervalDomainConsts int

const (
	interval_bot IntervalDomainConsts = iota  // no values
	interval_top // any value
)

func (IntervalDomainConsts) isIntervalDomainElem()

type IntervalDomain struct {
	lower, upper int64
}

func (l1 IntervalDomain) join(l2 IntervalDomainElem) IntervalDomainElem {
	switch l2_ty := l2.(type) {
		case 
	}
}

func (l1 IntervalDomainConsts) join(l2 IntervalDomainElem) IntervalDomainElem {
	switch l1 {
	case interval_bot:
		return interval_bot
	case interval_top:
		return interval_top
	}
	panic("unreachable")
}
