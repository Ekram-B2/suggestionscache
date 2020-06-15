package cachemanager

import (
	"testing"

	rankManager "github.com/Ekram-B2/rankmanager/rank"
)

func Test_encodeRaw(t *testing.T) {
	// 1. Define input and expected output to compare against (Arrange)
	bytesStream := []byte{123, 34, 110, 97, 109, 101, 34,
		58, 34, 84, 111, 114, 111, 110, 116, 111, 34, 44, 34,
		114, 97, 110, 107, 34, 58, 48, 46, 53, 125}

	inputRank := rankManager.Rank{Name: "Toronto", Rank: 0.5}
	// 2. Apply operation to get output (Act)
	actualStream, err := encodeRaw(inputRank)
	if err != nil {
		t.Fatalf("the operation failed to produce a stream given valid input")
	}
	if len(actualStream) != len(bytesStream) {
		t.Fatalf("expected did not match actual; expected was %v and actual was %v", bytesStream, actualStream)
	}
	if !isSameSlice(actualStream, bytesStream) {
		t.Fatalf("expected did not match actual; expected was %v and actual was %v", bytesStream, actualStream)
	}
}

func Test_encodeGzip(t *testing.T) {
	// 1. Define input and expected output to compare against (Arrange)
	expectedBytesStream := []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 170, 86, 202, 75, 204,
		77, 85, 178, 82, 10, 201, 47, 202, 207, 43, 201, 87, 210, 81, 42, 74, 204, 203, 86,
		178, 50, 208, 51, 173, 5, 4, 0, 0, 255, 255, 205, 209, 249, 51, 29, 0, 0, 0}

	inputRank := rankManager.Rank{Name: "Toronto", Rank: 0.5}
	// 2. Apply input to get output of operation (Act)
	actualStream, err := encodeGzip(inputRank)(inputRank)
	// 3. Check to see if the expected matches returned (Assert)
	if err != nil {
		t.Fatalf("an err occured instead with valid input")
	}
	if len(actualStream) != len(expectedBytesStream) {
		t.Fatalf("expected did not match actual; expected is %v and actual is %v", expectedBytesStream, actualStream)
	}
	if !isSameSlice(actualStream, expectedBytesStream) {
		t.Fatalf("expected did not match actual; expected is %v and actual is %v", expectedBytesStream, actualStream)
	}

}
