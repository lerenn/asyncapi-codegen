package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils/template"

type processable interface {
	Process(string, Specification) error
}

func processMap[T processable](spec Specification, m map[string]T, suffix string) error {
	for name, entity := range m {
		if err := entity.Process(template.Namify(name)+suffix, spec); err != nil {
			return err
		}
	}

	return nil
}
