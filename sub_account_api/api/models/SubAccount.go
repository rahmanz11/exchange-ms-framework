package models

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

// SubAccount model
type SubAccount struct {
	SubAccountId   uuid.UUID       `gorm:"primary_key;unique" json:"sub_account_id"`
	AccountNumber  string          `gorm:"size:255;not null;unique" json:"account_number"`
	Balance        float32         `json:"balance" gorm:"default:0.0"`
	Status         string          `gorm:"size:100;not null" json:"status"`
	Credential     string          `gorm:"size:250;not null" json:"credential"`
	LinkedAccounts []LinkedAccount `json:"linked_accounts" gorm:"foreignKey:SubAccountId"`
	CreatedAt      time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type LinkedAccount struct {
	AccountNumber string    `gorm:"primary_key;unique" json:"account_number"`
	SubAccountId  uuid.UUID `gorm:"unique" json:"sub_account_id"`
}

// Hash password using bcrypt
func Hash(credential string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(credential), bcrypt.DefaultCost)
}

// Verify a hashed password
func VerifyCredential(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeCreate Gorm Hook. This runs automatically before creating a new sub_account
func (subAccount *SubAccount) BeforeCreate(db *gorm.DB) error {
	// hashed the credential
	hashedCredential, err := Hash(subAccount.Credential)
	if err != nil {
		return err
	}
	// set the hashed credential and a new uuid for the sub_account
	subAccount.SubAccountId = uuid.New()
	subAccount.Credential = string(hashedCredential)
	return nil
}

// BeforeUpdate Gorm Hook. This runs automatically before updating a sub_account
func (subAccount *SubAccount) BeforeUpdate(db *gorm.DB) error {
	subAccount.UpdatedAt = time.Now()
	return nil
}

// Prepare the model for database insertion
func (subAccount *SubAccount) Prepare() {
	subAccount.CreatedAt = time.Now()
	subAccount.UpdatedAt = time.Now()
}

// Validate the model
func (subAccount *SubAccount) Validate(action string) error {
	// validate the model before creating a new sub_account
	switch strings.ToLower(action) {
	case "update":
		if subAccount.Balance <= 0 {
			return errors.New("balance invalid")
		}
		if subAccount.Credential == "" {
			return errors.New("required credential")
		}
		return nil

	case "login":
		if subAccount.Credential == "" {
			return errors.New("required credential")
		}
		if subAccount.AccountNumber == "" {
			return errors.New("required account number")
		}
		return nil

	default:
		if subAccount.Credential == "" {
			return errors.New("required credential")
		}
		if subAccount.AccountNumber == "" {
			return errors.New("required account number")
		}
	}
	return nil
}

// CreateSubAccount Create a new sub_account
func (subAccount *SubAccount) CreateSubAccount(db *gorm.DB) (*SubAccount, error) {
	// validate the model before creating a new sub_account
	subAccount.Prepare()
	err := subAccount.Validate("")
	if err != nil {
		return nil, err
	}
	// create a new sub_account
	err = db.Debug().Create(&subAccount).Error
	if err != nil {
		return &SubAccount{}, err
	}
	return subAccount, nil
}

// Find a sub_account by account number
func (subAccount *SubAccount) FindSubAccountByAccountNumber(db *gorm.DB, accountNumber string) (*SubAccount, error) {
	// find a sub_account by account number
	err := db.Debug().Model(&SubAccount{}).Preload("LinkedAccounts").Where("account_number = ?", accountNumber).Take(&subAccount).Error
	if err != nil {
		return &SubAccount{}, err
	}
	// if not found return error
	if gorm.ErrRecordNotFound == err {
		return &SubAccount{}, errors.New("sub account not found")
	}
	// return the sub_account
	return subAccount, err
}

// Find a sub_account by account id
func (subAccount *SubAccount) FindSubAccountByAccountId(db *gorm.DB, accountId string) (*SubAccount, error) {
	// find a sub_account by account id
	err := db.Debug().Model(&SubAccount{}).Preload("LinkedAccounts").Where("sub_account_id = ?", accountId).Take(&subAccount).Error
	if err != nil {
		return &SubAccount{}, err
	}
	// if not found return error
	if gorm.ErrRecordNotFound == err {
		return &SubAccount{}, errors.New("sub account not found")
	}
	// return the sub_account
	return subAccount, err
}

// Update a sub_account
func (subAccount *SubAccount) UpdateSubAccount(db *gorm.DB) (*SubAccount, error) {
	// Update the sub account
	// Only balance and status update is allowed for seurity reasons
	err := db.Debug().Model(&SubAccount{}).Where("account_number = ?", subAccount.AccountNumber).Updates(SubAccount{
		Balance: subAccount.Balance,
		Status:  subAccount.Status,
	}).Error
	if err != nil {
		return &SubAccount{}, err
	}
	return subAccount, nil
}

// Link an account to a sub_account
func (subAccount *SubAccount) LinkAccount(db *gorm.DB, accountNumber string) (*SubAccount, error) {
	// Create a new linked account
	// It links automatically using Foreign Key
	err := db.Debug().Model(&LinkedAccount{}).Create(&LinkedAccount{
		AccountNumber: accountNumber,
		SubAccountId:  subAccount.SubAccountId,
	}).Error
	if err != nil {
		return &SubAccount{}, err
	}
	return subAccount, nil
}
