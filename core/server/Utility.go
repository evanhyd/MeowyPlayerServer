package server

import (
	"log"
	"net/http"
)

/*
A helper function that logs the error locally, and then send it over the http response
*/
func sendError(resp http.ResponseWriter, statusCode int, errorMsg string) {
	log.Print(errorMsg)
	http.Error(resp, errorMsg, statusCode)
}
