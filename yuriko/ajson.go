package yuriko

import "encoding/json"

func Parse[T any](s string) (res T, err error) {
	if err = json.Unmarshal([]byte(s), &res); err != nil {
		return
	}
	return res, nil
}

func ParseResult[T any](s string) Result[T] {
	res, err := Parse[T](s)
	if err != nil {
		return Err[T](err)
	}
	return Ok(res)
}
