package model

import "time"

// Credit Score Rating Model
type Model struct {
	ID              uint32          `gorm:"column:model_id;primary_key;auto_increment" json:"-"`
	Name            string          `gorm:"size:255;not null;uniqueIndex:company_model" json:"name"`
	CompanyID       uint32          `gorm:"not null;uniqueIndex:company_model" json:"-"`
	Description     string          `gorm:"size:255;not null" json:"description"`
	InterceptValue  float64         `gorm:"not null" json:"interceptValue"`
	Active          bool            `gorm:"not null" json:"active"`
	Variables       []ModelVariable `gorm:"-" json:"variables,omitempty"` // The list of model variables
	CustomerCount   int64           `gorm:"-" json:"customerCount"`
	CreatedBy       string          `gorm:"-" json:"createdBy,omitempty"`
	CreatedById     uint32          `gorm:"null" json:"-"`
	CreatedAt       time.Time       `gorm:"not null" json:"createdOn,omitempty"`
	LastUpdatedBy   string          `gorm:"-" json:"lastUpdatedBy,omitempty"`
	LastUpdatedById uint32          `gorm:"null" json:"-"`
	LastUpdatedAt   time.Time       `gorm:"null" json:"lastUpdatedOn,omitempty"`
}

func (Model) TableName() string {
	return "model"
}
