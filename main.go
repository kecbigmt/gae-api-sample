package main

import (
  "fmt"
  "net/http"
  "encoding/json"

  "google.golang.org/appengine"
  "google.golang.org/appengine/log"
  "google.golang.org/appengine/urlfetch"

  "google.golang.org/appengine/taskqueue"
)

const url = "https://www.google.com/"

type Response struct {
  Status string `json:"status"`
  Message string `json:"message"`
}

func main() {
  http.HandleFunc("/fore", helloForeground)
  http.HandleFunc("/back", helloBackground)
  http.HandleFunc("/back/worker", helloBackgroundWorker)
  appengine.Main()
}

// Task Queueなしで処理するAPI（60秒でタイムアウト）
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

// レスポンスだけ速攻で返して時間のかかる処理はTask Queueに渡すAPI（10分 or 24時間でタイムアウト）
func helloBackground(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r) // App Engineのcontext作成

  // タスクをTaskQueueに投げる
  t := taskqueue.NewPOSTTask("/back/worker", map[string][]string{})
  if _, err := taskqueue.Add(ctx, t, ""); err != nil {
    http.Error(w, "internal server error", http.StatusInternalServerError)
    return
  }

  // レスポンス
  json.NewEncoder(w).Encode(Response{
    Status: "ok",
    Message: "processing in background.",
  })
}

// URLにリクエストして返ってきたステータスをログに残すworker
func helloBackgroundWorker(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r) // App Engineのcontext作成
  client := urlfetch.Client(ctx) // contextを使ってhttp.Clientを作成
  resp, err := client.Get(url)
  if err != nil {
    log.Errorf(ctx, "TEST: helloBackground: %s", err)
    return
  }
  log.Infof(ctx, "'%s' returns %s", url, resp.Status)

  return
}
