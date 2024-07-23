package sms

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/mervick/aes-everywhere/go/aes256"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"

	"sms-sorter/config"
)

var c *mongo.Collection

const (
	CName = "sms"

	KeyID          = "_id"
	KeyFromNumber  = "from_number"
	KeyContactName = "contact_name"
	KeyText        = "text"
	KeyOccurredAt  = "occurred_at"
	KeyToNumber    = "to_number"
	KeyReceivedAt  = "received_at"
	KeyTextType    = "text_type"
	KeyMessageHash = "message_hash"
	KeyCreatedAt   = "created_at"

	TextTypeNotSorted        = 0
	TextTypeSMS              = 201
	TextTypeLMS              = 202
	TextTypeMMS              = 203
	TextTypeSpamNotSpecified = 600
	TextTypeSpamCommonNumber = 601
)

func SetCollection(db *mongo.Database) {
	c = db.Collection(CName)
}

func Create(s *Sms) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	if s.ID.IsZero() {
		s.ID = primitive.NewObjectID()
	}
	s.CreatedAt = time.Now().Unix()

	_, err := c.InsertOne(ctx, *s)
	return err
}

func GetOneBy(by bson.M, opt ...*options.FindOneOptions) (*Sms, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	sr := c.FindOne(ctx, by, opt...)
	if sr.Err() != nil {
		return nil, sr.Err()
	}

	s := New()
	if err := sr.Decode(s); err != nil {
		return nil, fmt.Errorf("sr.Decode(s): %s", err)
	}

	return s, nil
}

func GetBy(by bson.M, opt ...*options.FindOptions) ([]Sms, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	ss := make([]Sms, 0)
	cur, err := c.Find(ctx, by, opt...)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ss, nil
		}
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		s := New()
		if err1 := cur.Decode(s); err1 != nil {
			continue
		}

		ss = append(ss, *s)
	}

	return ss, nil
}

func CountBy(by bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctx.Done()
	defer cancel()

	return c.CountDocuments(ctx, by)
}

func GetByID(id primitive.ObjectID) (*Sms, error) {
	return GetOneBy(bson.M{KeyID: id})
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

func IsExistHash(hash string) bool {
	s, _ := CountBy(bson.M{KeyMessageHash: hash})

	return s > 0
}

// ///////////////////////////////////////////////////
// January 20, 2020 at 01:18AM
func ParseFromJSONString(jsonStr string) (*Sms, error) {
	s := New()
	err := json.Unmarshal([]byte(jsonStr), s)
	if err != nil {
		return nil, err
	}
	if strings.Contains(s.Text, "FW\r\n") {
		s.ToNumber = config.PhoneForwarded
		s.Text = strings.ReplaceAll(s.Text, "FW\r\n", "")
	} else {
		s.ToNumber = config.PhoneDirect
	}

	// Parse Time
	at, err := time.Parse("January 2, 2006 at 03:04PM", s.OccurredAt)
	if err != nil {
		return nil, err
	}
	s.ReceivedAt = at.Add(-9 * time.Hour).Unix()
	s.MessageHash = GetSHA512(jsonStr)

	return s, err
}
