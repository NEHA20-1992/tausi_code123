package model

// Contains the enumeration text presented to the user for selection and its internal value representation
type EnumerationItem struct {
	ModelVariableID uint32  `gorm:"not null;uniqueIndex:model_variable_item" json:"-"`
	Text            string  `gorm:"size:255;not null;uniqueIndex:model_variable_item" json:"text"`
	Value           float64 `gorm:"not null" json:"value"`
}

func (EnumerationItem) TableName() string {
	return "enumeration_item"
}
