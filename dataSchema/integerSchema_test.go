package dataSchema

import (
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type IntegerSchemaTestSuite struct {
	suite.Suite
	schema Data
}

func Test_IntegerSchemaTestSuite(t *testing.T) {
	suite.Run(t, &IntegerSchemaTestSuite{})
}

func (ts *IntegerSchemaTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	min := 1
	max := 9
	ts.schema, _ = NewInteger(5, "%", &min, &max)
}

func (ts *IntegerSchemaTestSuite) Test_IntegerSchemaNew() {
	ts.Equal(5, ts.schema.Default)
	ts.Equal("integer", ts.schema.Type)
	ts.Equal("%", ts.schema.Unit)
	ds := ts.schema.DataSchema.(Integer)
	ts.Equal(1, *ds.Minimum)
	ts.Equal(9, *ds.Maximum)
}

func (ts *IntegerSchemaTestSuite) Test_IntegerSchemaFromString1() {
	result, err := ts.schema.FromString("6")
	ts.Nil(err)
	ts.Equal(6, result)
}

func (ts *IntegerSchemaTestSuite) Test_IntegerSchemaFromString2() {
	result, err := ts.schema.FromString("A")
	ts.Equal(0, result)
	ts.EqualError(err, "strconv.Atoi: parsing \"A\": invalid syntax")
}

func (ts *IntegerSchemaTestSuite) Test_IntegerSchemaValidate1() {
	var d interface{} = 6
	err := ts.schema.Validate(d)
	ts.Nil(err)
}

func (ts *IntegerSchemaTestSuite) Test_IntegerSchemaValidate2() {
	var d interface{} = "A"
	err := ts.schema.Validate(d)
	ts.EqualError(err, "incorrect integer value type")
}

func (ts *IntegerSchemaTestSuite) Test_IntegerSchemaValidate3() {
	err := ts.schema.Validate(10)
	ts.EqualError(err, "value is greater than maximum")
}

func (ts *IntegerSchemaTestSuite) Test_IntegerSchemaValidate4() {
	err := ts.schema.Validate(0)
	ts.EqualError(err, "value is less than minimum")
}

func (ts *IntegerSchemaTestSuite) Test_IntegerSchemaJson() {
	min := 0
	max := 9
	i, _ := NewInteger(5, "%", &min, &max)
	result, err := json.Marshal(&i)
	ts.Nil(err)
	ts.Equal(`{"default":5,"unit":"%","type":"integer","minimum":0,"maximum":9}`, string(result))
}
