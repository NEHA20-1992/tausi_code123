package service

import (
	"errors"
	"net/http"
	"time"

	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

var ErrCompanyAlreadyExists = errors.New("company name already exist")

type CompanyService interface {
	Create(claim *auth.AuthenticatedClaim, company *model.Company, logoContents []byte) (*model.Company, error)
	Update(claim *auth.AuthenticatedClaim, existingCompanyName string, company *model.Company, logoContents []byte) (*model.Company, error)
	Get(claim *auth.AuthenticatedClaim, companyName string) (*model.Company, error)
	GetById(claim *auth.AuthenticatedClaim, companyID uint32) (*model.Company, error)
	GetCompanyLogo(claim *auth.AuthenticatedClaim, companyName string) (*model.CompanyFile, error)
	GetAll(claim *auth.AuthenticatedClaim) ([]model.Company, error)
	GetEntityCount(claim *auth.AuthenticatedClaim) (*model.EnitityCount, error)
	GetCraEntityCount(claim *auth.AuthenticatedClaim, companyName string) (*model.CRACount, error)
}

type CompanyServiceImpl struct {
	db                 *gorm.DB
	companyTypeService CompanyTypeService
	countryService     CountryService
}

func GetCompanyService(db *gorm.DB) (service CompanyService) {
	return CompanyServiceImpl{
		db:                 db,
		countryService:     GetCountryService(db),
		companyTypeService: GetCompanyTypeService(db)}
}

func (m CompanyServiceImpl) Create(claim *auth.AuthenticatedClaim, company *model.Company, logoContents []byte) (result *model.Company, err error) {
	if company == nil {
		return
	}
	var existingRecord model.Company
	err = m.db.Model(&model.Company{}).Select("company_id").Where("name = ?", company.Name).Find(&existingRecord).Error
	if err != nil || existingRecord.ID > 0 {
		err = ErrCompanyAlreadyExists
		return
	}

	var currentUserId uint32 = 1

	currentUserId, err = m.prepare(claim, company)
	if err != nil {
		return
	}
	company.CreatedById = currentUserId

	err = m.db.Model(&company).Omit("last_updated_by_id", "last_updated_at").Create(&company).Error
	if err != nil {
		return
	}

	if len(logoContents) > 0 {
		mimeType := http.DetectContentType(logoContents)
		err = m.db.
			Model(&model.CompanyFile{}).
			Create(&model.CompanyFile{
				ID:       company.ID,
				Logo:     logoContents,
				MimeType: mimeType,
			}).
			Error

		if err != nil {
			return
		}
	}

	// This is the display the updated user
	createdRecord := model.Company{}
	err = m.db.Debug().Model(&createdRecord).Where("company_id = ?", company.ID).Take(&createdRecord).Error
	if err != nil {
		return
	}
	result = &createdRecord
	err = m.updateMeta(claim, result)
	if err != nil {
		return
	}

	return
}

func (m CompanyServiceImpl) Update(claim *auth.AuthenticatedClaim, existingCompanyName string, company *model.Company, logoContents []byte) (result *model.Company, err error) {
	var existingRecord model.Company
	err = m.db.Model(&model.Company{}).Select("company_id").Where("name = ?", existingCompanyName).First(&existingRecord).Error
	if err != nil {
		return
	}
	company.ID = existingRecord.ID

	var currentUserId uint32 = 1

	currentUserId, err = m.prepare(claim, company)
	if err != nil {
		return
	}

	company.CreatedById = currentUserId
	m.db.Debug().Model(&model.Company{}).Where("company_id = ?", company.ID).Take(&model.Company{}).UpdateColumns(
		map[string]interface{}{
			"name":               company.Name,
			"company_type_id":    company.CompanyTypeID,
			"address":            company.Address,
			"region_county":      company.RegionCounty,
			"country_id":         company.CountryId,
			"contact_number":     company.ContactNumber,
			"email_address":      company.EmailAddress,
			"active":             company.Active,
			"last_updated_by_id": currentUserId,
			"last_updated_at":    time.Now(),
		})

	if len(logoContents) > 0 {
		mimeType := http.DetectContentType(logoContents)

		var existingCompanyFile model.CompanyFile
		m.db.Model(&model.CompanyFile{}).Select("company_id").Where("company_id = ?", company.ID).First(&existingCompanyFile)

		if existingCompanyFile.ID > 0 {
			m.db.Debug().Model(&model.CompanyFile{}).
				Where("company_id = ?", company.ID).
				Take(&model.CompanyFile{}).
				UpdateColumns(
					map[string]interface{}{
						"logo":      logoContents,
						"mime_type": mimeType,
					})
			err = m.db.Error
			if err != nil {
				return
			}
		} else {
			mimeType := http.DetectContentType(logoContents)
			err = m.db.
				Model(&model.CompanyFile{}).
				Create(&model.CompanyFile{
					ID:       company.ID,
					Logo:     logoContents,
					MimeType: mimeType,
				}).
				Error

			if err != nil {
				return
			}
		}
	}
	// This is the display the updated user
	updatedRecord := model.Company{}
	err = m.db.Debug().Model(&updatedRecord).Where("company_id = ?", company.ID).Take(&updatedRecord).Error
	if err != nil {
		return
	}
	result = &updatedRecord
	err = m.updateMeta(claim, result)
	if err != nil {
		return
	}
	return
}

func (m CompanyServiceImpl) GetCompanyLogo(claim *auth.AuthenticatedClaim, companyName string) (result *model.CompanyFile, err error) {
	var company model.Company
	var companyLogo model.CompanyFile

	err = m.db.Model(&company).Where("name = ?", companyName).First(&company).Error
	if err != nil {
		return
	}

	err = m.db.Model(&companyLogo).Where("company_id = ?", company.ID).First(&companyLogo).Error
	if err != nil {
		return
	}

	result = &companyLogo
	return
}

func (m CompanyServiceImpl) GetById(claim *auth.AuthenticatedClaim, companyID uint32) (result *model.Company, err error) {
	var company model.Company

	err = m.db.Model(&company).Where("company_id = ?", companyID).First(&company).Error
	if err != nil {
		return
	}

	result = &company
	err = m.updateMeta(claim, result)
	if err != nil {
		return
	}

	return
}

func (m CompanyServiceImpl) Get(claim *auth.AuthenticatedClaim, companyName string) (result *model.Company, err error) {
	var company model.Company
	err = m.db.Model(&company).Where("name = ?", companyName).First(&company).Error
	if err != nil {
		return
	}

	result = &company
	err = m.updateMeta(claim, result)
	if err != nil {
		return
	}

	return
}

func (m CompanyServiceImpl) GetCraEntityCount(claim *auth.AuthenticatedClaim, companyName string) (result *model.CRACount, err error) {
	var craCount, projectCount int64
	var allCount model.CRACount

	var company model.Company
	err = m.db.Model(&company).Where("name = ?", companyName).First(&company).Error
	if err != nil {
		return
	}

	if company.ID > 0 {
		var models = []model.Model{}
		err = m.db.Where("company_id = ?", company.ID).
			Find(&models).
			Count(&projectCount).
			Error
		if err != nil {
			return
		}
	}

	if company.ID > 1 {
		var models = []model.User{}
		err = m.db.Where("company_id = ?", company.ID).
			Find(&models).
			Count(&craCount).
			Error
		if err != nil {
			return
		}
	}

	allCount.UserCount = craCount
	allCount.ProjectCount = projectCount
	result = &allCount
	return
}

func (m CompanyServiceImpl) GetEntityCount(claim *auth.AuthenticatedClaim) (result *model.EnitityCount, err error) {
	var companyCount, craCount, projectCount int64
	var allCount model.EnitityCount

	err = m.db.Table("company").Count(&companyCount).Error
	if err != nil {
		companyCount = 0
	}

	err = m.db.Table("user_user_role").Count(&craCount).Error
	if err != nil {
		craCount = 0
	}

	err = m.db.Table("model").Count(&projectCount).Error
	if err != nil {
		projectCount = 0
	}
	allCount.CompanyCount = companyCount
	allCount.UserCount = craCount
	allCount.ProjectCount = projectCount
	result = &allCount
	return
}

func (m CompanyServiceImpl) GetAll(claim *auth.AuthenticatedClaim) (result []model.Company, err error) {

	result = []model.Company{}
	err = m.db.Model(&model.Company{}).Limit(1000).Find(&result).Error
	if err != nil {
		result = nil
	}

	var companyList []model.Company = make([]model.Company, len(result))
	for companyIndex, aCompany := range result {
		var newCompany = &aCompany
		err = m.updateMeta(claim, newCompany)
		if err != nil {
			return
		}
		companyList[companyIndex] = *newCompany
	}

	result = companyList

	return
}

func (m CompanyServiceImpl) prepare(claim *auth.AuthenticatedClaim, company *model.Company) (currentUserId uint32, err error) {
	currentUserId = 1
	if claim != nil {
		currentUserId = claim.UserId
	}

	country, err := m.countryService.Get(company.Country)
	if err != nil {
		return
	}
	company.CountryId = country.ID

	companyType, err := m.companyTypeService.Get(company.Type_)
	if err != nil {
		return
	}
	company.CompanyTypeID = companyType.ID

	return
}

func (m CompanyServiceImpl) updateMeta(claim *auth.AuthenticatedClaim, company *model.Company) (err error) {
	var userMap map[int](*string) = make(map[int]*string)
	err = m.loadCountry(claim, company)
	if err != nil {
		return
	}
	err = m.loadCompanyType(claim, company)
	if err != nil {
		return
	}

	// companyLogo := model.CompanyFile{}
	// m.db.Model(&model.CompanyFile{}).Where("company_id = ?", company.ID).Find(&companyLogo)
	// company.Logo = companyLogo

	userName, err := getUserName(m.db, claim, userMap, company.CreatedById)
	if err != nil {
		return err
	}
	company.CreatedBy = userName
	if err != nil {
		return err
	}

	if company.LastUpdatedById > 0 {
		userName, err = getUserName(m.db, claim, userMap, company.LastUpdatedById)
		if err != nil {
			return err
		}
		company.LastUpdatedBy = userName
	}

	return
}

func (m CompanyServiceImpl) loadCountry(claim *auth.AuthenticatedClaim, company *model.Company) (err error) {
	var country *model.Country
	country, err = m.countryService.GetById(company.CountryId)
	if err != nil {
		return
	}

	company.Country = country.Name

	return
}

func getUserName(db *gorm.DB, claim *auth.AuthenticatedClaim, userMap map[int](*string), userId uint32) (userName string, err error) {
	if userMap == nil {
		userMap = make(map[int]*string)
	}
	var uService UserService = GetUserService(db)
	existingUserName := userMap[int(userId)]

	if existingUserName == nil {
		userName, err = uService.GetUserNameById(claim, userId)
		if err == nil {
			userMap[int(userId)] = &userName
		}
	} else {
		userName = *existingUserName
	}

	return
}

func (m CompanyServiceImpl) loadCompanyType(claim *auth.AuthenticatedClaim, company *model.Company) (err error) {
	var reference *model.CompanyType
	reference, err = m.companyTypeService.GetById(company.CompanyTypeID)
	if err != nil {
		return
	}

	company.Type_ = reference.Name

	return
}
