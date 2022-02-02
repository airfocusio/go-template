package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemVerString(t *testing.T) {
	assert := assert.New(t)
	test := func(input string, expectedOutput string, env map[string]string) {
		for key, value := range env {
			os.Setenv(key, value)
		}
		actualOutput, err := Render(".", RenderOptions{}, input)
		if assert.NoError(err) {
			assert.Equal(expectedOutput, *actualOutput)
		}
		for key := range env {
			os.Unsetenv(key)
		}
	}
	testError := func(input string, env map[string]string) {
		for key, value := range env {
			os.Setenv(key, value)
		}
		_, err := Render(".", RenderOptions{}, input)
		assert.Error(err)
		for key := range env {
			os.Unsetenv(key)
		}
	}

	test("", "", map[string]string{})
	test("a", "a", map[string]string{})
	test("a\nb\n", "a\nb\n", map[string]string{})
	test(`Hello, {{ .Env.NAME }}!`, "Hello, Tom!", map[string]string{"NAME": "Tom"})

	test(`Value = {{ .Env.VALUE | default "fallback" }}`, "Value = fallback", map[string]string{})
	test(`Value = {{ .Env.VALUE | default "fallback" }}`, "Value = provided", map[string]string{"VALUE": "provided"})

	test(`Value = {{ "Hello World" | replace "Hello" "Bye" }}`, "Value = Bye", map[string]string{})

	test(`Value = {{ "foo" | length }}`, "Value = 3", map[string]string{})

	test(`Value = {{ "FOO" | lower }}`, "Value = foo", map[string]string{})

	test(`Value = {{ "foo" | upper }}`, "Value = FOO", map[string]string{})

	test(`Value = {{ "  foo  " | trim }}`, "Value = foo", map[string]string{})

	test(`Value = {{ "0" | bool }}`, "Value = false", map[string]string{})
	test(`Value = {{ "1" | bool }}`, "Value = true", map[string]string{})
	test(`Value = {{ "no" | bool }}`, "Value = false", map[string]string{})
	test(`Value = {{ "yes" | bool }}`, "Value = true", map[string]string{})
	test(`Value = {{ "YES" | bool }}`, "Value = true", map[string]string{})
	test(`Value = {{ "false" | bool }}`, "Value = false", map[string]string{})
	test(`Value = {{ "true" | bool }}`, "Value = true", map[string]string{})
	test(`Value = {{ "TRUE" | bool }}`, "Value = true", map[string]string{})

	test(`Value = {{ if .Env.ENABLE | bool }}yes{{ else }}no{{ end }}`, "Value = no", map[string]string{"ENABLE": "0"})
	test(`Value = {{ if .Env.ENABLE | bool }}yes{{ else }}no{{ end }}`, "Value = yes", map[string]string{"ENABLE": "1"})

	test(`Value = {{ require .Env.VALUE }}`, "Value = foo", map[string]string{"VALUE": "foo"})
	testError(`Value = {{ require .Env.MISSING }}`, map[string]string{"VALUE": "foo"})
}
