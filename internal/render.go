package internal

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
)

type RenderOptions struct{}

type RenderData struct {
	Now time.Time
	Env map[string]string
}

func Render(dir string, opts RenderOptions, input string) (*string, error) {
	tmpl, err := template.New("template").Parse(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template: %w", err)
	}
	buf := new(bytes.Buffer)
	if buf == nil {
		return nil, fmt.Errorf("unable to initialize render buffer")
	}
	err = tmpl.Execute(buf, BuildRenderData())
	if err != nil {
		return nil, fmt.Errorf("unable to render template: %w", err)
	}
	output := buf.String()
	return &output, nil
}

func BuildRenderData() RenderData {
	env := map[string]string{}
	for _, envEntry := range os.Environ() {
		envEntryParts := strings.SplitN(envEntry, "=", 2)
		key := envEntryParts[0]
		value := envEntryParts[1]
		env[key] = value
	}
	return RenderData{Now: time.Now(), Env: env}
}
