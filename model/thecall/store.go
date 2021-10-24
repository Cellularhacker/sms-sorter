package thecall

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/mervick/aes-everywhere/go/aes256"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sms-sorter/config"
	"time"
)

var c *mongo.Collection

const (
	CName = "thecall"

	KeyID          = "_id"
	KeyPhoneNumber = "phone_number"
	KeySubject     = "subject"
	KeyIsWhiteList = "is_white_list"
	KeyCreatedAt   = "created_at"
	KeyUpdatedAt   = "updated_at"
	KeyAskCount    = "ask_count"
)

func SetCollection(db *mongo.Database) {
	c = db.Collection(CName)
}

func Create(tc *TheCall) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	if tc.ID.IsZero() {
		tc.ID = primitive.NewObjectID()
	}
	tc.CreatedAt = time.Now().Unix()

	_, err := c.InsertOne(ctx, *tc)
	return err
}

func GetOneBy(by bson.M, opt ...*options.FindOneOptions) (*TheCall, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	sr := c.FindOne(ctx, by, opt...)
	if sr.Err() != nil {
		return nil, sr.Err()
	}

	tc := New()
	if err := sr.Decode(tc); err != nil {
		return nil, fmt.Errorf("sr.Decode(tc): %s", err)
	}

	return tc, nil
}

func GetBy(by bson.M, opt ...*options.FindOptions) ([]TheCall, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	tcs := make([]TheCall, 0)
	cur, err := c.Find(ctx, by, opt...)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return tcs, nil
		}
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		tc := New()
		if err1 := cur.Decode(tc); err1 != nil {
			continue
		}

		tcs = append(tcs, *tc)
	}

	return tcs, nil
}

func UpdateSet(by, set bson.M, opt ...*options.UpdateOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	_, err := c.UpdateOne(ctx, by, bson.M{"$set": set}, opt...)
	return err
}

func CountBy(by bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	return c.CountDocuments(ctx, by)
}

func GetByID(id primitive.ObjectID) (*TheCall, error) {
	return GetOneBy(bson.M{KeyID: id})
}

func GetByPhoneNumber(phoneNumber string) (*TheCall, error) {
	return GetOneBy(bson.M{KeyPhoneNumber: phoneNumber})
}

func GetEncrypted(privateKey string) string {
	return aes256.Encrypt(privateKey, config.EncryptionSecret)
}

func GetDecrypted(encryptedText string) string {
	return aes256.Decrypt(encryptedText, config.EncryptionSecret)
}

func GetSHA512(content string) string {
	sha := sha512.New()                        // SHA512 해시 인스턴스 생성
	sha.Write([]byte(content))                 // 해시 인스턴스에 데이터 추가
	sha.Write([]byte(config.EncryptionSecret)) // 해시 인스턴스에 데이터 추가
	h2 := sha.Sum(nil)                         // 해시 인스턴스에 저장된 데이터의 SHA512 해시 값 추출

	return fmt.Sprintf("%x", h2)
}

func IsSpam(phoneNumber string) bool {
	cnt, _ := CountBy(bson.M{KeyPhoneNumber: phoneNumber, KeyIsWhiteList: false})

	return cnt > 0
}

func IsExist(phoneNumber string) bool {
	cnt, _ := CountBy(bson.M{KeyPhoneNumber: phoneNumber})

	return cnt > 0
}
