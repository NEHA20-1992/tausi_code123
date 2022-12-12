package model

// Customer data element for the Credit Score Rating
type CustomerInformationItem struct {
	CustomerInformationID uint32 `gorm:"not null;primary_key" json:"-"`
	ModelVariableID       uint32 `gorm:"not null;primary_key" json:"-"`
	// ModelVariable         []ModelVariable `gorm:"-" json:"items,omitempty"`
	Name              string  `gorm:"-" json:"name,omitempty"`
	Value             float64 `gorm:"not null" json:"value"`
	PreprocessedValue float64 `gorm:"null" json:"preprocessedValue"`
}

func (CustomerInformationItem) TableName() string {
	return "customer_information_item"
}
