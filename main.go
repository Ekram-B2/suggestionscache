package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"github.com/Ekram-B2/suggestionscache/cachemanager"
)

func main() {

	// 1. Set up router object
	r := chi.NewRouter()

	r.Get("/", HandleRoot)

	r.Get("/rank", cachemanager.HandleRequestForSuggestions)

	// 2. Define catch all endpoint to help determine how to recover from the error case
	r.Get("/*", handleCatchAll)

	var bindingPort string
	if os.Getenv("DEPLOYMENT_TYPE") == "1" {
		// Hardcoded the port number in development mode
		bindingPort = ":8082"
	} else {
		bindingPort = ":" + os.Getenv("PORT")
	}

	// 3. Start the web application process and bind the application to a port

	http.ListenAndServe(bindingPort, r)

}

// HandleRoot is a handler function for the root server that is used for testing
func HandleRoot(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Received root ping from user"))
}

func handleCatchAll(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("The endpoint referenced is not currently supported"))
}
