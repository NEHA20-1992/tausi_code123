package model

// Country details
type Country struct {
	ID   uint32 `gorm:"column:country_id;primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null" json:"name,omitempty"`
}

func (Country) TableName() string {
	return "country"
}
