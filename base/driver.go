package core

import (
	"io/ioutil"
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

func MakeFileStubWriter(baseDir string) (*FileStubWriter, error) {
	return &FileStubWriter{
		baseDir,
		[]error{},
	}, nil
}
