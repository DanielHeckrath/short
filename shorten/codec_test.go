package shorten

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var encodeData = []struct {
	encoded string
	decoded int64
}{
	{"1", 1},     // check first
	{"Z", 61},    // check last single char
	{"10", 62},   // check wrap around
	{"g8", 1000}, // a few other to be sure
	{"2Bi", 10000},
	{"q0U", 100000},
	{"4c92", 1000000},
	{"FXsk", 10000000},
	{"6LAze", 100000000},
	{"23", math.MaxInt8},
	{"8wv", math.MaxInt16},
	{"2lkCB1", math.MaxInt32},
	{"aZl8N0y58M7", math.MaxInt64},
}

func TestEncode(t *testing.T) {
	a := assert.New(t)

	for _, data := range encodeData {
		encoded := encode(data.decoded)
		a.Equal(data.encoded, encoded, "Encoded value does not match expectation")
	}
}

func TestDecode(t *testing.T) {
	a := assert.New(t)

	for _, data := range encodeData {
		decoded := decode(data.encoded)
		a.Equal(data.decoded, decoded, "Decoded value does not match expectation")
	}
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encode(int64(i))
	}
}

func BenchmarkDecode(b *testing.B) {
	count := len(encodeData)
	for i := 0; i < b.N; i++ {
		index := i % count
		decode(encodeData[index].encoded)
	}
}
