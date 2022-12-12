package validator

import (
	"errors"

	"github.com/NEHA20-1992/tausi_code/api/model"
)

var ErrCompanyNameRequired = errors.New("company name required")
var ErrCompanyAddressRequired = errors.New("company address required")
var ErrRegionCountyRequired = errors.New("region county required")
var ErrEmailRequired = errors.New("email address required")
var ErrContactNumberRequired = errors.New("contact number required")
var ErrCompanyTypeRequired = errors.New("company type required")
var ErrCompanyStatusRequired = errors.New("company status required")
var ErrCompanyCountryRequired = errors.New("country required")

var ErrFileNameRequired = errors.New("File name required")
var ErrFileDescriptionRequired = errors.New("File description required")

func ValidateCompany(data *model.Company) (err error) {
	if data == nil {
		return
	}

	if data.Name == "" {
		err = ErrCompanyNameRequired
		return
	}

	if data.Address == "" {
		err = ErrCompanyAddressRequired
		return
	}

	if data.RegionCounty == "" {
		err = ErrRegionCountyRequired
		return
	}

	if data.EmailAddress == "" {
		err = ErrEmailRequired
		return
	}

	if data.ContactNumber == "" {
		err = ErrContactNumberRequired
	}

	if data.Type_ == "" {
		err = ErrCompanyTypeRequired
		return
	}
	return
}

func ValidateCompanyDataFile(data *model.CompanyDataFile) (err error) {
	if data == nil {
		return
	}

	if data.Name == "" {
		err = ErrFileNameRequired
		return
	}

	if data.Description == "" {
		err = ErrFileDescriptionRequired
		return
	}
	return
}
