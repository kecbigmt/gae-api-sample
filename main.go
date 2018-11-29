package main

import (
  "fmt"
  "net/http"
  "encoding/json"

  "google.golang.org/appengine"
  "google.golang.org/appengine/log"
  "google.golang.org/appengine/urlfetch"
)

const url = "https://www.google.com/"

type Response struct {
  Status string `json:"status"`
  Message string `json:"message"`
}

func main() {
  http.HandleFunc("/fore", helloForeground)
  http.HandleFunc("/back", helloBackground)
  appengine.Main()
}

// OKパターン
func helloForeground(w http.ResponseWriter, r *http.Request) {
  // URLにリクエストして返ってきたステータスをログに残す
  ctx := appengine.NewContext(r) // App Engineのcontext作成
  client := urlfetch.Client(ctx) // contextを使ってhttp.Clientを作成
  resp, err := client.Get(url)
  if err != nil {
    log.Errorf(ctx, "TEST: helloForeground: %s", err)
    return
  }
  log.Infof(ctx, "'%s' returns %s", url, resp.Status)

  // レスポンス
  json.NewEncoder(w).Encode(Response{
    Status: "ok",
    Message: fmt.Sprintf("'%s' returns %s", url, resp.Status),
  })
}

// NGパターン
func helloBackground(w http.ResponseWriter, r *http.Request) {
  // URLにリクエストして返ってきたステータスをログに残すgoroutine
  go func() {
    ctx := appengine.NewContext(r) // App Engineのcontext作成
    client := urlfetch.Client(ctx) // contextを使ってhttp.Clientを作成
    resp, err := client.Get(url)
    if err != nil {
      log.Errorf(ctx, "TEST: helloBackground: %s", err)
      return
    }
    log.Infof(ctx, "'%s' returns %s", url, resp.Status)
  }()

  // レスポンス
  json.NewEncoder(w).Encode(Response{
    Status: "ok",
    Message: "processing in background.",
  })
}
