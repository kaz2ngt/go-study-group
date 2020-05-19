package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/apbgo/go-study-group/chapter7/kadai/model"
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
	"/user_fortune": func(w http.ResponseWriter, r *http.Request) {
		// リクエストボディの取得
		defer r.Body.Close()

		// リクエストBodyの内容を取得
		var req model.Request
		// []byteに変換
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err.Error())
			return
		}
		err = json.Unmarshal(data, &req)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(req)

		// レスポンスの作成
		response := model.Response{
			Status: http.StatusOK,
			Data:   fmt.Sprintf("ID:%vの%sさんの運勢は%sです！", req.UserID, req.Name, runFortune(false)),
		}
		log.Println(response)

		var res bytes.Buffer
		enc := json.NewEncoder(&res)
		if err = enc.Encode(response); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")

		w.Write(res.Bytes())
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
