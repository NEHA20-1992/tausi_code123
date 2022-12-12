package service

import (
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

type CountryService interface {
	GetAll() ([]model.Country, error)
	Get(name string) (*model.Country, error)
	GetById(countryId uint32) (*model.Country, error)
}

type CountryServiceImpl struct {
	db *gorm.DB
}

func GetCountryService(db *gorm.DB) CountryService {
	return &CountryServiceImpl{db: db}
}

func (m *CountryServiceImpl) GetAll() (result []model.Country, err error) {
	result = []model.Country{}
	err = m.db.Model(&model.Country{}).Limit(1000).Find(&result).Error
	if err != nil {
		result = nil
	}

	return
}

func (m CountryServiceImpl) Get(name string) (result *model.Country, err error) {
	var data = model.Country{}
	err = m.db.Model(&model.Country{}).Where("name = ?", name).First(&data).Error

	if err == nil {
		result = &data
	}

	return
}

func (m CountryServiceImpl) GetById(countryId uint32) (result *model.Country, err error) {
	var data = model.Country{}
	err = m.db.Model(&model.Country{}).Where("country_id = ?", countryId).First(&data).Error

	if err == nil {
		result = &data
	}

	return
}
