package util

import (
	"fmt"
	"log"
	"sms-sorter/service/mqtt"
	"sms-sorter/service/telegram"
	"sms-sorter/util/uTime"
	"time"
)

func SendFailed(location string, err error) {
	t := uTime.GetKST(nil)
	msg := fmt.Sprintf("[ERROR/%s]\n=> %s", location, err)

	telegram.SendFailedMsg(msg, t)
	mqtt.SendFailedMessage(msg, t)
	log.Println(msg, t.Format(time.RFC822))
}

func SendStarted(hostname, localIP, pubIP string) {
	telegram.SendStarted(hostname, localIP, pubIP)
	mqtt.SendStarted(hostname, localIP, pubIP)
}

func SendNormalStopped(hostname, localIP, pubIP string) {
	telegram.SendStopped(hostname, localIP, pubIP)
	mqtt.SendStopped(hostname, localIP, pubIP)
}