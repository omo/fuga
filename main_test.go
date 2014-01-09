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

func TestCppGenerator(t *testing.T) {
	wri := MakeTestingStubWriter()
	gen := FindGenerator("cpp")
	WriteStub(wri, gen)
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

func TestListBuildUnits(t *testing.T) {
	entries := ListBuildUnits("./testroot")
	listed := []string{}
	for _, e := range entries {
		listed = append(listed, e.PrimaryFile())
	}

	// FIXME: add files from other generators.
	expectTrue(1 == len(listed), "len(listed)", t)
	expect(listed[0], "testroot/2014/01042256-c/foo.c", t)

	e0 := entries.Pick(0)
	expectTrue(e0.IsValid(), "e0.IsValid()", t)
	e1 := entries.Pick(uint(len(listed)))
	expectTrue(!e1.IsValid(), "e1.IsValid()", t)
}

func TestPickBuildUnitFromScrachDirr(t *testing.T) {
	entry := PickBuildUnitFromScrachDir("./testroot/2014/01042256-c/")
	expect(entry.PrimaryFile(), "testroot/2014/01042256-c/foo.c", t)

	shouldBeEmpty := PickBuildUnitFromScrachDir("./testroot/")
	expectTrue(!shouldBeEmpty.IsValid(), "shouldBeEmpty", t)
}

func TestMakeEditorCommandArgs(t *testing.T) {
	args1 := makeEditorCommandArgs("vi", "foo")
	expectTrue(2 == len(args1), "len(args1)", t)
	expect(args1[0], "vi", t)
	expect(args1[1], "foo", t)

	args2 := makeEditorCommandArgs("emacs -nw", "bar")
	expectTrue(3 == len(args2), "len(args2)", t)
	expect(args2[0], "emacs", t)
	expect(args2[1], "-nw", t)
	expect(args2[2], "bar", t)
}

func TestFindLanguageSuffix(t *testing.T) {
	expect(findLanguageSuffix("testroot/2014/01042256-c"), "c", t)
	expect(findLanguageSuffix("testroot/2014/01042256-go"), "go", t)
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
