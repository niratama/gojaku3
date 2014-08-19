package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"flag"
	"github.com/darkhelmet/twitterstream"
	"github.com/rakyll/globalconf"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

// WebSocketで返すTweetデータ
type WSTweet struct {
	Text            string `json:"text"`
	Name            string `json:"name"`
	ScreenName      string `json:"screen_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

// flag/globalconfから取得するデータ
var (
	accessToken       = flag.String("access_token", "", "Twitter access token")
	accessTokenSecret = flag.String("access_token_secret", "", "Twitter access token secret")
	consumerKey       = flag.String("consumer_key", "", "Twitter consumer key")
	consumerSecret    = flag.String("consumer_secret", "", "Twitter consumer secret")
	listen            = flag.String("listen", ":3000", "HTTP Listen port")
	timeoutString     = flag.String("timeout", "30m", "Connection timeout")
)
var timeout time.Duration

// Twitter検索
func twitterSearchHandler(ws *websocket.Conn) {
	defer ws.Close()

	// 検索キーワードの取得
	req := ws.Request() // http.Requestが返る
	query := req.FormValue("q")
	log.Printf("query: %s", query)

	// Twitter streaming APIに接続
	client := twitterstream.NewClientTimeout(
		*consumerKey,
		*consumerSecret,
		*accessToken,
		*accessTokenSecret,
		timeout,
	)
	conn, err := client.Track(query)
	if err != nil {
		log.Printf("Tracking failed: %s", err)
		return
	}
	defer conn.Close()

	for {
		if tweet, err := conn.Next(); err == nil {
			// Tweetが公式Retweetだった場合はなにもしない
			if tweet.RetweetedStatus != nil {
				continue
			}

			// Websocketに流すJSONを作成
			data := WSTweet{
				tweet.Text,
				tweet.User.Name,
				tweet.User.ScreenName,
				tweet.User.ProfileImageUrl,
			}
			json, _ := json.Marshal(data)

			// Websocketに流す
			_, err = ws.Write(json)
			if err != nil {
				log.Printf("Writing to Websocket failed: %s", err)
				return
			}
		} else {
			log.Printf("Decoding tweet failed: %s", err)
			return
		}
	}
}

func init() {
	// ホームディレクトリ以下の設定を読み込む
	conf, err := globalconf.New("twstream")
	if err != nil {
		log.Fatalf("Can't load config: %s", err)
	}
	conf.ParseAll()

	// タイムアウトの値をtime.Durationに変換する
	timeout, err = time.ParseDuration(*timeoutString)
	if err != nil {
		log.Fatalf("Can't parse timeout(%s) :%s", *timeoutString, err)
	}
}

func main() {
	// staticなファイルの置き場
	pwd, _ := os.Getwd()
	staticPath := path.Join(pwd, "static")

	// Twitter検索
	http.Handle("/search", websocket.Handler(twitterSearchHandler))
	// それ以外はstaticなファイル
	http.Handle("/", http.FileServer(http.Dir(staticPath)))

	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}
