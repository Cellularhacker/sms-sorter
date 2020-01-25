package config

import (
	"log"
	"os"
	"strings"
)

const (
	ModeProduction  = "Production"
	ModeDevelopment = "Development"

	ServerName = "sms-saver"
)

var SqlURL = ""
var TokyoMongoURL = ""
var TokyoMongoAddr []string
var TokyoMongoUsername = ""
var TokyoMongoAuthDB = ""
var TokyoMongoPass = ""
var EncryptionSecret = ""
var PushoverSecret = ""
var PushoverDeviceID = ""

var MqttURL = ""
var MqttClientID = ""

var Mode = ModeDevelopment

func init() {
	EncryptionSecret = os.Getenv("SMS_ENCRYPT")
	if EncryptionSecret == "" {
		log.Fatalln("SMS_ENCRYPT missing")
	}

	MqttURL = os.Getenv("SMS_ADMIN_MQTT_URL")
	MqttClientID = os.Getenv("SMS_ADMIN_MQTT_CLIENT_ID")

	Mode = os.Getenv("SMS_MODE")
	if IsProductionMode() {
		log.Println("Running SMS_ in Production Mode")
	} else {
		log.Println("Running SMS_ in Development Mode")
	}

	if IsProductionMode() {
		SqlURL = os.Getenv("SMS_MYSQL_URL_APP_ENGINE")
	} else {
		SqlURL = os.Getenv("SMS_MYSQL_URL_IP")
	}

	TokyoMongoURL = os.Getenv("SMS_MONGO_URL")
	TokyoMongoAddr = strings.Split(os.Getenv("SMS_MONGO_ADDR"), ",")
	TokyoMongoUsername = os.Getenv("SMS_MONGO_USERNAME")
	TokyoMongoPass = os.Getenv("SMS_MONGO_PASS")
	TokyoMongoAuthDB = os.Getenv("SMS_MONGO_AUTH_DB")
	SqlURL = os.Getenv("SMS_MYSQL_URL")

	PushoverSecret = os.Getenv("SMS_PUSHOVER_SECRET")
	PushoverDeviceID = os.Getenv("SMS_PUSHOVER_DEVICE_ID")
}

func IsProductionMode() bool {
	return Mode == ModeProduction
}

const PhoneForwarded = "010-6514-6909"
const PhoneDirect = "010-3254-6909"