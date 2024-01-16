package interaction

import (
	"encoding/json"
	"testing"

	"github.com/project-eria/go-wot/dataSchema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type PropertyAffordanceTestSuite struct {
	suite.Suite
}

func Test_PropertyAffordanceTestSuite(t *testing.T) {
	suite.Run(t, &PropertyAffordanceTestSuite{})
}

func (ts *PropertyAffordanceTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func (ts *PropertyAffordanceTestSuite) Test_PropertyAffordanceNew() {
	d, _ := dataSchema.NewBoolean(
		dataSchema.BooleanDefault(false),
	)
	result := NewProperty("A", "B", "C", false, false, true, nil, d)
	ts.Equal("A", result.Key)
	ts.Equal("B", result.Title)
	ts.Equal("C", result.Description)
	ts.Equal(false, result.ReadOnly)
	ts.Equal(false, result.WriteOnly)
	ts.Equal(true, result.Observable)
	ts.Nil(result.UriVariables)
	ts.Equal(d, result.Data)
}

func (ts *PropertyAffordanceTestSuite) Test_PropertyAffordanceJsonMarshal() {
	d, _ := dataSchema.NewInteger(
		dataSchema.IntegerDefault(5),
		dataSchema.IntegerUnit("%"),
		dataSchema.IntegerMin(1),
		dataSchema.IntegerMax(9),
	)
	p := NewProperty("A", "B", "C", false, false, true, nil, d)
	result, err := json.Marshal(p)
	ts.Nil(err)
	ts.Equal(`{"title":"B","description":"C","forms":[],"default":5,"unit":"%","type":"integer","minimum":1,"maximum":9,"readOnly":false,"writeOnly":false,"observable":true}`, string(result))
}

func (ts *PropertyAffordanceTestSuite) Test_PropertyAffordanceJsonUnmarshal() {
	j := []byte(`{"title":"B","description":"C","forms":[],"default":5,"unit":"%","type":"integer","minimum":1,"maximum":9,"readOnly":false,"writeOnly":false,"observable":true}`)
	var result Property
	err := json.Unmarshal(j, &result)
	ts.Nil(err)
	ts.Equal("", result.Key)
	ts.Equal("B", result.Title)
	ts.Equal("C", result.Description)
	ts.Equal(5, result.Data.Default)
	ts.Equal("%", result.Data.Unit)
	ts.Equal("integer", result.Data.Type)
	ts.Equal(1, *result.Data.DataSchema.(dataSchema.Integer).Minimum)
	ts.Equal(9, *result.Data.DataSchema.(dataSchema.Integer).Maximum)
	ts.Equal(false, result.ReadOnly)
	ts.Equal(false, result.WriteOnly)
	ts.Equal(true, result.Observable)
}
