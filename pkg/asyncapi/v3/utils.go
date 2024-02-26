package asyncapiv3

type processable interface {
	Process(string, Specification)
}

func processMap[T processable](spec Specification, m map[string]T) {
	for name, entity := range m {
		entity.Process(name, spec)
	}
}
