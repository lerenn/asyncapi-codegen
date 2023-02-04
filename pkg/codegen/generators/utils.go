package generators

import "github.com/lerenn/asyncapi-codegen/pkg/asyncapi"

func getCorrelationIDsLocationsByChannel(spec asyncapi.Specification) map[string]string {
	locations := make(map[string]string)

	for k, v := range spec.Channels {
		locations[k] = v.CorrelationIDLocation(spec)
	}

	return locations
}
