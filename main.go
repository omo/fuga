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

func dropNonFlagArgs(args []string) (nonFlags, flags []string) {
	nonFlags = []string{}

	for ; 0 < len(args) && -1 == strings.Index(args[0], "-"); args = args[1:] {
		nonFlags = append(nonFlags, args[0])
	}

	return nonFlags, args
}

func parseAndBuildArgs() ([]string, error) {
	// Setup extra flags from the dot file.
	dotText, err := readDotFile()
	if nil != err {
		return nil, err
	}

	dotFlags := parseDotFileToArgs(dotText)
	nonFlags, flags := dropNonFlagArgs(os.Args[1:])
	os.Args = append(os.Args[:1], append(dotFlags, flags...)...)

	flag.Parse()

	return append(nonFlags, flag.Args()...), nil
}

// Common Flags
var givenWorkspace = flag.String("workspace", filepath.Join("~", ".fuga"),
	"The directory where fuga generates stubs")

var logVerbose = flag.Bool("verbose", false,
	"Emit verbose log for diagnosing problem")

// Bootstrap
func main() {
	args, err := parseAndBuildArgs()
	panicIfError(err)

	if len(args) <= 0 {
		// FIXME: Use flag.Usage
		fail(fmt.Sprintf("Give one of these commands: %v", ListCommands()))
	}

	if *logVerbose {
		EnableVerboseLog()
	}

	commandName := args[0]
	command := FindCommand(commandName)
	if nil == command {
		fail(fmt.Sprintf("Command %s not found", commandName))
	}

	wd, err := os.Getwd()
	if err != nil {
		fail(err.Error())
	}

	workspace := resolveHome(*givenWorkspace)
	if err := command.Run(args[1:], CommandSettings{Workspace: workspace, Wd: wd}); nil != err {
		fail(err.Error())
	}
}
