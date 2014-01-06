package langs

import (
	base "github.com/omo/fuga/base"
)

type TemplateMap map[string]string

func (self TemplateMap) WriteTo(writer base.StubWriter, name string) {
	writer.WriteFile(name, self[name])
}
