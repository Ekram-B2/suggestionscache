package cachemanager

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	l4g "github.com/alecthomas/log4go"

	"github.com/Ekram-B2/rankmanager/rank/rankclient"
	"github.com/Ekram-B2/suggestionsmanager/config"
)

//HandleRequestForSuggestions is the handler meant to process the GET request made for cached suggestions
func HandleRequestForSuggestions(rw http.ResponseWriter, req *http.Request) {

	// 1. Check to see if the required query parameters are provided
	searchTerm := req.URL.Query().Get("searchTerm")
	if searchTerm == "" {
		l4g.Error("SYSTEM-ERROR: unable to find required query parameter 'searchTerm'")
		http.Error(rw, "There was an error retreiving one of the required query parameters; plese include a parameter for 'searchTerm' in your request.", http.StatusBadRequest)
		return
	}

	realTerm := req.URL.Query().Get("realTerm")
	if searchTerm == "" {
		l4g.Error("SYSTEM-ERROR: unable to find required query parameter 'realTerm'")
		http.Error(rw, "There was an error retreiving one of the required query parameters; plese include a parameter for 'realTerm' in your request.", http.StatusBadRequest)
		return
	}

	// 2. Check if latitudes and longitudes for the realTerm and searchTerm are provided. If these params aren't
	// provided, then this does not mean that there is an error as is it possible to query without this information
	searchTermLat := req.URL.Query().Get("searchTermLat")

	searchTermLng := req.URL.Query().Get("searchTermLng")

	realTermLat := req.URL.Query().Get("realTermLat")

	realTermLng := req.URL.Query().Get("realTermLng")

	// 3. Load configuration

	loadedConfig, err := config.LoadConfiguration(config.GetConfigPath(os.Getenv("CONFIG_OPERATION_TYPE")))
	if err != nil {
		l4g.Error("OPERATION-ERROR: there was an error loading the configuration object: %s", err.Error())
		http.Error(rw, "We were unable to retreive a rank for the real term provided; please try again after waiting some time.", http.StatusInternalServerError)
		return
	}

	// 4. Create a cache manager through which to proxy cache contents and commit requests to a rank manager
	cacheManager := GetCacheManager(realTerm,
		searchTerm,
		realTermLat,
		realTermLng,
		searchTermLat,
		searchTermLng,
		loadedConfig)

	// 5. Check the cache to see if rank can be found there

	// generate a cache key given the input
	cacheKey := generateCacheKey(loadedConfig.CacheKeyType)(realTerm,
		searchTerm,
		realTermLat,
		realTermLng,
		searchTermLat,
		searchTermLng,
	)

	// apply cache key to determine if there is a cache hit or miss

	isCacheHit, returnedRank, err := cacheManager.getRankFromCache(cacheKey, getByteDecoder(loadedConfig.CacheKeyType))

	if err != nil {
		// this case is different from a cache miss. this means that there was an error with interacting with a cache
		l4g.Error("SYSTEM-ERROR: was unable to retreive rank from cache: %s", err.Error())
		http.Error(rw, "There was an error a rank for the realterm.", http.StatusInternalServerError)
		return
	}

	// 6. Based on whether there was a cache hit or miss, set the X-Cache header
	if isCacheHit {
		rw.Header().Set("X-Cache", "HIT")
	} else {
		rw.Header().Set("X-Cache", "MISS")
	}
	// 7. Based on whether the value was in the cache or not, we can proceed along two different control flows
	// A) If we haven't found the entry in the cache, then we will have to go to the rank manager service to draw content, and then store that content.
	// B) If we have found an entry in the cache, then we return that back to the suggestions manager

	if !isCacheHit {
		// Set the header to state that there was a miss
		// There was a cache miss and thus, apply the upstream rank manager service to retreive a rank
		returnedRank, err = rankclient.GetRank(searchTerm,
			realTerm,
			realTermLat,
			searchTermLat,
			realTermLng,
			searchTermLng)
		if err != nil {
			l4g.Error("SYSTEM-ERROR: there was an error retrieving the from the rank manager service")
			http.Error(rw, "There was an error retrieving the rank.", http.StatusInternalServerError)
			return
		}
		// Store the retreived rank into the cache
		err = cacheManager.putInCache(cacheKey, returnedRank, getByteEncoder(loadedConfig.ByteEncoderType, returnedRank))
		if err != nil {
			l4g.Error("SYSTEM-ERROR: there was an error storing into the cache")
			http.Error(rw, "There was an error retrieving the rank.", http.StatusInternalServerError)
			return
		}
	}

	// 7. Return the response back to the user

	// Set the response writer interface to specify content type
	rw.Header().Add("Content-Type", "application/json; charset=UTF-8")
	// Set the status code to OK since this is the desired result
	rw.WriteHeader(http.StatusOK)
	// Combine the return object with information concerning whether the cache was hit or not within a JSON array

	b := &bytes.Buffer{}
	err = json.NewEncoder(b).Encode(returnedRank)
	if err != nil {
		l4g.Error("OPERATION-ERROR: unable to marshall the response out to a byte stream: %s", err.Error())
		http.Error(rw, "There was an error a rank for the realterm.", http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(b.Bytes())
	if err != nil {
		l4g.Error("OPERATION-ERROR: unable to write stream out to reponse writer: %s", err.Error())
		http.Error(rw, "There was an error a rank for the realterm.", http.StatusInternalServerError)
		return
	}

}
