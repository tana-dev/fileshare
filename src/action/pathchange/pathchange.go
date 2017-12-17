package pathchange

import (
	"fmt"
	"html/template"
	"net/http"
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

	tmpl, _ := template.ParseFiles("./resources/view/pathchange/index.html")
	tmpl.Execute(w, h)

}
