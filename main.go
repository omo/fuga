package main

import (
	"flag"
	"fmt"
	. "github.com/omo/mkfoo/base"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

type JavaGenerator struct{}

func (*JavaGenerator) Generate(writer StubWriter) error {
	return nil
}

type CGenerator struct{}

const sourceTempalte = `
#include <stdio.h>

int main(int argc, char* argv[]) {
  printf("Hello, World!\n");
  return 0;
}
`

const makefileTemplate = `
CC=gcc
TARGET=./foo
SOURCE=./foo.c

run : ${TARGET}
	${TARGET}

clean :
	-rm ${TARGET}

${TARGET} : ${SOURCE}
	${CC} $^ -o $@

.PHONY : run clean
`

func (*CGenerator) Generate(writer StubWriter) error {
	writer.WriteFile("foo.c", sourceTempalte)
	writer.WriteFile("Makefile", makefileTemplate)
	return nil
}

type GeneratorTable map[string]StubGenerator

func makeTable() GeneratorTable {
	return GeneratorTable{
		"java": &JavaGenerator{},
		"c":    &CGenerator{},
	}
}

func makeGenerator(suffix string) StubGenerator {
	return makeTable()[suffix]
}

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

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) <= 0 {
		// FIXME: Use flag.Usage
		fmt.Printf("Specify prefix\n")
		os.Exit(1)
	}

	params := makeParameters(args)
	err := EnsureDir(params.Workspace)
	panicIfError(err)
	writer, err := MakeFileStubWriter(makeBaseDir(params))
	panicIfError(err)
	gen := makeGenerator(params.Suffix)
	err = gen.Generate(writer)
	panicIfError(err)
	panicIfError(writer.LastError())
}
