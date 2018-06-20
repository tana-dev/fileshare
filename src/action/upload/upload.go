package upload

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	//	"strings"
	"../../lib"
)

type Html struct {
	User       string
	Ip         string
	Download   map[string]string
	Upload     string
	Pathchange string
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var ip string
	var user string
	var url string
	var download map[string]string
	var upload string
	var pathchange string

	// ユーザー設定情報取得
	userConfig, err := appconfig.Parse("./config/user.json")
	if err != nil {
		fmt.Println("error ")
	}

	// ユーザー情報セット
	ip = userConfig.Host + ":" + userConfig.Port
	url = userConfig.Protocol + "://" + ip
	user = userConfig.Username

	// downloadセット
	download = map[string]string{}
	for i, v := range userConfig.Download {
		download[i] = url + "/download" + v
	}

	// uploadセット
	upload = url + "/upload"

	// pathchangeセット
	pathchange = url + "/pathchange"

	h := Html{
		User:       user,
		Ip:         ip,
		Download:   download,
		Upload:     upload,
		Pathchange: pathchange,
	}

	tmpl, _ := template.ParseFiles("./resources/view/upload/index.html")
	tmpl.Execute(w, h)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {

	// ユーザー設定情報取得
	userConfig, err := appconfig.Parse("./config/user.json")
	if err != nil {
		fmt.Println("error ")
	}
	upload := userConfig.Upload
	//	upload = `\` + strings.Replace(upload, "/", `\`, -1) // 1.Windows

	if r.Method != "POST" {
		http.Error(w, "Allowed POST method only", http.StatusMethodNotAllowed)
		return
	}

	err = r.ParseMultipartForm(32 << 20) // maxMemory
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("upload_files")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	f, err := os.Create(upload + "\\" + handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)
	http.Redirect(w, r, "/upload/", http.StatusFound)

}

func SaveFileHandler(w http.ResponseWriter, r *http.Request) {

	// check Method
	if r.Method != "POST" {
		http.Error(w, "Allowed POST method only", http.StatusMethodNotAllowed)
		return
	}


	// check maxMemory
	err := r.ParseMultipartForm(32 << 20) // maxMemory
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get current directory (from post data)
	currentDir := r.FormValue("currentDirectory")
	//currentDir = `\` + strings.Replace(currentDir, "/", `\`, -1)  // 1.Windows


	// get file data (from post data)
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// get filename
	targetFile := currentDir + "/" + handler.Filename
	//targetFile = `\` + strings.Replace(targetFile, "/", `\`, -1)  // 1.Windows

	// check exited file
	_, err = os.Stat(targetFile)
	if err == nil {
		fmt.Println("同名ファイルが存在してます")
		http.Redirect(w, r, "/download/"+currentDir, http.StatusFound)
	}

	// open taget file
	f, err := os.Create(targetFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()


	io.Copy(f, file)
	http.Redirect(w, r, "/download/"+currentDir, http.StatusFound)

}
