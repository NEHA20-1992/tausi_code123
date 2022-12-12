package service

import (
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

type CompanyFileTypeService interface {
	Get(id uint32) (*model.CompanyDataFileType, error)
	GetById(companyFileTypeId uint32) (*model.CompanyDataFileType, error)
}

type CompanyFileTypeServiceImpl struct {
	db *gorm.DB
}

func GetCompanyFileTypeService(db *gorm.DB) (service CompanyFileTypeService) {
	//return CompanyFileTypeServiceImpl{db: db}

	return
}

func (m CompanyFileTypeServiceImpl) Get(id uint32) (result *model.CompanyDataFileType, err error) {
	var data = model.CompanyDataFileType{}
	err = m.db.Model(&model.CompanyDataFileType{}).Where("company_data_file_type_id = ?", id).First(&data).Error

	if err == nil {
		result = &data
	}

	return
}
