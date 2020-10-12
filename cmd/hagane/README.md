# hagane

hagane is a template base code generator.

```sh
$ hagane -template template.go.tmpl -o sample_mock.go -data '{"type":"DB"}' sample.go
```

* `-o`: output file path (default stdout)
* `-f`: template format (default "{{.}}")
* `-template`: template file (data use `-f` option)
* `-data`: extra data as JSON format

See [the example](../../_examples/hagane/).
