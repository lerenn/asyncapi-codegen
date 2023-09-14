package asyncapi

import (
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	"github.com/stretchr/testify/suite"
)

func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessageSuite))
}

type MessageSuite struct {
	suite.Suite
}

func (suite *MessageSuite) TestIsCorrelationIDRequired() {
	cases := []struct {
		Message  Message
		Required bool
	}{
		{
			Message: Message{
				Headers: &Schema{
					Required: []string{"correlationId"},
				},
				CorrelationID: &CorrelationID{
					Location: "$message.header#/correlationId",
				},
			},
			Required: true,
		},
		{
			Message: Message{
				Payload: &Schema{
					Required: []string{"correlationId"},
				},
				CorrelationID: &CorrelationID{
					Location: "$message.payload#/correlationId",
				},
			},
			Required: true,
		},
		{
			Message: Message{
				Headers: &Schema{
					Properties: map[string]*Schema{
						"toto": utils.ToPointer(Schema{
							Required: []string{"correlationId"},
						}),
					},
				},
				CorrelationID: &CorrelationID{
					Location: "$message.header#/toto/correlationId",
				},
			},
			Required: true,
		},
		{
			Message: Message{
				Headers: &Schema{},
				CorrelationID: &CorrelationID{
					Location: "$message.header#/correlationId",
				},
			},
			Required: false,
		},
		{
			Message:  Message{},
			Required: false,
		},
	}

	for i, c := range cases {
		suite.Require().Equal(c.Required, c.Message.isCorrelationIDRequired(), i)
	}
}
