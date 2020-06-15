# cachemanager
--
    import "github.com/Ekram-B2/suggestionscache/cachemanager"


## Usage

#### func  GetCacheManager

```go
func GetCacheManager(realTerm, searchTerm, realTermLat, realTermLng, searchTermLat, searchTermLng string, config config.Config) cacheManager
```
GetCacheManager is a factory applied to return a cache manager

#### func  HandleRequestForSuggestions

```go
func HandleRequestForSuggestions(rw http.ResponseWriter, req *http.Request)
```
HandleRequestForSuggestions is the handler meant to process the GET request made
for cached suggestions
