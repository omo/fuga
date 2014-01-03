package main

import (
	"flag"
	"fmt"
	. "github.com/omo/fuga/base"
	_ "github.com/omo/fuga/langs"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
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

func readDotFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ".fugarc"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", err
	}

	return string(data), nil
}

func parseDotFileToArgs(text string) []string {
	splitPattern := regexp.MustCompile(`\r?\n`)
	emptyPattern := regexp.MustCompile(`^$`)
	commentPattern := regexp.MustCompile(`^(.*)#`)

	splittedString := splitPattern.Split(text, -1)
	ret := []string{}
	for _, v := range splittedString {

		if mayIncludeComment := commentPattern.FindString(v); mayIncludeComment != "" {
			v = mayIncludeComment[0 : len(mayIncludeComment)-1]
		}

		v = strings.Trim(v, " \t")
		if emptyPattern.MatchString(v) {
			continue
		}

		ret = append(ret, v)
	}

	return ret
}

// Common Flags
var givenWorkspace = flag.String("workspace", filepath.Join("~", ".fuga"),
	"The directory where fuga generates stubs")

// Bootstrap
func main() {
	// Setup extra flags from the dot file.
	dotText, err := readDotFile()
	panicIfError(err)
	dotFlags := parseDotFileToArgs(dotText)
	os.Args = append(os.Args[:1], append(dotFlags, os.Args[1:]...)...)

	flag.Parse()
	args := flag.Args()

	if len(args) <= 0 {
		// FIXME: Use flag.Usage
		fail(fmt.Sprintf("Give one of these commands: %v", ListCommands()))
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
