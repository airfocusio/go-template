package internal

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"
)

type RenderData struct {
	Now time.Time
	Val map[string]string
	Env map[string]string
}

func Render(data RenderData, input string) (*string, error) {
	tmpl, err := template.New("template").Funcs(funcMap).Parse(input)
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
	"require": func(value interface{}) (interface{}, error) {
		if value == nil {
			return nil, fmt.Errorf("value is missing")
		}
		return value, nil
	},
}
