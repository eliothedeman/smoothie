package smoothie

import "testing"

func TestDataFrameInsert(t *testing.T) {
	df := NewDataFrame(make([]float64, 100))

	for i := 0; i < 101; i++ {
		df.Push(float64(i))
	}

	if df.Index(0) != 1 {
		t.Fail()
	}
}

func TestGrow(t *testing.T) {
	df := NewDataFrame(make([]float64, 10))

	for i := 0; i < 15; i++ {
		df.Push(float64(i))
	}
	df.Grow(5)
	df.Push(15)

	if df.Len() != 15 {
		t.Fail()
	}

	if df.Index(14) != 15 {
		t.Fail()
	}
}

func TestShrink(t *testing.T) {
	df := NewDataFrame(make([]float64, 10))

	df.Shrink(2)

	if df.Len() != 8 {
		t.Fail()
	}
}
