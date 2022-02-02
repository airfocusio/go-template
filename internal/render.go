package internal

import (
	"bytes"
	"fmt"
	"math/rand"
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
	"first": func(value interface{}) interface{} {
		defer recovery()

		v := reflect.ValueOf(value)

		switch v.Kind() {
		case reflect.String:
			return string([]rune(v.String())[0])
		case reflect.Slice, reflect.Array:
			return v.Index(0).Interface()
		}

		return ""
	},
	"last": func(value interface{}) interface{} {
		defer recovery()

		v := reflect.ValueOf(value)

		switch v.Kind() {
		case reflect.String:
			str := []rune(v.String())
			return string(str[len(str)-1])
		case reflect.Slice, reflect.Array:
			return v.Index(v.Len() - 1).Interface()
		}

		return ""
	},
	"slice": func(start int, end int, value interface{}) interface{} {
		defer recovery()

		v := reflect.ValueOf(value)

		if start < 0 {
			start = 0
		}

		switch v.Kind() {
		case reflect.String:
			str := []rune(v.String())

			if end > len(str) {
				end = len(str)
			}

			return string(str[start:end])
		case reflect.Slice:
			return v.Slice(start, end).Interface()
		}
		return ""
	},
	"join": func(arg string, value []string) string {
		defer recovery()

		return strings.Join(value, arg)
	},
	"random": func(value interface{}) interface{} {
		defer recovery()

		rand.Seed(time.Now().UTC().UnixNano())

		v := reflect.ValueOf(value)

		switch v.Kind() {
		case reflect.String:
			str := []rune(v.String())
			return string(str[rand.Intn(len(str))])
		case reflect.Slice, reflect.Array:
			return v.Index(rand.Intn(v.Len())).Interface()
		}

		return ""
	},
}
