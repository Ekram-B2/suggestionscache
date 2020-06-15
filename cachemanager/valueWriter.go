package cachemanager

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	l4g "github.com/alecthomas/log4go"

	rankManager "github.com/Ekram-B2/rankmanager/rank"
)

// byteEncoder defines the operation to encode a rank into a byte stream
type byteEncoder func(rankManager.Rank) ([]byte, error)

// getByteEncoder is a factory applied to get the appropriate encoder algorithm
func getByteEncoder(compressorType string, rank rankManager.Rank) byteEncoder {
	switch compressorType {
	case "gzip":
		return encodeGzip(rank)
	default:
		return encodeRaw
	}
}

// encodeGzip is applied to return a value from the cache given a cache key
func encodeGzip(rank rankManager.Rank) byteEncoder {
	return func(rankManager.Rank) (compressedBytes []byte, err error) {
		// 1. Check cache to see if value exists given the cache key
		rawStream, err := encodeRaw(rank)
		if err != nil {
			l4g.Error("OPERATION-ERROR: unable to get uncompressed byte stream from the rank")
			return nil, err
		}

		// 2. Apply stream to back a writer interface
		var b bytes.Buffer
		writer := gzip.NewWriter(&b)

		// 3. Convert to a gzip writer to apply gzip decompression
		_, err = writer.Write(rawStream)
		if err != nil {
			l4g.Error("OPERATION-ERROR: unable write compressed bytes to backing byte slice")
			return nil, err
		}
		writer.Close()
		// 5. Return compressed byte stream
		return b.Bytes(), nil
	}
}

// encodeRaw is applied to return a value from the cache given a cache key
func encodeRaw(rank rankManager.Rank) ([]byte, error) {
	// 1. Check cache to see if value exists given the cache key
	valueStream, err := json.Marshal(rank)
	if err != nil {
		l4g.Error("OPERATION-ERROR: unable to marshall rank into a byte stream")
		return nil, err
	}
	return valueStream, nil
}
