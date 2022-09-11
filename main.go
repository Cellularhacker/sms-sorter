package main

import (
	"github.com/robfig/cron"
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
		logger.L.Fatal(err)
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
