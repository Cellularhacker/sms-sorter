package main

import (
	"github.com/robfig/cron"
	"log"
	"net"
	"os"
	"os/signal"
	"sms-sorter/data"
	"sms-sorter/model/sms"
	"sms-sorter/model/thecall"
	"sms-sorter/service/renewSpamDB"
	"sms-sorter/util"
	"syscall"

	"github.com/chyeh/pubip"
)

var c *cron.Cron
var localIP = ""
var pubIP = ""
var hostname = ""

func init() {
	c = cron.New()
	//telegram.Init()
	//mqtt.Init()
	//
	// Initializing Data...
	data.Init()
	//
	// Set Context
	sms.SetStore(data.NewSmsStore())
	thecall.SetStore(data.NewTheCallStore())
	//finefss.SetStore(data.NewFineFssStore())
	//finefssCategory.SetStore(data.NewFineFssCategoryStore())
	// Service...
}

func main() {
	err := renewSpamDB.TheCall()
	if err != nil {
		log.Fatalln(err)
	}
	return
	//res := `{"from_number": "01065146909","contact_name": "Cellularhacker@DEXEOS","text": "[Web발신]정확하고 안전하게 !전문가와 함께 진행 !하루 평균 200% 순이익https://bit.ly/38jCWnA","occurred_at": "January 18, 2020 at 07:20PM"}`
	//t := &temp{}
	//
	//err := json.Unmarshal([]byte(res), t)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//at, err := time.Parse("January 2, 2006 at 03:04PM", t.OccurredAtStr)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//loc, _ := time.LoadLocation("Asia/Seoul")
	//at = at.Add(-9 * time.Hour)
	//
	//log.Printf("occurred_at: %v\n", at.In(loc))
	//
	//r := regexp.MustCompile("([0-9]{3})-?([0-9]{4})-?([0-9]{4})")
	////r := regexp.MustCompile("^(01[016789]{1}|02|0[3-9]{1}[0-9]{1})-?[0-9]{3,4}-?[0-9]{4}$")
	//list := r.FindAllStringSubmatch(t.FromNumber, -1)
	//
	//for i, e := range list {
	//	for j, f := range e {
	//		log.Printf("[%d:%d] %v\n", i, j, f)
	//	}
	//}
	//
	//return
	// Temporary backup.
	//err := pushover.Backup()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//return

	// Send Startup Message
	go func() {
		hostname, _ = os.Hostname()
		localIP = GetOutboundIP().String()
		pubIPRes, _ := pubip.Get()
		pubIP = pubIPRes.String()
		// Send a telegramMessage to notice server has been started.
		util.SendStarted(hostname, localIP, pubIP)
	}()

	cronJobs()
	handleServerStop()

	select {} // block forever
}

func cronJobs() {
	//_, _ = c.AddFunc("0 * * * *", updateMenu)
	c.Start()
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func handleServerStop() {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		if sig == syscall.SIGTERM || sig == syscall.SIGINT {
			util.SendNormalStopped(hostname, localIP, pubIP)
			os.Exit(0)
		}
	}()
}
