package cachemanager

import (
	"bytes"
	"time"

	l4g "github.com/alecthomas/log4go"
	rcache "github.com/patrickmn/go-cache"

	rankManager "github.com/Ekram-B2/rankmanager/rank"
	"github.com/Ekram-B2/suggestionsmanager/config"
)

type ramCacheManager struct {
	searchTerm    string
	realTerm      string
	searchTermLat string
	searchTermLng string
	realTermLat   string
	realTermLng   string
	config        config.Config
}

type ramCacheclient interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, d time.Duration)
}

var (
	rmclient ramCacheclient
)

func init() {
	rmclient = rcache.New(5*time.Minute, 10*time.Minute)
}

func (rm ramCacheManager) getBytesFromCache(key string) (isCacheHit bool, valuestream []byte) {
	// 1. Check cache to see if value exists given the cache key
	valueStream, isFound := rmclient.Get(key)
	// 2. Return a state based on whether there was a cache hit or cache miss
	if isFound {
		// a cache miss does not not mean an error
		return isFound, valueStream.([]byte)
	}
	return isFound, nil
}

// getFromCache is applied to return a value from the cache given a cache key
func (rm ramCacheManager) getRankFromCache(key, byteDecoderType string) (isCacheHit bool, returnedRank rankManager.Rank, err error) {

	// 1. Check cache to see if value exists given the cache key
	isCacheHit, valueStream := rm.getBytesFromCache(key)

	if !isCacheHit {
		// there was a cache miss - this is not considered an error
		return false, rankManager.Rank{}, nil
	}

	// 2. Convert valueStream to an io reader
	reader := bytes.NewBuffer(valueStream)
	// 3. Apply byte decoder to return a rank
	returnedRank, err = getByteDecoder(rm.config.ByteDecoderType)(reader)
	if err != nil {
		l4g.Error("OPERATION-ERROR: unable to get rank from the decoder: %s", err.Error())
		return true, rankManager.Rank{}, nil
	}
	return true, returnedRank, nil
}

func (rm ramCacheManager) putInCache(key string, value rankManager.Rank, byteEncoder byteEncoder) error {
	// 1. Convert rank into byte stream to put into cache
	valueStream, err := getByteEncoder(rm.config.ByteEncoderType, value)(value)
	if err != nil {
		l4g.Error("OPERATION-ERROR: unable to get byte stream from the encoder: %s", err.Error())
		return err
	}
	// 2. Apply client to put into the cache
	rmclient.Set(key, valueStream, rcache.DefaultExpiration)
	// 4. Return error if such exists
	return nil

}
