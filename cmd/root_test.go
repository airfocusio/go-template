package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/airfocusio/go-template/internal"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	input := "Hello, {{ .Val.message }}!"
	inputReader := strings.NewReader(input)
	var outputBytes bytes.Buffer
	outputWriter := io.Writer(&outputBytes)
	err := Run(inputReader, outputWriter, internal.RenderData{
		Val: map[string]interface{}{
			"message": "World",
		},
	})
	assert.NoError(t, err)
	output := outputBytes.String()
	assert.Equal(t, "Hello, World!", output)
}

func TestBuildRenderData(t *testing.T) {
	valueFlags := arrayFlags{"some=thing"}

	valueFile, err := ioutil.TempFile(os.TempDir(), "go-template-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(valueFile.Name())

	err = ioutil.WriteFile(valueFile.Name(), []byte(`
persons:
- name: Adam
  age: 42
- name: Eve
  age: 41
apple: pie
`), 0o664)
	assert.NoError(t, err)
	valueFileFlags := arrayFlags{valueFile.Name()}

	os.Clearenv()
	os.Setenv("FOO", "BAR")

	data, err := BuildRenderData(valueFlags, valueFileFlags)
	assert.NoError(t, err)

	assert.Equal(t, internal.RenderData{
		Val: map[string]interface{}{
			"some": "thing",
			"persons": []interface{}{
				map[string]interface{}{
					"name": "Adam",
					"age":  42,
				},
				map[string]interface{}{
					"name": "Eve",
					"age":  41,
				},
			},
			"apple": "pie",
		},
		Env: map[string]string{
			"FOO": "BAR",
		},
	}, *data)
}
