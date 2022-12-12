package model

type UserRole struct {
	ID   uint32 `gorm:"column:user_role_id;primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}

func (UserRole) TableName() string {
	return "user_role"
}
