package internal

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
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
	tmpl, err := template.New("template").Funcs(funcMap).Parse(input)
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

// recovery will silently swallow all unexpected panics.
func recovery() {
	recover()
}

// inspired by https://github.com/leekchan/gtf
var funcMap = template.FuncMap{
	"default": func(arg interface{}, value interface{}) interface{} {
		defer recovery()

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Invalid:
			return arg
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			if v.Len() == 0 {
				return arg
			}
		case reflect.Bool:
			if !v.Bool() {
				return arg
			}
		default:
			return value
		}

		return value
	},
	"replace": func(s1 string, s2 string, s3 string) string {
		defer recovery()

		return strings.Replace(s2, s1, s3, -1)
	},
	"length": func(value interface{}) int {
		defer recovery()

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return v.Len()
		case reflect.String:
			return len([]rune(v.String()))
		}

		return 0
	},
	"lower": func(s string) string {
		defer recovery()

		return strings.ToLower(s)
	},
	"upper": func(s string) string {
		defer recovery()

		return strings.ToUpper(s)
	},
	"trim": func(s string) string {
		defer recovery()

		return strings.TrimSpace(s)
	},
	"bool": func(value interface{}) bool {
		defer recovery()

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Bool:
			return v.Bool()
		case reflect.String:
			lower := strings.ToLower(v.String())
			return lower == "1" || lower == "yes" || lower == "true"
		}

		return false
	},
}
