package account

import "math"

func roundFraction(numerator uint64, denominator uint64) uint64 {
	return (numerator + denominator/2) / denominator
}

func floatToUint64(x float64) uint64 {
	return uint64(math.Round(x))
}

func floatMulToUint64(a, b float64) uint64 {
	return floatToUint64(a * b)
}
