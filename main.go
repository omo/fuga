package main

import (
	"flag"
	"fmt"
	. "github.com/omo/fuga/base"
	_ "github.com/omo/fuga/langs"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
)

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

func resolveHome(path string) string {
	usr, err := user.Current()
	panicIfError(err)
	pattern := regexp.MustCompile(`^~`)
	return pattern.ReplaceAllString(path, usr.HomeDir)
}

// Common Flags
var givenWorkspace = flag.String("workspace", filepath.Join("~", ".fuga"),
	"The directory where fuga generates stubs")

// Bootstrap
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

	workspace := resolveHome(*givenWorkspace)
	if err := command.Run(args[1:], CommandSettings{Workspace: workspace}); nil != err {
		fail(err.Error())
	}
}
