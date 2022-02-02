package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/airfocusio/go-template/internal"
)

func run(dir string, opts internal.RenderOptions) error {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("unable to read stdin: %w", err)
	}
	output, err := internal.Render(dir, opts, input)
	if err != nil {
		return fmt.Errorf("unable to generate version: %w", err)
	}
	os.Stdout.Write([]byte(*output))
	return nil
}

func Execute(version FullVersion) error {
	flag.Parse()
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	opts := internal.RenderOptions{}
	err = run(dir, opts)
	if err != nil {
		return err
	}
	return nil
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
