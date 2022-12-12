package model

import (
	"time"
)

type CompanyDataFile struct {
	ID                    uint32    `gorm:"column:file_id;primary_key;auto_increment" json:"id"`
	Name                  string    `gorm:"size:255;not null" json:"name,omitempty"`
	CompanyID             uint32    `gorm:"not null" json:"-"`
	CompanyDataFileTypeID uint32    `gorm:"not null" json:"-"`
	ModelID               uint32    `gorm:"null" json:"-"`
	Description           string    `gorm:"size:255;not null" json:"description,omitempty"`
	CreatedBy             string    `gorm:"-" json:"createdBy,omitempty"`
	CreatedById           uint32    `gorm:"null" json:"-"`
	CreatedAt             time.Time `gorm:"not null" json:"createdOn,omitempty"`
}

func (CompanyDataFile) TableName() string {
	return "company_data_file"
}

var companyDataFileArray = []*CompanyDataFile{}
