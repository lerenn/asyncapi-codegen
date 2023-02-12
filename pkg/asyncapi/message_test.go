package asyncapi

import (
	"testing"

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
				Headers: &Any{
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
				Payload: &Any{
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
				Headers: &Any{},
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
