package main

import (
	. "github.com/omo/fuga/base"
)

type ListCommand struct{}

func (self *ListCommand) Run(args []string, settings CommandSettings) error {
	return nil
}

func (self *ListCommand) Name() string {
	return "list"
}

func init() {
	AddCommand(&ListCommand{})
}
