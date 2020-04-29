package skltn_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/cloverrose/skltn"
)

const testTemplate = `func (d *{{.ReceiverTyp}}) {{.MethodName}}({{range $idx, $arg := .Args}}{{if $idx}},{{end}}{{$arg.Name}} {{$arg.Typ}}{{end}}) ({{range $idx, $ret := .Returns}}{{if $idx}},{{end}}{{$ret}}{{end}}) {}`

func TestIntegration(t *testing.T) {
	got, err := skltn.Generate([]byte(`func (d *db) Update(x *a.D, y []a) (*a.DTO, error) {`), testTemplate, true)
	if err != nil {
		t.Fatalf("unexpected error '%s'", err.Error())
	}
	exp := `func (d *db) Update(x *a.D, y []a) (*a.DTO, error) {}`
	if diff := cmp.Diff(exp, string(got)); diff != "" {
		t.Errorf("unexpected unit test created:\n%s", diff)
	}
}
