package validator

import (
	"errors"

	"github.com/NEHA20-1992/tausi_code/api/model"
)

var ErrFirstNameRequired = errors.New("first name required")
var ErrlastNameRequired = errors.New("last name required")
var ErrProofNumberRequired = errors.New("proof number required")
var ErrContactNumRequired = errors.New("contact number required")
var ErrCityRequired = errors.New("city required")

func ValidateCustomer(data *model.CustomerInformation) (err error) {

	if data == nil {
		return
	}
	if data.FirstName == "" {
		err = ErrFirstNameRequired
		return
	}
	if data.LastName == "" {
		err = ErrlastNameRequired
		return
	}

	if data.CustomerIdProofNumber == "" {
		err = ErrProofNumberRequired
		return
	}

	if data.ContactNumber == "" {
		err = ErrContactNumRequired
		return
	}
	if data.City == "" {
		err = ErrCityRequired
		return
	}

	return
}
