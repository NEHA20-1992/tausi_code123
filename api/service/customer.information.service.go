package service

import (
	"errors"
	"math"

	// "strconv"
	"time"

	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/model"

	// "github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var ErrAllModelVariableRequired = errors.New("value for all model variable required")
var ErrModelVariableNotFound = errors.New("model variable not found")
var ErrInvalidCustomerIDProofType = errors.New("invalid customer id proof type")

const (
	SqlSearch                         string = `Select from customer_information where customer_information_id=?`
	SqlPagination                     string = ` LIMIT %d OFFSET %d ",  perPage, (page-1)*perPage`
	SqlComputeCustomerInformationItem string = `UPDATE customer_information_item cii
	INNER JOIN
model_variable mv ON mv.model_variable_id = cii.model_variable_id 
SET 
cii.preprocessed_value = IF(mv.standard_deviation_value = 0,
	ROUND(cii.value, 9),
	ROUND((cii.value - mv.mean_value) / mv.standard_deviation_value,
			9))
WHERE
cii.customer_information_id = ?;
`
	SqlComputeCustomerInformation string = `UPDATE customer_information ci
	INNER JOIN
(SELECT 
	customer_information_id,
		ROUND(EXP(ppv) / (1 + EXP(ppv)), 2) probability_of_default_percentage,
		100 - (ROUND(EXP(ppv) / (1 + EXP(ppv)), 4) * 100) group_score
FROM
	(SELECT 
	cii.customer_information_id,
		ROUND(SUM(cii.preprocessed_value * mv.coefficient_value) + m.intercept_value, 9) ppv
FROM
	customer_information_item cii
INNER JOIN model_variable mv ON mv.model_variable_id = cii.model_variable_id
INNER JOIN model m ON m.model_id = mv.model_id
WHERE
	cii.customer_information_id = ?
GROUP BY cii.customer_information_id, m.intercept_value) b) a ON a.customer_information_id = ci.customer_information_id 
SET 
ci.probability_of_default_percentage = a.probability_of_default_percentage,
ci.group_score = TRUNCATE(a.group_score, 2)`
)

type CustomerInformationService interface {
	Create(claim *auth.AuthenticatedClaim, companyName string, modelName string, CustomerInformation *model.CustomerInformation) (*model.CustomerInformation, error)
	Get(claim *auth.AuthenticatedClaim, companyName string, modelName string, ID uint32) (*model.CustomerInformation, error)
	GetAllCompayCustomer(claim *auth.AuthenticatedClaim, companyName string, modelName string) ([]model.CustomerInformation, error)
	// GetAll(claim *auth.AuthenticatedClaim, companyName string) ([]model.CustomerInformation, error)
	GetAll(claim *auth.AuthenticatedClaim, companyName string, modelName string, pageNumber uint32, size uint32, cid uint32, search string, CreditScore float32) ([]model.CustomerInformation, error)
	// pagination([]model.CustomerInformation)([]model.CustomerInformation,error)
	GetCreditScore(claim *auth.AuthenticatedClaim, companyName string, modelName string, ID uint32) (*model.CustomerCreditScore, error)
}

type CustomerInformationServiceImpl struct {
	db                         *gorm.DB
	customerIdProofTypeService CustomerIdProofTypeService
	modelService               ModelService
	modelVariableService       ModelVariableService
	companyService             CompanyService
}

func GetCustomerInformationService(db *gorm.DB) CustomerInformationService {
	return &CustomerInformationServiceImpl{
		db:                   db,
		modelService:         GetModelService(db),
		companyService:       GetCompanyService(db),
		modelVariableService: GetModelVariableService(db)}
}

// func (c *fiber.Ctx)pagination(custInfo []model.CustomerInformation)(result []model.CustomerInformation,err error){
// sql := "Select * from custinfo"
// page := strconv.Atoi(m.Query("page","1"))
// m.db.Raw(sql).Scan(&custInfo)
// return
// }
func (m *CustomerInformationServiceImpl) Create(claim *auth.AuthenticatedClaim, companyName string, modelName string, customerInformation *model.CustomerInformation) (result *model.CustomerInformation, err error) {
	var modelValue *model.Model
	modelValue, err = m.modelService.GetDetails(claim, companyName, modelName)
	if err != nil {
		return
	}
	customerInformation.CompanyID = modelValue.CompanyID
	customerInformation.ModelID = modelValue.ID
	customerInformation.CreatedById = claim.UserId
	customerInformation.CreatedAt = time.Now()

	var cIdProofType = model.CustomerIdProofType{}
	err = m.db.Model(&model.CustomerIdProofType{}).Where("name = ?", customerInformation.CustomerIDProofType).First(&cIdProofType).Error
	if err != nil {
		err = ErrInvalidCustomerIDProofType
		return
	}

	customerInformation.CustomerIDProofTypeID = cIdProofType.ID
	if len(modelValue.Variables) != len(customerInformation.Items) {
		err = ErrAllModelVariableRequired
		return
	}

	var variableMap map[string]*model.ModelVariable = make(map[string]*model.ModelVariable)
	for _, mv := range modelValue.Variables {
		variableMap[mv.Name] = &mv
	}

	for _, ci := range customerInformation.Items {
		if variableMap[ci.Name] == nil {
			err = ErrModelVariableNotFound
			return
		}
	}

	tx := m.db.Begin()
	err = tx.
		Model(&customerInformation).
		Create(&customerInformation).
		Error
	if err != nil {
		return
	}

	var ciList []model.CustomerInformationItem = make([]model.CustomerInformationItem, len(customerInformation.Items))
	for indx, ci := range customerInformation.Items {
		var newRecord *model.CustomerInformationItem = &ci
		newRecord.CustomerInformationID = customerInformation.ID
		ciList[indx] = *newRecord

		var variableDef *model.ModelVariable
		variableDef, err = m.modelVariableService.Get(claim, companyName, modelName, ci.Name)
		if err != nil {
			return
		}
		newRecord.ModelVariableID = variableDef.ID

		err = tx.
			Model(&model.CustomerInformationItem{}).
			Create(newRecord).
			Error
		if err != nil {
			return
		}
	}
	tx.Exec(SqlComputeCustomerInformationItem, customerInformation.ID)
	tx.Exec(SqlComputeCustomerInformation, customerInformation.ID)
	customerInformation.Items = ciList

	var createdRecord model.CustomerInformation
	err = tx.Model(&model.CustomerInformation{}).
		Where("customer_information_id = ?", customerInformation.ID).
		Find(&createdRecord).
		Error
	if err != nil {
		tx.Rollback()
		return
	}

	updatedRecord, err := m.updatedRecord(tx, claim, &createdRecord)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()

	result = updatedRecord

	return
}

func (m *CustomerInformationServiceImpl) Get(claim *auth.AuthenticatedClaim, companyName string, modelName string, ID uint32) (result *model.CustomerInformation, err error) {
	var modelValue *model.Model
	modelValue, err = m.modelService.Get(claim, companyName, modelName)
	if err != nil {
		return
	}
	result, err = m.get(claim, modelValue.ID, ID)
	return
}

func (m *CustomerInformationServiceImpl) GetCreditScore(claim *auth.AuthenticatedClaim, companyName string, modelName string, ID uint32) (result *model.CustomerCreditScore, err error) {
	var modelValue *model.Model
	modelValue, err = m.modelService.Get(claim, companyName, modelName)
	if err != nil {
		return
	}
	var ci model.CustomerInformation
	err = m.db.Model(&model.CustomerInformation{}).
		Where("model_id = ? AND customer_information_id = ?", modelValue.ID, ID).
		Find(&ci).
		Error
	if err != nil {
		return
	}
	var ciCreditScore model.CustomerCreditScore

	ciCreditScore.ID = ci.ID
	ciCreditScore.ModelID = ci.ModelID
	ciCreditScore.FirstName = ci.FirstName
	ciCreditScore.LastName = ci.LastName
	ciCreditScore.ContactNumber = ci.ContactNumber
	ciCreditScore.City = ci.City
	ciCreditScore.ProbabilityOfDefaultPercentage = math.Round(ci.ProbabilityOfDefaultPercentage * 100)
	ciCreditScore.GroupScore = math.Round(ci.GroupScore*100) / 100
	//ciCreditScore.GroupScore = ci.GroupScore

	var grouScoreRoundValue int64 = int64(ci.GroupScore)

	//var ciDefaultPercentage float64 = ci.ProbabilityOfDefaultPercentage
	switch {
	case grouScoreRoundValue >= 90 && grouScoreRoundValue <= 100:
		ciCreditScore.CreditScore = "A1"
	case grouScoreRoundValue >= 80 && grouScoreRoundValue <= 89:
		ciCreditScore.CreditScore = "A2"
	case grouScoreRoundValue >= 70 && grouScoreRoundValue <= 79:
		ciCreditScore.CreditScore = "B1"
	case grouScoreRoundValue >= 60 && grouScoreRoundValue <= 69:
		ciCreditScore.CreditScore = "B2"
	case grouScoreRoundValue >= 50 && grouScoreRoundValue <= 59:
		ciCreditScore.CreditScore = "C1"
	case grouScoreRoundValue >= 40 && grouScoreRoundValue <= 49:
		ciCreditScore.CreditScore = "C2"
	case grouScoreRoundValue >= 30 && grouScoreRoundValue <= 39:
		ciCreditScore.CreditScore = "D1"
	case grouScoreRoundValue >= 20 && grouScoreRoundValue <= 29:
		ciCreditScore.CreditScore = "D2"
	case grouScoreRoundValue >= 10 && grouScoreRoundValue <= 19:
		ciCreditScore.CreditScore = "E1"
	case grouScoreRoundValue >= 0 && grouScoreRoundValue <= 9:
		ciCreditScore.CreditScore = "E2"
	}

	result = &ciCreditScore

	return
}

func (m *CustomerInformationServiceImpl) get(claim *auth.AuthenticatedClaim, modelId uint32, ID uint32) (result *model.CustomerInformation, err error) {
	var existingRecord model.CustomerInformation
	err = m.db.Model(&model.CustomerInformation{}).
		Where("model_id = ? AND customer_information_id = ?", modelId, ID).
		Find(&existingRecord).
		Error
	if err != nil {
		return
	}

	result, err = m.updatedRecord(m.db, claim, &existingRecord)
	if err != nil {
		return
	}

	return
}

func (m *CustomerInformationServiceImpl) getx(claim *auth.AuthenticatedClaim, ID uint32) (result *model.CustomerInformation, err error) {
	var existingRecord model.CustomerInformation
	err = m.db.Model(&model.CustomerInformation{}).
		Where(" customer_information_id = ?", ID).
		Find(&existingRecord).
		Error
	if err != nil {
		return
	}

	result, err = m.updatedRecord(m.db, claim, &existingRecord)
	if err != nil {
		return
	}

	return
}
func (m *CustomerInformationServiceImpl) GetAllCompayCustomer(claim *auth.AuthenticatedClaim, companyName string, modelName string) (result []model.CustomerInformation, err error) {
	modelValue, err := m.modelService.Get(claim, companyName, modelName)
	if err != nil {
		return
	}

	list := []model.CustomerInformation{}
	err = m.db.Model(&model.CustomerInformation{}).
		Select("customer_information_id").
		Where("model_id = ?", modelValue.ID).
		Find(&list).
		Error
	if err != nil {
		return
	}

	var resultList []model.CustomerInformation = make([]model.CustomerInformation, len(list))
	for ciIndx, ci := range list {
		var ciValue, updatedRecord *model.CustomerInformation
		ciValue, err = m.get(claim, modelValue.ID, ci.ID)
		if err != nil {
			return
		}

		updatedRecord, err = m.updatedRecord(m.db, claim, ciValue)
		if err != nil {
			return
		}

		resultList[ciIndx] = *updatedRecord
	}

	result = resultList

	return
}

func (m *CustomerInformationServiceImpl) GetAll(claim *auth.AuthenticatedClaim, companyName string, modelName string, pageNumber uint32, size uint32, cid uint32, search string, CreditScore float32) (result []model.CustomerInformation, err error) {
	var result1 []model.CustomerInformation
	if companyName == "" {
		list := []model.CustomerInformation{}
		err = m.db.Model(&model.CustomerInformation{}).
			Select("customer_information_id").Limit(int(size)).Offset(int((size) * (pageNumber - 1))).
			Find(&list).
			Error
		if err != nil {
			return
		}
		var resultList []model.CustomerInformation = make([]model.CustomerInformation, len(list))
		for ciIndx, ci := range list {
			var ciValue, updatedRecord *model.CustomerInformation
			ciValue, err = m.getx(claim, ci.ID)
			if err != nil {
				return
			}

			updatedRecord, err = m.updatedRecord(m.db, claim, ciValue)
			if err != nil {
				return
			}

			resultList[ciIndx] = *updatedRecord

		}
		result1 = resultList

	} else if modelName == "" {
		modelValue, _ := m.companyService.Get(claim, companyName)

		list := []model.CustomerInformation{}
		err = m.db.Model(&model.CustomerInformation{}).
			Select("customer_information_id").
			Where("company_id = ? and customer_information_id <= ? and customer_information_id >?", modelValue.ID, size*(pageNumber), size*(pageNumber-1)).
			Find(&list).
			Error
		if err != nil {
			return
		}

		var resultList []model.CustomerInformation = make([]model.CustomerInformation, len(list))
		for ciIndx, ci := range list {
			var ciValue, updatedRecord *model.CustomerInformation
			ciValue, err = m.getx(claim, ci.ID)
			if err != nil {
				return
			}

			updatedRecord, err = m.updatedRecord(m.db, claim, ciValue)
			if err != nil {
				return
			}

			resultList[ciIndx] = *updatedRecord
		}
		result1 = resultList
	} else {
		modelValue, _ := m.modelService.Get(claim, companyName, modelName)

		list := []model.CustomerInformation{}
		err = m.db.Model(&model.CustomerInformation{}).
			Select("customer_information_id").
			Where("model_id = ? and company_id =? and customer_information_id <= ? and customer_information_id >?", modelValue.ID, modelValue.CompanyID, size*pageNumber, size*(pageNumber-1)).
			Find(&list).
			Error
		if err != nil {
			return
		}

		var resultList []model.CustomerInformation = make([]model.CustomerInformation, len(list))
		for ciIndx, ci := range list {
			var ciValue, updatedRecord *model.CustomerInformation
			ciValue, err = m.getx(claim, ci.ID)
			if err != nil {
				return
			}

			updatedRecord, err = m.updatedRecord(m.db, claim, ciValue)
			if err != nil {
				return
			}

			resultList[ciIndx] = *updatedRecord
		}
		result1 = resultList
	}
	if cid > 0 {

		err = m.db.Model(&model.CustomerInformation{}).
			Select("customer_information_id").
			Where(" customer_information_id =? and customer_information_id <= ? and customer_information_id >?", cid, size*pageNumber, size*(pageNumber-1)).
			Find(&result1).
			Error
		if err != nil {
			return
		}
		var resultList1 []model.CustomerInformation = make([]model.CustomerInformation, len(result1))
		for ciIndx, ci := range result1 {
			var ciValue, updatedRecord *model.CustomerInformation
			ciValue, err = m.getx(claim, ci.ID)
			if err != nil {
				return
			}

			updatedRecord, err = m.updatedRecord(m.db, claim, ciValue)
			if err != nil {
				return
			}

			resultList1[ciIndx] = *updatedRecord
		}
		result1 = resultList1
	}
	if search != "" {

		err = m.db.Model(&model.CustomerInformation{}).
			Select("customer_information_id").
			Where(" first_name =? and customer_information_id <= ? and customer_information_id >?", search, size*pageNumber, size*(pageNumber-1)).
			Find(&result1).
			Error
		if err != nil {
			return
		}
		var resultList1 []model.CustomerInformation = make([]model.CustomerInformation, len(result1))
		for ciIndx, ci := range result1 {
			var ciValue, updatedRecord *model.CustomerInformation
			ciValue, err = m.getx(claim, ci.ID)
			if err != nil {
				return
			}

			updatedRecord, err = m.updatedRecord(m.db, claim, ciValue)
			if err != nil {
				return
			}

			resultList1[ciIndx] = *updatedRecord
		}
		result1 = resultList1
	}
	if CreditScore > 0 {

		err = m.db.Model(&model.CustomerInformation{}).
			Select("customer_information_id").
			Where(" group_score =>? and customer_information_id <= ? and customer_information_id >?", CreditScore, size*pageNumber, size*(pageNumber-1)).
			Find(&result1).
			Error
		if err != nil {
			return
		}
		var resultList1 []model.CustomerInformation = make([]model.CustomerInformation, len(result1))
		for ciIndx, ci := range result1 {
			var ciValue, updatedRecord *model.CustomerInformation
			ciValue, err = m.getx(claim, ci.ID)
			if err != nil {
				return
			}

			updatedRecord, err = m.updatedRecord(m.db, claim, ciValue)
			if err != nil {
				return
			}

			resultList1[ciIndx] = *updatedRecord
		}
		result1 = resultList1
	}

	result = result1
	return
}
func (m CustomerInformationServiceImpl) updatedRecord(tx *gorm.DB, claim *auth.AuthenticatedClaim, createdRecord *model.CustomerInformation) (result *model.CustomerInformation, err error) {
	createdRecord.CreatedBy, err = getUserName(m.db, claim, nil, createdRecord.CreatedById)
	if err != nil {
		return
	}

	createdItemList := []model.CustomerInformationItem{}
	err = tx.Model(&model.CustomerInformationItem{}).
		Where("customer_information_id = ?", createdRecord.ID).
		Find(&createdItemList).
		Error

	var newItemList []model.CustomerInformationItem = make([]model.CustomerInformationItem, len(createdItemList))
	for indx, ci := range createdItemList {
		var mv = &model.ModelVariable{}
		err = m.db.
			Model(mv).
			Where("model_variable_id = ?", ci.ModelVariableID).
			First(&mv).
			Error
		if err != nil {
			return
		}
		ci.Name = mv.Name
		newItemList[indx] = ci
	}
	createdRecord.Items = newItemList

	result = createdRecord

	return
}

func (m CustomerInformationServiceImpl) prepare(claim *auth.AuthenticatedClaim, customerInformation *model.CustomerInformation) (currentUserId uint32, err error) {
	currentUserId = 1
	if claim != nil {
		currentUserId = claim.UserId
	}

	customerIDProofType, err := m.customerIdProofTypeService.Get(customerInformation.CustomerIDProofType)
	if err != nil {
		return
	}
	customerInformation.CustomerIDProofTypeID = customerIDProofType.ID
	return
}

func (m CustomerInformationServiceImpl) CustomerIdProofType(claim *auth.AuthenticatedClaim, customerinformation *model.CustomerInformation) (err error) {
	var reference *model.CustomerIdProofType
	reference, err = m.customerIdProofTypeService.GetById(customerinformation.CustomerIDProofTypeID)
	if err != nil {
		return
	}

	customerinformation.CustomerIDProofType = reference.Name

	return
}
func Round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}
