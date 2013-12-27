package main

import (
	"fmt"
	. "github.com/omo/fuga/base"
	_ "github.com/omo/fuga/langs"
)

//
// Hello command is made just for debugging. There is no real use.
//
type HelloCommand struct{}

func (self *HelloCommand) Run(args []string) error {
	fmt.Printf("Hello args: %v\n", args)
	return nil
}

func (self *HelloCommand) Name() string {
	return "hello"
}

func init() {
	AddCommand(&HelloCommand{})
}
