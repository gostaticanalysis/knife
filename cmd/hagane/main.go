package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"

	"github.com/gostaticanalysis/knife"
)

var (
	flagOut       string
	flagFormat    string
	flagTemplate  string
	flagExtraData string
)

func init() {
	flag.StringVar(&flagOut, "o", "", "output file path")
	flag.StringVar(&flagFormat, "f", "{{.}}", "output format")
	flag.StringVar(&flagTemplate, "template", "", "template file")
	flag.StringVar(&flagExtraData, "data", "", "extra data as JSON format")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (rerr error) {
	knifeOpt := &knife.KnifeOption{Tests: true}
	k, err := knife.New(knifeOpt, flag.Args()[1:]...)
	if err != nil {
		return fmt.Errorf("cannot create knife: %w", err)
	}

	var opt knife.ExecuteOption
	if flagExtraData != "" {
		err := json.Unmarshal([]byte(flagExtraData), &opt.ExtraData)
		if err != nil {
			return fmt.Errorf("cannot parse JSON: %w", err)
		}
	}

	pkgs := k.Packages()
	if len(pkgs) == 0 {
		return errors.New("does not find package")
	}

	var tmpl any = flagFormat
	if flagTemplate != "" {
		tmpl, err = os.ReadFile(flagTemplate)
		if err != nil {
			return fmt.Errorf("cannot read template: %w", err)
		}
	}

	pkg := pkgs[0]
	var buf bytes.Buffer
	if err := k.Execute(&buf, pkg, tmpl, &opt); err != nil {
		return fmt.Errorf("cannot knife execute: %w", err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("cannot format: %w", err)
	}

	if len(bytes.TrimSpace(src)) == 0 {
		return nil
	}

	var w io.Writer = os.Stdout
	if flagOut != "" {
		f, err := os.Create(flagOut)
		if err != nil {
			return fmt.Errorf("cannot create file: %w", err)
		}
		defer func() {
			if err := f.Close(); err != nil && rerr == nil {
				rerr = err
			}
		}()
		w = f
	}

	if _, err := fmt.Fprintln(w, string(src)); err != nil {
		return fmt.Errorf("cannot output source: %w", err)
	}

	return nil
}
