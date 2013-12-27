package main

import (
	"errors"
	"fmt"
	. "github.com/omo/fuga/base"
	_ "github.com/omo/fuga/langs"
	"strings"
)

//
// Hello command is made just for debugging. There is no real use.
//
type HelloCommand struct{}

func (self *HelloCommand) Run(args []string) error {
	fmt.Printf("Hello args: %v\n", args)
	if "error" == args[0] {
		return errors.New(fmt.Sprintf("Hello error: %s", strings.Join(args[1:], " ")))
	}

	return nil
}

func (self *HelloCommand) Name() string {
	return "hello"
}

func init() {
	AddCommand(&HelloCommand{})
}
