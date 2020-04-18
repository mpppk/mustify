package lib

import (
	"errors"
	"math"
)

// Sqrt returns the square root of x
func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, errors.New("invalid value")
	}
	return math.Sqrt(x), nil
}

// Sqrt returns the square root of x without error
func SqrtWithoutError(x float64) float64 {
	return math.Sqrt(x)
}

// SumAndSub returns sum and sub of arguments
func SumAndSub(v1, v2 int) (int, int, error) {
	return v1 + v2, v1 - v2, nil
}

func unexportedFunc() error {
	return nil
}
