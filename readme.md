# mustify
[![CircleCI](https://circleci.com/gh/mpppk/mustify.svg?style=svg)](https://circleci.com/gh/mpppk/mustify)
[![Build status](https://ci.appveyor.com/api/projects/status/qv1fyq6fm8ni4cne?svg=true)](https://ci.appveyor.com/project/mpppk/mustify)
![GitHub Actions](https://github.com/mpppk/mustify/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/mpppk/mustify/branch/master/graph/badge.svg)](https://codecov.io/gh/mpppk/mustify)
[![GoDoc](https://godoc.org/github.com/mpppk/mustify?status.svg)](https://godoc.org/github.com/mpppk/mustify)

mustify is CLI tool for generate MustXXX methods from go source which includes methods that return error.

In Go, there is a culture where a method that occur a panic instead of returning an error has Must as a prefix.  
For example, `regexp.MustCompile` is the wrapper for `regexp.Compile` which occure panic if `regexp.Compile` return error.
mustify generate the wrapper automatically from all functions which return error.

## Getting Started
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

func noErrorReturnFunc() int {
    return 42
}

func unexportedFunc() error {
	return nil
}
```

Then execute mustify:

```shell script
$ mustify lib | xargs goimports -w
$ cat lib/must-math.go
package lib

func MustSqrt(x float64) float64 {
	v, err := Sqrt(x)
	if err != nil {
		panic(err)
	}
	return v
}
```
