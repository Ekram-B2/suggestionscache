package cachemanager

import (
	"strings"

	rankManager "github.com/Ekram-B2/rankmanager/rank"
)

// cacheManager enacts as a application to store/retreive from a cache, or commit a request to a rank manager
type cacheManager interface {
	// putInCache puts a key/value pairing into a cache store
	putInCache(string, rankManager.Rank, byteEncoder) error
	// getFromCache gets a rank from the cache if it exists in the cache, or an empty rank if no key/value pairing exists
	getRankFromCache(string, string) (bool, rankManager.Rank, error)
	// getRankFromCache returns a byte stream from the cache
	getBytesFromCache(string) (bool, []byte)
}

// GetCacheManager is a factory applied to return a cache manager
func GetCacheManager(realTerm, searchTerm, realTermLat, realTermLng, searchTermLat, searchTermLng, cacheType string) cacheManager {
	switch cacheType {
	case "ramcache":
		return ramCacheManager{realTerm: realTerm,
			searchTerm:    searchTerm,
			realTermLat:   realTermLat,
			realTermLng:   realTermLng,
			searchTermLat: searchTermLat,
			searchTermLng: searchTermLng,
		}
	default:

		return ramCacheManager{realTerm: realTerm,
			searchTerm:    searchTerm,
			realTermLat:   realTermLat,
			realTermLng:   realTermLng,
			searchTermLat: searchTermLat,
			searchTermLng: searchTermLng,
		}
	}
}

// generateCacheKeyDefault is applied to generate a hash that servces as the cache key
func generateCacheKeyDefault(realTerm, searchTerm, searchTermLat, searchTermLng, realTermLat, realTermLng string) string {
	// 1. Return the cache key for a rank
	if searchTermLat == "" || searchTermLng == "" || realTermLat == "" || realTermLng == "" {
		// given the application is meant to work with the suggestions manager, we can conclude that deriving
		// this sort of output is appropriate
		return strings.ToLower(searchTerm) + "-" + strings.ToLower(realTerm)
	}
	return strings.ToLower(searchTerm) + "-" + strings.ToLower(realTerm) + "lat-lng"
}

// createCacheKey defines the operation required to create cache key
type createCacheKey func(string, string, string, string, string, string) string

// generateCacheKey is a factory that generates the op used to create a cache key
func generateCacheKey(opType string) createCacheKey {
	switch opType {
	case "default":
		return generateCacheKeyDefault
	default:
		return generateCacheKeyDefault
	}

}
