package main

import (
	"fmt"
	. "github.com/omo/fuga/base"
	"testing"
	"time"
)

var _ = fmt.Printf

func TestHello(t *testing.T) {
	// Do nothing.
}

func makeTestParameters(t *testing.T) *Parameters {
	now, err := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Dec  5, 2013 at 5:04pm (JST)")
	expectOK(err, t)
	return &Parameters{"/root", now, "java"}
}

func TestMakeBaseDir(t *testing.T) {
	actual := makeBaseDir(makeTestParameters(t))
	expect(actual, "/root/2013/12051704-java", t)
}

func TestCGenerator(t *testing.T) {
	wri := MakeTestingStubWriter()
	gen := FindGenerator("c")
	WriteStub(wri, gen)
	expectTrue(wri.IsWritten("foo.c"), "foo.c", t)
	expectTrue(wri.IsWritten("Makefile"), "Makefile", t)
}

func TestFindCommand(t *testing.T) {
	fullnameCommand := FindCommand("generate")
	expect(fullnameCommand.Name(), "generate", t)
	abbreviatedCommand := FindCommand("g")
	expect(abbreviatedCommand.Name(), "generate", t)
}

func TestNonConflictingName(t *testing.T) {
	expect(nonConflictingName("/foo/1234-c"), "/foo/1234-c-001", t)
	expect(nonConflictingName("/foo/1234-c-001"), "/foo/1234-c-002", t)
}

func TestParseDotFileToArgs(t *testing.T) {
	flagString := `
--flaga=x
--flagb=y # Should be skipped
   
--flagc=z
`
	flags := parseDotFileToArgs(flagString)
	expectTrue(len(flags) == 3, "len(flags)", t)
	expect(flags[0], "--flaga=x", t)
	expect(flags[1], "--flagb=y", t)
	expect(flags[2], "--flagc=z", t)
}

func TestListPrimaryFiles(t *testing.T) {
	listed := []string{}
	listPrimaryFiles("./testroot", func(path string) {
		listed = append(listed, path)
	})

	// FIXME: add files from other generators.
	expectTrue(1 == len(listed), "len(listed)", t)
	expect(listed[0], "testroot/2014/01042256-c/foo.c", t)
}

// Copied from github.com/eknkc/amber/amber_test.go
func expect(cur, expected string, t *testing.T) {
	if cur != expected {
		t.Fatalf("Expected {%s} got {%s}.", expected, cur)
	}
}

func expectOK(err error, t *testing.T) {
	if err != nil {
		t.Fatal("Should be OK")
	}
}

func expectTrue(ok bool, subject string, t *testing.T) {
	if !ok {
		t.Fatalf("%s should be OK", subject)
	}
}
