package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/airfocusio/go-template/internal"
)

var valueFlags arrayFlags

func Execute(version FullVersion) error {
	flag.Var(&valueFlags, "value", "Provide values like key=value (can be repeated)")
	flag.Parse()

	opts := buildRenderData()
	err := run(opts)
	if err != nil {
		return err
	}
	return nil
}

func run(data internal.RenderData) error {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("unable to read stdin: %w", err)
	}
	output, err := internal.Render(data, string(input))
	if err != nil {
		return fmt.Errorf("unable to generate version: %w", err)
	}
	os.Stdout.Write([]byte(*output))
	return nil
}

func buildRenderData() internal.RenderData {
	val := map[string]string{}
	for _, entry := range valueFlags {
		segments := strings.SplitN(entry, "=", 2)
		key := segments[0]
		value := ""
		if len(segments) == 2 {
			value = segments[1]
		}
		val[key] = value
	}

	env := map[string]string{}
	for _, envEntry := range os.Environ() {
		envEntryParts := strings.SplitN(envEntry, "=", 2)
		key := envEntryParts[0]
		value := envEntryParts[1]
		env[key] = value
	}
	return internal.RenderData{Now: time.Now(), Val: val, Env: env}
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
