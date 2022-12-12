package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

var ErrNotUploaded = errors.New("ERROR:file not uploaded")

type FileService interface {
	GetAllRawFile(claim *auth.AuthenticatedClaim, companyName string) ([]model.CompanyDataFile, error)
	GetAllProcessedFile(claim *auth.AuthenticatedClaim, companyName string, modelName string) ([]model.CompanyDataFile, error)
	UploadProcessedFile(claim *auth.AuthenticatedClaim, companyName string, modelName string, customerinfo []model.CustomerInformation, fileheader *multipart.FileHeader) (*model.CompanyDataFile, error)
	UploadRawFile(claim *auth.AuthenticatedClaim, companyDataFile *model.CompanyDataFile) (*model.CompanyDataFile, error)
}

type FileServiceImpl struct {
	db                  *gorm.DB
	modelService        ModelService
	fileTypeService     CompanyFileTypeService
	customerInfoService CustomerInformationService
	companyService      CompanyService
	//userRoleService UserRoleService
}

func GetFileService(db *gorm.DB) FileService {
	return FileServiceImpl{
		db:                  db,
		modelService:        GetModelService(db),
		fileTypeService:     GetCompanyFileTypeService(db),
		customerInfoService: GetCustomerInformationService(db),
		companyService:      GetCompanyService(db),

		//userRoleService: GetUserRoleService(db)
	}
}

func (m FileServiceImpl) GetAllRawFile(claim *auth.AuthenticatedClaim, companyName string) (result []model.CompanyDataFile, err error) {
	var company *model.Company
	company, err = m.companyService.Get(claim, companyName)
	if err != nil {
		return
	}

	// fileType, err := m.fileTypeService.Get(fileTypeId)
	// if err != nil {
	// 	return
	// }

	result = []model.CompanyDataFile{}
	err = m.db.Model(&model.CompanyDataFile{}).Limit(1000).Where("company_id = ? AND company_data_file_type_id = ?", company.ID, 1).Find(&result).Error
	if err != nil {
		result = nil
	}

	var resultList []model.CompanyDataFile = make([]model.CompanyDataFile, len(result))
	for indxValue, aRecord := range result {
		var newRecord = &aRecord
		resultList[indxValue] = *newRecord
	}
	result = resultList
	return
}

func (m FileServiceImpl) GetAllProcessedFile(claim *auth.AuthenticatedClaim, companyName string, modelName string) (result []model.CompanyDataFile, err error) {
	var company *model.Company

	company, err = m.companyService.Get(claim, companyName)
	if err != nil {
		return
	}
	modelValue, err := m.modelService.Get(claim, companyName, modelName)
	if err != nil {
		return
	}
	// fileType, err := m.fileTypeService.Get(fileTypeId)
	// if err != nil {
	// 	return
	// }

	result = []model.CompanyDataFile{}
	err = m.db.Model(&model.CompanyDataFile{}).Limit(1000).Where("company_id = ? AND model_id=? AND company_data_file_type_id = ?", company.ID, modelValue.ID, 2).Find(&result).Error
	if err != nil {
		result = nil
	}

	var resultList []model.CompanyDataFile = make([]model.CompanyDataFile, len(result))
	for indxValue, aRecord := range result {
		var newRecord = &aRecord
		resultList[indxValue] = *newRecord
	}
	result = resultList
	return
}

func (m FileServiceImpl) UploadProcessedFile(claim *auth.AuthenticatedClaim, companyName string, modelName string, customerinfo []model.CustomerInformation, fileheader *multipart.FileHeader) (result *model.CompanyDataFile, err error) {

	if customerinfo == nil {
		return
	}
	for _, req := range customerinfo {
		result, _ := m.customerInfoService.Create(claim, companyName, modelName, &req)
		fmt.Println(result)
	}
	companyDataFile := model.CompanyDataFile{}
	value, err := m.companyService.Get(claim, companyName)
	value1, err := m.modelService.Get(claim, companyName, modelName)

	companyDataFile.CompanyID = value.ID
	companyDataFile.ModelID = value1.ID
	companyDataFile.CreatedById = claim.UserId
	companyDataFile.Name = fileheader.Filename
	companyDataFile.CompanyDataFileTypeID = 2
	companyDataFile.CreatedAt = time.Now()
	companyDataFile.Description = "Processed"

	err = m.db.Model(&companyDataFile).Create(&companyDataFile).Error
	if err != nil {
		return
	}
	createdRecord := model.CompanyDataFile{}
	err = m.db.Debug().Model(&createdRecord).Where("file_id = ?", companyDataFile.ID).Take(&createdRecord).Error
	if err != nil {
		return
	}
	result = &createdRecord

	return
}

func (m FileServiceImpl) UploadRawFile(claim *auth.AuthenticatedClaim, companyDataFile *model.CompanyDataFile) (result *model.CompanyDataFile, err error) {
	if companyDataFile == nil {
		return
	}
	var existingRecord model.CompanyDataFile
	err = m.db.Model(&model.CompanyDataFile{}).Select("company_id").Where("name = ?", companyDataFile.Name).Find(&existingRecord).Error
	if err != nil || existingRecord.ID > 0 {
		err = ErrCompanyAlreadyExists
		return
	}

	companyDataFile.CreatedById = claim.UserId

	err = m.db.Model(&companyDataFile).Omit("created_by_id").Create(&companyDataFile).Error
	if err != nil {
		return
	}
	// This is the display the updated user
	createdRecord := model.CompanyDataFile{}
	err = m.db.Debug().Model(&createdRecord).Where("file_id = ?", companyDataFile.ID).Take(&createdRecord).Error
	if err != nil {
		return
	}
	result = &createdRecord

	return
}
