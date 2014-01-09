package main

import (
	"errors"
	"fmt"
	. "github.com/omo/fuga/base"
	"regexp"
)

var _ = fmt.Printf

type RunCommand struct{}

var filenameRe = regexp.MustCompile("([[:alnum:]]+)$")

func findLanguageSuffix(path string) string {
	matched := filenameRe.FindAllString(path, -1)
	return matched[0]
}

func (self *RunCommand) Run(args []string, settings CommandSettings) error {
	picked, err := PickBuildUnitFromArgs(settings, args)
	if err != nil {
		return err
	}

	suffix := findLanguageSuffix(picked.Dir())
	runner := FindRunner(suffix)
	if nil == runner {
		return errors.New(fmt.Sprintf("No scratch runner for %s", picked.Dir()))
	}

	return runner.Run(BuildRunnerParams{Unit: picked})
}

func (self *RunCommand) Name() string {
	return "run"
}

func init() {
	AddCommand(&RunCommand{})
}
