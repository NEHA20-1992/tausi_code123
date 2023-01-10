package model

import "time"

// Customer Information to be analyzed for the Credit Score Rating
type CustomerInformation struct {
	ID                             uint32                    `gorm:"column:customer_information_id;primary_key;auto_increment" json:"id"`
	ModelID                        uint32                    `gorm:"not null" json:"-"`
	FirstName                      string                    `gorm:"size:255;not null" json:"firstName,omitempty"`
	LastName                       string                    `gorm:"size:255;not null" json:"lastName,omitempty"`
	CustomerIdProofNumber          string                    `gorm:"size:255;not null" json:"customerIdProofNumber,omitempty"`
	CustomerIDProofTypeID          uint32                    `gorm:"not null" json:"-"`
	CustomerIDProofType            string                    `gorm:"-" json:"type,omitempty"`
	ContactNumber                  string                    `gorm:"size:25;not null" json:"contactNumber,omitempty"`
	City                           string                    `gorm:"size:255;not null" json:"city,omitempty"`
	AccountID                      string                    `gorm:"size:255;not null"`
	ProbabilityOfDefaultPercentage float64                   `gorm:"null" json:"probabilityOfDefaultPercentage,omitempty"`
	GroupScore                     float64                   `gorm:"null" json:"groupScore,omitempty"`
	Items                          []CustomerInformationItem `gorm:"-" json:"items"` // The values for the Credit Score Rating
	CreatedBy                      string                    `gorm:"-" json:"createdBy,omitempty"`
	CompanyID                      uint32                    `gorm:"not null" json:"-"`
	CreatedById                    uint32                    `gorm:"null" json:"-"`
	CreatedAt                      time.Time                 `gorm:"not null" json:"createdOn,omitempty"`
}

type CustomerCreditScore struct {
	ID                             uint32  `gorm:"customer_information_id;" json:"id"`
	ModelID                        uint32  `gorm:"null" json:"-"`
	FirstName                      string  `gorm:"size:255;not null" json:"firstName"`
	LastName                       string  `gorm:"size:255;not null" json:"lastName"`
	ContactNumber                  string  `gorm:"size:25;not null" json:"contactNumber"`
	City                           string  `gorm:"size:255;not null" json:"city"`
	ProbabilityOfDefaultPercentage float64 `gorm:"null" json:"probabilityOfDefaultPercentage"`
	GroupScore                     float64 `gorm:"null" json:"groupScore"`
	CreditScore                    string  `gorm:"size:25;not null" json:"creditScore"`
}

type CustomerFilterRequest struct {
	CompanyName   string
	ModelName     string
	City          string
	MinGroupScore float64
	MaxGroupScore float64
	MinPercentage float64
	MaxPercentage float64
	Sort          string
	PageNumber    uint32
	Size          uint32
	CID           uint32
	MID           uint32
}

func (CustomerInformation) TableName() string {
	return "customer_information"
}
