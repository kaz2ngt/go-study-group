package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var fortuneList = [...]string{
	"大吉",
	"中吉",
	"吉",
	"凶",
}

func runFortune(isCheat bool) string {
	rand.Seed(time.Now().UnixNano())
	if isCheat {
		// 大吉を返却
		return fortuneList[0]
	}
	fortune := fortuneList[rand.Intn(len(fortuneList))]
	return fortune
}

// 処理ハンドラマップ
var handlerMap = map[string]http.HandlerFunc{
	"/": func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, server.")
	},
	"/fortune": func(w http.ResponseWriter, r *http.Request) {
		p := r.FormValue("p")
		isCheat := p == "cheat"
		fortune := runFortune(isCheat)
		fmt.Fprint(w, fortune)
	},
}

func main() {
	// ハンドラをエントリポイントと紐付け
	for path, handler := range handlerMap {
		http.HandleFunc(path, handler)
	}

	// サーバをlocalhost:8080で起動
	http.ListenAndServe(":8080", nil)
}
