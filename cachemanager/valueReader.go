package cachemanager

import (
	"compress/gzip"
	"encoding/json"
	"io"

	l4g "github.com/alecthomas/log4go"

	rankManager "github.com/Ekram-B2/rankmanager/rank"
)

// byteDecoder defines the operation for directly reading from the cache
type byteDecoder func(io.Reader) (rankManager.Rank, error)

// getByteDecoder is a factory that returns the byte reader operation to apply from reading a
// stream from a cache service
func getByteDecoder(compressorType string) byteDecoder {
	switch compressorType {
	case "gzip":
		return decodeGzip
	default:
		return decodeRaw
	}
}

// decodeGzip is applied to decode the raw stream using the gzip algorithm
func decodeGzip(reader io.Reader) (rankManager.Rank, error) {

	decompressedValueReader, err := gzip.NewReader(reader)
	if err != nil {
		l4g.Error("unable to decompress representation backing reader drawn from cache: %s", err.Error())
		return rankManager.Rank{}, err
	}

	rank := rankManager.Rank{}
	err = json.NewDecoder(decompressedValueReader).Decode(&rank)
	if err != nil {
		l4g.Error("unable to decode reader to go structure: %s", err.Error())
		return rankManager.Rank{}, err
	}
	// 3. Close the gzip.NewReader given what's it stated as a caller responsibility in its description
	decompressedValueReader.Close()
	return rank, nil

}

// decodeRaw is applied to return a value from the cache given a cache key
func decodeRaw(reader io.Reader) (rankManager.Rank, error) {
	// 1. Check cache to see if value exists given the cache key

	rank := rankManager.Rank{}
	err := json.NewDecoder(reader).Decode(&rank)
	if err != nil {
		l4g.Error(err)
		return rankManager.Rank{}, err
	}
	// 2. Return failure case if the software was unable to pipe data into byte slice
	if err != nil {
		l4g.Error("unable to pipe representation backing reader into a byte slice: %s", err.Error())
		return rankManager.Rank{}, err
	}
	return rank, nil
}
