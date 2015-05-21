package smoothie

import (
	"log"
	"testing"
)

func TestDataFrameInsert(t *testing.T) {
	df := NewDataFrame(make([]float64, 100))

	for i := 0; i < 1000; i++ {
		df.Push(float64(i))
		log.Println(df.data)
		log.Println(df.Data())
	}
}
