package telegram

import (
	"fmt"
	"log"
	"sms-sorter/config"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var to *MonitorRoom

var initialized = false

type MonitorRoom struct{}
type ChatRoom struct{}

func (*MonitorRoom) Recipient() string {
	return config.TelegramChatID
}

var bot *tb.Bot

func Init() {
	log.Println("Initializing telegram bot..")
	var err error
	bot, err = tb.NewBot(tb.Settings{
		Token:  config.TelegramAccessToken,
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	})
	if err != nil {
		log.Fatalln(err)
	}

	to = &MonitorRoom{}
	initialized = true
}

func SendMessage(message string) {
	loc, _ := time.LoadLocation("Asia/Seoul")

	SendMessageAt(message, time.Now().In(loc))
}

func SendMessageAt(message string, t time.Time) {
	if !config.IsProductionMode() || !initialized {
		return
	}

	msg := fmt.Sprintf("%s\n%s", message, t.Format(time.RFC822))
	log.Println("  Sending telegram Message...")
	_, err := bot.Send(to, msg)
	if err != nil {
		log.Println("   Failed to send Monitor", err)
	}
	log.Println("  [Telegram] message sent.")
}

func SendStarted(hostname string, localIP string, pubIP string) {
	log.Println("SendStarted()")
	msg := fmt.Sprintf("<%s> started successfully\nHostname:%s\nLocal IP:%s\nPublic IP:%s\n", config.ServerName, hostname, localIP, pubIP)
	SendMessage(msg)
}

func SendStopped(hostname string, localIP, pubIP string) {
	msg := fmt.Sprintf("<%s> stopping normally\nHostname:%s\nLocal IP:%s\nPublic IP:%s", config.ServerName, hostname, localIP, pubIP)
	SendMessage(msg)
}

func SendFailed(location string, err error, at time.Time) {
	msg := fmt.Sprintf("[ERROR/%s]\n=> %s", location, err)
	SendMessageAt(msg, at)
}

func SendFailedMsg(message string, at time.Time) {
	SendMessageAt(message, at)
}
