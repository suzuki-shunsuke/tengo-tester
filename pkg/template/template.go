package template

import (
	"bytes"
	"text/template"
)

type Template struct {
	Template *template.Template
	Text     string
}

func Compile(tpl string) (Template, error) {
	tmpl, err := template.New("text").Parse(tpl)
	if err != nil {
		return Template{}, err
	}
	return Template{
		Template: tmpl,
		Text:     tpl,
	}, nil
}

func (tpl Template) Render(params interface{}) (string, error) {
	if tpl.Template == nil {
		return "", nil
	}
	buf := &bytes.Buffer{}
	if err := tpl.Template.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}
