package model

// City details
type City struct {
	ID   uint32 `gorm:"column:city_id;primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null" json:"name,omitempty"`
}

func (City) TableName() string {
	return "city"
}
