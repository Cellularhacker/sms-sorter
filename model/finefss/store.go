package finefss

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
	CName = "finefss"

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

func Create(ff *FineFss) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	if ff.ID.IsZero() {
		ff.ID = primitive.NewObjectID()
	}
	if ff.CategoryID.IsZero() {
		return fmt.Errorf("invalid 'category_id'")
	}
	ff.CreatedAt = time.Now().Unix()

	_, err := c.InsertOne(ctx, *ff)
	return err
}

func GetOneBy(by bson.M, opt ...*options.FindOneOptions) (*FineFss, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	sr := c.FindOne(ctx, by, opt...)
	if sr.Err() != nil {
		return nil, sr.Err()
	}

	ff := New()
	if err := sr.Decode(ff); err != nil {
		return nil, fmt.Errorf("sr.Decode(ff): %s", err)
	}

	return ff, nil
}

func GetBy(by bson.M, opt ...*options.FindOptions) ([]FineFss, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	ffs := make([]FineFss, 0)
	cur, err := c.Find(ctx, by, opt...)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ffs, nil
		}
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		ff := New()
		if err1 := cur.Decode(ff); err1 != nil {
			continue
		}

		ffs = append(ffs, *ff)
	}

	return ffs, nil
}

func CountBy(by bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	return c.CountDocuments(ctx, by)
}

func GetByPhoneNumber(phoneNumber string) (*FineFss, error) {
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
