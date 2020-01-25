package sms

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/mervick/aes-everywhere/go/aes256"
	"strings"
	"time"

	"sms-sorter/config"
)

var store Store

const (
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

	DefaultTel   = "010-3254-6909"
	DefaultFWTel = "010-6514-6909"

	TextTypeNotSorted        = 0
	TextTypeSMS              = 201
	TextTypeLMS              = 202
	TextTypeMMS              = 203
	TextTypeSpamNotSpecified = 600
	TextTypeSpamCommonNumber = 601
)

type Store interface {
	Create(a *Sms) error
	GetBy(bson.M) ([]Sms, error)
	GetByDesc(bson.M) ([]Sms, error)
	GetBySortLimit(query bson.M, sort string, limit int) ([]Sms, error)
	UpdateSet(what, set bson.M) error
	Delete(what bson.M, all bool) error
	CountBy(bson.M) (int, error)
	GetBySortLimitSkip(query bson.M, sort string, limit, skip int) ([]Sms, error)
}

func SetStore(s Store) {
	store = s
}

func GetByID(id bson.ObjectId) (*Sms, error) {
	as, err := store.GetByDesc(bson.M{KeyID: id})
	if err != nil {
		return nil, err
	}
	if len(as) > 0 {
		return &as[0], nil
	}
	return nil, nil
}

func DeleteBy(by bson.M) error {
	return store.Delete(by, true)
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
	s, _ := store.GetBy(bson.M{KeyMessageHash: hash})

	return len(s) > 0
}

/////////////////////////////////////////////////////
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
