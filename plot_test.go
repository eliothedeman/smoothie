package smoothie

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgsvg"
)

func testPlot(df *DataFrame, name string, mod func(*DataFrame) *DataFrame) {

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	newdf := mod(df)

	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	err = plotutil.AddLinePoints(p, "raw", df.PlotPoints(), "smooth", newdf.PlotPoints())
	if err != nil {
		log.Fatal(err)
	}

	c := vgsvg.New(16*vg.Inch, 9*vg.Inch)

	can := draw.New(c)

	p.Draw(can)
	p.Save(16*vg.Inch, 9*vg.Inch, fmt.Sprintf("%s.png", name))
	f, err := os.Create(fmt.Sprintf("%s.svg", name))
	if err != nil {
		log.Fatal(err)
	}

	c.WriteTo(f)

}

func randDF(size int) *DataFrame {
	df := EmptyDataFrame(size)

	for i := 0; i < df.Len(); i++ {
		if i == 0 {
			df.Insert(i, math.Sin(float64(i)/float64(df.Len())))
		} else {
			df.Insert(i, rand.Float64()+math.Sin(16*math.Pi*float64(i)/float64(df.Len()))+math.Sin(8*math.Pi*float64(i)/float64(df.Len())))
		}
	}
	return df
}

type mod func(df *DataFrame) *DataFrame

var (
	test_mods = map[string]mod{
		"moving_average": func(df *DataFrame) *DataFrame {
			return df.MovingAverage(10)
		},
		"weighted_average": func(df *DataFrame) *DataFrame {
			return df.WeightedMovingAverage(10, LinearWeighting)
		},
		"double_smooth": func(df *DataFrame) *DataFrame {
			return df.DoubleExponentialSmooth(0.2, 0.3)
		},
		"single_smooth": func(df *DataFrame) *DataFrame {
			return df.SingleExponentialSmooth(0.3)
		},
	}
)

func TestPlotDF(t *testing.T) {
	rand := randDF(200)

	for k, v := range test_mods {
		testPlot(rand, k, v)
	}
}
