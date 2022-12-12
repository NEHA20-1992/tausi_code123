package service

import (
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

type UserRoleService interface {
	GetAll() ([]model.UserRole, error)
	Get(name string) (*model.UserRole, error)
}

type UserRoleServiceImpl struct {
	db *gorm.DB
}

func GetUserRoleService(db *gorm.DB) (service UserRoleService) {
	return UserRoleServiceImpl{db: db}
}

func (m UserRoleServiceImpl) GetAll() (result []model.UserRole, err error) {
	result = []model.UserRole{}
	err = m.db.Model(&model.UserRole{}).Limit(1000).Find(&result).Error

	if err != nil {
		result = []model.UserRole{}
	}

	return
}

func (m UserRoleServiceImpl) Get(name string) (result *model.UserRole, err error) {
	var data = model.UserRole{}
	err = m.db.Model(&model.UserRole{}).Where("name = ?", name).First(&data).Error

	if err == nil {
		result = &data
	}

	return
}
