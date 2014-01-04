package main

import (
	"fmt"
	. "github.com/omo/fuga/base"
	"os"
	"path/filepath"
	"regexp"
)

var _ = fmt.Printf

type ListCommand struct{}

type visitListItem func(string)

func listPrimaryFiles(workspace string, visit visitListItem) {
	digitDirPattern := regexp.MustCompile(`^(\d{4}|\d{8})`)
	fooPattern := regexp.MustCompile(`(?i)^foo\.[[:alnum:]]+$`)

	filepath.Walk(workspace,
		func(path string, info os.FileInfo, err error) error {
			if workspace == path {
				return nil
			}

			basename := filepath.Base(path)

			if info.Mode().IsDir() {
				// Goes into seemingly generated directories only.
				if nil == digitDirPattern.FindStringIndex(basename) {
					return filepath.SkipDir
				}

				return nil
			}

			if nil == fooPattern.FindStringIndex(basename) {
				return nil
			}

			visit(path)
			return nil
		})
}

func (self *ListCommand) Run(args []string, settings CommandSettings) error {
	listPrimaryFiles(settings.Workspace,
		func(path string) {
			fmt.Printf("%s\n", path)
		})

	return nil
}

func (self *ListCommand) Name() string {
	return "list"
}

func init() {
	AddCommand(&ListCommand{})
}