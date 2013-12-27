package main

import (
	"fmt"
	. "github.com/omo/fuga/base"
	_ "github.com/omo/fuga/langs"
	"os/user"
	"path/filepath"
	"time"
)

func makeBaseDir(param *Parameters) string {
	return filepath.Join(
		param.Workspace,
		fmt.Sprintf("%4d", param.Now.Year()),
		fmt.Sprintf("%2d%2d%2d%2d-%s", param.Now.Month(), param.Now.Day(), param.Now.Hour(), param.Now.Minute(), param.Suffix))
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

type GenerateCommand struct{}

func (self *GenerateCommand) Run(args []string) error {
	params := makeParameters(args)

	err := EnsureDir(params.Workspace)
	if nil != err {
		return err
	}

	writer, err := MakeFileStubWriter(makeBaseDir(params))
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
