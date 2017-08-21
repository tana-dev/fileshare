package main

import (
  "fmt"
  "net/http"
  // "net/url"
  "io/ioutil"
//  "io"
  "os"
  "strings"
  "path/filepath"
)

func handler(w http.ResponseWriter, r *http.Request) {
  // pathを取るにはr.URL.Pathで受け取文末のスラッシュを削除
  fpath := strings.TrimRight(r.URL.Path, "/")
  fname := filepath.Base(fpath)

  // ファイル存在チェック
  finfo, err := os.Stat(fpath)
  if err != nil {
    fmt.Fprintf(w, "error")
    return
  }

  if(finfo.IsDir()){
    // ディレクトリ配下のファイル一覧を取得
    fmt.Fprintf(w, fpath)
  } else {
    // ファイルをコピー
    cpath := "/tmp/copy/" + fname
//    _ = os.Link(fpath, cpath)
//    _, err = io.Copy(fpath, cpath)
    if err := os.Rename(fpath, cpath); err != nil {
      fmt.Println(err)
    }
    fmt.Fprintf(w, cpath)
    fmt.Fprintf(w, fpath)
  }

  // ak, _ := url.QueryUnescape(path)
  // ak := Query()path.Encode()
  // fmt.Fprintf(w, "Hello!%s", path)
  // fmt.Fprintf(w, ak)
}

func main() {
  http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
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
