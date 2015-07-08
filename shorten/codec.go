package shorten

import (
	"math"
	"strings"
)

const (
	// symbols are all characters used for short-urls
	symbols = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base    = int64(len(symbols))
)

// encode converts a number into our *base* representation
func encode(number int64) string {
	rest := number % base
	// strings are a bit weird in go...
	result := string(symbols[rest])
	if number-rest != 0 {
		newnumber := (number - rest) / base
		result = encode(newnumber) + result
	}
	return result
}

// decode takes a string in our encoding and returns the decimal integer.
func decode(input string) int64 {
	const floatbase = float64(base)
	l := len(input)
	var sum = 0
	for index := l - 1; index > -1; index-- {
		current := string(input[index])
		pos := strings.Index(symbols, current)
		sum = sum + (pos * int(math.Pow(floatbase, float64((l-index-1)))))
	}
	return int64(sum)
}
