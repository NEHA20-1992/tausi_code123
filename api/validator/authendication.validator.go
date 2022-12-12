package validator

import (
	"errors"

	"github.com/NEHA20-1992/tausi_code/api/model"
)

var ErrAuthEmailRequired = errors.New("Email is required")
var ErrAuthResetCodeRequired = errors.New("Reset code is required")
var ErrAuthNewPasswordRequired = errors.New("New password is required")
var ErrAuthConfirmPasswordRequired = errors.New("Confirm password is required")

func ValidateResetPassword(data *model.ResetPasswordRequest) (err error) {
	if data == nil {
		return
	}

	if data.ResetCode == "" {
		err = ErrAuthResetCodeRequired
		return
	}

	if data.Email == "" {
		err = ErrAuthEmailRequired
		return
	}

	if data.NewPassword == "" {
		err = ErrAuthNewPasswordRequired
		return
	}

	if data.ConfirmPassword == "" {
		err = ErrAuthConfirmPasswordRequired
		return
	}

	return
}
