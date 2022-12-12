package service

import (
	"errors"

	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

var ErrVariableAlreadyExists = errors.New("variable already exist")
var ErrModelUsageFound = errors.New("model has been in use")
var ErrModelVariableItemNotFound = errors.New("model variable item not found")

type ModelVariableService interface {
	Create(claim *auth.AuthenticatedClaim, companyName string, modelName string, ModelVariable *model.ModelVariable) (*model.ModelVariable, error)
	Update(claim *auth.AuthenticatedClaim, companyName string, modelName string, existingVariableName string, ModelVariable *model.ModelVariable) (*model.ModelVariable, error)
	Delete(claim *auth.AuthenticatedClaim, companyName string, modelName string, existingVariableName string) (string, error)
	Get(claim *auth.AuthenticatedClaim, companyName string, modelName string, modelVariableName string) (*model.ModelVariable, error)
	GetAll(claim *auth.AuthenticatedClaim, companyName string, modelName string) ([]model.ModelVariable, error)
}

type ModelVariableServiceImpl struct {
	db           *gorm.DB
	modelService ModelService
}

func GetModelVariableService(db *gorm.DB) ModelVariableService {
	return &ModelVariableServiceImpl{
		db:           db,
		modelService: GetModelService(db)}
}

func (m *ModelVariableServiceImpl) Create(claim *auth.AuthenticatedClaim, companyName string, modelName string, modelVariable *model.ModelVariable) (result *model.ModelVariable, err error) {
	var modelValue *model.Model
	modelValue, err = m.modelService.Get(claim, companyName, modelName)
	if err != nil {
		return
	}
	modelVariable.ModelID = modelValue.ID

	var customerInformation model.CustomerInformation
	err = m.db.Model(&model.CustomerInformation{}).
		Where("model_id = ?", modelVariable.ModelID).
		Limit(1).
		FirstOrInit(&customerInformation, &model.CustomerInformation{ID: 0}).
		Error
	if err != nil {
		return
	}
	if customerInformation.ID > 0 {
		err = ErrModelUsageFound
		return
	}

	var existingRecord model.Model
	err = m.db.Model(&model.ModelVariable{}).
		Select("model_variable_id").
		Where("model_id = ? AND name = ?", modelValue.ID, modelVariable.Name).
		Find(&existingRecord).
		Error
	if err != nil || existingRecord.ID > 0 {
		err = ErrVariableAlreadyExists
		return
	}

	err = m.db.
		Model(&modelVariable).
		Create(&modelVariable).
		Error
	if err != nil {
		return
	}

	err = m.updateEnumerations(modelVariable)
	if err != nil {
		return
	}

	result = modelVariable

	return
}

func (m *ModelVariableServiceImpl) Update(claim *auth.AuthenticatedClaim, companyName string, modelName string, existingVariableName string, modelVariable *model.ModelVariable) (result *model.ModelVariable, err error) {
	var existingRecord *model.ModelVariable
	existingRecord, err = m.Get(claim, companyName, modelName, existingVariableName)
	if err != nil {
		return
	}

	modelVariable.ID = existingRecord.ID
	modelVariable.ModelID = existingRecord.ModelID

	var customerInformation model.CustomerInformation
	err = m.db.Model(&model.CustomerInformation{}).
		Where("model_id = ?", modelVariable.ModelID).
		Limit(1).
		FirstOrInit(&customerInformation, &model.CustomerInformation{ID: 0}).
		Error
	if err != nil {
		return
	}
	if customerInformation.ID > 0 {
		err = ErrModelUsageFound
		return
	}

	if modelVariable.Name != existingVariableName {
		var existingVariableRecord model.ModelVariable
		err = m.db.Model(&model.ModelVariable{}).
			Select("model_variable_id").
			Where("model_id = ? AND name = ?", modelVariable.ModelID, modelVariable.Name).
			Find(&existingVariableRecord).
			Error
		if err != nil || existingVariableRecord.ID > 0 {
			err = ErrVariableAlreadyExists
			return
		}
	}

	err = m.db.Debug().
		Model(&model.ModelVariable{}).
		Where("model_variable_id = ?", modelVariable.ID).
		Take(&model.ModelVariable{}).
		UpdateColumns(
			map[string]interface{}{
				"name":                     modelVariable.Name,
				"data_type":                modelVariable.DataType,
				"coefficient_value":        modelVariable.CoefficientValue,
				"mean_value":               modelVariable.MeanValue,
				"standard_deviation_value": modelVariable.StandardDeviationValue,
			}).
		Error
	if err != nil {
		return
	}

	err = m.updateEnumerations(modelVariable)
	if err != nil {
		return
	}

	result = modelVariable

	return
}

func (m *ModelVariableServiceImpl) Delete(claim *auth.AuthenticatedClaim, companyName string,
	modelName string, existingVariableName string) (result string, err error) {
	var existingRecord *model.ModelVariable
	existingRecord, err = m.Get(claim, companyName, modelName, existingVariableName)
	if err != nil {
		err = ErrModelVariableNotFound
		return
	}

	var customerInformation model.CustomerInformation
	err = m.db.Model(&model.CustomerInformation{}).
		Where("model_id = ?", existingRecord.ModelID).
		Limit(1).
		FirstOrInit(&customerInformation, model.CustomerInformation{ID: 0}).
		Error
	if err != nil {
		return
	}

	if customerInformation.ID > 0 {
		err = ErrModelUsageFound
		return
	}

	if len(existingRecord.Enumerations) > 0 {
		err = m.db.Debug().Exec("DELETE FROM enumeration_item WHERE model_variable_id", existingRecord.ID).Error
		if err != nil {
			err = ErrModelVariableItemNotFound
			return
		}

		result = "variable item deleted"
	}

	err = m.db.Debug().Exec("DELETE FROM model_variable WHERE model_id = ? AND model_variable_id = ?", existingRecord.ModelID, existingRecord.ID).Error
	if err != nil {
		return
	}
	result = "variable deleted"
	return
}

func (m *ModelVariableServiceImpl) updateEnumerations(modelVariable *model.ModelVariable) (err error) {
	err = m.db.Debug().
		Delete(&model.EnumerationItem{}, &model.EnumerationItem{ModelVariableID: modelVariable.ID}).
		Error
	if err != nil {
		return
	}

	for _, item := range modelVariable.Enumerations {
		item.ModelVariableID = modelVariable.ID
		err = m.db.
			Model(item).
			Create(item).
			Error
		if err != nil {
			return
		}
	}

	return
}

func (m *ModelVariableServiceImpl) Get(claim *auth.AuthenticatedClaim, companyName string, modelName string, modelVariableName string) (result *model.ModelVariable, err error) {
	var modelValue *model.Model
	modelValue, err = m.modelService.Get(claim, companyName, modelName)
	if err != nil {
		return
	}

	var existingRecord model.ModelVariable
	err = m.db.Model(&model.ModelVariable{}).
		Where("model_id = ? AND name = ?", modelValue.ID, modelVariableName).
		Find(&existingRecord).
		Error
	if err != nil {
		return
	}

	if existingRecord.DataType == 2 {
		eiList := []model.EnumerationItem{}
		err = m.db.
			Model(&model.EnumerationItem{}).
			Limit(1000).
			Where("model_variable_id = ?", existingRecord.ID).
			Find(&eiList).
			Error
		if err != nil {
			return
		}
		existingRecord.Enumerations = eiList
	}

	result = &existingRecord

	return
}

func (m *ModelVariableServiceImpl) GetAll(claim *auth.AuthenticatedClaim, companyName string, modelName string) (result []model.ModelVariable, err error) {
	var modelValue *model.Model
	modelValue, err = m.modelService.Get(claim, companyName, modelName)
	if err != nil {
		return
	}

	result, err = getAllModelVariables(m.db, claim, modelValue.ID)

	return
}

func getAllModelVariables(db *gorm.DB, claim *auth.AuthenticatedClaim, modelID uint32) (result []model.ModelVariable, err error) {
	result = []model.ModelVariable{}
	err = db.
		Model(&model.ModelVariable{}).
		Limit(1000).
		Where("model_id = ?", modelID).
		Find(&result).
		Error
	if err != nil {
		result = nil
	}

	var resultList []model.ModelVariable = make([]model.ModelVariable, len(result))
	for indxValue, aRecord := range result {
		var newRecord = &aRecord

		// Load the Enumerator Items
		if newRecord.DataType == 2 {
			eiList := []model.EnumerationItem{}
			err = db.Debug().
				Model(&model.EnumerationItem{}).
				Limit(1000).
				Where("model_variable_id = ? ", newRecord.ID).
				Find(&eiList).
				Error
			if err != nil {
				return
			}
			newRecord.Enumerations = eiList
		}

		resultList[indxValue] = *newRecord
	}

	result = resultList

	return
}
