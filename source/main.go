package main

import (
	"flag"
	"log"
	"meowyplayerserver/api"
	"net/http"
)

func main() {
	isHttp := flag.Bool("http", false, "Use http insteaed of https.")
	flag.Parse()

	apiManager := api.MakeAPIManager()
	if err := apiManager.Initialize(); err != nil {
		log.Fatalln(err)
	}

	var err error
	if *isHttp {
		err = http.ListenAndServe(":80", nil)
	} else {
		//openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout cert.key -out cert.crt
		err = http.ListenAndServeTLS(":443", "cert.crt", "cert.key", nil)
	}
	log.Println(err)
}
