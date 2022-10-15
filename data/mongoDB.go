package data

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sms-sorter/config"
	"sms-sorter/util/logger"
	"time"
)

const (
	SmsDBName  = "sms"
	SpamDBName = "spamDB"

	CWhosCall = "whoscall"
)

var mongoDBClient *mongo.Client

func connectMongoDB() {
	logger.L.Info("Connecting to LoT Admin Mongo DB...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt := options.Client().ApplyURI(config.MongoURL)
	opt.SetAppName(config.ServerName)
	opt.SetReadPreference(readpref.SecondaryPreferred())
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		logger.L.Fatal("Failed to connect LoT Admin Mongo DB! =>", err)
	}
	err = client.Ping(context.Background(), readpref.SecondaryPreferred())
	if err != nil {
		logger.L.Fatal("Failed to ping LoT Admin Mongo DB! =>", err)
	}

	mongoDBClient = client
}

// InitMongo initiates our database session
func InitMongo() {
	connectMongoDB()
}
func GetSmsDB() *mongo.Database {
	return mongoDBClient.Database(SmsDBName)
}
func GetSpamDB() *mongo.Database {
	return mongoDBClient.Database(SpamDBName)
}
