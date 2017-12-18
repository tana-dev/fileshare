package download

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"../../lib"
)

type Html struct {
	FileinfoList [][]string
	Breadcrumbs  map[string]string
	User         string
	Ip           string
	Download     map[string]string
	Upload       string
	Pathchange   string
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var ip string
	var user string
	var url string
	var fileinfoList [][]string
	var breadcrumbs map[string]string
	var fpath string
	var fname string
	var download map[string]string
	var upload string
	var pathchange string

	// ユーザー設定情報取得
	userConfig, err := appconfig.Parse("./config/user.json")
	if err != nil {
		fmt.Println("error ")
	}

	// ユーザー情報セット
	ip = userConfig.Host + ":"+ userConfig.Port
	url = userConfig.Protocol + "://"+ ip
	user = userConfig.Username

	// downloadセット
	download = map[string]string{}
	for i,v := range userConfig.Download {
		download[i] = url + "/download" + v
	}

	// uploadセット
	upload = url + "/upload"

	// pathchangeセット
	pathchange = url + "/pathchange"

	fpath = r.URL.Path
	fpath = strings.TrimLeft(fpath, "/download/")
	fpath = strings.TrimRight(fpath, "/")
	fpath = "/" + fpath

	// pathを取るにはr.URL.Pathで受け取文末のスラッシュを削除
	fpath = `\` + strings.Replace(fpath, "/", `\`, -1) // 1.Windows
	fpath = strings.TrimRight(fpath, `\`)                   // 1.Windows
	fname = filepath.Base(fpath)

	// ファイル存在チェック
	fi, err := os.Stat(fpath)
	if err != nil {
		fmt.Fprintf(w, "ファイル、もしくはディレクトが存在しません")
		return
	}

	// breadcrumbs create
	dirs_list := strings.Split(strings.TrimLeft(fpath, "\\\\"), "\\") // 1.Windows
//	dirs_list := strings.Split(strings.TrimLeft(fpath, "/"), "/") // 2.Linux
	breadcrumbs = map[string]string{}
	var indexs map[int]string
	indexs = map[int]string{}
	for i := 0; i < len(dirs_list); i++ {
		for l := 0; l <= i; l++ {
			if l == 0 {
				indexs[i] = "/" + dirs_list[l] + "/"
			} else {
				indexs[i] = indexs[i] + dirs_list[l] + "/"
			}
		}
		index := url + "/download" + indexs[i]
		breadcrumbs[index] = dirs_list[i]
	}

	if fi.IsDir() {
		fpaths := dirwalk(fpath)
		for _, fp := range fpaths {
			var fileinfo []string
			var dir string
			link := strings.Replace(fp, `\`, "/", -1)      // 1.Windows
			link = url + "/download" + strings.Replace(link, "/", "", 2) // 1.Windows
			//link := url + "/download" + fp // 2.Linux
			name := filepath.Base(fp)
			f, _ := os.Stat(fp)
			if f.IsDir() {
				dir = "fa-folder"
			} else {
				dir = "fa-file-o"
			}

			if err != nil {
				fmt.Fprintf(w, "ファイルの読み込みに失敗しました")
				return
			}
			updatetime_tmp := f.ModTime()
			updatetime := updatetime_tmp.Format("2006-01-02 15:04:05")

			fileinfo = append(fileinfo, link)
			fileinfo = append(fileinfo, name)
			fileinfo = append(fileinfo, updatetime)
			fileinfo = append(fileinfo, dir)
			fileinfoList = append(fileinfoList, fileinfo)
		}
		// sort.Sort(fileinfoList)

	} else {
		ext := fname[strings.LastIndex(fname, "."):]
		out := readfile(fpath)
		ctype := createContentType(ext)
		w.Header().Set("Content-Disposition", "attachment; filename="+fname)
		w.Header().Set("Content-Type", ctype)
		// w.Header().Set("Content-Length", string(len(out)))
		w.Write(out)
		return
	}

	h := Html{
		FileinfoList: fileinfoList,
		Breadcrumbs:  breadcrumbs,
		User:         user,
		Ip:           ip,
		Download:     download,
		Upload:       upload,
		Pathchange:   pathchange,
	}

	tmpl, _ := template.ParseFiles("./resources/view/download/index.html")
	tmpl.Execute(w, h)

}

func readfile(srcpath string) []byte {

	src, err := os.Open(srcpath)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	contents, _ := ioutil.ReadAll(src)

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

	// w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(name)))

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
		// ctype = "application/vnd.ms-excel" // EXCELファイル
		ctype = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" // EXCELファイル
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
	var dpaths []string
	var fpaths []string
	for _, file := range files {
		if 0 != strings.Index(file.Name(), ".") && 0 != strings.Index(file.Name(), "~$") && 0 != strings.Index(file.Name(), "Thumbs.db") {

			f := filepath.Join(dir, file.Name())

			// ファイル存在チェック
			fi, _ := os.Stat(f)
			if fi.IsDir() {
				dpaths = append(dpaths, filepath.Join(dir, file.Name()))
			} else {
				fpaths = append(fpaths, filepath.Join(dir, file.Name()))
			}
		}
	}

	if nil == dpaths && nil != fpaths {
		paths = fpaths
	} else if nil != dpaths && nil == fpaths {
		paths = dpaths
	} else {
		paths = append(dpaths, fpaths...)
	}

	return paths
}
