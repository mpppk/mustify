# mustify
![GitHub Actions](https://github.com/mpppk/mustify/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/mpppk/mustify/branch/master/graph/badge.svg)](https://codecov.io/gh/mpppk/mustify)
[![GoDoc](https://godoc.org/github.com/mpppk/mustify?status.svg)](https://godoc.org/github.com/mpppk/mustify)

mustify is CLI tool for generate MustXXX methods from go source which includes methods that return error.

In Go, there is a culture where a method that occur a panic instead of returning an error has `Must` as a prefix.  
For example, `regexp.MustCompile` is the wrapper for `regexp.Compile` which occure panic if `regexp.Compile` return error.
mustify generate the wrapper automatically from all functions which return error.

## Usage
Assume you have below functions in `lib/math.go`

```go
package lib

import (
	"errors"
	"math"
)

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, errors.New("invalid value")
	}
	return math.Sqrt(x), nil
}

func SqrtWithoutError(x float64) float64 {
	return math.Sqrt(x)
}

func SumAndSub(v1, v2 int) (int, int, error) {
	return v1 + v2, v1 - v2, nil
}

func unexportedFunc() error {
	return nil
}
```

Then execute mustify:

```shell script
$ mustify < lib/math.go > lib/must-math.go
$ cat lib/must-math.go
package lib

func MustSqrt(x float64) float64 {
	_v0, _err := Sqrt(x)
	if _err != nil {
		panic(_err)
	}
	return _v0
}
func MustSumAndSub(v1, v2 int) (int, int) {
	_v0, _v1, _err := SumAndSub(v1, v2)
	if _err != nil {
		panic(_err)
	}
	return _v0, _v1
}
```
