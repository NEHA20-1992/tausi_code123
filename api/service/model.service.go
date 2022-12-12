package service

import (
	"errors"
	"time"

	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

var ErrModelAlreadyExists = errors.New("project already exist")

type ModelService interface {
	Create(claim *auth.AuthenticatedClaim, companyName string, model *model.Model) (*model.Model, error)
	Update(claim *auth.AuthenticatedClaim, companyName string, existingModelName string, model *model.Model) (*model.Model, error)
	Get(claim *auth.AuthenticatedClaim, companyName string, modelName string) (*model.Model, error)
	GetById(claim *auth.AuthenticatedClaim, companyID uint32, modelID uint32) (*model.Model, error)

	GetDetails(claim *auth.AuthenticatedClaim, companyName string, modelName string) (*model.Model, error)
	GetAll(claim *auth.AuthenticatedClaim, companyName string) ([]model.Model, error)
	GetCustomerCount(claim *auth.AuthenticatedClaim, companyName string, modelName string) (result int64, err error)
}

type ModelServiceImpl struct {
	db             *gorm.DB
	companyService CompanyService
}

func GetModelService(db *gorm.DB) ModelService {
	return &ModelServiceImpl{
		db:             db,
		companyService: GetCompanyService(db)}
}

func (m *ModelServiceImpl) Create(claim *auth.AuthenticatedClaim, companyName string, modelValue *model.Model) (result *model.Model, err error) {
	if modelValue == nil {
		return
	}

	var currentUserId uint32 = 1
	currentUserId, err = m.prepare(claim, companyName, modelValue)
	if err != nil {
		return
	}
	modelValue.CreatedById = currentUserId

	var existingRecord model.Model
	err = m.db.Model(&model.Model{}).
		Select("model_id").
		Where("company_id = ? AND name = ?", modelValue.CompanyID, modelValue.Name).
		Find(&existingRecord).
		Error
	if err != nil || existingRecord.ID > 0 {
		err = ErrModelAlreadyExists
		return
	}

	err = m.db.
		Model(&modelValue).
		Omit("last_updated_by_id", "last_updated_at").
		Create(&modelValue).
		Error
	if err != nil {
		return
	}

	// This is the display the updated user
	createdRecord := model.Model{}
	err = m.db.Debug().
		Model(&createdRecord).
		Where("model_id = ?", modelValue.ID).
		Take(&createdRecord).
		Error
	if err != nil {
		return
	}

	result = &createdRecord
	err = m.updateMeta(claim, nil, result)
	if err != nil {
		return
	}
	return
}

func (m *ModelServiceImpl) Update(claim *auth.AuthenticatedClaim, companyName string, existingModelName string, modelValue *model.Model) (result *model.Model, err error) {
	var currentUserId uint32 = 1

	currentUserId, err = m.prepare(claim, companyName, modelValue)
	if err != nil {
		return
	}
	modelValue.CreatedById = currentUserId

	if modelValue.Name != existingModelName {
		var existingModelRecord model.Model
		err = m.db.Model(&model.Model{}).
			Select("model_id").
			Where("company_id = ? AND name = ?", modelValue.CompanyID, modelValue.Name).
			Find(&existingModelRecord).
			Error
		if err != nil || existingModelRecord.ID > 0 {
			err = ErrModelAlreadyExists
			return
		}
	}

	var existingRecord model.Model
	err = m.db.
		Model(&model.Model{}).
		Select("model_id").
		Where("company_id = ? AND name = ?", modelValue.CompanyID, existingModelName).
		First(&existingRecord).
		Error
	if err != nil {
		return
	}

	modelValue.ID = existingRecord.ID

	err = m.db.Debug().
		Model(&model.Model{}).
		Where("model_id = ?", modelValue.ID).
		Take(&model.Model{}).
		UpdateColumns(
			map[string]interface{}{
				"name":               modelValue.Name,
				"description":        modelValue.Description,
				"intercept_value":    modelValue.InterceptValue,
				"active":             modelValue.Active,
				"last_updated_by_id": currentUserId,
				"last_updated_at":    time.Now(),
			}).
		Error
	if err != nil {
		return
	}

	// This is the display the updated user
	updatedRecord := model.Model{}
	err = m.db.Debug().
		Model(&updatedRecord).
		Where("model_id = ?", modelValue.ID).
		Take(&updatedRecord).
		Error
	if err != nil {
		return
	}
	result = &updatedRecord
	err = m.updateMeta(claim, nil, result)
	if err != nil {
		return
	}

	return
}

func (m *ModelServiceImpl) GetCustomerCount(claim *auth.AuthenticatedClaim, companyName string, modelName string) (result int64, err error) {
	var existingRecord *model.Model
	existingRecord, err = m.Get(claim, companyName, modelName)
	if err != nil {
		err = ErrModelVariableNotFound
		return
	}

	if existingRecord.ID > 0 {
		var customerInformations = []model.CustomerInformation{}
		err = m.db.Where("model_id = ?", existingRecord.ID).
			Find(&customerInformations).
			Count(&result).
			Error
		if err != nil {
			return
		}
	}

	return
}

func (m *ModelServiceImpl) Get(claim *auth.AuthenticatedClaim, companyName string, modelName string) (result *model.Model, err error) {
	var company *model.Company
	company, err = m.companyService.Get(claim, companyName)
	if err != nil {
		return
	}

	var modelValue model.Model
	err = m.db.
		Model(&modelValue).
		Where("company_id = ? AND name = ?", company.ID, modelName).
		First(&modelValue).
		Error
	if err != nil {
		return
	}
	err = m.updateMeta(claim, nil, &modelValue)
	if err != nil {
		return
	}

	result = &modelValue

	return
}
func (m *ModelServiceImpl) GetById(claim *auth.AuthenticatedClaim, companyID uint32, modelID uint32) (result *model.Model, err error) {
	var company *model.Company
	company, err = m.companyService.GetById(claim, companyID)
	if err != nil {
		return
	}

	var modelValue model.Model
	err = m.db.
		Model(&modelValue).
		Where("company_id = ? AND model_id = ?", company.ID, modelID).
		First(&modelValue).
		Error
	if err != nil {
		return
	}
	err = m.updateMeta(claim, nil, &modelValue)
	if err != nil {
		return
	}

	result = &modelValue

	return
}

func (m *ModelServiceImpl) GetDetails(claim *auth.AuthenticatedClaim, companyName string, modelName string) (result *model.Model, err error) {
	result, err = m.Get(claim, companyName, modelName)
	if err != nil {
		return
	}

	mvList, err := getAllModelVariables(m.db, claim, result.ID)

	result.Variables = mvList

	return
}

func (m *ModelServiceImpl) GetAll(claim *auth.AuthenticatedClaim, companyName string) (result []model.Model, err error) {
	var company *model.Company
	company, err = m.companyService.Get(claim, companyName)
	if err != nil {
		return
	}

	result = []model.Model{}
	err = m.db.
		Model(&model.Model{}).
		Limit(1000).
		Where("company_id = ?", company.ID).
		Find(&result).
		Error
	if err != nil {
		result = nil
	}

	var userMap map[int](*string) = make(map[int]*string)
	var resultList []model.Model = make([]model.Model, len(result))
	for indxValue, aRecord := range result {
		var newRecord = &aRecord
		var count int64
		if aRecord.ID > 0 {
			var customerInformations = []model.CustomerInformation{}
			err = m.db.Where("model_id = ?", aRecord.ID).
				Find(&customerInformations).
				Count(&count).
				Error
			if err != nil {
				return
			}
		}
		aRecord.CustomerCount = count
		err = m.updateMeta(claim, userMap, newRecord)
		if err != nil {
			return
		}

		resultList[indxValue] = *newRecord
	}

	result = resultList

	return
}

func (m ModelServiceImpl) prepare(claim *auth.AuthenticatedClaim, companyName string, modelValue *model.Model) (currentUserId uint32, err error) {
	currentUserId = 1
	if claim != nil {
		currentUserId = claim.UserId
	}

	var company *model.Company
	company, err = m.companyService.Get(claim, companyName)
	if err != nil {
		return
	}
	modelValue.CompanyID = company.ID

	return
}

func (m ModelServiceImpl) updateMeta(claim *auth.AuthenticatedClaim, userMap map[int](*string), modelValue *model.Model) (err error) {
	if userMap == nil {
		userMap = make(map[int]*string)
	}

	userName, err := getUserName(m.db, claim, userMap, modelValue.CreatedById)
	if err != nil {
		return err
	}
	modelValue.CreatedBy = userName
	if err != nil {
		return err
	}

	if modelValue.LastUpdatedById > 0 {
		userName, err = getUserName(m.db, claim, userMap, modelValue.LastUpdatedById)
		if err != nil {
			return err
		}
		modelValue.LastUpdatedBy = userName
	}

	return
}
