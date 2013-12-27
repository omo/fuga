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
	now, err := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Dec 25, 2013 at 5:34pm (JST)")
	expectOK(err, t)
	return &Parameters{"/root", now, "java"}
}

func TestMakeBaseDir(t *testing.T) {
	actual := makeBaseDir(makeTestParameters(t))
	expect(actual, "/root/2013/12251734-java", t)
}

func TestCGenerator(t *testing.T) {
	wri := MakeTestingStubWriter()
	gen := makeGenerator("c")
	WriteStub(wri, gen)
	expectTrue(wri.IsWritten("foo.c"), "foo.c", t)
	expectTrue(wri.IsWritten("Makefile"), "Makefile", t)
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
