package rabbit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePublicProperty(t *testing.T) {
	params := []Param{
		{Type: "public", Name: "paramName", Value: "paramValue"},
	}

	c := Config{Job: "testJob"}
	generatedParams := c.generateProperty(params...)

	// Expected serialized parameters
	expectedParams := "s:9:\"paramName\";s:10:\"paramValue\";"

	// Assertions
	assert.Equal(t, expectedParams, generatedParams)
}

func TestGenerateInvalidProperty(t *testing.T) {
	params := []Param{
		{Type: "invalid", Name: "paramName", Value: "paramValue"},
	}

	c := Config{Job: "testJob"}
	generatedParams := c.generateProperty(params...)

	// Expected serialized parameters
	expectedParams := "s:9:\"paramName\";s:10:\"paramValue\";"

	// Assertions
	assert.Equal(t, expectedParams, generatedParams)
}

func TestGenerateProtectedProperty(t *testing.T) {
	params := []Param{
		{Type: "protected", Name: "paramName", Value: "paramValue"},
	}

	c := Config{Job: "testJob"}
	generatedParams := c.generateProperty(params...)

	// Expected serialized parameters
	expectedParams := fmt.Sprintf("s:%d:\"\u0000*\u0000%s\";s:%d:\"%s\";", len("paramName")+3, "paramName", len("paramValue"), "paramValue")

	// Assertions
	assert.Equal(t, expectedParams, generatedParams)
}

func TestGeneratePrivateProperty(t *testing.T) {
	params := []Param{
		{Type: "private", Name: "paramName", Value: "paramValue"},
	}

	c := Config{Job: "testJob"}
	generatedParams := c.generateProperty(params...)

	// Expected serialized parameters
	expectedParams := fmt.Sprintf("s:%d:\"\u0000%s\u0000%s\";s:%d:\"%s\";", len(c.Job)+2, c.Job, "paramName", len("paramValue"), "paramValue")

	// Assertions
	assert.Equal(t, expectedParams, generatedParams)
}

func TestDelay(t *testing.T) {
	c := Config{Delay: 5}

	// Expected response
	expected := fmt.Sprintf("i:%d;", c.Delay)

	// Assertions
	assert.Equal(t, expected, c.generateDelay())
}

func TestNoDelay(t *testing.T) {
	c := Config{Delay: 0}

	// Expected response
	expected := "N;"

	// Assertions
	assert.Equal(t, expected, c.generateDelay())
}
