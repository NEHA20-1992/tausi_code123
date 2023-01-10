package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	// "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/helper"
	"github.com/NEHA20-1992/tausi_code/api/validator"

	// "github.com/northbright/pathhelper"
	// "github.com/northbright/xls2csv-go/xls2csv"

	"github.com/xuri/excelize/v2"
	//  "https://github.com/tealeg/xlsx2csv"
	//  "github.com/NEHA20-1992/tausi_code/api/helper"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"gorm.io/gorm"
)

var ErrNotUploaded = errors.New("ERROR:file not uploaded")
var ErrFileNameNotValid = errors.New("enter valid file name.")

const (
	rawTypeFile                             = 1
	processedTypeFile                       = 2
	SqlLoadInFileCustomerInformation string = ``
)

type FileService interface {
	// GetAllRawFile(claim *auth.AuthenticatedClaim, companyName string) ([]model.CompanyDataFile, error)
	GetAllProcessedFile(claim *auth.AuthenticatedClaim, companyName string, modelName string) ([]model.CompanyDataFile, error)
	UploadProcessedFile(claim *auth.AuthenticatedClaim, companyName string, modelName string, fileheader string, file *excelize.File) (*model.CompanyDataFile, error)
	UploadRawFile(claim *auth.AuthenticatedClaim, companyName string, modelName string, fileheader string, file *excelize.File) (*model.CompanyDataFile, error)
	DownloadRawDataFile(claim *auth.AuthenticatedClaim, companyName string, fileName string) (*excelize.File, []model.CompanyDataFile, error)
	// UploadRawFile(claim *auth.AuthenticatedClaim, companyDataFile *model.CompanyDataFile) (*model.CompanyDataFile, error)
}

type FileServiceImpl struct {
	db                  *gorm.DB
	modelService        ModelService
	fileTypeService     CompanyFileTypeService
	customerInfoService CustomerInformationService
	userService         UserService
	companyService      CompanyService
	variableService     ModelVariableService
	//userRoleService UserRoleService
}

func GetFileService(db *gorm.DB) FileService {
	return FileServiceImpl{
		db:                  db,
		modelService:        GetModelService(db),
		fileTypeService:     GetCompanyFileTypeService(db),
		userService:         GetUserService(db),
		customerInfoService: GetCustomerInformationService(db),
		companyService:      GetCompanyService(db),
		variableService:     GetModelVariableService(db),

		//userRoleService: GetUserRoleService(db)
	}
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

	result = []model.CompanyDataFile{}
	err = m.db.Model(&model.CompanyDataFile{}).Where("company_id = ? AND model_id=? AND company_data_file_type_id = ?", company.ID, modelValue.ID, 2).Find(&result).Error
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

func (m FileServiceImpl) UploadProcessedFile(claim *auth.AuthenticatedClaim, companyName string, modelName string, fileheader string, file *excelize.File) (result *model.CompanyDataFile, err error) {
	if file == nil {
		return
	}
	companyDataFile := model.CompanyDataFile{}
	value, err := m.companyService.Get(claim, companyName)
	value1, err := m.modelService.Get(claim, companyName, modelName)

	companyDataFile.CompanyID = value.ID
	companyDataFile.ModelID = value1.ID
	companyDataFile.CreatedById = claim.UserId
	companyDataFile.Name = fileheader
	companyDataFile.CompanyDataFileTypeID = processedTypeFile
	companyDataFile.CreatedAt = time.Now()
	companyDataFile.Description = "Processed Data"

	err = m.db.Model(&companyDataFile).Create(&companyDataFile).Error
	if err != nil {
		return
	}

	req := []*model.CustomerInformation{}
	list, err := m.variableService.GetAll(claim, companyName, modelName)
	for i := 0; i < file.SheetCount; i++ {
		sheet1 := file.WorkBook.Sheets.Sheet[i].Name

		rows, _ := file.GetRows(sheet1)
		for _, row := range rows[1:] {

			request := model.CustomerInformation{}

			request.FirstName = row[0]
			request.LastName = row[1]
			request.CustomerIdProofNumber = row[2]
			request.CustomerIDProofType = row[3]
			request.ContactNumber = row[4]
			request.City = row[5]
			request.AccountID = row[6]
			items := []model.CustomerInformationItem{}
			for i, _ := range row[7:] {
				item := model.CustomerInformationItem{}
				list1 := model.ModelVariable{}
				err = m.db.Model(&list).
					Select("*").
					Where("name =?", rows[0][i+7]).
					Find(&list1).
					Error
				if err != nil {
					return
				}

				if reflect.ValueOf(list1).IsZero() {
					continue
				} else {
					item.CustomerInformationID = request.ID
					item.ModelVariableID = list1.ModelID
					item.Name = list1.Name
					ID, _ := strconv.ParseFloat(row[i+7], 64)
					item.Value = ID

				}
				items = append(items, item)
			}
			request.Items = items
			err = validator.ValidateCustomer(&request)
			if err != nil {
				return
			}
			req = append(req, &request)
		}

	}

	fmt.Println(len(req))
	for _, req1 := range req {
		_, err1 := m.customerInfoService.Create(claim, companyName, modelName, req1)
		if err1 != nil {
			panic(err1)
			return
		}
	}
	createdRecord := model.CompanyDataFile{}
	err = m.db.Debug().Model(&createdRecord).Where("file_id = ?", companyDataFile.ID).Take(&createdRecord).Error
	if err != nil {
		return
	}
	result = &createdRecord
	return
}

func (m FileServiceImpl) UploadRawFile(claim *auth.AuthenticatedClaim, companyName string, modelName string, fileheader string, file *excelize.File) (result *model.CompanyDataFile, err error) {
	if file == nil {
		return
	}

	companyDataFile := model.CompanyDataFile{}
	value, err := m.companyService.Get(claim, companyName)
	value1, err := m.modelService.Get(claim, companyName, modelName)

	companyDataFile.CompanyID = value.ID
	companyDataFile.ModelID = value1.ID
	companyDataFile.CreatedById = claim.UserId
	companyDataFile.Name = fileheader
	companyDataFile.CompanyDataFileTypeID = rawTypeFile
	companyDataFile.CreatedAt = time.Now()
	companyDataFile.Description = "Raw Data"

	err = m.db.Model(&companyDataFile).Create(&companyDataFile).Error
	if err != nil {
		return
	}
	createdRecord := model.CompanyDataFile{}
	err = m.db.Debug().Model(&createdRecord).Where("file_id = ?", companyDataFile.ID).Take(&createdRecord).Error
	if err != nil {
		return
	}
	user, err := m.userService.GetById(claim, 1)
	if err != nil {
		return
	}
	if file != nil {

		htmlBody := "<h3>Hi " + user.FirstName + "</h3>" +
			"<div style='border: 7px solid lightgray;margin-left: 100px;padding: 25px;width: 400px;text-align: center;'>" +

			"<div style='margin: 20px 0px 20px 0px;'>" +
			"<div>New file is uploaded By user of by Company </div>" +
			"<div>" + value.Name + "</div>" + "<div>For model</div>" +
			"<div>" + value1.Name + "</div>" +
			"</div>" +

			"<div>Thankyou for visiting Tausi App Website.</div>" +
			"</div>" +
			"</div>"
		var subject string = "Welcome to Tausi - Credit Scoring Engine"

		err = helper.SendEmailServiceSmtp1(user.Email, user.FirstName, subject, htmlBody, user, 1)

		if err != nil {
			return
		}

	}

	result = &createdRecord

	return
}

func (m FileServiceImpl) DownloadRawDataFile(claim *auth.AuthenticatedClaim, companyName string, fileName string) (result *excelize.File, result1 []model.CompanyDataFile, err error) {
	var company *model.Company
	company, err = m.companyService.Get(claim, companyName)
	if err != nil {
		return
	}
	result1 = []model.CompanyDataFile{}
	err = m.db.Model(&model.CompanyDataFile{}).Where("company_id = ? AND company_data_file_type_id = ?", company.ID, rawTypeFile).Find(&result1).Error
	if err != nil {
		panic(err)
		return
	}

	if fileName != "" {
		var name string
		err = m.db.Model(&result1).Limit(1000).Select("name").Where("name = ? ", fileName).Find(&name).Error
		if name == "" {

			err = ErrFileNameNotValid
		}
	}
	path := filepath.Join("./files/raw_data/", fileName)
	result, _ = excelize.OpenFile(path)
	// user, err := m.userService.GetByEmail(claim, claim.Email)
	// if err != nil {
	// 	return
	// }
	// if result != nil {

	// 	htmlBody := "<h3>Hi " + user.FirstName + "</h3>" +
	// 		"<div style='border: 7px solid lightgray;margin-left: 100px;padding: 25px;width: 400px;text-align: center;'>" +

	// 		"<div style='margin: 20px 0px 20px 0px;'>" +
	// 		"<div>File  downloaded Succesfully </div>" +

	// 		"</div>" +

	// 		"<div>Thankyou for visiting Tausi App Website.</div>" +
	// 		"</div>" +
	// 		"</div>"
	// 	var subject string = "Welcome to Tausi - Credit Scoring Engine"

	// 	err = helper.SendEmailServiceSmtp1(user.Email, user.FirstName, subject, htmlBody, user, 1)
	// 	if err != nil {
	// 		return
	// 	}
	// }

	return

}
