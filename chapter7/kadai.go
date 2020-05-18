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

// 処理ハンドラマップ
var handlerMap = map[string]http.HandlerFunc{
	"/": func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, server.")
	},
	"/fortune": func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UnixNano())
		fortune := fortuneList[rand.Intn(len(fortuneList))]
		if p := r.FormValue("p"); p == "cheat" {
			// p=cheatが指定されているときは大吉
			fortune = "大吉"
		}
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
