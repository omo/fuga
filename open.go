package main

import (
	"flag"
	"fmt"
	. "github.com/omo/fuga/base"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var _ = fmt.Printf

type OpenCommand struct{}

func makeEditorCommandArgs(given, filename string) []string {
	// FIXME: should support quoted blank in better way.
	blankPattern := regexp.MustCompile(`\s+`)
	trimmed := strings.Trim(given, " \"")
	return append(blankPattern.Split(trimmed, -1), filename)
}

func OpenWithEditor(filename string) error {
	params := makeEditorCommandArgs(*givenEditor, filename)
	cmd := exec.Command(params[0], params[1:]...)
	// http://stackoverflow.com/questions/8875038/
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return err
	}

	_ = cmd.Wait()
	return nil
}

func (self *OpenCommand) Run(args []string, settings CommandSettings) error {

	picked, err := PickBuildUnitFromArgs(settings, args)
	if err != nil {
		return err
	}

	return OpenWithEditor(picked.PrimaryFile())
}

func (self *OpenCommand) Name() string {
	return "open"
}

var givenEditor = flag.String("editor", "vi",
	"open, --open: Editor program to open the file.")

func SetEditor(name string) {
	*givenEditor = name
}

func init() {
	AddCommand(&OpenCommand{})
}
