package yuriko

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestResult(t *testing.T) {
	a := 10
	res := Ok[int](a)
	assert.Equal(t, res.Unwrap(), 10)
}

func TestUnwrapErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("not panic")
		}
		t.Log("panic!")
	}()

	res := Err[int](errors.New("faq"))
	res.Unwrap()
}

func TestMap(t *testing.T) {
	s := "12345"
	res := Ok[string](s)
	r2 := Map[string, int](res, func(s string) int {
		i, _ := strconv.Atoi(s)
		return i
	})
	assert.Equal(t, r2, Ok[int](12345))
}

type Foo struct {
	Dd int
}

func TestParse(t *testing.T) {
	s := `{"Dd": 123}`
	res, _ := Parse[map[string]int](s)
	assert.Equal(t, res, map[string]int{
		"Dd": 123,
	})
	res2, _ := Parse[Foo](s)
	assert.Equal(t, res2, Foo{
		Dd: 123,
	})
}

func TestFlatMap(t *testing.T) {
	res := Ok[string](`{"Dd": 123}`)
	r2 := AndThen[string, Foo](res, ParseResult[Foo])
	assert.Equal(t, r2.Unwrap(), Foo{Dd: 123})
}
