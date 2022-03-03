package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v3"
)

type RenderData struct {
	Val map[string]interface{}
	Env map[string]string
}

func Render(data RenderData, input string) (*string, error) {
	tmpl, err := template.New("template").Funcs(FuncMap()).Parse(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template: %w", err)
	}
	buf := new(bytes.Buffer)
	if buf == nil {
		return nil, fmt.Errorf("unable to initialize render buffer")
	}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, fmt.Errorf("unable to render template: %w", err)
	}
	output := buf.String()
	return &output, nil
}

func FuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	extra := template.FuncMap{
		"required":      required,
		"toJson":        toJson,
		"toPrettyJson":  toPrettyJson,
		"fromJson":      fromJson,
		"fromJsonArray": fromJsonArray,
		"toYaml":        toYaml,
		"fromYaml":      fromYaml,
		"fromYamlArray": fromYamlArray,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func logError(err error) {
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("Warn: %v\n", err)))
	}
}

func required(warn string, val interface{}) (interface{}, error) {
	if val == nil {
		return val, fmt.Errorf(warn)
	} else if val, ok := val.(string); ok && val == "" {
		return val, fmt.Errorf(warn)
	}
	return val, nil
}

func toJson(v interface{}) string {
	data, err := json.Marshal(v)
	logError(err)
	return strings.TrimSuffix(string(data), "\n")
}

func toPrettyJson(v interface{}) string {
	output, err := json.MarshalIndent(v, "", "  ")
	logError(err)
	return string(output)
}

func fromJson(str string) map[string]interface{} {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &data)
	logError(err)
	return data
}

func fromJsonArray(str string) []interface{} {
	data := []interface{}{}
	err := json.Unmarshal([]byte(str), &data)
	logError(err)
	return data
}

func toYaml(v interface{}) string {
	data, err := yaml.Marshal(v)
	logError(err)
	return strings.TrimSuffix(string(data), "\n")
}

func fromYaml(str string) map[string]interface{} {
	data := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(str), &data)
	logError(err)
	return data
}

func fromYamlArray(str string) []interface{} {
	data := []interface{}{}
	err := yaml.Unmarshal([]byte(str), &data)
	logError(err)
	return data
}
