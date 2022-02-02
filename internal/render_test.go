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

	test("", "", map[string]string{})
	test("a", "a", map[string]string{})
	test("a\nb\n", "a\nb\n", map[string]string{})
	test(`Hello, {{ .Env.NAME }}!`, "Hello, Tom!", map[string]string{"NAME": "Tom"})
}
