package langs

import (
	"bytes"
	base "github.com/omo/fuga/base"
	"text/template"
)

type TemplateMap map[string]string

func (self TemplateMap) WriteTo(writer base.StubWriter, name string) {
	writer.WriteFile(name, self[name])
}

func (self TemplateMap) WriteToWith(writer base.StubWriter, name string, data interface{}) {
	doc := &bytes.Buffer{}
	templateText := self[name]
	tmpl, err := template.New("cppTemplate").Parse(templateText)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(doc, data)
	if err != nil {
		panic(err)
	}

	writer.WriteFile(name, doc.String())
}

type MakefileRunner struct{}

func (*MakefileRunner) Run(params base.BuildRunnerParams) error {
	return runProgram("make", []string{}, params.Unit.Dir())
}
