package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	expectedTemplate = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Nam eget tempor ${params.placeholder_1}, ut ultricies leo. Nunc sit amet orci iaculis, euismod neque sit amet, bibendum magna.
Nullam eget interdum risus.
Aenean euismod, metus sed fermentum tincidunt, ex ${params.placeholder_2} sollicitudin lorem, eu aliquam arcu tortor ac massa.
Proin ultrices eget ligula ut ${params.placeholder_1}.
Sed tincidunt euismod urna vitae tincidunt.
Cras nec aliquam velit.
Praesent non dui ac ${global.placeholder_1} velit ${global.placeholder_2} laoreet vel in turpis.
Donec euismod tellus vel sem vulputate ${params.placeholder_3}.
Suspendisse feugiat quam quis sagittis fringilla.
Nullam dictum vehicula libero nec ${global.placeholder_3}.
Sed suscipit vitae justo vel porttitor.`
	expectedTemplateResulst = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Nam eget tempor bulbasaur, ut ultricies leo. Nunc sit amet orci iaculis, euismod neque sit amet, bibendum magna.
Nullam eget interdum risus.
Aenean euismod, metus sed fermentum tincidunt, ex charmander sollicitudin lorem, eu aliquam arcu tortor ac massa.
Proin ultrices eget ligula ut bulbasaur.
Sed tincidunt euismod urna vitae tincidunt.
Cras nec aliquam velit.
Praesent non dui ac bulbasaur velit charmander laoreet vel in turpis.
Donec euismod tellus vel sem vulputate blastoise.
Suspendisse feugiat quam quis sagittis fringilla.
Nullam dictum vehicula libero nec blastoise.
Sed suscipit vitae justo vel porttitor.`
	placeholder_1       = "bulbasaur"
	placeholder_2       = "charmander"
	placeholder_3       = "blastoise"
	templateGoldenPath  = "../../test/lorem.golden.tpl"
	templateGoldenPath2 = "../../test/ipsum.golden.tpl"
)

func Test_Load(t *testing.T) {
	template, err := LoadTemplate(templateGoldenPath)

	assert.Nil(t, err)
	assert.Equal(t, expectedTemplate, template)

	template, err = LoadTemplate(templateGoldenPath + "foo")

	assert.NotNil(t, err)
	assert.Equal(t, "", template)
}

func Test_ExecuteTemplate(t *testing.T) {
	template, err := LoadTemplate(templateGoldenPath)
	assert.Nil(t, err)

	global := map[string]string{
		"placeholder_1": "bulbasaur",
		"placeholder_2": "charmander",
		"placeholder_3": "blastoise",
	}

	vars := map[string]string{
		"placeholder_1": "bulbasaur",
		"placeholder_2": "charmander",
		"placeholder_3": "blastoise",
	}

	result := ExecuteTemplate(template, global, vars)

	assert.Equal(t, expectedTemplateResulst, result)
}
