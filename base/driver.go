package core

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

type TestingScratchWriter struct {
	writtenFiles map[string]string
}

func (self *TestingScratchWriter) WriteFile(filename, content string) {
	if content == "" {
		log.Panic("Empty content is given for %s", filename)
	}

	self.writtenFiles[filename] = content
}

func (self *TestingScratchWriter) LastError() error {
	return nil
}

func (self *TestingScratchWriter) PrimaryFileName() string {
	return ""
}

func (self *TestingScratchWriter) IsWritten(filename string) bool {
	_, ok := self.writtenFiles[filename]
	return ok
}

func WriteStub(writer ScratchWriter, generator StubGenerator) error {
	return generator.Generate(writer)
}

func MakeTestingScratchWriter() *TestingScratchWriter {
	return &TestingScratchWriter{
		map[string]string{},
	}
}

type FileScratchWriter struct {
	baseDir         string
	primaryFileName string
	errors          []error
}

func (self *FileScratchWriter) WriteFile(filename, content string) {
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

func (self *FileScratchWriter) LastError() error {
	if 0 == len(self.errors) {
		return nil
	}

	return self.errors[len(self.errors)-1]
}

func (self *FileScratchWriter) PrimaryFileName() string {
	return self.primaryFileName
}

func MakeFileScratchWriter(baseDir string) (*FileScratchWriter, error) {
	return &FileScratchWriter{
		baseDir: baseDir,
		errors:  []error{},
	}, nil
}
