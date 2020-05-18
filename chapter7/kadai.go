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

// 処理ハンドラ
func handler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	fortune := fortuneList[rand.Intn(len(fortuneList))]
	if p := r.FormValue("p"); p == "cheat" {
		// p=cheatが指定されているときは大吉
		fortune = "大吉"
	}
	fmt.Fprint(w, fortune)
}

func main() {
	// ハンドラをエントリポイントと紐付け
	http.HandleFunc("/fortune", handler)

	// サーバをlocalhost:8080で起動
	http.ListenAndServe(":8080", nil)
}
