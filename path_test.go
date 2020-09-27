package knife_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gostaticanalysis/knife"
)

func ExampleAt() {
	type Bar struct {
		N []int
	}

	type Foo struct {
		Bar *Bar
	}

	f := &Foo{
		Bar: &Bar{
			N: []int{100},
		},
	}

	v, err := knife.At(f, `Bar.N[0]`)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output:
	// 100
}

func ExamplePath_Eval() {
	type Bar struct {
		N []int
	}

	type Foo struct {
		Bar *Bar
	}

	f := &Foo{
		Bar: &Bar{
			N: []int{100},
		},
	}

	p, err := knife.NewPath(`Bar.N[0]`)
	if err != nil {
		panic(err)
	}

	var v int
	if err := p.Eval(f, &v); err != nil {
		panic(err)
	}

	fmt.Println(v)
	// Output:
	// 100
}

func TestAt(t *testing.T) {
	t.Parallel()

	type C struct{ N int }
	type B struct{ C *C }
	type A struct{ B *B }

	cases := []struct {
		v      interface{}
		expr   string
		want   interface{}
		hasErr bool
	}{
		{struct{ A int }{100}, "A", 100, false},
		{struct{ A []int }{[]int{100, 200}}, "A[1]", 200, false},
		{struct{ A map[string][]int }{map[string][]int{"foo": {100, 200}}}, `A["foo"][1]`, 200, false},
		{struct{ A map[int][]int }{map[int][]int{200: {100, 200}}}, `A[200][1]`, 200, false},
		{struct{ A map[string]int }{map[string]int{"B": 100}}, `A.B`, 100, false},
		{struct{ A struct{ B int } }{struct{ B int }{100}}, `A.B`, 100, false},
		{struct{ A struct{ B int } }{struct{ B int }{100}}, `A.C`, nil, true},
		{&A{&B{&C{100}}}, `B.C`, &C{100}, false},
		{struct{ N []int }{[]int{100}}, `N[0]`, 100, false},
		{nil, `import "fmt"`, nil, true},
		{nil, `Call()`, nil, true},
		{struct{ N map[int]int }{map[int]int{-1: 100}}, `N[-1]`, 100, false},
		{struct{ N map[int]int }{map[int]int{0: 100}}, `N[1-1]`, 100, false},
		{struct{ N map[int]int }{map[int]int{0: 100}}, `N[(0)]`, 100, false},
		{struct{ N map[bool]int }{map[bool]int{true: 100}}, `N[true]`, 100, false},
		{struct{ N map[bool]int }{map[bool]int{false: 100}}, `N[false]`, 100, false},
		{struct{ N map[bool]int }{map[bool]int{true: 100}}, `N[100 > 0]`, 100, false},
		{struct{ N map[bool]int }{map[bool]int{true: 100}}, `N[true]`, 100, false},
		{struct{ N map[bool]int }{map[bool]int{true: 100}}, `N[1 + f()]`, nil, true},
		{struct{ N map[bool]int }{map[bool]int{true: 100}}, `N[T]`, nil, true},
		{struct{ N map[bool]int }{map[bool]int{true: 100}}, `N[-T]`, nil, true},
		{struct{ N map[bool]int }{map[bool]int{true: 100}}, `N[T - 1]`, nil, true},
		{struct{ N map[bool]int }{map[bool]int{true: 100}}, `N["key" + 1]`, nil, true},
		{struct{ N map[int]int }{map[int]int{10: 100}}, `N[99999999999999999999999999999]`, nil, true},
		{struct{ N []int }{[]int{100}}, `N[99999999999999999999999999999]`, nil, true},
		{struct{ N map[float64]int }{map[float64]int{1.5: 100}}, `N[1.5]`, 100, false},
		{struct{ N map[float64]int }{map[float64]int{1.5: 100}}, `N[99999999999999999999999999999999999.0]`, nil, true},
		{struct{ N map[int]int }{map[int]int{100: 100}}, `N[0]`, nil, false},
		{100, `N[0]`, nil, true},
		{100, `A.B[0]`, nil, true},
		{struct{ N int }{100}, `N[0]`, nil, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.expr, func(t *testing.T) {
			t.Parallel()
			got, err := knife.At(tt.v, tt.expr)
			switch {
			case tt.hasErr && err == nil:
				t.Fatal("expected error has not occur")
			case !tt.hasErr && err != nil:
				t.Fatal("unexpected error:", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
