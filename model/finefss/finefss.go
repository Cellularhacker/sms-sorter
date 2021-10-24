package finefss

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FineFss struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	CategoryID primitive.ObjectID `bson:"category_id" json:"category_id"`

	SequenceNo int `bson:"sequence_no" json:"sequence_no"`

	CompanyName                string   `bson:"company_name" json:"company_name"`
	CompanyNameEnglish         string   `bson:"company_name_english" json:"company_name_english"`
	CompanyPhoneNumber         string   `bson:"company_phone_number" json:"company_phone_number"`
	PhoneNumber                string   `bson:"phone_number" json:"phone_number"`
	PhoneNumberNote            string   `bson:"phone_number_note" json:"phone_number_note"`
	Address                    string   `bson:"address" json:"address"`
	Homepage                   string   `bson:"homepage" json:"homepage"`
	DepositInsurance           bool     `bson:"deposit_insurance" json:"deposit_insurance"`
	SupervisoryAuthority       string   `bson:"supervisory_authority" json:"supervisory_authority"`
	LicensedRegisteredBusiness []string `bson:"licensed_registered_business" json:"licensed_registered_business"`
	Etc                        string   `bson:"etc" json:"etc"`
	Note                       string   `bson:"note" json:"note"`

	CreatedAt int64 `bson:"created_at" json:"created_at"`
	UpdatedAt int64 `bson:"updated_at" json:"updated_at"`
}

func New() *FineFss {
	return &FineFss{}
}

func (c *FineFss) Upsert() error {
	ps, err := GetByPhoneNumber(c.PhoneNumber)
	if err != nil {
		return nil
	}

	// Create a new.
	if ps == nil {
		return Create(c)
	}

	// Check Update
	//if ps.Subject != c.Subject && ps.IsWhiteList == c.IsWhiteList {
	//	return store.UpdateSet(bson.M{KeyID: ps.ID}, bson.M{KeySubject: c.Subject})
	//}

	// Nothing to update.
	return nil
}
