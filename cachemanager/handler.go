package cachemanager

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	l4g "github.com/alecthomas/log4go"

	"github.com/Ekram-B2/rankmanager/rank/rankclient"
)

//HandleRequestForSuggestions is the handler meant to process the GET request made for cached suggestions
func HandleRequestForSuggestions(rw http.ResponseWriter, req *http.Request) {

	// 1. Check to see if the required query parameters are provided
	searchTerm := req.URL.Query().Get("searchTerm")
	if searchTerm == "" {
		l4g.Error("unable to find required query parameter 'searchTerm'")
		http.Error(rw, "there was an error retreiving one of the required query parameters; plese include a parameter for 'searchTerm' in your request :)", http.StatusBadRequest)
		return
	}

	realTerm := req.URL.Query().Get("realTerm")
	if searchTerm == "" {
		l4g.Error("unable to find required query parameter 'realTerm'")
		http.Error(rw, "there was an error retreiving one of the required query parameters; plese include a parameter for 'realTerm' in your request :)", http.StatusBadRequest)
		return
	}

	// 2. Check if latitudes and longitudes for the realTerm and searchTerm are provided. If these params aren't
	// provided, then this does not mean that there is an error as is it possible to query without this information
	searchTermLat := req.URL.Query().Get("searchTermLat")

	searchTermLng := req.URL.Query().Get("searchTermLng")

	realTermLat := req.URL.Query().Get("realTermLat")

	realTermLng := req.URL.Query().Get("realTermLng")

	// 3. Create a cache manager through which to proxy cache contents and commit requests to a rank manager

	cacheManager := GetCacheManager(realTerm,
		searchTerm,
		realTermLat,
		realTermLng,
		searchTermLat,
		searchTermLng,
		"test")

	// 4. Check the cache to see if rank can be found there

	// generate a cache key given the input
	cacheKey := generateCacheKey(os.Getenv("OPERATION_TYPE"))(realTerm,
		searchTerm,
		searchTermLat,
		searchTermLng,
		realTermLat,
		realTermLng)

	// apply cache key to determine if there is a cache hit or miss

	isCacheHit, returnedRank, err := cacheManager.getRankFromCache(cacheKey, os.Getenv("OPERATION_TYPE"))
	if err != nil {
		// this case is different from a cache miss. this means that there was an error with interacting with a cache
		l4g.Error("was unable to retreive rank from cache: %s", err.Error())
		http.Error(rw, "There was an error retrieving ranks from the rank manager service.", http.StatusInternalServerError)
		return
	}

	// 5. Based on whether the value was in the cache or not, we can proceed along two different control flows
	// A) If we haven't found the entry in the cache, then we will have to go to the rank manager service to draw content, and then store that content.
	// B) If we have found an entry in the cache, then we return that back to the suggestions manager

	if !isCacheHit {
		// There was a cache miss and thus, apply the downstream rank manager service to retreive a rank
		returnedRank, err = rankclient.GetRank(searchTerm,
			realTerm,
			realTermLat,
			searchTermLat,
			realTermLng,
			searchTermLng)
		if err != nil {
			l4g.Error("there was an error retrieving the from the rank manager service")
			http.Error(rw, "There was an error retrieving ranks from the rank manager service.", http.StatusInternalServerError)
			return
		}
		// Store the retreived rank into the cache
		err = cacheManager.putInCache(cacheKey, returnedRank, getByteEncoder(os.Getenv("OPERATION_TYPE"), returnedRank))
		if err != nil {
			l4g.Error("there was an error storing into the cache.")
			http.Error(rw, `There was error an storing into the cache.`, http.StatusInternalServerError)
			return
		}
	}

	// 6. Return the response back to the user

	// Set the response writer interface to specify content type
	rw.Header().Add("Content-Type", "application/json; charset=UTF-8")
	// Set the status code to OK since this is the desired result
	rw.WriteHeader(http.StatusOK)
	// Combine the return object with information concerning whether the cache was hit or not within a JSON array

	b := &bytes.Buffer{}
	err = json.NewEncoder(b).Encode(returnedRank)
	if err != nil {
		l4g.Error("error in marshalling to response to a output stream: %s", err.Error())
		http.Error(rw, "there was an error marshalling to the expected output; please try again later after waiting some time :)", http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(b.Bytes())
	if err != nil {
		l4g.Error("error in writing steam out to reponse content: %s", err.Error())
		http.Error(rw, "there was an writing to replying content to the response; please try again later after waiting some time :)", http.StatusInternalServerError)
		return
	}

}
