package main

import (
	"flag"
	"fmt"
	. "github.com/omo/fuga/base"
	_ "github.com/omo/fuga/langs"
	"os"
)

type JavaGenerator struct{}

func (*JavaGenerator) Generate(writer StubWriter) error {
	return nil
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func fail(message string) {
	fmt.Printf(message)
	fmt.Printf("\n")
	os.Exit(1)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) <= 0 {
		// FIXME: Use flag.Usage
		fail("Specify prefix")
	}

	commandName := args[0]
	command := FindCommand(commandName)
	if nil == command {
		fail(fmt.Sprintf("Command %s not found", commandName))
	}

	command.Run(args[1:])
}
