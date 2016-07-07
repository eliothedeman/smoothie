package smoothie

func (d *DataFrame) KernalConvolve(kernal *DataFrame) *DataFrame {
	output := NewDataFrame(d.Len())
	tmp := NewDataFrame(kernal.Len())
	middle := kernal.Len() / 2
	for i := 0; i < d.Len(); i++ {
		for x := 0; x < kernal.Len(); x++ {
			if i+x >= d.Len() || i+x-middle < 0 {
				continue
			}
			tmp.Insert(x, kernal.Index(x)*d.Index(i+x-middle))
		}

		output.Insert(i, tmp.Avg()/kernal.Avg())
	}

	return output.Reduce(output.Len() - (middle * 2))
}
