package core

import (
	"path/filepath"
	"time"
)

//
// For generators
//

// FIXME: rename to GeneratorParameters
type Parameters struct {
	Workspace string
	Now       time.Time
	Suffix    string
}

type ScratchWriter interface {
	WriteFile(filename, content string)
	LastError() error
	PrimaryFileName() string
}

type StubGenerator interface {
	Generate(ScratchWriter) error
}

type BuildUnit struct {
	primaryFile string
}

func (self BuildUnit) PrimaryFile() string {
	return self.primaryFile
}

func (self BuildUnit) Dir() string {
	return filepath.Dir(self.primaryFile)
}

func (self BuildUnit) PrimaryBase() string {
	return filepath.Base(self.primaryFile)
}

func (self BuildUnit) IsValid() bool {
	return "" != self.PrimaryFile()
}

func MakeBuildUnit(file string) BuildUnit {
	return BuildUnit{primaryFile: file}
}

type ScratchRunnerParams struct {
	Unit BuildUnit
}

type ScratchRunner interface {
	Run(params ScratchRunnerParams) error
}

type Language interface {
	MakeGenerator() StubGenerator
	MakeRunner() ScratchRunner
}

type LanguageTable map[string]Language

var theLanguageTable = LanguageTable{}

func AddLanguage(suffix string, lang Language) {
	theLanguageTable[suffix] = lang
}

func FindLanguage(suffix string) Language {
	return theLanguageTable[suffix]
}

func FindGenerator(suffix string) StubGenerator {
	return FindLanguage(suffix).MakeGenerator()
}

func FindRunner(suffix string) ScratchRunner {
	return FindLanguage(suffix).MakeRunner()
}

//
// For commands
//
type CommandSettings struct {
	Workspace string
	Wd        string
}

type Command interface {
	Name() string
	Run(args []string, settings CommandSettings) error
}

type CommandList []Command

var theCommandList = CommandList{}

func AddCommand(command Command) {
	theCommandList = append(theCommandList, command)
}

func FindCommand(name string) Command {
	for _, c := range theCommandList {
		if c.Name() == name || c.Name()[0:1] == name {
			return c
		}
	}

	return nil
}

func ListCommands() []string {
	ret := []string{}
	for _, c := range theCommandList {
		ret = append(ret, c.Name())
	}

	return ret
}
