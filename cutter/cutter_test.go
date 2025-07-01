package cutter_test

import (
	"testing"

	"github.com/gostaticanalysis/knife/cutter"
)

func TestNew(t *testing.T) {
	cutterOpt := &cutter.CutterOption{Tests: true}
	c, err := cutter.New(cutterOpt, "fmt")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if got := c.KnifePackages(); len(got) == 0 {
		t.Error("cutter.New must creates knife.Package")
	}
}
