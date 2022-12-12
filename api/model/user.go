package model

import (
	"time"
)

// User details
type User struct {
	ID              uint32     `gorm:"column:user_id;primary_key;auto_increment" json:"-"`
	FirstName       string     `gorm:"size:255;not null" json:"firstName"`
	LastName        string     `gorm:"size:255;not null" json:"lastName"`
	Email           string     `gorm:"size:255;not null;unique" json:"email"`
	Password        string     `gorm:"size:255;null" json:"password,omitempty"`
	ResetCode       string     `gorm:"size:255;null" json:"resetCode,omitempty"`
	ContactNumber   string     `gorm:"size:25;not null" json:"contactNumber"`
	CompanyId       uint32     `gorm:"null" json:"-"`
	CompanyName     string     `gorm:"null" json:"company"`
	Roles           []string   `gorm:"-" json:"roles,omitempty"`
	UserRoles       []UserRole `gorm:"many2many:user_user_role;" json:"-"`
	Admin           bool       `gorm:"-" json:"admin"`
	EmailVerified   bool       `gorm:"default:false;not null" json:"emailVerified"` // Set to true if the user's email has been verified.
	Active          bool       `gorm:"default:false;not null" json:"active"`
	CreatedBy       string     `gorm:"-" json:"createdBy,omitempty"`
	CreatedById     uint32     `gorm:"null" json:"-"`
	CreatedAt       time.Time  `gorm:"not null" json:"createdOn,omitempty"`
	LastUpdatedBy   string     `gorm:"-" json:"lastUpdatedBy,omitempty"`
	LastUpdatedById uint32     `gorm:"null" json:"-"`
	LastUpdatedAt   time.Time  `gorm:"null" json:"lastUpdatedOn,omitempty"`
}

// chage password request details
type ChangePasswordRequest struct {
	Password        string `gorm:"size:255;null" json:"password,omitempty"`
	NewPassword     string `gorm:"size:255;null" json:"newPassword,omitempty"`
	ConfirmPassword string `gorm:"size:255;null" json:"confirmPassword,omitempty"`
}

// chage password request details
type ForgetPasswordRequest struct {
	CompanyName string `gorm:"size:255;null" json:"companyName,omitempty"`
	Email       string `gorm:"size:255;null" json:"email,omitempty"`
}

func (User) TableName() string {
	return "user"
}
