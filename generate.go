package main

import (
	"errors"
	"fmt"
	. "github.com/omo/fuga/base"
	_ "github.com/omo/fuga/langs"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

type EnsureDirOption int

const (
	NeedsFresh      EnsureDirOption = iota
	DoesntNeedFresh EnsureDirOption = iota
)

type MkdirAllError struct {
	error string
}

func (self *MkdirAllError) Error() string {
	return self.error
}

func EnsureDir(dirname string, option EnsureDirOption) error {
	stat, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dirname, 0777); err != nil {
			return &MkdirAllError{err.Error()}
		}

		return nil
	}

	if !stat.IsDir() || option == NeedsFresh {
		return errors.New(fmt.Sprintf("File %s is already exist!", dirname))
	}

	return nil
}

func nonConflictingName(name string) string {
	pattern := regexp.MustCompile(`[0-9]+$`)
	ordinal := pattern.FindString(name)

	if ordinal == "" {
		return name + "-001"
	}

	ordinalAsNumber, _ := strconv.ParseUint(ordinal, 10, 32)
	return pattern.ReplaceAllString(name, fmt.Sprintf("%03d", ordinalAsNumber+1))
}

func makeBaseDir(param *Parameters) string {
	return filepath.Join(
		param.Workspace,
		fmt.Sprintf("%04d", param.Now.Year()),
		fmt.Sprintf("%02d%02d%02d%02d-%s", param.Now.Month(), param.Now.Day(), param.Now.Hour(), param.Now.Minute(), param.Suffix))
}

func defaultWorkspace() string {
	usr, err := user.Current()
	panicIfError(err)
	return filepath.Join(usr.HomeDir, "work", "foos")
}

func makeParameters(args []string) *Parameters {
	return &Parameters{
		defaultWorkspace(),
		time.Now(),
		args[0],
	}
}

func ensureBaseDir(name string) error {
	for made := false; !made; {
		if err := EnsureDir(name, NeedsFresh); err != nil {
			switch err.(type) {
			case *MkdirAllError:
				return err
			default:
				name = nonConflictingName(name)
			}
		} else {
			made = true
		}
	}

	return nil
}

type GenerateCommand struct{}

func (self *GenerateCommand) Run(args []string) error {
	params := makeParameters(args)

	err := EnsureDir(params.Workspace, DoesntNeedFresh)
	if nil != err {
		return err
	}

	baseDir := makeBaseDir(params)
	if err := ensureBaseDir(baseDir); err != nil {
		return err
	}

	writer, err := MakeFileStubWriter(baseDir)
	if nil != err {
		return err
	}

	gen := FindGenerator(params.Suffix)
	if err := gen.Generate(writer); nil != err {
		return err
	}

	if nil != writer.LastError() {
		return writer.LastError()
	}

	return nil
}

func (self *GenerateCommand) Name() string {
	return "generate"
}

func init() {
	AddCommand(&GenerateCommand{})
}
