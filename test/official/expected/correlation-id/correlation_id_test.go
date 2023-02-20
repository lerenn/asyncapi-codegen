package correlationID

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestCorrelationIDSuite(t *testing.T) {
	suite.Run(t, new(CorrelationIDSuite))
}

type CorrelationIDSuite struct {
	suite.Suite
}

func (suite *CorrelationIDSuite) TestDateTimeCorrectUnmarshal() {
	type structWithTime struct {
		SentAtSchema SentAtSchema
	}

	originalTime := time.Unix(60, 0)

	// Marshal
	strt := structWithTime{
		SentAtSchema: SentAtSchema(originalTime),
	}
	b, err := json.Marshal(strt)
	suite.Require().NoError(err)

	// Unmarshal
	var newStrt structWithTime
	suite.Require().NoError(json.Unmarshal(b, &newStrt))

	// Compare
	suite.Require().Equal(originalTime, time.Time(newStrt.SentAtSchema))
}
