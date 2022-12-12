package model

// Model Variable to be captured for the Model
type ModelVariable struct {
	ID                     uint32            `gorm:"column:model_variable_id;primary_key;auto_increment" json:"-"`
	ModelID                uint32            `gorm:"not null;uniqueIndex:model_variable_name" json:"-"`
	Name                   string            `gorm:"size:255;not null;uniqueIndex:model_variable_name" json:"name"` // Name of the Variable (Unique within the Model)
	DataType               uint16            `gorm:"not null" json:"dataType"`                                      // 1: Numeric Value 2: Enumerated Value
	CoefficientValue       float64           `gorm:"not null" json:"coefficientValue"`
	MeanValue              float64           `gorm:"not null" json:"meanValue"`
	StandardDeviationValue float64           `gorm:"not null" json:"standardDeviationValue"`
	Enumerations           []EnumerationItem `gorm:"-" json:"enumerations,omitempty"` // The enumerated values for the model variable (dataType: 2)
}

func (ModelVariable) TableName() string {
	return "model_variable"
}
