package skltn

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parser_parseF(t *testing.T) {
	tests := []struct {
		name    string
		src string
		want    *Method
		wantErr bool
	}{
		{
			name: "simple",
			src: "func (d db) Save() error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Save",
				Args: []Arg{},
				Returns: []string{"error"},
			},
		},
		{
			name: "pointer receiver",
			src: "func (d *db) Save() error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Save",
				Args: []Arg{},
				Returns: []string{"error"},
			},
		},
		{
			name: "no arg",
			src: "func (d db) Ping() error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Ping",
				Args: []Arg{},
				Returns: []string{"error"},
			},
		},
		{
			name: "single arg single return",
			src: "func (d *db) Save(x int) error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Save",
				Args: []Arg{{"x", "int"}},
				Returns: []string{"error"},
			},
		},
		{
			name: "multiple args single return",
			src: "func (d *db) Save(x int, y value) error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Save",
				Args: []Arg{{"x", "int"},{"y", "value"}},
				Returns: []string{"error"},
			},
		},
		{
			name: "single arg multiple return",
			src: "func (d *db) Get(x int) (int, string) {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Get",
				Args: []Arg{{"x", "int"}},
				Returns: []string{"int", "string"},
			},
		},
		{
			name: "pointer arg type",
			src: "func (d db) Get(x *string) error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Get",
				Args: []Arg{{"x", "*string"}},
				Returns: []string{"error"},
			},
		},
		{
			name: "long arg type",
			src: "func (d db) Get(x a.B) error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Get",
				Args: []Arg{{"x", "a.B"}},
				Returns: []string{"error"},
			},
		},
		{
			name: "array arg type",
			src: "func (d *db) Get(x []string) error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Get",
				Args: []Arg{{"x", "[]string"}},
				Returns: []string{"error"},
			},
		},
		{
			name: "long pointer arg type",
			src: "func (d db) Get(x *a.B) error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Get",
				Args: []Arg{{"x", "*a.B"}},
				Returns: []string{"error"},
			},
		},
		{
			name: "long array arg type",
			src: "func (d *db) Update(x []a.B) error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{{"x", "[]a.B"}},
				Returns: []string{"error"},
			},
		},
		{
			name: "pointer array arg type",
			src: "func (a db) Update(x *[]a.B) error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{{"x", "*[]a.B"}},
				Returns: []string{"error"},
			},
		},
		{
			name: "array pointer arg type",
			src: "func (a db) Update(x []*a.B) error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{{"x", "[]*a.B"}},
				Returns: []string{"error"},
			},
		},
		{
			name: "has varargs",
			src: "func (d db) Update(x a.B, y int, z ...string) error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{{"x", "a.B"}, {"y", "int"}, {"z", "...string"}},
				Returns: []string{"error"},
			},
		},
		{
			name: "pointer return type",
			src: "func (a db) Update() *error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{},
				Returns: []string{"*error"},
			},
		},
		{
			name: "array return type",
			src: "func (a db) Update() []error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{},
				Returns: []string{"[]error"},
			},
		},
		{
			name: "pointer array return type",
			src: "func (a db) Update() *[]error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{},
				Returns: []string{"*[]error"},
			},
		},
		{
			name: "array pointer return type",
			src: "func (a db) Update() []*error {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{},
				Returns: []string{"[]*error"},
			},
		},
		{
			name: "long array pointer return type",
			src: "func (a db) Update() []*a.B {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{},
				Returns: []string{"[]*a.B"},
			},
		},
		{
			name: "complex return type 1",
			src: "func (a db) Update() (a.B, []error) {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{},
				Returns: []string{"a.B", "[]error"},
			},
		},
		{
			name: "complex return type 2",
			src: "func (a db) Update() (a.B, *int, error) {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{},
				Returns: []string{"a.B", "*int", "error"},
			},
		},
		{
			name: "complex return type 3",
			src: "func (a db) Update() ([]*a.B, *[]int, error) {",
			want: &Method{
				ReceiverTyp: "db",
				MethodName: "Update",
				Args: []Arg{},
				Returns: []string{"[]*a.B", "*[]int", "error"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newParser([]byte(tt.src))
			got, err := p.parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("parseF() got = %v, want %v, diff=%v", got, tt.want, diff)
			}
		})
	}
}


