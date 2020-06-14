package cachemanager

import (
	"bytes"
	"encoding/json"
	"testing"

	rankManager "github.com/Ekram-B2/rankmanager/rank"
)

func Test_decodeRaw(t *testing.T) {
	// 1. Define input and expected output to compare against (Arrange)
	inputRank := rankManager.Rank{Name: "Toronto", Rank: 0.5}
	bytesStream, _ := json.Marshal(inputRank)
	inputReader := bytes.NewBuffer(bytesStream)
	// 2. Apply operation to get output (Act)
	actualRank, err := decodeRaw(inputReader)
	// 3. Check to see if the actual matches the expected (Assert)
	if err != nil {
		t.Fatalf("the operation failed to produce a stream given valid input")
	}
	if actualRank.Name != inputRank.Name {
		t.Fatalf("expected did not match actual; expected was %v and actual was %v", actualRank.Name, inputRank.Name)
	}
	if actualRank.Rank != inputRank.Rank {
		t.Fatalf("expected did not match actual; expected was %v and actual was %v", actualRank.Rank, inputRank.Rank)
	}

}

func Test_decodeGzip(t *testing.T) {
	// 1. Define the expected output and input args
	inputRank := rankManager.Rank{Name: "Toronto", Rank: 0.5}
	valueStream := []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 170, 86, 202, 75, 204,
		77, 85, 178, 82, 10, 201, 47, 202, 207, 43, 201, 87, 210, 81, 42, 74, 204, 203, 86,
		178, 50, 208, 51, 173, 5, 4, 0, 0, 255, 255, 205, 209, 249, 51, 29, 0, 0, 0}
	inputReader := bytes.NewBuffer(valueStream)
	// 2. Apply operation to get output (Act)
	actualRank, err := decodeGzip(inputReader)
	// 3. Check to see if the actual matches the expected (Assert)
	if err != nil {
		t.Fatalf("failed to produce stream given valid output")
	}
	if actualRank.Name != inputRank.Name {
		t.Fatalf("expected did not match actual; expected was %v and actual was %v", actualRank.Name, inputRank.Name)
	}
	if actualRank.Rank != inputRank.Rank {
		t.Fatalf("expected did not match actual; expected was %v and actual was %v", actualRank.Rank, inputRank.Rank)
	}

}
