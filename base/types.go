package core

import "time"

//
// For generators
//

// FIXME: rename to GeneratorParameters
type Parameters struct {
	Workspace string
	Now       time.Time
	Suffix    string
}

type StubWriter interface {
	WriteFile(filename, content string)
	LastError() error
	PrimaryFileName() string
}

type StubGenerator interface {
	Generate(StubWriter) error
}

type GeneratorTable map[string]StubGenerator

var theGeneratorTable = GeneratorTable{}

func AddGenerator(suffix string, generator StubGenerator) {
	theGeneratorTable[suffix] = generator
}

func FindGenerator(suffix string) StubGenerator {
	return theGeneratorTable[suffix]
}

//
// For commands
//
type CommandSettings struct {
	Workspace string
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
