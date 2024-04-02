//go:generate go run ../../../../cmd/asyncapi-codegen -p issue97 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue97

// CheckGeneration is just to test that the generation is correct.
func CheckGeneration() {
	var msg ReferencePayloadArrayMessage
	msg.Payload = ArraySchema{
		"content1",
		"content2",
		"content3",
	}
}
