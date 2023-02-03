package generators

import "bytes"

type BrokerNATSGenerator struct {
}

func (bng BrokerNATSGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		brokerNATSTemplatePath,
	)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmplt.Execute(buf, bng); err != nil {
		return "", err
	}

	return buf.String(), nil
}
