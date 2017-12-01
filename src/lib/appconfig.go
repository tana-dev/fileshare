package appconfig

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
)

// 設定ファイルの値を表現する構造体
type Config struct {
	Username string
	Protocol string
	Host string
	Port string
	Bookmarks map[string]string
}

// 設定ファイルを読み込む
func Parse(filename string) (Config, error) {
	var c Config
	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error: t")
		return c, err
	}
	err = json.Unmarshal(jsonString, &c)
	if err != nil {
		fmt.Println("error: t1")
		return c, err
	}
	return c, nil
}
