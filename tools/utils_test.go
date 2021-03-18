package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildInsertValuePlaceHolder(t *testing.T) {
	str := BuildInsertValuePlaceHolder(2, 3)
	exp := "(?, ?), (?, ?), (?, ?)"
	assert.Equal(t, exp, str)
}

func TestBuildInsertValuePlaceHolderV1(t *testing.T) {
	str := buildInsertValuePlaceHolderV1(2, 3)
	exp := "(?, ?), (?, ?), (?, ?)"
	assert.Equal(t, exp, str)
}
func TestBuildInsertValuePlaceHolderV2(t *testing.T) {
	str := buildInsertValuePlaceHolderV2(2, 3)
	exp := "(?, ?), (?, ?), (?, ?)"
	assert.Equal(t, exp, str)
}
func TestBuildInsertValuePlaceHolderV3(t *testing.T) {
	str := buildInsertValuePlaceHolderV3(3, 3)
	exp := "(?, ?, ?), (?, ?, ?), (?, ?, ?)"
	assert.Equal(t, exp, str)

	str = buildInsertValuePlaceHolderV3(3, 4)
	exp = "(?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?)"
	assert.Equal(t, exp, str)
}

func TestBuildInsertValuePlaceHolderV4(t *testing.T) {
	str := buildInsertValuePlaceHolderV4(3, 3)
	exp := "(?, ?, ?), (?, ?, ?), (?, ?, ?)"
	assert.Equal(t, exp, str)

	str = buildInsertValuePlaceHolderV4(3, 4)
	exp = "(?, ?, ?), (?, ?, ?), (?, ?, ?), (?, ?, ?)"
	assert.Equal(t, exp, str)
}

func BenchmarkBuildInsertValuePlaceHolder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildInsertValuePlaceHolder(200, 1000)
	}
}

func BenchmarkBuildInsertValuePlaceHolderV1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildInsertValuePlaceHolderV1(200, 1000)
	}
}

func BenchmarkBuildInsertValuePlaceHolderV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildInsertValuePlaceHolderV2(200, 1000)
	}
}

func BenchmarkBuildInsertValuePlaceHolderV3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildInsertValuePlaceHolderV3(200, 1000)
	}
}

func BenchmarkBuildInsertValuePlaceHolderV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildInsertValuePlaceHolderV4(200, 1000)
	}
}
