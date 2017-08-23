package main

import (
	"fmt"
	"net/http"
	// "net/url"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Person struct {
	Name string
	From string
}

func handler(w http.ResponseWriter, r *http.Request) {

	// コピー先
	cpath := `/tmp/copy/`
	// cpath := `\\gndomain\MacShare\企画開発本部\開発部門\個人\tanaka-shu\copy\`

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
		// fpaths := dirwalk(fpath)
		dirwalk(fpath)
		// for _, fp := range fpaths {
		// 	fmt.Fprintln(w, fp)
		// }

	} else {
		cpath = cpath + fname
		ext := fname[strings.LastIndex(fname, "."):]
		// copyfile(fpath, cpath)
		out := readfile(fpath)
		ctype := createContentType(ext)
		w.Header().Set("Content-Disposition", "attachment; filename="+fname)
		w.Header().Set("Content-Type", ctype)
		// w.Header().Set("Content-Length", string(len(out)))
		w.Write(out)
	}

	// fmt.Fprintln(w, fpath)
	// fmt.Fprintln(w, cpath)
	p := Person{
		Name: "sekky",
		From: "埼玉",
	}

	tmpl := template.Must(template.ParseFiles("./view/index.html"))
	tmpl.Execute(w, p)

}

func main() {
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
    http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func readfile(srcpath string) []byte {

	src, err := os.Open(srcpath)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	contents, _ := ioutil.ReadAll(src)
	// if err != nil {
	// 	panic(err)
	// }
	// defer contents.Close()

	return contents
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

func createContentType(ext string) string {

	var ctype string

	switch ext {
	case ".txt":
		ctype = "text/plain"
	case ".csv":
		ctype = "text/csv" // CSVファイル
	case ".html":
		ctype = "text/html" // HTMLファイル
	case ".css":
		ctype = "text/css" // CSSファイル
	case ".js":
		ctype = "text/javascript" // JavaScriptファイル
	case ".exe":
		ctype = "application/octet-stream" // EXEファイルなどの実行ファイル
	case ".pdf":
		ctype = "application/pdf" // PDFファイル
	case ".xlsx":
		ctype = "application/vnd.ms-excel" // EXCELファイル
	case ".ppt":
		ctype = "application/vnd.ms-powerpoint" // PowerPointファイル
	case ".docx":
		ctype = "application/msword" // WORDファイル
	case ".jpeg", ".jpg":
		ctype = "image/jpeg" // JPEGファイル(.jpg, .jpeg)
	case ".png":
		ctype = "image/png" // PNGファイル
	case ".gif":
		ctype = "image/gif" // GIFファイル
	case ".bmp":
		ctype = "image/bmp" // Bitmapファイル
	case ".zip":
		ctype = "application/zip" // Zipファイル
	case ".lzh":
		ctype = "application/x-lzh" // LZHファイル
	case ".tar":
		ctype = "application/x-tar" // tarファイル/tar&gzipファイル
	case ".mp3":
		ctype = "audio/mpeg" // MP3ファイル
	case ".mp4":
		ctype = "audio/mp4" // MP4ファイル
	case ".mpeg":
		ctype = "video/mpeg" // MPEGファイル（動画）
	default:
		ctype = "text/plain"
	}

	return ctype
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if 0 != strings.Index(file.Name(), ".") {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return paths
}
