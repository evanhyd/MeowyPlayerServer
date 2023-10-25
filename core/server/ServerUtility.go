package server

import (
	"fmt"
	"log"
	"net/http"

	"meowyplayerserver.com/utility/assert"
)

/*
A helper function that logs the error locally, and then send it over the http response
*/
func sendError(resp http.ResponseWriter, statusCode int, err error, message string) {
	log.Printf("%v: %v\n", message, err)
	resp.WriteHeader(statusCode)
	_, err = fmt.Fprintf(resp, "%v: %v", message, err)
	assert.NoErr(err, "failed to send error over http")
}
