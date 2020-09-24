package a

//go:generate go run ../../main.go -o a_mock.go DB a.go

type DB interface {
	Get(id string) int
	Set(id string, v int)
}
