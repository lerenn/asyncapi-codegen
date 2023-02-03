package generators

import "github.com/lerenn/asyncapi-codegen/pkg/asyncapi"

func getCorrelationIdsLocationsByChannel(spec asyncapi.Specification) map[string]string {
	locations := make(map[string]string)

	for k, v := range spec.Channels {
		locations[k] = v.CorrelationIdLocation(spec)
	}

	return locations
}
