package main

import (
	"fmt"
	"net/http"
	// "net/url"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {

	// コピー先
	cpath := "/tmp/copy/"
	// cpath := `\\gndomain\MacShare\企画開発本部\開発部門\個人\tanaka-shu\`

	// pathを取るにはr.URL.Pathで受け取文末のスラッシュを削除
	fpath := strings.TrimRight(r.URL.Path, "/")
	// fpath = strings.TrimLeft(fpath, "/")
	fname := filepath.Base(fpath)

	// ファイル存在チェック
	finfo, err := os.Stat(fpath)
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	}

	if finfo.IsDir() {
		// ディレクトリ配下のファイル一覧を取得
		fmt.Fprintln(w, dirwalk(fpath))
		return
	} else {
		// ファイルをコピー
		cpath = cpath + fname
		copyfile(fpath, cpath)
	}

	fmt.Fprintln(w, fpath)
	fmt.Fprintln(w, cpath)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func copyfile(srcpath string, dstpath string) {

	src, err := os.Open(srcpath)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create(dstpath)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		panic(err)
	}
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
