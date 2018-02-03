package main

import (
	"net/http"
	"../action/download"
	"../action/upload"
	"../action/pathchange"
)

func main() {

	// ユーザー設定情報取得
//	userConfig, err := appconfig.Parse("./config/user.json")
//	if err != nil {
//		fmt.Println("error ")
//	}

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))

	http.HandleFunc("/download/", download.Handler)
	http.HandleFunc("/upload/", upload.Handler)
	http.HandleFunc("/uploadsave/", upload.SaveHandler)
	http.HandleFunc("/pathchange/", pathchange.Handler)

	http.ListenAndServe(":8080", nil)
}
