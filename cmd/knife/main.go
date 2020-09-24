package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gostaticanalysis/knife"
)

var (
	flagFormat    string
	flagTemplate  string
	flagExtraData string
	flagXPath     string
)

func init() {
	flag.StringVar(&flagFormat, "f", "{{.}}", "output format")
	flag.StringVar(&flagTemplate, "template", "", "template file")
	flag.StringVar(&flagExtraData, "data", "", "extra data (key:value,key:value)")
	flag.StringVar(&flagXPath, "xpath", "", "A XPath expression for an AST node")
	flag.Parse()
}

func main() {
	if err := run(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	k, err := knife.New(args...)
	if err != nil {
		return err
	}

	var w io.Writer = os.Stdout

	opt := &knife.FormatOption{
		XPath: flagXPath,
	}

	if flagExtraData != "" {
		extraData, err := parseExtraData(flagExtraData)
		if err != nil {
			return err
		}
		opt.ExtraData = extraData
	}

	pkgs := k.Packages()
	for i, pkg := range pkgs {
		switch {
		case flagTemplate != "":
			if err := k.FormatWithTemplate(w, pkg, flagTemplate, opt); err != nil {
				return err
			}
		default:
			if err := k.Format(w, pkg, flagFormat, opt); err != nil {
				return err
			}
		}

		if i != len(pkgs)-1 {
			fmt.Fprintln(w)
		}
	}

	return nil
}

func parseExtraData(extraData string) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	kvs := strings.Split(extraData, ",")
	for i := range kvs {
		kv := strings.Split(strings.TrimSpace(kvs[i]), ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid extra data: %s", kvs[i])
		}
		m[kv[0]] = kv[1]
	}
	return m, nil
}
