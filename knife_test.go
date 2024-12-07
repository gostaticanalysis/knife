package knife

import (
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
