package smoothie

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestCat(t *testing.T) {
	one := NewDataFrame[int](100)
	for i := 0; i < 100; i++ {
		one.Push(i)
	}
	two := NewDataFrame[int](100)
	for i := 0; i < 100; i++ {
		two.Push(i)
	}

	three := one.Cat(two)
	if three.Len() != 200 {
		t.Fatal()
	}

	for i := 0; i < 200; i++ {
		if three.Index(i) != i%100 {
			t.Fatal(three.Index(i), i%100)
		}
	}
}

func TestDataFrameInsert(t *testing.T) {
	df := NewDataFrame[int](100)
	for i := 0; i < 101; i++ {
		df.Push(i)
	}

	if df.Index(0) != 1 {
		t.Fail()
	}
}

func TestGrow(t *testing.T) {
	df := NewDataFrame[int](10)
	for i := 0; i < 15; i++ {
		df.Push(i)
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
	df := NewDataFrame[int](10)
	df.Shrink(2)

	if df.Len() != 8 {
		t.Fail()
	}
}

func TestIntegerOperations(t *testing.T) {
	df := NewDataFrameFromSlice([]int{3, 1, 4, 1, 5})

	if df.Min() != 1 || df.Max() != 5 || df.Sum() != 14 || df.Median() != 3 {
		t.Fatal("unexpected integer statistics")
	}

	derivative := df.Derivative().Data()
	if len(derivative) != 4 || derivative[0] != -2 || derivative[3] != 4 {
		t.Fatal("unexpected integer derivative", derivative)
	}
}
