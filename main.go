package main

import (
	"github.com/robfig/cron"
	"log"
	"net"
	"os"
	"os/signal"
	"sms-sorter/config"
	"sms-sorter/data"
	"sms-sorter/model/finefss"
	"sms-sorter/model/finefssCategory"
	"sms-sorter/model/sms"
	"sms-sorter/model/thecall"
	"sms-sorter/service/telegram"
	"sms-sorter/util"
	"sms-sorter/util/logger"
	"syscall"

	"github.com/chyeh/pubip"
)

var c *cron.Cron
var localIP = ""
var pubIP = ""
var hostname = ""

func init() {
	logger.Init(config.IsProductionMode())
	// Initializing Data...
	data.Init()

	c = cron.New()

	// Set Collection
	sms.SetCollection(data.GetSmsDB())
	thecall.SetCollection(data.GetSmsDB())
	finefss.SetCollection(data.GetSpamDB())
	finefssCategory.SetCollection(data.GetSpamDB())
	//finefss.SetStore(data.NewFineFssStore())
	//finefssCategory.SetStore(data.NewFineFssCategoryStore())

	// Service...
	telegram.Init()

}

func main() {
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
