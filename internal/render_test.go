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
	testError := func(input string, data RenderData) {
		_, err := Render(data, input)
		assert.Error(err)
	}

	test("", "", RenderData{})
	test("a", "a", RenderData{})
	test("a\nb\n", "a\nb\n", RenderData{})

	test(`Hello, {{ .Env.NAME }}!`, "Hello, Tom!", RenderData{Env: map[string]string{"NAME": "Tom"}})
	test(`Hello, {{ .Val.name }}!`, "Hello, Tom!", RenderData{Val: map[string]interface{}{"name": "Tom"}})
	test(`Hello, {{ .Val.person.name }}!`, "Hello, Tom!", RenderData{Val: map[string]interface{}{"person": map[string]interface{}{"name": "Tom"}}})

	test(`Value = {{ .Env.VALUE | default "fallback" }}`, "Value = fallback", RenderData{})
	test(`Value = {{ .Env.VALUE | default "fallback" }}`, "Value = provided", RenderData{Env: map[string]string{"VALUE": "provided"}})

	test(`Value = {{ "Hello World" | replace "Hello" "Bye" }}`, "Value = Bye", RenderData{})

	test(`Value = {{ "foo" | length }}`, "Value = 3", RenderData{})

	test(`Value = {{ "FOO" | lower }}`, "Value = foo", RenderData{})

	test(`Value = {{ "foo" | upper }}`, "Value = FOO", RenderData{})

	test(`Value = {{ "  foo  " | trim }}`, "Value = foo", RenderData{})

	test(`Value = {{ "0" | bool }}`, "Value = false", RenderData{})
	test(`Value = {{ "1" | bool }}`, "Value = true", RenderData{})
	test(`Value = {{ "no" | bool }}`, "Value = false", RenderData{})
	test(`Value = {{ "yes" | bool }}`, "Value = true", RenderData{})
	test(`Value = {{ "YES" | bool }}`, "Value = true", RenderData{})
	test(`Value = {{ "false" | bool }}`, "Value = false", RenderData{})
	test(`Value = {{ "true" | bool }}`, "Value = true", RenderData{})
	test(`Value = {{ "TRUE" | bool }}`, "Value = true", RenderData{})

	test(`Value = {{ if .Env.ENABLE | bool }}yes{{ else }}no{{ end }}`, "Value = no", RenderData{Env: map[string]string{"ENABLE": "0"}})
	test(`Value = {{ if .Env.ENABLE | bool }}yes{{ else }}no{{ end }}`, "Value = yes", RenderData{Env: map[string]string{"ENABLE": "1"}})

	test(`Value = {{ require .Env.VALUE }}`, "Value = foo", RenderData{Env: map[string]string{"VALUE": "foo"}})
	testError(`Value = {{ require .Env.MISSING }}`, RenderData{Env: map[string]string{"VALUE": "foo"}})
}
