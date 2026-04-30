package dataSchema

import (
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
	stringType, _ := NewString(
		StringMinLength(16),
	)
	ts.schema, _ = NewObject(
		ObjectProperty("mystring", &stringType),
		ObjectRequired([]string{"mystring"}),
	)
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaNew() {
	//ts.Equal("test", ts.schema.Default)
	ts.Equal("object", ts.schema.Type)
	ds := ts.schema.DataSchema.(Object)
	ts.Equal([]string{"mystring"}, ds.Required)
	ts.Len(ds.Properties, 1, "Properties map should have 1 element")
	ts.Contains(ds.Properties, "mystring", "Properties should contain 'mystring'")
	ts.Equal("string", ds.Properties["mystring"].Type, "mystring should have Type 'string'")
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaIsNullOrEmpty1() {
	var d interface{} = map[string]string{
		"a": "b",
	}
	res := ts.schema.IsNullOrEmpty(d)
	ts.False(res)
}
func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaIsNullOrEmpty2() {
	var d interface{} = map[string]string{
		"a": "",
	}
	res := ts.schema.IsNullOrEmpty(d)
	ts.False(res)
}
func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaIsNullOrEmpty3() {
	var d interface{} = ""
	res := ts.schema.IsNullOrEmpty(d)
	ts.True(res)
}

func (ts *ObjectSchemaTestSuite) Test_ObjectSchemaIsNullOrEmpty4() {
	var d interface{} = map[string]string{}
	res := ts.schema.IsNullOrEmpty(d)
	ts.True(res)
}
