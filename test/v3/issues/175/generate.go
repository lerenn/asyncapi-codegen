//go:generate go run ../../../../cmd/asyncapi-codegen -p issue175 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue175

import (
	"github.com/TheSadlig/asyncapi-codegen/pkg/utils"
)

// CheckGeneration is just to test that the generation is correct.
func CheckGeneration() {
	var msg Type1Message
	msg.Payload = []ItemFromType1MessagePayload{
		{
			Age:   utils.ToPointer(int64(1)),
			Email: utils.ToPointer("email1"),
			Name:  utils.ToPointer("name1"),
		},
	}

	var msg2 Type2Message
	msg2.Payload = ArrayPayloadSchema{
		{
			Age:   utils.ToPointer(int64(1)),
			Email: utils.ToPointer("email1"),
			Name:  utils.ToPointer("name1"),
		},
	}

	var msg3 Type3Message
	msg3.Payload = []string{
		"content1",
		"content2",
		"content3",
	}
}
