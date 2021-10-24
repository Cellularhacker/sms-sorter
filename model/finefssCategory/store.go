package finefssCategory

import (
	"crypto/sha512"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/mervick/aes-everywhere/go/aes256"
	"sms-sorter/config"
)

var store Store

const (
	KeyID        = "_id"
	KeyKey       = "key"
	KeyValue     = "value"
	KeyText      = "text"
	KeyCreatedAt = "created_at"
)

type Store interface {
	Create(a *FineFssCategory) error
	GetBy(bson.M) ([]FineFssCategory, error)
	GetByDesc(bson.M) ([]FineFssCategory, error)
	GetBySortLimit(query bson.M, sort string, limit int) ([]FineFssCategory, error)
	UpdateSet(what, set bson.M) error
	Delete(what bson.M, all bool) error
	CountBy(bson.M) (int, error)
	GetBySortLimitSkip(query bson.M, sort string, limit, skip int) ([]FineFssCategory, error)
}

func SetStore(s Store) {
	store = s
}

func GetOneBy(by bson.M) (*FineFssCategory, error) {
	ts, err := store.GetBy(by)
	if err != nil {
		return nil, err
	}
	if len(ts) > 0 {
		return &ts[0], nil
	}
	return nil, nil
}

func GetByID(id bson.ObjectId) (*FineFssCategory, error) {
	as, err := store.GetByDesc(bson.M{KeyID: id})
	if err != nil {
		return nil, err
	}
	if len(as) > 0 {
		return &as[0], nil
	}
	return nil, nil
}

func GetByKey(key string) (*FineFssCategory, error) {
	return GetOneBy(bson.M{KeyKey: key})
}
func GetByValue(value string) (*FineFssCategory, error) {
	return GetOneBy(bson.M{KeyValue: value})
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
