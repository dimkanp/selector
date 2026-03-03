package selector

type Selector map[string]any

func ToSelector(s any) Selector {
	switch v := s.(type) {
	case Selector:
		return v
	default:
		return BaseSelector()
	}
}

func BaseSelector() Selector {
	return Selector{
		"_system_base_fields": struct{}{},
	}
}
