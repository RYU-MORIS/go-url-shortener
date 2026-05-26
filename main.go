package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os" // 💡 クラウドの環境変数を読み込むための道具を追加
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var urlDatabase = make(map[string]string)

func generateShortKey(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTリクエストのみ受け付けます", http.StatusMethodNotAllowed)
		return
	}
	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URLが入力されていません", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey(6)
	urlDatabase[shortKey] = longURL

	// 💡 転送先のドメイン（localhost）は後ほどクラウド用に自動対応させるため、
	// ここではとりあえず相対パスやホスト名に対応できるようにキーだけを返す形式などに対応させますが、
	// 一旦、動いているサーバーのホスト名を自動で取得するように賢く書き換えます。
	response := map[string]string{
		"short_url": fmt.Sprintf("http://%s/%s", r.Host, shortKey),
		"original":  longURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[1:]

	// 💡 ここを追加：カギが空っぽなら、UI画面（index.html）を表示する！
	if shortKey == "" {
		http.ServeFile(w, r, "index.html")
		return
	}

	longURL, exists := urlDatabase[shortKey]
	if !exists {
		http.Error(w, "URLが見つかりません", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	// 💡【重要】GCPが指定するポート番号（待ち受け番号）を取得。無ければ8080にする。
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 サーバー起動ポート: %s\n", port)
	// 指定されたポートで待ち受けを開始
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
