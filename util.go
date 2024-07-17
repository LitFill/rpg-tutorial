package main

func clamp[N ~int | ~float64](value, min, max N) N {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
