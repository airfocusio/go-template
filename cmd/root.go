package cmd

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/airfocusio/go-template/pkg"
	"gopkg.in/yaml.v3"
)

func Execute(version FullVersion) error {
	var valueFlags arrayFlags
	flag.Var(&valueFlags, "value", "Provide values like key=value (can be repeated)")
	var valueFileFlags arrayFlags
	flag.Var(&valueFileFlags, "value-file", "Provide value files with yaml (can be repeated)")
	flag.Parse()

	data, err := BuildRenderData(valueFlags, valueFileFlags)
	if err != nil {
		return err
	}
	err = Run(os.Stdin, os.Stdout, *data)
	if err != nil {
		return err
	}
	return nil
}

func Run(stdin io.Reader, stdout io.Writer, data pkg.RenderData) error {
	input, err := ioutil.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("unable to read stdin: %w", err)
	}
	output, err := pkg.Render(data, string(input))
	if err != nil {
		return fmt.Errorf("unable to render: %w", err)
	}
	stdout.Write([]byte(*output))
	return nil
}

func BuildRenderData(valueFlags arrayFlags, valueFileFlags arrayFlags) (*pkg.RenderData, error) {
	val := map[string]interface{}{}
	for _, entry := range valueFlags {
		segments := strings.SplitN(entry, "=", 2)
		key := segments[0]
		value := ""
		if len(segments) == 2 {
			value = segments[1]
		}
		val[key] = value
	}

	for _, file := range valueFileFlags {
		yamlBytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("unable to build render data: %w", err)
		}
		var yamlValue map[string]interface{}
		err = yaml.Unmarshal(yamlBytes, &yamlValue)
		if err != nil {
			return nil, fmt.Errorf("unable to build render data: %w", err)
		}

		for k, v := range yamlValue {
			val[k] = v
		}
	}

	env := map[string]string{}
	for _, envEntry := range os.Environ() {
		envEntryParts := strings.SplitN(envEntry, "=", 2)
		key := envEntryParts[0]
		value := envEntryParts[1]
		env[key] = value
	}
	return &pkg.RenderData{Val: val, Env: env}, nil
}

type FullVersion struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

func (v FullVersion) ToString() string {
	result := v.Version
	if v.Commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, v.Commit)
	}
	if v.Date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, v.Date)
	}
	if v.BuiltBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, v.BuiltBy)
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		result = fmt.Sprintf("%s\nmodule version: %s, checksum: %s", result, info.Main.Version, info.Main.Sum)
	}
	return result
}
