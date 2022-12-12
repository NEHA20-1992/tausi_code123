package model

import (
	"time"
)

// Company details
type Company struct {
	ID              uint32    `gorm:"column:company_id;primary_key;auto_increment" json:"-"`
	Name            string    `gorm:"size:255;not null" json:"name,omitempty"`
	Type_           string    `gorm:"-" json:"type,omitempty"`
	CompanyTypeID   uint32    `gorm:"not null" json:"-"`
	CompanyFileId   uint32    `gorm:"null" json:"-"`
	Address         string    `gorm:"size:255;not null" json:"address,omitempty"`
	RegionCounty    string    `gorm:"size:255;not null" json:"regionCounty,omitempty"`
	Country         string    `gorm:"-" json:"country,omitempty"`
	CountryId       uint32    `gorm:"not null" json:"-"`
	ContactNumber   string    `gorm:"size:25;not null" json:"contactNumber,omitempty"`
	EmailAddress    string    `gorm:"size:255;not null" json:"emailAddress,omitempty"`
	Active          bool      `gorm:"not null" json:"active"`
	CreatedBy       string    `gorm:"-" json:"createdBy,omitempty"`
	CreatedById     uint32    `gorm:"null" json:"-"`
	CreatedAt       time.Time `gorm:"not null" json:"createdOn,omitempty"`
	LastUpdatedBy   string    `gorm:"-" json:"lastUpdatedBy,omitempty"`
	LastUpdatedById uint32    `gorm:"null" json:"-"`
	LastUpdatedAt   time.Time `gorm:"null" json:"lastUpdatedOn,omitempty"`
}

type EnitityCount struct {
	CompanyCount int64
	UserCount    int64
	ProjectCount int64
}

type CRACount struct {
	UserCount    int64
	ProjectCount int64
}

func (Company) TableName() string {
	return "company"
}
