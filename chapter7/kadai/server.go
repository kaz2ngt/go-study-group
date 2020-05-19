package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	mux := http.NewServeMux()
	// ハンドラをエントリポイントと紐付け
	for path, handler := range handlerMap {
		mux.HandleFunc(path, handler)
	}

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// OSからのシグナルを待つ
	go func() {
		// SIGTERM: コンテナが終了する時に送信されるシグナル
		// SIGINT: Ctrl+c
		sigCh := make(chan os.Signal, 1)
		// 受け取るシグナルを指定
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		// チャネルでの待受、シグナルを受け取るまで以降は処理されない
		<-sigCh

		log.Println("start graceful shutdown server.")
		// タイムアウトのコンテキストを設定（後述）
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// Graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Println(err)
			// 接続されたままのコネクションも明示的に切る
			srv.Close()
		}
		log.Println("HTTPServer shutdown.")
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Print(err)
	}
}
