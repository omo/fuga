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

func (self *TestingStubWriter) PrimaryFileName() string {
	return ""
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
	baseDir         string
	primaryFileName string
	errors          []error
}

func (self *FileStubWriter) WriteFile(filename, content string) {
	path := filepath.Join(self.baseDir, filename)
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if nil != err {
		self.errors = append(self.errors, err)
		return
	}

	if self.primaryFileName == "" {
		self.primaryFileName = path
	}
}

func (self *FileStubWriter) LastError() error {
	if 0 == len(self.errors) {
		return nil
	}

	return self.errors[len(self.errors)-1]
}

func (self *FileStubWriter) PrimaryFileName() string {
	return self.primaryFileName
}

func MakeFileStubWriter(baseDir string) (*FileStubWriter, error) {
	return &FileStubWriter{
		baseDir: baseDir,
		errors:  []error{},
	}, nil
}
