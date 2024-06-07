//go:generate go run ../../../../cmd/asyncapi-codegen -g types -p issue192 -i asyncapi.yaml,openapi.yaml -o ./asyncapi.gen.go

package issue192

import "github.com/TheSadlig/asyncapi-codegen/pkg/utils"

// This is just to test that the generation is correct

// CheckGeneration is a function to check the generation of the code.
func CheckGeneration() {
	// Check local schema
	var local LocalSchema
	local.Data.Hello = utils.ToPointer("hello")
	local.Data.World = utils.ToPointer("world")

	// Check distant schema
	var distant DistantSchema
	distant.Data.Hello = utils.ToPointer("hello")
	distant.Data.World = utils.ToPointer("world")
}
