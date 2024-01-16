package interaction

import (
	"encoding/json"
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type ActionAffordanceTestSuite struct {
	suite.Suite
}

func Test_ActionAffordanceTestSuite(t *testing.T) {
	suite.Run(t, &ActionAffordanceTestSuite{})
}

func (ts *ActionAffordanceTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func (ts *ActionAffordanceTestSuite) Test_ActionAffordanceNew() {
	i, _ := dataSchema.NewInteger(
		dataSchema.IntegerDefault(0),
		dataSchema.IntegerUnit("%"),
		dataSchema.IntegerMin(0),
		dataSchema.IntegerMax(100),
	)
	o, _ := dataSchema.NewString()
	result := NewAction("A", "B", "C", &i, &o)
	ts.Equal("A", result.Key)
	ts.Equal("B", result.Title)
	ts.Equal("C", result.Description)
	ts.Equal(i, *result.Input)
	ts.Equal(o, *result.Output)
}

func (ts *ActionAffordanceTestSuite) Test_ActionAffordanceJsonMarshal() {
	i, _ := dataSchema.NewInteger(
		dataSchema.IntegerDefault(0),
		dataSchema.IntegerUnit("%"),
		dataSchema.IntegerMin(0),
		dataSchema.IntegerMax(100),
	)
	o, _ := dataSchema.NewString()
	a := NewAction("A", "B", "C", &i, &o)
	result, err := json.Marshal(a)
	ts.Nil(err)
	ts.Equal(`{"input":{"default":0,"unit":"%","type":"integer","minimum":0,"maximum":100},"output":{"type":"string"},"title":"B","description":"C","forms":[]}`, string(result))
}

func (ts *ActionAffordanceTestSuite) Test_ActionAffordanceJsonUnmarshal() {
	j := []byte(`{"input":{"default":0,"unit":"%","type":"integer","minimum":0,"maximum":100},"output":{"type":"string"},"title":"B","description":"C","forms":[]}`)
	var result Action
	err := json.Unmarshal(j, &result)
	ts.Nil(err)
	ts.Equal("", result.Key)
	ts.Equal("B", result.Title)
	ts.Equal("C", result.Description)
	ts.Equal(0, result.Input.Default)
	ts.Equal("%", result.Input.Unit)
	ts.Equal("integer", result.Input.Type)
	ts.Equal(0, *result.Input.DataSchema.(dataSchema.Integer).Minimum)
	ts.Equal(100, *result.Input.DataSchema.(dataSchema.Integer).Maximum)
	ts.Equal("string", result.Output.Type)
}
