package algebra

import (
	"testing"
)

func Test_ExtInt(t *testing.T) {
	t.Log("0 * 100 =", ExtInt_Zero().Add(ExtInt_Finite(100)))
	t.Log("0 * ∞ =", ExtInt_Zero().Add(ExtInt_Infty()))
	t.Log("-5 + ∞ =", ExtInt_Finite(-5).Add(ExtInt_Infty()))
	t.Log("-∞ * ∞ =", ExtInt_NegInfty().Mul(ExtInt_Infty()))
	t.Log("-∞ * -∞ =", ExtInt_NegInfty().Mul(ExtInt_NegInfty()))
	t.Log("-∞ * -99 =", ExtInt_NegInfty().Mul(ExtInt_Finite(-99)))
	t.Log("5 * 8 =", ExtInt_Finite(5).Mul(ExtInt_Finite(8)))

	t.Log("5 <= ∞ ==", ExtInt_Finite(5).Leq(ExtInt_Infty()))
	t.Log("5 <= -∞ ==", ExtInt_Finite(5).Leq(ExtInt_NegInfty()))
	t.Log("5 <= 4 ==", ExtInt_Finite(5).Leq(ExtInt_Finite(4)))

	t.Log("min(5, -3, ∞) =", ExtInt_Finite(5).Min(ExtInt_Finite(-3), ExtInt_Infty()))

}
