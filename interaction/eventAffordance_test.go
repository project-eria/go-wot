package interaction

import (
	"encoding/json"
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type EventAffordanceTestSuite struct {
	suite.Suite
}

func Test_EventAffordanceTestSuite(t *testing.T) {
	suite.Run(t, &EventAffordanceTestSuite{})
}

func (ts *EventAffordanceTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func (ts *EventAffordanceTestSuite) Test_EventAffordanceNew() {
	d, _ := dataSchema.NewInteger(
		dataSchema.IntegerDefault(0),
		dataSchema.IntegerUnit("%"),
		dataSchema.IntegerMin(0),
		dataSchema.IntegerMax(100),
	)
	result := NewEvent("A", "B", "C", &d)
	ts.Equal("A", result.Key)
	ts.Equal("B", result.Title)
	ts.Equal("C", result.Description)
	ts.Equal(d, *result.Data)
}

func (ts *EventAffordanceTestSuite) Test_EventAffordanceJsonMarshal() {
	d, _ := dataSchema.NewInteger(
		dataSchema.IntegerDefault(0),
		dataSchema.IntegerUnit("%"),
		dataSchema.IntegerMin(0),
		dataSchema.IntegerMax(100),
	)
	e := NewEvent("A", "B", "C", &d)
	result, err := json.Marshal(e)
	ts.Nil(err)
	ts.Equal(`{"data":{"default":0,"unit":"%","type":"integer","minimum":0,"maximum":100},"title":"B","description":"C","forms":[]}`, string(result))
}

func (ts *EventAffordanceTestSuite) Test_EventAffordanceJsonUnmarshal() {
	j := []byte(`{"data":{"default":0,"unit":"%","type":"integer","minimum":0,"maximum":100},"title":"B","description":"C","forms":[]}`)
	var result Event
	err := json.Unmarshal(j, &result)
	ts.Nil(err)
	ts.Equal("", result.Key)
	ts.Equal("B", result.Title)
	ts.Equal("C", result.Description)
	ts.Equal(0, result.Data.Default)
	ts.Equal("%", result.Data.Unit)
	ts.Equal("integer", result.Data.Type)
	ts.Equal(0, *result.Data.DataSchema.(dataSchema.Integer).Minimum)
	ts.Equal(100, *result.Data.DataSchema.(dataSchema.Integer).Maximum)
}
