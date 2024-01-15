package dataSchema

import (
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type BooleanSchemaTestSuite struct {
	suite.Suite
	schema Data
}

func Test_BooleanSchemaTestSuite(t *testing.T) {
	suite.Run(t, &BooleanSchemaTestSuite{})
}

func (ts *BooleanSchemaTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.schema, _ = NewBoolean(true)
}

func (ts *BooleanSchemaTestSuite) Test_BooleanSchemaNew() {
	ts.Equal(true, ts.schema.Default)
	ts.Equal("boolean", ts.schema.Type)
	ds := ts.schema.DataSchema.(Boolean)
	ts.NotNil(ds)
}

func (ts *BooleanSchemaTestSuite) Test_BooleanSchemaFromString1() {
	result, err := ts.schema.FromString("true")
	ts.Nil(err)
	ts.Equal(true, result)
}

func (ts *BooleanSchemaTestSuite) Test_BooleanSchemaFromString2() {
	result, err := ts.schema.FromString("A")
	ts.Nil(result)
	ts.EqualError(err, "incorrect boolean value type")
}

func (ts *BooleanSchemaTestSuite) Test_BooleanSchemaValidate1() {
	var d interface{} = true
	err := ts.schema.Validate(d)
	ts.Nil(err)
}

func (ts *BooleanSchemaTestSuite) Test_BooleanSchemaValidate2() {
	var d interface{} = "A"
	err := ts.schema.Validate(d)
	ts.EqualError(err, "incorrect boolean value type")
}

func (ts *BooleanSchemaTestSuite) Test_BooleanSchemaJsonMarshal() {
	result, err := json.Marshal(&ts.schema)
	ts.Nil(err)
	ts.Equal(`{"default":true,"type":"boolean"}`, string(result))
}

func (ts *BooleanSchemaTestSuite) Test_BooleanSchemaJsonUnmarshal() {
	j := []byte(`{"default":true,"type":"boolean"}`)
	var result Data
	err := json.Unmarshal(j, &result)
	ts.Nil(err)
	ts.Equal(true, result.Default)
	ts.Equal("boolean", result.Type)
}
