package main

import (
	"fmt"
	"net/http"

	"meowyplayerserver.com/core/resource"
	"meowyplayerserver.com/core/server"
	"meowyplayerserver.com/utility/assert"
	"meowyplayerserver.com/utility/logger"
)

func main() {
	logger.Initiate()
	resource.MakeNecessaryPath()

	http.HandleFunc("/stats", server.GetInstance().ServerStats)
	http.HandleFunc("/list", server.GetInstance().ServerRequestList)
	http.HandleFunc("/upload", server.GetInstance().ServerRequestUpload)
	http.HandleFunc("/download", server.GetInstance().ServerRequestDownload)

	fmt.Println("meowyplayer server is running...")
	assert.NoErr(http.ListenAndServe(":80", nil), "meowyplayer server has crashed")
}
