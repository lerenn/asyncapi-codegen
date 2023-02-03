package generators

import "bytes"

type BrokerControllerGenerator struct {
}

func (bcg BrokerControllerGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		brokerControllerTemplatePath,
	)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmplt.Execute(buf, bcg); err != nil {
		return "", err
	}

	return buf.String(), nil
}
