package mqtt

import (
	"fmt"
	json "github.com/json-iterator/go"
	"sms-sorter/config"
	"sms-sorter/util/logger"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type Status struct {
	ServerName string `json:"server_name"`
	Status     string `json:"status"`
	Timestamp  int64  `json:"timestamp"`
	Period     string `json:"period"`
}

var TopicStatus = func() string {
	return "Server/Status/SMS-SAVER"
}

var mqttClient mqtt.Client
var initialized = false

func Init() {
	logger.L.Info("MQTT Initializing...")

	if config.MqttURL == "" {
		logger.L.Fatal("SMS_ADMIN_MQTT_URL missing")
	}
	if config.MqttClientID == "" {
		logger.L.Fatal("SMS_ADMIN_MQTT_CLIENT_ID missing")
	}

	options := mqtt.NewClientOptions()
	options.AddBroker(config.MqttURL)
	options.SetAutoReconnect(true)
	options.SetClientID(config.MqttClientID)
	options.SetUsername(config.MqttClientID)

	mqttClient = mqtt.NewClient(options)
	t := mqttClient.Connect()
	if t.Wait() {
		if t.Error() != nil {
			logger.L.Fatal("Mqtt Error", t.Error())
		} else {
			logger.L.Info("MQTT Connected")
			initialized = true
		}
	}
}

func Begin() {
	SendStatus("startup", "startup")

	go func() {
		for range time.NewTicker(1 * time.Minute).C {
			SendStatus("alive", "1m")
		}
	}()
}

func SendStatus(status string, period string) {
	str := &Status{
		ServerName: config.ServerName,
		Status:     status,
		Timestamp:  time.Now().Unix(),
		Period:     period,
	}
	payload, err := json.Marshal(str)
	if err != nil {
		logger.L.Info("mqtt.SendStatus", err)
		return
	}

	if !initialized {
		logger.L.Info("[MQTT/Not-Initialized] payload => ", string(payload))
		return
	}
	publish(TopicStatus(), 1, false, payload)
	logger.L.Info("[MQTT] Sent status.")
}

func SendFailed(location string, err error, at time.Time) {
	status := fmt.Sprintf("{\"locatoin\":\"%s\",\"error\":\"%s\"", location, err.Error())
	SendFailedMessage(status, at)
}

func SendStarted(hostname, localIP, pubIP string) {
	resp := fmt.Sprintf("{\"hostname\":\"%s\", \"local_ip\":\"%s\", \"pub_ip\": \"%s\"}", hostname, localIP, pubIP)
	SendStatus(resp, "startup")
}

func SendStopped(hostname, localIP, pubIP string) {
	resp := fmt.Sprintf("{\"hostname\":\"%s\", \"local_ip\":\"%s\", \"pub_ip\": \"%s\"}", hostname, localIP, pubIP)
	SendStatus(resp, "normal-stop")
}

func SendFailedMessage(status string, at time.Time) {
	str := &Status{ServerName: config.ServerName, Status: status, Timestamp: at.Unix(), Period: "error"}
	payload, err := json.Marshal(str)
	if err != nil {
		logger.L.Info("mqtt.SendFailed - SendFailedMessage()", err)
		return
	}

	publish(TopicStatus(), 1, false, payload)
	logger.L.Info("[MQTT] SendFailedMessage() - Success.")
}

func publish(topic string, qos byte, retained bool, payload interface{}) {
	if !config.IsProductionMode() || !initialized {
		return
	}
	mqttClient.Publish(topic, qos, retained, payload)
}
