package main

import (
	"fmt"
	"net/http"
	"github.com/tana-dev/fileshare/action/download"
	"github.com/tana-dev/fileshare/action/upload"
	"github.com/tana-dev/fileshare/action/pathchange"
	"github.com/tana-dev/fileshare/lib"
)

func main() {

	// get user infomartion
	userConfig, err := appconfig.Parse("./config/user.json")
	if err != nil {
		fmt.Println("error ")
	}

	// Resources
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	// HandleFunc
	http.HandleFunc("/download/", download.Handler)
	http.HandleFunc("/upload/", upload.Handler)
	http.HandleFunc("/uploadsave/", upload.SaveHandler)
	http.HandleFunc("/uploadfile/", upload.SaveFileHandler)
	http.HandleFunc("/pathchange/", pathchange.Handler)

	// Listen
	http.ListenAndServe(":" + userConfig.Port, nil)
}
