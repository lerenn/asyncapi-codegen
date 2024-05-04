//go:generate go run ../../../../cmd/asyncapi-codegen -g types -p issue185 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue185

import "time"

// This is just to test that the generation is correct

// CheckGeneration is just to test that the generation is correct.
func CheckGeneration() {
	var event EventPayloadSchema

	// Check that each field exists and that the type is correct
	event.Time = time.Now()
	event.Data = ContentDataSchema{
		ContentId: "content_id",
	}
	event.Id = "id"
}
