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

func TestRegexpInTemplate(t *testing.T) {
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
			err = parsed.Execute(&buf, nil)
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
