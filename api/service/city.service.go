package service

import (
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

type CityService interface {
	GetAll() ([]model.City, error)
	Get(name string) (*model.City, error)
	GetById(cityId uint32) (*model.City, error)
}

type CityServiceImpl struct {
	db *gorm.DB
}

func GetCityService(db *gorm.DB) CityService {
	return &CityServiceImpl{db: db}
}

func (m *CityServiceImpl) GetAll() (result []model.City, err error) {
	result = []model.City{}
	err = m.db.Model(&model.City{}).Limit(1000).Find(&result).Error
	if err != nil {
		result = nil
	}

	return
}

func (m CityServiceImpl) Get(name string) (result *model.City, err error) {
	var data = model.City{}
	err = m.db.Model(&model.City{}).Where("name = ?", name).First(&data).Error

	if err == nil {
		result = &data
	}

	return
}

func (m CityServiceImpl) GetById(cityId uint32) (result *model.City, err error) {
	var data = model.City{}
	err = m.db.Model(&model.City{}).Where("city_id = ?", cityId).First(&data).Error

	if err == nil {
		result = &data
	}

	return
}
