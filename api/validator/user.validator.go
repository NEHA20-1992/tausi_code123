package validator

import (
	"errors"

	"github.com/NEHA20-1992/tausi_code/api/model"
)

var ErrUserFirstNameRequired = errors.New("first name is required")
var ErrUserLastNameRequired = errors.New("last name is required")

func ValidateUser(data *model.User) (err error) {
	if data == nil {
		return
	}

	if data.FirstName == "" {
		err = ErrUserFirstNameRequired
		return
	}

	if data.LastName == "" {
		err = ErrUserLastNameRequired
		return
	}

	return
}
