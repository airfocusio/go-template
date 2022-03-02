package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemVerString(t *testing.T) {
	assert := assert.New(t)
	test := func(input string, expectedOutput string, data RenderData) {
		actualOutput, err := Render(data, input)
		if assert.NoError(err) {
			assert.Equal(expectedOutput, *actualOutput)
		}
	}
	testError := func(input string, data RenderData, errorMessage string) {
		_, err := Render(data, input)
		assert.EqualError(err, errorMessage)
	}

	test("", "", RenderData{})
	test("a", "a", RenderData{})
	test("a\nb\n", "a\nb\n", RenderData{})

	test(`Hello, {{ .Env.NAME }}!`, "Hello, Tom!", RenderData{Env: map[string]string{"NAME": "Tom"}})
	test(`Hello, {{ .Val.name }}!`, "Hello, Tom!", RenderData{Val: map[string]interface{}{"name": "Tom"}})
	test(`Hello, {{ .Val.person.name }}!`, "Hello, Tom!", RenderData{Val: map[string]interface{}{"person": map[string]interface{}{"name": "Tom"}}})

	test(`Value = {{ .Env.VALUE | default "fallback" }}`, "Value = fallback", RenderData{})
	test(`Value = {{ .Env.VALUE | default "fallback" }}`, "Value = provided", RenderData{Env: map[string]string{"VALUE": "provided"}})

	test(`Value = {{ .Env.VALUE | required "missing" }}`, "Value = foo", RenderData{Env: map[string]string{"VALUE": "foo"}})
	testError(`Value = {{ .Env.MISSING | required "missing" }}`, RenderData{Env: map[string]string{"VALUE": "foo"}}, "unable to render template: template: template:1:26: executing \"template\" at <required \"missing\">: error calling required: missing")
}
