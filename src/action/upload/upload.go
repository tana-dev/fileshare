package upload

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"../../lib"
)

type Html struct {
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

	h := Html{
		User:         user,
		Ip:           ip,
		Download:     download,
		Upload:       upload,
		Pathchange:   pathchange,
	}

//	templ_file, err := Asset("../resources/view/download/index.html")
//	tmpl, _ := template.New("tmpl").Parse(string(templ_file))
//		templates := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/base.html", "templates/view.html"))

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

    f, err := os.Create( upload + "/" + handler.Filename )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer f.Close()

    io.Copy(f, file)
    http.Redirect(w, r, "/upload/", http.StatusFound)

}
