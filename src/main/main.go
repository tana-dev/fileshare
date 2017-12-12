package main

import (
	"net/http"
	"../downloader"
)

func main() {
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))

	http.HandleFunc("/", downloader.DownloaderHandler)

	http.ListenAndServe(":8080", nil)
}

