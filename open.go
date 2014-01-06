package main

import (
	"errors"
	"flag"
	"fmt"
	. "github.com/omo/fuga/base"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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

	var nth uint64 = 0
	if 0 < len(args) {
		n, err := strconv.ParseUint(args[0], 10, 32)
		if nil != err {
			// FIXME: Could give better error message
			return err
		}

		nth = n
	}

	picked := ListPrimaryFiles(settings.Workspace).Pick(uint(nth))
	if !picked.IsValid() {
		return errors.New("Cannot find valid files.")
	}

	return OpenWithEditor(picked.PrimaryFile)
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
