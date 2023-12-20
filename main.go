package main

import (
	"flag"
	"log"
	"net/http"

	"meowyplayerserver.com/core/analytics"
	"meowyplayerserver.com/core/authentication"
	"meowyplayerserver.com/core/collection"
	"meowyplayerserver.com/core/server"
	"meowyplayerserver.com/utility/logger"
)

var isHttps bool
var isRegister bool
var username string
var password string

func init() {
	flag.BoolVar(&isHttps, "https", true, "uses https, requires a valid certificate")
	flag.BoolVar(&isRegister, "register", false, "register a user, -username(must) and -password(optional)")
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
}

func main() {
	logger.Initiate()
	analytics.Initialize()
	authentication.Initialize()
	collection.Initialize()

	if isRegister {
		authentication.RegisterAccount(username, password)
		return
	}

	http.HandleFunc("/stats", server.Instance().ServerStats)
	http.HandleFunc("/list", server.Instance().ServerList)
	http.HandleFunc("/register", server.Instance().ServerRegister)
	http.HandleFunc("/upload", server.Instance().ServerUpload)
	http.HandleFunc("/download", server.Instance().ServerDownload)

	log.Println("meowyplayer server is running...")

	var err error
	if isHttps {
		//openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout keyFile.key -out certFile.crt
		err = http.ListenAndServeTLS(":443", "certFile.crt", "keyFile.key", nil)
	} else {
		err = http.ListenAndServe(":80", nil)
	}
	logger.Error(err, 0)
}
