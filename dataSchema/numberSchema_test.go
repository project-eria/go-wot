package dataSchema

import (
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type NumberSchemaTestSuite struct {
	suite.Suite
	schema Data
}

func Test_NumberSchemaTestSuite(t *testing.T) {
	suite.Run(t, &NumberSchemaTestSuite{})
}

func (ts *NumberSchemaTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.schema, _ = NewNumber(
		NumberDefault(5.5),
		NumberUnit("%"),
		NumberMin(1),
		NumberMax(9),
	)
}

func (ts *NumberSchemaTestSuite) Test_NumberSchemaNew() {
	ts.Equal(5.5, ts.schema.Default)
	ts.Equal("number", ts.schema.Type)
	ts.Equal("%", ts.schema.Unit)
	ds := ts.schema.DataSchema.(Number)
	ts.Equal(1, *ds.Minimum)
	ts.Equal(9, *ds.Maximum)
}

func (ts *NumberSchemaTestSuite) Test_NumberSchemaFromString1() {
	result, err := ts.schema.FromString("6.6")
	ts.Nil(err)
	ts.Equal(6.6, result)
}

func (ts *NumberSchemaTestSuite) Test_NumberSchemaFromString2() {
	result, err := ts.schema.FromString("A")
	ts.Equal(0.0, result)
	ts.EqualError(err, "strconv.ParseFloat: parsing \"A\": invalid syntax")
}

func (ts *NumberSchemaTestSuite) Test_NumberSchemaValidate1() {
	var d interface{} = 6.6
	err := ts.schema.Validate(d)
	ts.Nil(err)
}

func (ts *NumberSchemaTestSuite) Test_NumberSchemaValidate2() {
	var d interface{} = "A"
	err := ts.schema.Validate(d)
	ts.EqualError(err, "incorrect number value type")
}

func (ts *NumberSchemaTestSuite) Test_NumberSchemaValidate3() {
	err := ts.schema.Validate(10.0)
	ts.EqualError(err, "value is greater than maximum")
}

func (ts *NumberSchemaTestSuite) Test_NumberSchemaValidate4() {
	err := ts.schema.Validate(0.0)
	ts.EqualError(err, "value is less than minimum")
}

func (ts *NumberSchemaTestSuite) Test_NumberSchemaJsonMarshal() {
	result, err := json.Marshal(&ts.schema)
	ts.Nil(err)
	ts.Equal(`{"default":5.5,"unit":"%","type":"number","minimum":1,"maximum":9}`, string(result))
}

func (ts *NumberSchemaTestSuite) Test_NumberSchemaJsonUnmarshal() {
	j := []byte(`{"default":5.5,"unit":"%","type":"number","minimum":1,"maximum":9}`)
	var result Data
	err := json.Unmarshal(j, &result)
	ts.Nil(err)
	ts.Equal(5.5, result.Default)
	ts.Equal("%", result.Unit)
	ts.Equal("number", result.Type)
	ts.Equal(1, *result.DataSchema.(Number).Minimum)
	ts.Equal(9, *result.DataSchema.(Number).Maximum)
}
