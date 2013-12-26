package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TestingStubWriter struct {
	writtenFiles map[string]string
}

func (self *TestingStubWriter) WriteFile(filename, content string) {
	self.writtenFiles[filename] = content
}

func (self *TestingStubWriter) LastError() error {
	return nil
}

func (self *TestingStubWriter) IsWritten(filename string) bool {
	_, ok := self.writtenFiles[filename]
	return ok
}

func WriteStub(writer StubWriter, generator StubGenerator) error {
	return generator.Generate(writer)
}

func MakeTestingStubWriter() *TestingStubWriter {
	return &TestingStubWriter{
		map[string]string{},
	}
}

type FileStubWriter struct {
	baseDir string
	errors  []error
}

func (self *FileStubWriter) WriteFile(filename, content string) {
	err := ioutil.WriteFile(filepath.Join(self.baseDir, filename), []byte(content), 0644)
	if nil != err {
		self.errors = append(self.errors, err)
		return
	}
}

func (self *FileStubWriter) LastError() error {
	if 0 == len(self.errors) {
		return nil
	}

	return self.errors[len(self.errors)-1]
}

func EnsureDir(dirname string) error {
	stat, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return os.MkdirAll(dirname, 0777)
	}

	if !stat.IsDir() {
		return errors.New(fmt.Sprintf("File %s is already exist!", dirname))
	}

	return nil
}

func MakeFileStubWriter(baseDir string) (*FileStubWriter, error) {
	// FIXME: This should ensure the freshness.
	if err := EnsureDir(baseDir); err != nil {
		return nil, err
	}

	return &FileStubWriter{
		baseDir,
		[]error{},
	}, nil
}
