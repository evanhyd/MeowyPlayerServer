package main

import (
	"fmt"
	"net/http"

	"meowyplayerserver.com/core/resource"
	"meowyplayerserver.com/core/server"
	"meowyplayerserver.com/utility/logger"
)

func main() {
	logger.Initiate()
	resource.MakeNecessaryPath()

	http.HandleFunc("/stats", server.GetInstance().ServerStats)
	http.HandleFunc("/list", server.GetInstance().ServerRequestList)
	http.HandleFunc("/upload", server.GetInstance().ServerRequestUpload)
	http.HandleFunc("/download", server.GetInstance().ServerRequestDownload)
	http.HandleFunc("/remove", server.GetInstance().ServerRequestRemove)

	fmt.Println("meowyplayer server is running...")
	err := http.ListenAndServe(":80", nil)
	logger.Error(err, 0)
}
