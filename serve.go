package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
)

func serve() {
	fmt.Println("This server should only be used as a test server!")

	// Get the address from the SITE_ADDRESS environment variable,
	//  defaulting to localhost:8080
	address := os.Getenv("SITE_ADDRESS")
	if address == "" {
		address = "localhost:8080"
	}

	// Create static file server on docs
	mux := http.NewServeMux()

	dir := http.Dir("docs")
	fileServer := http.FileServer(dir)

	mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		// Check for 404 and append '.html' if 404 then serve
		buf := httptest.NewRecorder()
		fileServer.ServeHTTP(buf, request)
		if buf.Result().StatusCode == 404 {
			request.URL.Path += ".html"
			fileServer.ServeHTTP(response, request)
		} else {
			fileServer.ServeHTTP(response, request)
		}
	})

	fmt.Printf("Hosting on http://%s/\n", address)

	check(http.ListenAndServe(address, mux))
}
