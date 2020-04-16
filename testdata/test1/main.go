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
