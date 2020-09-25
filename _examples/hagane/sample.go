package sample

//go:generate hagane -template template.go.tmpl -o sample_mock.go -data {"type":"DB"} sample.go

type DB interface {
	Get(id string) int
	Set(id string, v int)
}

type db struct {}

func (db) Get(id string) int {
	return 0
}

func (db) Set(id string, v int) {}
