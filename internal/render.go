package internal

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

type RenderData struct {
	Val map[string]interface{}
	Env map[string]string
}

func Render(data RenderData, input string) (*string, error) {
	tmpl, err := template.New("template").Funcs(funcMap()).Parse(input)
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

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	extra := template.FuncMap{
		"required": fnRequired,
		"bool":     fnBool,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f

}

func fnRequired(warn string, val interface{}) (interface{}, error) {
	if val == nil {
		return val, fmt.Errorf(warn)
	} else if val, ok := val.(string); ok && val == "" {
		return val, fmt.Errorf(warn)
	}
	return val, nil
}

func fnBool(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.String:
		lower := strings.ToLower(v.String())
		return lower == "1" || lower == "yes" || lower == "true"
	}

	return false
}
