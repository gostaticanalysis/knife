package knife

import (
	"go/types"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	want := strings.TrimSpace(version)
	got := Version()
	if want != got {
		t.Errorf("knife.Version does not match: (got,wan) = (%q, %q)", got, want)
	}
}

func TestTemplateFunctions(t *testing.T) {
	cases := []struct {
		name     string
		template string
		want     string
	}{
		{
			name:     "regex match",
			template: `{{regexp "^[A-Z]" "Test"}}`,
			want:     "true",
		},
		{
			name:     "regex no match",
			template: `{{regexp "^[0-9]+$" "abc"}}`,
			want:     "false",
		},
		{
			name:     "len function with slice",
			template: `{{len .Types}}`,
			want:     "0",
		},
		{
			name:     "len function with string",
			template: `{{len "hello"}}`,
			want:     "5",
		},
		{
			name:     "last function with string",
			template: `{{last "hello"}}`,
			want:     "o",
		},
		{
			name:     "data function",
			template: `{{data "testkey"}}`,
			want:     "testvalue",
		},
		{
			name:     "br function",
			template: `line1{{br}}line2`,
			want:     "line1\nline2",
		},
		{
			name:     "pkg function",
			template: `{{.Name}}`,
			want:     "test",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			td := &TempalteData{
				Pkg: types.NewPackage("test", "test"),
				Extra: map[string]any{
					"testkey": "testvalue",
				},
			}
			tmpl := NewTemplate(td)

			parsed, err := tmpl.Parse(tt.template)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			var buf strings.Builder
			pkg := NewPackage(td.Pkg)
			err = parsed.Execute(&buf, pkg)
			if err != nil {
				t.Fatalf("failed to execute template: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			if got != tt.want {
				t.Errorf("template execution result = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGodocFunction(t *testing.T) {
	cases := []struct {
		name     string
		template string
		contains string
	}{
		{
			name:     "godoc fmt package",
			template: `{{godoc "fmt"}}`,
			contains: "Package fmt implements formatted I/O",
		},
		{
			name:     "godoc fmt.Println function",
			template: `{{godoc "fmt.Println"}}`,
			contains: "func Println",
		},
		{
			name:     "godoc with -src flag",
			template: `{{godoc "-src" "fmt.Println"}}`,
			contains: "func Println",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			td := &TempalteData{
				Pkg: types.NewPackage("test", "test"),
			}
			tmpl := NewTemplate(td)

			parsed, err := tmpl.Parse(tt.template)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			var buf strings.Builder
			pkg := NewPackage(td.Pkg)
			err = parsed.Execute(&buf, pkg)
			if err != nil {
				t.Fatalf("failed to execute template: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			if !strings.Contains(got, tt.contains) {
				t.Errorf("template execution result = %q, want to contain %q", got, tt.contains)
			}
		})
	}
}

func TestTypeConversionFunctions(t *testing.T) {
	cases := []struct {
		name     string
		template string
		want     string
	}{
		{
			name:     "basic type conversion",
			template: `{{with basic (typeof "int")}}{{.Name}}{{end}}`,
			want:     "int",
		},
		{
			name:     "slice type conversion returns nil for non-slice",
			template: `{{$slice := slice (typeof "int")}}{{if $slice}}has slice{{else}}no slice{{end}}`,
			want:     "no slice",
		},
		{
			name:     "under function",
			template: `{{$t := typeof "int"}}{{with $t}}{{$u := under .TypesType}}{{if $u}}has underlying{{else}}no underlying{{end}}{{end}}`,
			want:     "has underlying",
		},
		{
			name:     "identical function",
			template: `{{identical (typeof "int") (typeof "int")}}`,
			want:     "true",
		},
		{
			name:     "identical function false",
			template: `{{identical (typeof "int") (typeof "string")}}`,
			want:     "false",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			td := &TempalteData{
				Pkg: types.NewPackage("test", "test"),
			}
			tmpl := NewTemplate(td)

			parsed, err := tmpl.Parse(tt.template)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			var buf strings.Builder
			pkg := NewPackage(td.Pkg)
			err = parsed.Execute(&buf, pkg)
			if err != nil {
				t.Fatalf("failed to execute template: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			if got != tt.want {
				t.Errorf("template execution result = %q, want %q", got, tt.want)
			}
		})
	}
}
