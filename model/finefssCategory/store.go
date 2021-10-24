package finefssCategory

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
	CName = "finefssCategory"

	KeyID        = "_id"
	KeyKey       = "key"
	KeyValue     = "value"
	KeyText      = "text"
	KeyCreatedAt = "created_at"
)

func SetCollection(db *mongo.Database) {
	c = db.Collection(CName)
}

func Create(ffc *FineFssCategory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	if ffc.ID.IsZero() {
		ffc.ID = primitive.NewObjectID()
	}
	ffc.CreatedAt = time.Now().Unix()

	_, err := c.InsertOne(ctx, *ffc)
	return err
}

func GetOneBy(by bson.M, opt ...*options.FindOneOptions) (*FineFssCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	sr := c.FindOne(ctx, by, opt...)
	if sr.Err() != nil {
		return nil, sr.Err()
	}

	ffc := New()
	if err := sr.Decode(ffc); err != nil {
		return nil, fmt.Errorf("sr.Decode(ffc): %s", err)
	}

	return ffc, nil
}

func GetBy(by bson.M, opt ...*options.FindOptions) ([]FineFssCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	ffcs := make([]FineFssCategory, 0)
	cur, err := c.Find(ctx, by, opt...)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ffcs, nil
		}
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		ffc := New()
		if err1 := cur.Decode(ffc); err1 != nil {
			continue
		}

		ffcs = append(ffcs, *ffc)
	}

	return ffcs, nil
}

func CountBy(by bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	return c.CountDocuments(ctx, by)
}

func GetByID(id primitive.ObjectID) (*FineFssCategory, error) {
	return GetOneBy(bson.M{KeyID: id})
}

func GetByKey(key string) (*FineFssCategory, error) {
	return GetOneBy(bson.M{KeyKey: key})
}
func GetByValue(value string) (*FineFssCategory, error) {
	return GetOneBy(bson.M{KeyValue: value})
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
