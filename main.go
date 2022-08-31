// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"github.com/Handson-peng/LineBotAccounting/sheet"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var bot *linebot.Client
var service sheet.Service
var	zone = time.FixedZone("UTC+8", 8*60*60)
type cmfunc func([]string) string

var command map[string]cmfunc = map[string]cmfunc{
	"記帳": Record,
	"總計": GetSum,
}

func main() {
	var err error

	ctx := context.Background()

	srv, err := sheets.NewService(ctx)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	sheet.SpreadsheetId = os.Getenv("SpreadsheetId")
	service = sheet.Service(*srv)

	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)

	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				textSlice := strings.Split(message.Text, " ")
				runfunc, ok := command[textSlice[0]]
				if !ok {
					return
				}
				result := runfunc(textSlice[1:])
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Get:"+message.Text+" ,resulte:"+result)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func Record(text []string) string {
	now := time.Now().In(zone)
	textSlice := make([]string, len(text)+1)
	textSlice[0] = now.Format("01/02 15:04")
	copy(textSlice[1:], text)
	service.AppendRow(now.Format("2006/01"), textSlice)
	return "紀錄成功"
}
func GetSum(text []string) string {
	now := time.Now().In(zone)
	res := service.ValueGet(now.Format("2006/01"), "G1")
	return fmt.Sprintf("%v", res[0][0])
}
