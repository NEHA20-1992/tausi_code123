package model

type CompanyType struct {
	ID   uint32 `gorm:"column:company_type_id;primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null" json:"name,omitempty"`
}

func (CompanyType) TableName() string {
	return "company_type"
}
