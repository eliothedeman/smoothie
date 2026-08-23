# smoothie

Generic numeric data-frame operations for integer and floating-point values.

```go
integers := smoothie.NewDataFrameFromSlice([]int{1, 2, 3})
floats := smoothie.NewFloatDataFrameFromSlice([]float64{1, 2, 3})

sum := integers.Sum()
average := floats.Avg()
```

`FloatDataFrame` provides floating-point-only operations such as averages,
moving averages, exponential smoothing, Fourier transforms, and convolution.
