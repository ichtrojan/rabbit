package rabbit

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePublicProperty(t *testing.T) {
	params := []Param{
		{Type: "public", Name: "paramName", Value: "paramValue"},
	}

	c := Config{Job: "testJob"}
	generatedParams := c.generateProperty(params...)

	expectedParams := "s:9:\"paramName\";s:10:\"paramValue\";"

	assert.Equal(t, expectedParams, generatedParams)
}

func TestGenerateInvalidProperty(t *testing.T) {
	params := []Param{
		{Type: "invalid", Name: "paramName", Value: "paramValue"},
	}

	c := Config{Job: "testJob"}
	generatedParams := c.generateProperty(params...)

	expectedParams := "s:9:\"paramName\";s:10:\"paramValue\";"

	assert.Equal(t, expectedParams, generatedParams)
}

func TestGenerateProtectedProperty(t *testing.T) {
	params := []Param{
		{Type: "protected", Name: "paramName", Value: "paramValue"},
	}

	c := Config{Job: "testJob"}
	generatedParams := c.generateProperty(params...)

	expectedParams := fmt.Sprintf(
		"s:%d:\"\u0000*\u0000%s\";s:%d:\"%s\";",
		len("paramName")+3, "paramName", len("paramValue"),
		"paramValue",
	)

	assert.Equal(t, expectedParams, generatedParams)
}

func TestGeneratePrivateProperty(t *testing.T) {
	params := []Param{
		{Type: "private", Name: "paramName", Value: "paramValue"},
	}

	c := Config{Job: "testJob"}
	generatedParams := c.generateProperty(params...)

	expectedParams := fmt.Sprintf(
		"s:%d:\"\u0000%s\u0000%s\";s:%d:\"%s\";",
		len(c.Job)+2, c.Job, "paramName", len("paramValue"),
		"paramValue",
	)

	assert.Equal(t, expectedParams, generatedParams)
}

func TestDelay(t *testing.T) {
	c := Config{Delay: 5}

	expected := fmt.Sprintf("i:%d;", c.Delay)

	assert.Equal(t, expected, c.generateDelay())
}

func TestNoDelay(t *testing.T) {
	c := Config{Delay: 0}

	expected := "N;"

	assert.Equal(t, expected, c.generateDelay())
}

func TestEmptyAppName(t *testing.T) {
	params := []Param{
		{Type: "protected", Name: "paramName", Value: "paramValue"},
	}

	c := Config{}

	generatedParams := c.Dispatch(params...)

	expectedParams := errors.New("AppName is required in config")

	assert.Equal(t, expectedParams, generatedParams)
}
