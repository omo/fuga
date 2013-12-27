package langs

import (
	base "github.com/omo/fuga/base"
)

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

func (*CGenerator) Generate(writer base.StubWriter) error {
	writer.WriteFile("foo.c", sourceTempalte)
	writer.WriteFile("Makefile", makefileTemplate)
	return nil
}

func init() {
	base.AddGenerator("c", &CGenerator{})
}
