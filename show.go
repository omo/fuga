package main

import (
	. "github.com/omo/fuga/base"
)

type ShowCommand struct{}

func (self *ShowCommand) Run(args []string, settings CommandSettings) error {
	SetEditor("cat")
	return FindCommand("open").Run(args, settings)
}

func (self *ShowCommand) Name() string {
	return "show"
}

func init() {
	AddCommand(&ShowCommand{})
}
