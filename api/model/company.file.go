package model

// Company File details
type CompanyFile struct {
	ID       uint32 `gorm:"column:company_id;primary_key;auto_increment" json:"id"`
	Logo     []byte `gorm:"not null" json:"logo"`
	MimeType string `gorm:"null;size:512" json:"mimeType"`
}

func (CompanyFile) TableName() string {
	return "company_file"
}
