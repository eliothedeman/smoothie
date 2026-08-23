package smoothie

import (
	"math"
	"testing"
)

func TestFloatDataFrameStatistics(t *testing.T) {
	df := NewFloatDataFrameFromSlice([]float64{1, 2, 3, 4})

	if df.Avg() != 2.5 || df.StdDev() != 1 || df.Median() != 3 {
		t.Fatal("unexpected floating-point statistics")
	}

	if got := df.MovingAverage(2).Data(); !reflectFloatSlice(got, []float64{1.5, 2.5, 3.5, 4}) {
		t.Fatal("unexpected moving average", got)
	}
}

func TestFloatDataFrameFunctions(t *testing.T) {
	df := NewFloatDataFrameFromSlice([]float64{1, 2, 3, 4})

	if df.AddConst(1).Index(0) != 2 {
		t.Fatal("wrapper did not preserve dataframe operations")
	}
	if df.DoubleSmoothPredictN(2, 0.5, 0.5).Len() != 6 {
		t.Fatal("unexpected prediction length")
	}
	if len(df.FFT()) != df.Len() {
		t.Fatal("unexpected transform length")
	}

	triple := NewFloatDataFrameFromSlice([]float64{1, 2, 1, 2, 1, 2, 1, 2}).TripleSmooth(0.5, 0.3, 0.5, 2)
	if triple.Len() != 8 {
		t.Fatal("unexpected triple smoothing length")
	}

	if !math.IsNaN(NewFloatDataFrameFromSlice([]float64{math.NaN()}).Avg()) {
		t.Fatal("expected NaN for an empty floating-point sample")
	}
}

func reflectFloatSlice(left, right []float64) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
