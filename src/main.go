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
	// pathを取るにはr.URL.Pathで受け取文末のスラッシュを削除
	fpath := strings.TrimRight(r.URL.Path, "/")
	fpath = strings.TrimLeft(fpath, "/")
	fname := filepath.Base(fpath)

	// ファイル存在チェック
	finfo, err := os.Stat(fpath)
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	}

	if finfo.IsDir() {
		// ディレクトリ配下のファイル一覧を取得
		fmt.Fprintf(w, fpath)
	} else {
		// ファイルをコピー
		// cpath := "/tmp/copy/" + fname
		cpath := `\\gndomain\MacShare\企画開発本部\開発部門\個人\tanaka-shu\` + fname

		src, err := os.Open(fpath)
		if err != nil {
			panic(err)
		}
		defer src.Close()

		dst, err := os.Create(cpath)
		if err != nil {
			panic(err)
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(w, cpath)
		fmt.Fprintln(w, fpath)
	}

	// ak, _ := url.QueryUnescape(path)
	// ak := Query()path.Encode()
	// fmt.Fprintf(w, "Hello!%s", path)
	// fmt.Fprintf(w, ak)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
