package smoothie

// KernalConvolve applies a centered kernel to the dataframe. The misspelling
// is retained for compatibility with the original API.
func (d *FloatDataFrame) KernalConvolve(kernel *FloatDataFrame) *FloatDataFrame {
	if kernel.Len() == 0 {
		panic("FloatDataFrame: convolution kernel cannot be empty")
	}

	output := NewFloatDataFrame(d.Len())
	middle := kernel.Len() / 2
	divisor := kernel.Sum()
	for i := 0; i < d.Len(); i++ {
		var total float64
		for x := 0; x < kernel.Len(); x++ {
			index := i + x - middle
			if index < 0 || index >= d.Len() {
				continue
			}
			total += kernel.Index(x) * d.Index(index)
		}
		output.Insert(i, total/divisor)
	}
	return output
}

// KernelConvolve is the correctly-spelled form of KernalConvolve.
func (d *FloatDataFrame) KernelConvolve(kernel *FloatDataFrame) *FloatDataFrame {
	return d.KernalConvolve(kernel)
}
