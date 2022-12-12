package service

import (
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

type CustomerIdProofTypeService interface {
	GetAll() ([]model.CustomerIdProofType, error)
	Get(name string) (*model.CustomerIdProofType, error)
	GetById(CustomerIdProofTypeId uint32) (*model.CustomerIdProofType, error)
}

type CustomerIdProofTypeServiceImpl struct {
	db *gorm.DB
}

func GetCustomerIdProofTypeService(db *gorm.DB) (service CustomerIdProofTypeService) {
	return CustomerIdProofTypeServiceImpl{db: db}
}

func (m CustomerIdProofTypeServiceImpl) GetAll() (result []model.CustomerIdProofType, err error) {
	result = []model.CustomerIdProofType{}
	err = m.db.Model(&model.CustomerIdProofType{}).Limit(1000).Find(&result).Error
	if err != nil {
		result = []model.CustomerIdProofType{}
	}

	return
}

func (m CustomerIdProofTypeServiceImpl) Get(name string) (result *model.CustomerIdProofType, err error) {
	var data = model.CustomerIdProofType{}
	err = m.db.Model(&model.CustomerIdProofType{}).Where("name = ?", name).First(&data).Error

	if err == nil {
		result = &data
	}

	return
}

func (m CustomerIdProofTypeServiceImpl) GetById(customerIdProofTypeId uint32) (result *model.CustomerIdProofType, err error) {
	var data = model.CustomerIdProofType{}
	err = m.db.Model(&model.CustomerIdProofType{}).Where("customer_id_proof_type_Id = ?", customerIdProofTypeId).First(&data).Error

	if err == nil {
		result = &data
	}

	return
}
