package model

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidHashValue = errors.New("invalid Hash Value")

func (u *User) Hash() error {
	passwordText := fmt.Sprintf("HASH:%d:%s", u.ID, u.Password)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwordText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = fmt.Sprintf("HASH:%s", hex.EncodeToString(passwordHash))

	return nil
}

func (u *User) VerifyPassword(password string) error {
	parts := strings.Split(u.Password, ":")
	if parts != nil && len(parts) != 2 {
		return ErrInvalidHashValue
	}
	passwordText := fmt.Sprintf("HASH:%d:%s", u.ID, password)
	hashValue, err := hex.DecodeString(parts[1])
	if err != nil {
		return ErrInvalidHashValue
	}

	newArray := make([]byte, len(hashValue)) // Avoid the array with 120 capacity and make a 60 byte length array
	copy(newArray, hashValue)

	return bcrypt.CompareHashAndPassword(hashValue, []byte(passwordText))
}
