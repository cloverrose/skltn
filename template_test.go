package skltn

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const tmpl = `{{.ReceiverTyp}}
{{.MethodName}}
{{range $arg := .Args}}{{$arg.Name}} -> {{$arg.Typ}}{{end}}
{{range $ret := .Returns}}{{$ret}}{{end}}`

func TestTemplate(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   *Method
		exp  string
	}{
		{
			name: "base case",
			exp: "db\nUpdate\nx -> string\nerror",
			in: &Method{
				ReceiverTyp: "db",
				MethodName:  "Update",
				Args: []Arg{
					{
						Name: "x",
						Typ:  "string",
					},
				},
				Returns: []string{
					"error",
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			b, err := renderTemplate(tt.in, tmpl)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			if diff := cmp.Diff(tt.exp, string(b)); diff != "" {
				t.Errorf("got different output than what was expected:\n%s", diff)
			}
		})
	}
}
