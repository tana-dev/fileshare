package main

import (
	"net/http"
	"../downloader"
	"../uploader"
)

func main() {

	// ユーザー設定情報取得
//	userConfig, err := appconfig.Parse("./config/user.json")
//	if err != nil {
//		fmt.Println("error ")
//	}

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))

	http.HandleFunc("/downloader/", downloader.Handler)
	http.HandleFunc("/uploader/", uploader.Handler)
	http.HandleFunc("/uploadersave/", uploader.SaveHandler)

	http.ListenAndServe(":8080", nil)
}

