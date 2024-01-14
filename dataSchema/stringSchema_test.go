package dataSchema

import (
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type StringSchemaTestSuite struct {
	suite.Suite
	schema Data
}

func Test_StringSchemaTestSuite(t *testing.T) {
	suite.Run(t, &StringSchemaTestSuite{})
}

func (ts *StringSchemaTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	min := uint16(2)
	max := uint16(9)
	ts.schema, _ = NewString("test", &min, &max, `[A-Za-z]*`)
}

func (ts *StringSchemaTestSuite) Test_StringSchemaNew() {

	ts.Equal("test", ts.schema.Default)
	ts.Equal("string", ts.schema.Type)
	ds := ts.schema.DataSchema.(String)
	ts.Equal(uint16(2), *ds.MinLength)
	ts.Equal(uint16(9), *ds.MaxLength)
	ts.Equal(`^[A-Za-z]*$`, ds.Pattern)
}
func (ts *StringSchemaTestSuite) Test_StringSchemaNewError1() {
	min := uint16(2)
	max := uint16(9)
	result, err := NewString("test", &min, &max, `\o`)
	ts.Equal(Data{}, result)
	ts.EqualError(err, "invalid pattern")
}

func (ts *StringSchemaTestSuite) Test_StringSchemaFromString1() {
	result, err := ts.schema.FromString("t3st")
	ts.Nil(err)
	ts.Equal("t3st", result)
}

func (ts *StringSchemaTestSuite) Test_StringSchemaValidate1() {
	var d interface{} = "test"
	err := ts.schema.Validate(d)
	ts.Nil(err)
}

func (ts *StringSchemaTestSuite) Test_StringSchemaValidate2() {
	var d interface{} = 1
	err := ts.schema.Validate(d)
	ts.EqualError(err, "incorrect string value type")
}

func (ts *StringSchemaTestSuite) Test_StringSchemaValidate3() {
	err := ts.schema.Validate("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	ts.EqualError(err, "string too long")
}

func (ts *StringSchemaTestSuite) Test_StringSchemaValidate4() {
	err := ts.schema.Validate("A")
	ts.EqualError(err, "string too short")
}

func (ts *StringSchemaTestSuite) Test_StringSchemaValidate5() {
	err := ts.schema.Validate("A3C")
	ts.EqualError(err, "string does not match pattern")
}

func (ts *StringSchemaTestSuite) Test_StringSchemaJson() {
	result, err := json.Marshal(&ts.schema)
	ts.Nil(err)
	ts.Equal(`{"default":"test","type":"string","minLength":2,"maxLength":9,"pattern":"^[A-Za-z]*$"}`, string(result))
}
