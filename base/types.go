package core

import "time"

type Parameters struct {
	Workspace string
	Now       time.Time
	Suffix    string
}

type StubWriter interface {
	WriteFile(filename, content string)
	LastError() error
}

type StubGenerator interface {
	Generate(StubWriter) error
}
