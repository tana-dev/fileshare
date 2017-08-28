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

type Html struct {
	FileinfoList [][]string
	Breadcrumbs  map[string]string
}

func handler(w http.ResponseWriter, r *http.Request) {

	var url string
	var fileinfoList [][]string
	var breadcrumbs map[string]string
	var fpath string
	var fname string

	url = "http://10.27.145.100:8080/"
	// url = "http://192.168.33.22:8080/"
	fpath = r.URL.Path
	fpath1 := r.URL.Path
	fpath1 = strings.TrimRight(fpath1, "/")

	// pathを取るにはr.URL.Pathで受け取文末のスラッシュを削除
	fpath = `\` + strings.Replace(r.URL.Path, "/", `\`, -1) // 1.Windows
	fpath = strings.TrimRight(fpath, `\`) // 1.Windows
	// fpath = strings.TrimRight(fpath, "/") // 2. Linux
	fname = filepath.Base(fpath)

	// ファイル存在チェック
	fi, err := os.Stat(fpath)
	if err != nil {
		fmt.Fprintf(w, "ファイル、もしくはディレクトが存在しません")
		return
	}

	// breadcrumbs create
	dirs_list := strings.Split(strings.TrimLeft(fpath1, "/"), "/")
	breadcrumbs = map[string]string{}
	var indexs map[int]string
	indexs = map[int]string{}
	for i := 0; i < len(dirs_list); i++ {
		for l := 0; l <= i; l++ {
			if l == 0 {
				indexs[i] = dirs_list[l] + "/"
			} else {
				indexs[i] = indexs[i] + dirs_list[l] + "/"
			}
		}
		index := url + indexs[i]
		breadcrumbs[index] = dirs_list[i]
	}

	if fi.IsDir() {
		fpaths := dirwalk(fpath)
		for _, fp := range fpaths {
			var fileinfo []string
			link := strings.Replace(fp, `\`, "/", -1)       // 2.Windows
			link = url + strings.Replace(link, "/", "", 2)  // 2.Windows
			// link := url + strings.Replace(fp, "/", "", 1) // 2.Linux
			name := filepath.Base(fp)

			fi, err := os.Stat(fp)
			if err != nil {
				fmt.Fprintf(w, "ファイルの読み込みに失敗しました")
				return
			}
			updatetime_tmp := fi.ModTime()
			updatetime := updatetime_tmp.Format("2006-01-02 15:04:05")

			fileinfo = append(fileinfo, link)
			fileinfo = append(fileinfo, name)
			fileinfo = append(fileinfo, updatetime)
			fileinfoList = append(fileinfoList, fileinfo)
		}

	} else {
		ext := fname[strings.LastIndex(fname, "."):]
		out := readfile(fpath)
		ctype := createContentType(ext)
		w.Header().Set("Content-Disposition", "attachment; filename="+fname)
		w.Header().Set("Content-Type", ctype)
		// w.Header().Set("Content-Length", string(len(out)))
		w.Write(out)
	}

	// fmt.Fprintln(w, fpath)
	// fmt.Fprintln(w, cpath)
	h := Html{
		FileinfoList: fileinfoList,
		Breadcrumbs:  breadcrumbs,
	}

	// results := []ToDo{ToDo{5323, "foo", "bar"}, ToDo{632, "foo", "bar"}}
	// funcs := template.FuncMap{"add": add}
	// temp := template.Must(template.New("index.html").Funcs(funcs).ParseFiles(templateDir + "/index.html"))
	// temp.Execute(writer, results)

	// funcs := template.FuncMap{"add": add}
	// tmpl := template.Must(template.New("./view/index.html").Funcs(funcs).ParseFiles("./view/index.html"))
	// tmpl.Execute(w, h)

	tmpl := template.Must(template.ParseFiles("./view/index.html"))
	tmpl.Execute(w, h)

}

func main() {
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

func add(x, y int) int {
	return x + y
}
