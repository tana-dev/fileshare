package pathchange

import (
	"fmt"
	"html/template"
	"net/http"
	"../../lib"
)

type Html struct {
	User           string
	Ip             string
	Download       map[string]string
	DownloadBase   string
	Pathchange     string
	PathchangeLink string
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var ip string
	var user string
	var url string
	var download map[string]string
	var downloadBase string
	var pathchange string
	var pathchangeLink string

	// get user info
	userConfig, err := appconfig.Parse("./config/user.json")
	if err != nil {
		fmt.Println("error ")
	}

	// set static info
	ip = userConfig.Host + ":" + userConfig.Port
	url = userConfig.Protocol + "://" + ip
	user = userConfig.Username

	// downloadセット
	download = map[string]string{}
	for i,v := range userConfig.Download {
		download[i] = url + "/download" + v
	}
	downloadBase = url + "/download"

	// pathchangeセット
	pathchange = url + "/pathchange"
	pathchangeLink = url + "/download" + userConfig.Pathchange

	h := Html{
		User:           user,
		Ip:             ip,
		Download:       download,
		DownloadBase:   downloadBase,
		Pathchange:     pathchange,
		PathchangeLink: pathchangeLink,
	}

	tmpl, _ := template.ParseFiles("./resources/view/pathchange/index.html")
	tmpl.Execute(w, h)

}
