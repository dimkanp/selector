package selector

type PreparableSlice[T Preparable] struct {
	Preparable
	slice *[]T
}

func Slice[T Preparable](p *[]T) *PreparableSlice[T] {
	var t T
	l := &PreparableSlice[T]{
		Preparable: t,
		slice:      p,
	}

	return l
}
