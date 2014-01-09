package main

import (
	"errors"
	"flag"
	"fmt"
	. "github.com/omo/fuga/base"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

var _ = fmt.Printf

// FIXME: Move to separate file
type BuildUnitList []BuildUnit

func imin(x, y int) int {
	return int(math.Min(float64(x), float64(y)))
}

func (self BuildUnitList) Round(n int) int {
	return imin(n, len(self))
}

func (self BuildUnitList) Pick(n uint) BuildUnit {
	if 0 == len(self) || len(self) <= int(n) {
		return BuildUnit{}
	}

	return self[n]
}

func PickBuildUnitFromArgs(settings CommandSettings, args []string) (BuildUnit, error) {
	var nth uint64 = 0
	if 0 < len(args) {
		n, err := strconv.ParseUint(args[0], 10, 32)
		if nil != err {
			// FIXME: Could give better error message
			return BuildUnit{}, err
		}

		nth = n
	}

	picked := ListBuildUnits(settings.Workspace).Pick(uint(nth))
	if !picked.IsValid() {
		return BuildUnit{}, errors.New("Cannot find valid files.")
	}

	return picked, nil
}

// sort.Interface
func (a BuildUnitList) Len() int           { return len(a) }
func (a BuildUnitList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BuildUnitList) Less(i, j int) bool { return a[i].PrimaryFile() > a[j].PrimaryFile() }

func ListBuildUnits(workspace string) BuildUnitList {
	ret := BuildUnitList{}

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

			ret = append(ret, MakeBuildUnit(path))
			return nil
		})

	sort.Sort(ret)
	return ret
}

func exp10(num int) int {
	ret := 0
	for 0 < num {
		ret = ret + 1
		num = num / 10
	}

	return ret
}

type ListCommand struct{}

func (self *ListCommand) Run(args []string, settings CommandSettings) error {

	count := 0
	if 0 < len(args) {
		n, err := strconv.ParseUint(args[0], 10, 32)
		if nil != err {
			// FIXME Could give better error message
			return err
		}

		count = int(n)
	}

	list := ListBuildUnits(settings.Workspace)
	if count != 0 {
		list = list[0:list.Round(count)]
	}

	width := exp10(len(list))
	for i, e := range list {
		if *shouldPrintOridinalNumber {
			fmt.Printf("%"+strconv.FormatInt(int64(width), 10)+"d: %s\n", i, e.PrimaryFile())
		} else {
			fmt.Printf("%s\n", e.PrimaryFile())
		}
	}

	return nil
}

func (self *ListCommand) Name() string {
	return "list"
}

var shouldPrintOridinalNumber = flag.Bool("number", false,
	"list: Print ordinal numbers for each item")

func init() {
	AddCommand(&ListCommand{})
}
