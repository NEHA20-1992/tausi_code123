package model

type CustomerIdProofType struct {
	ID   uint32 `gorm:"column:customer_id_proof_type_id;primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null" json:"name,omitempty"`
}

func (CustomerIdProofType) TableName() string {
	return "customer_id_proof_type"
}
