package dataSchema

import (
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type ObjectSchemaTestSuite struct {
	suite.Suite
	schema Data
}

func Test_ObjectSchemaTestSuite(t *testing.T) {
	suite.Run(t, &ObjectSchemaTestSuite{})
}

func (ts *ObjectSchemaTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	// Object has no nested `properties`/`required` yet (see the TODO in
	// objectSchema.go), so only the generic Data options apply here
	ts.schema, _ = NewObject(
		ObjectDefault(map[string]interface{}{
			"mystring": "A",
		}),
	)
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaNew() {
	ts.Equal(map[string]interface{}{"mystring": "A"}, ts.schema.Default)
	ts.Equal("object", ts.schema.Type)
	ds := ts.schema.DataSchema.(Object)
	ts.NotNil(ds)
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaNewUnit() {
	schema, err := NewObject(
		ObjectUnit("celsius"),
	)
	ts.Nil(err)
	ts.Equal("celsius", schema.Unit)
	ts.Nil(schema.Default)
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaFromString() {
	result, err := ts.schema.FromString(`{"mystring":"A"}`)
	ts.Nil(result)
	ts.EqualError(err, "not implemented")
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaValidate1() {
	var d interface{} = map[string]interface{}{
		"a": "b",
	}
	err := ts.schema.Validate(d)
	ts.Nil(err)
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaValidate2() {
	// Only map[string]interface{} is accepted, a typed map is not
	var d interface{} = map[string]string{
		"a": "b",
	}
	err := ts.schema.Validate(d)
	ts.EqualError(err, "incorrect object value type")
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaValidate3() {
	var d interface{} = ""
	err := ts.schema.Validate(d)
	ts.EqualError(err, "incorrect object value type")
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaValidate4() {
	var d interface{}
	err := ts.schema.Validate(d)
	ts.EqualError(err, "missing value")
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaJsonMarshal() {
	result, err := json.Marshal(&ts.schema)
	ts.Nil(err)
	ts.Equal(`{"default":{"mystring":"A"},"type":"object"}`, string(result))
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaJsonUnmarshal() {
	j := []byte(`{"default":{"mystring":"A"},"type":"object"}`)
	var result Data
	err := json.Unmarshal(j, &result)
	ts.Nil(err)
	ts.Equal(map[string]interface{}{"mystring": "A"}, result.Default)
	ts.Equal("object", result.Type)
	ts.Equal(Object{}, result.DataSchema)
}
