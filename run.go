package main

import (
	"fmt"
	. "github.com/omo/fuga/base"
)

var _ = fmt.Printf

type RunCommand struct{}

func (self *RunCommand) Run(args []string, settings CommandSettings) error {
	picked, err := PickBuildUnitFromArgs(settings, args)
	if err != nil {
		return err
	}

	// FIXME: impl
	fmt.Printf("run: %s\n", picked.Dir())
	return nil
}

func (self *RunCommand) Name() string {
	return "run"
}

func init() {
	AddCommand(&RunCommand{})
}
