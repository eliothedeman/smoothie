package smoothie

import "testing"

func TestDerivative(t *testing.T) {
	df := NewDataFrame[int](10)

	for i := 0; i < 10; i++ {
		df.Push(i * 10)
	}

	der := df.Derivative()

	der.ForEach(func(v int, i int) {
		if v != 10 {
			t.Fail()
		}
	})
}

func TestNegitiveDerivitive(t *testing.T) {
	df := NewDataFrame[int](10)

	for i := 0; i < 10; i++ {
		df.Push(i * -10)
	}

	der := df.Derivative()

	der.ForEach(func(v int, i int) {
		if v != -10 {
			t.Fail()
		}
	})
}
