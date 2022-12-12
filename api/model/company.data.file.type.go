package model

type CompanyDataFileType struct {
	ID   uint32 `gorm:"column:company_data_file_type_id;primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null" json:"name,omitempty"`
}

func (CompanyDataFileType) TableName() string {
	return "company_data_file_type"
}
