package data

import (
	"crypto/tls"
	"log"
	"net"
	"sms-sorter/config"

	"github.com/globalsign/mgo"
)

const (
	SmsDBName  = "sms"
	SpamDBName = "spamDB"

	CSms = "sms"

	CTheCall         = "thecall"
	CFineFss         = "finefss"
	CFineFssCategory = "finefssCategory"
	CWhosCall        = "whoscall"
)

var tokyoSession *mgo.Session

func createTokyoDBSession() {
	log.Println("MongoDB Initializing..")
	var err error

	if len(config.TokyoMongoAddr) > 1 {
		log.Println("TokyoMongoAddr: ", config.TokyoMongoAddr)
		log.Println("TokyoMongoUsername: ", config.TokyoMongoUsername)
		log.Println("TokyoMongoPass: ", config.TokyoMongoPass)
		log.Println("TokyoMongoAuth: ", config.TokyoMongoAuthDB)

		tlsConfig := &tls.Config{}
		dialInfo := &mgo.DialInfo{
			Addrs:    config.TokyoMongoAddr,
			Database: config.TokyoMongoAuthDB,
			Username: config.TokyoMongoUsername,
			Password: config.TokyoMongoPass,
		}
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		tokyoSession, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			log.Fatal("TokyoMongoDB failed: ", err)
		}

	} else {
		log.Println("TokyoMongoURL: ", config.TokyoMongoURL)

		tokyoSession, err = mgo.Dial(config.TokyoMongoURL)
		if err != nil {
			log.Fatal("TokyoMongoDB failed: ", err)
		}
	}

	log.Println("TokyoMongoDB Connected")
}

func getTokyoSession() *mgo.Session {
	if tokyoSession == nil {
		createTokyoDBSession()
	}
	return tokyoSession.Clone()
}

//Init initiates our database session
func InitTokyoMongo() {
	createTokyoDBSession()
	//ensureTokyoIndex()
}

func ensureTokyoIndex() {
	admSess := getTokyoSession()
	defer admSess.Close()
	//
	//liveUserCount := mgo.Index{Key: []string{"live_user_count"}, Background: true, Sparse: false}
	//appInstalledCount := mgo.Index{Key: []string{"app_installed_count"}, Background: true, Sparse: false}
	//appUsingUsersCount := mgo.Index{Key: []string{"app_using_users_count"}, Background: true, Sparse: false}
	//mau := mgo.Index{Key: []string{"mau"}, Background: true, Sparse: false}
	//wau := mgo.Index{Key: []string{"wau"}, Background: true, Sparse: false}
	//dau := mgo.Index{Key: []string{"dau"}, Background: true, Sparse: false}
	//activatedAPIKeysCount := mgo.Index{Key: []string{"activated_api_keys_count"}, Background: true, Sparse: false}
	//dailyTradingVolumeInUSD := mgo.Index{Key: []string{"daily_trading_volume_in_usd"}, Background: true, Sparse: false}
	//currentAccumulatedBalanceInUSD := mgo.Index{Key: []string{"current_accumulated_balance_in_usd"}, Background: true, Sparse: false}
	//_type := mgo.Index{Key: []string{"type"}, Background: true, Sparse: false}
	//createdAt := mgo.Index{Key: []string{"-created_at"}, Background: true, Sparse: false}
	//lastTokyoStats := mgo.Index{Key: []string{"-created_at", "type"}, Background: true, Sparse: false}

	////TokyoStats
	//checkTokyoIndexError(admSess.DB(TokyoDBName).C(CTokyoStats).EnsureIndex(_type))
	//checkTokyoIndexError(admSess.DB(TokyoDBName).C(CTokyoStats).EnsureIndex(createdAt))
	//checkTokyoIndexError(admSess.DB(TokyoDBName).C(CTokyoStats).EnsureIndex(lastTokyoStats))
}

func checkTokyoIndexError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
