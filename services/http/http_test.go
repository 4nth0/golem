package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_normalizeRouteConfiguration(t *testing.T) {

	route := normalizeRouteConfiguration(HTTPHandler{
		BodyFile: templateGoldenPath,
	})

	template, err := LoadTemplate(templateGoldenPath)
	assert.Nil(t, err)

	assert.Equal(t, DefaultStatusCode, route.Code)
	assert.Equal(t, DefaultMethod, route.Method)
	assert.Equal(t, template, route.Body)
}

func Test_normalizeRouteConfigurationMultipleBodyFiles(t *testing.T) {

	route := normalizeRouteConfiguration(HTTPHandler{
		BodyFiles: []string{templateGoldenPath, templateGoldenPath2},
	})

	template_1, err := LoadTemplate(templateGoldenPath)
	assert.Nil(t, err)

	template_2, err := LoadTemplate(templateGoldenPath2)
	assert.Nil(t, err)

	assert.Equal(t, template_1, route.Bodies[0])
	assert.Equal(t, template_2, route.Bodies[1])
}

func Test_normalizeRouteConfigurationMultipleBodies(t *testing.T) {

	route := normalizeRouteConfiguration(HTTPHandler{
		BodyFiles: []string{templateGoldenPath, templateGoldenPath2},
	})

	template_1, err := LoadTemplate(templateGoldenPath)
	assert.Nil(t, err)

	template_2, err := LoadTemplate(templateGoldenPath2)
	assert.Nil(t, err)

	assert.Equal(t, template_1, route.Bodies[0])
	assert.Equal(t, template_2, route.Bodies[1])
}
