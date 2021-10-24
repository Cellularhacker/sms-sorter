package finefss

import (
	"crypto/sha512"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/mervick/aes-everywhere/go/aes256"
	"sms-sorter/config"
)

var store Store

const (
	KeyID          = "_id"
	KeyPhoneNumber = "phone_number"
	KeySubject     = "subject"
	KeyIsWhiteList = "is_white_list"
	KeyCreatedAt   = "created_at"
	KeyUpdatedAt   = "updated_at"
	KeyAskCount    = "ask_count"
)

type Store interface {
	Create(a *FineFss) error
	GetBy(bson.M) ([]FineFss, error)
	GetByDesc(bson.M) ([]FineFss, error)
	GetBySortLimit(query bson.M, sort string, limit int) ([]FineFss, error)
	UpdateSet(what, set bson.M) error
	Delete(what bson.M, all bool) error
	CountBy(bson.M) (int, error)
	GetBySortLimitSkip(query bson.M, sort string, limit, skip int) ([]FineFss, error)
}

func SetStore(s Store) {
	store = s
}

func GetOneBy(by bson.M) (*FineFss, error) {
	ts, err := store.GetBy(by)
	if err != nil {
		return nil, err
	}
	if len(ts) > 0 {
		return &ts[0], nil
	}
	return nil, nil
}

func GetByID(id bson.ObjectId) (*FineFss, error) {
	as, err := store.GetByDesc(bson.M{KeyID: id})
	if err != nil {
		return nil, err
	}
	if len(as) > 0 {
		return &as[0], nil
	}
	return nil, nil
}

func GetByPhoneNumber(phoneNumber string) (*FineFss, error) {
	return GetOneBy(bson.M{KeyPhoneNumber: phoneNumber})
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

func IsSpam(phoneNumber string) bool {
	cnt, _ := store.CountBy(bson.M{KeyPhoneNumber: phoneNumber, KeyIsWhiteList: false})

	return cnt > 0
}

func IsExist(phoneNumber string) bool {
	cnt, _ := store.CountBy(bson.M{KeyPhoneNumber: phoneNumber})

	return cnt > 0
}
