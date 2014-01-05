package main

import (
	"fmt"
	. "github.com/omo/fuga/base"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var _ = fmt.Printf

type ListCommand struct{}

type visitListItem func(string)

type ListEntry struct {
	PrimaryFile string
}

type ListEntryList []ListEntry

func (a ListEntryList) Len() int           { return len(a) }
func (a ListEntryList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ListEntryList) Less(i, j int) bool { return a[i].PrimaryFile > a[j].PrimaryFile }

func listPrimaryFiles(workspace string) []ListEntry {
	ret := ListEntryList{}

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

			ret = append(ret, ListEntry{PrimaryFile: path})
			return nil
		})

	sort.Sort(ret)
	return ret
}

func (self *ListCommand) Run(args []string, settings CommandSettings) error {
	for _, e := range listPrimaryFiles(settings.Workspace) {
		fmt.Printf("%s\n", e.PrimaryFile)
	}

	return nil
}

func (self *ListCommand) Name() string {
	return "list"
}

func init() {
	AddCommand(&ListCommand{})
}
