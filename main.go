package main

import (
	"log"
	"net/http"

	"meowyplayerserver.com/core/analytics"
	"meowyplayerserver.com/core/authentication"
	"meowyplayerserver.com/core/collection"
	"meowyplayerserver.com/core/server"
	"meowyplayerserver.com/utility/logger"
)

func main() {
	logger.Initiate()
	analytics.Initialize()
	authentication.Initialize()
	collection.Initialize()

	http.HandleFunc("/stats", server.Instance().ServerStats)
	http.HandleFunc("/list", server.Instance().ServerList)
	http.HandleFunc("/upload", server.Instance().ServerUpload)
	http.HandleFunc("/download", server.Instance().ServerDownload)

	log.Println("meowyplayer server is running...")
	err := http.ListenAndServe(":80", nil)
	logger.Error(err, 0)
}
