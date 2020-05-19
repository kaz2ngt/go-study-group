package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/apbgo/go-study-group/chapter7/kadai/model"
)

func main() {
	if err := postUserForm(); err != nil {
		log.Fatal(err)
	}
}

func postUserForm() error {
	// リクエストjsonデータの作成
	reqModel := &model.Request{
		UserID: 1,
		Name:   "user_fortune",
	}
	var body bytes.Buffer
	enc := json.NewEncoder(&body)
	if err := enc.Encode(reqModel); err != nil {
		return err
	}
	// ヘッダーデータの作成
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := request("POST", "http://localhost:8080/user_fortune", bytes.NewReader(body.Bytes()), headers)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var resModel model.Response
	dec := json.NewDecoder(res.Body)
	if err = dec.Decode(&resModel); err != nil {
		return err
	}

	fmt.Println(resModel)
	return nil
}

func request(method string, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	// クライアントを作成（タイムアウトを指定）
	client := &http.Client{Timeout: time.Duration(10) * time.Second}

	// タイムアウト・キャンセル用のコンテキストを作成
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Request を生成
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// headerがあれば設定
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// リクエストにタイムアウトを設定したコンテキストを持たせる
	req = req.WithContext(ctx)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
