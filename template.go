package skltn

import (
	"bytes"
	"text/template"
)

func renderTemplate(method *Method, tplStr string) ([]byte, error) {
	tmpl, err := template.New("unittest").Parse(tplStr)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, method); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
