package service

import (
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

type CompanyTypeService interface {
	GetAll() ([]model.CompanyType, error)
	Get(name string) (*model.CompanyType, error)
	GetById(companyTypeId uint32) (*model.CompanyType, error)
}

type CompanyTypeServiceImpl struct {
	db *gorm.DB
}

func GetCompanyTypeService(db *gorm.DB) (service CompanyTypeService) {
	return CompanyTypeServiceImpl{db: db}
}

func (m CompanyTypeServiceImpl) GetAll() (result []model.CompanyType, err error) {
	result = []model.CompanyType{}
	err = m.db.Model(&model.CompanyType{}).Limit(1000).Find(&result).Error
	if err != nil {
		result = []model.CompanyType{}
	}

	return
}

func (m CompanyTypeServiceImpl) Get(name string) (result *model.CompanyType, err error) {
	var data = model.CompanyType{}
	err = m.db.Model(&model.CompanyType{}).Where("name = ?", name).First(&data).Error

	if err == nil {
		result = &data
	}

	return
}

func (m CompanyTypeServiceImpl) GetById(companyTypeId uint32) (result *model.CompanyType, err error) {
	var data = model.CompanyType{}
	err = m.db.Model(&model.CompanyType{}).Where("company_type_id = ?", companyTypeId).First(&data).Error

	if err == nil {
		result = &data
	}

	return
}
