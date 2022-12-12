package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"

	// "github.com/xuri/excelize@latest"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"
	"github.com/NEHA20-1992/tausi_code/api/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var err1 = errors.New("not able to store file")
var err2 = errors.New("not able to unmarshal file")

type FileHandler interface {
	GetAllRawFile(w http.ResponseWriter, req *http.Request)
	GetAllProcessedFile(w http.ResponseWriter, req *http.Request)
	UploadRawFile(w http.ResponseWriter, req *http.Request)
	UploadProcessedFile(w http.ResponseWriter, req *http.Request)
}

type FileHandlerImpl struct {
	service             service.FileService
	companyService      service.CompanyService
	IdTypeService       service.CustomerIdProofTypeService
	VariableService     service.ModelVariableService
	customerInfoService service.CustomerInformationService
	modelService        service.ModelService
}

func GetFileHandlerInstance(db *gorm.DB) (handler FileHandler) {
	return FileHandlerImpl{
		service: service.GetFileService(db)}
}

func (h FileHandlerImpl) GetAllRawFile(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if claim.CompanyId != 1 {
		response.ERROR(w, http.StatusUnauthorized, ErrNotAdmin)
		return
	}

	var companyName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	// int1, _ := strconv.ParseUint(vars["fileTypeId"], 10, 32)
	// ui := uint32(int1)
	// fileTypeId = ui

	result, err := h.service.GetAllRawFile(claim, companyName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h FileHandlerImpl) GetAllProcessedFile(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if claim.CompanyId != 1 {
		response.ERROR(w, http.StatusUnauthorized, ErrNotAdmin)
		return
	}

	var companyName, modelName string
	// var fileTypeId uint32
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]
	// int1, _ := strconv.ParseUint(vars["fileTypeId"], 10, 32)
	// ui := uint32(int1)
	// fileTypeId = ui

	result, err := h.service.GetAllProcessedFile(claim, companyName, modelName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h FileHandlerImpl) UploadProcessedFile(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// Upto 10MB in memory and rest in temp files.
	err = req.ParseMultipartForm(1024 * 1024 * 10)
	if err != nil {
		panic(err)
	}

	var aFile multipart.File
	var aFileHeader *multipart.FileHeader
	aFile, aFileHeader, err = req.FormFile("ProcessedFile")
	if err != nil || aFileHeader == nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	f, err := os.OpenFile("C:/Users/Test/go/tausi_code123/analyzed_data/"+aFileHeader.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
	}
	defer f.Close()
	io.Copy(f, aFile)

	file1, err := excelize.OpenFile("C:/Users/Test/go/tausi_code123/analyzed_data/" + aFileHeader.Filename)
	if err != nil {

		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	sheet1 := file1.WorkBook.Sheets.Sheet[0].Name
	rows := file1.GetRows(sheet1)
	var companyName, modelName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err2)
		return
	}
	request1 := []model.CustomerInformation{}
	for _, row := range rows {
		request := model.CustomerInformation{}
		request.FirstName = row[0]
		request.LastName = row[1]
		request.CustomerIdProofNumber = row[2]
		request.CustomerIDProofType = row[3]
		request.ContactNumber = row[4]
		request.City = row[5]
		ID5, _ := strconv.ParseFloat(row[6], 64)
		ID6, _ := strconv.ParseFloat(row[7], 64)
		ID7, _ := strconv.ParseFloat(row[8], 64)
		ID8, _ := strconv.ParseFloat(row[9], 64)
		ID9, _ := strconv.ParseFloat(row[10], 64)
		ID10, _ := strconv.ParseFloat(row[10], 64)
		items := []model.CustomerInformationItem{
			{
				Name:  "Group Allocatable Income",
				Value: ID5,
			},

			{
				Name:  "Income",
				Value: ID6,
			},

			{
				Name:  "Income Ratio",
				Value: ID7,
			},

			{
				Name:  "No of dependents",
				Value: ID8,
			},

			{
				Name:  "Marital Status",
				Value: ID9,
			},
			{
				Name:  "Gender",
				Value: ID10,
			},
		}

		request.Items = items

		request1 = append(request1, request)
	}

	result, err := h.service.UploadProcessedFile(claim, companyName, modelName, request1, aFileHeader)

	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)

}

func (h FileHandlerImpl) UploadRawFile(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// Upto 10MB in memory and rest in temp files.
	err = req.ParseMultipartForm(1024 * 1024 * 10)
	if err != nil {
		panic(err)
	}

	var aFile multipart.File
	var aFileHeader *multipart.FileHeader
	aFile, aFileHeader, err = req.FormFile("rawFile")
	if err != nil || aFileHeader == nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var fileContents []byte
	if aFile != nil && aFileHeader != nil {
		defer aFile.Close()
		fileContents, err = ioutil.ReadAll(aFile)
		if err != nil || fileContents == nil {
			response.ERROR(w, http.StatusBadRequest, err)
			return
		}
	}

	aCompanyData := req.FormValue("rawData")
	if aCompanyData == "" {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// fmt.Println(req.FormValue("rawData"))
	fmt.Println(aCompanyData)

	aCompanyData, err = url.PathUnescape(aCompanyData)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println(aCompanyData)
	var companyDataBytes []byte = []byte(aCompanyData)
	fmt.Println(companyDataBytes)
	var request = model.CompanyDataFile{}
	err = json.Unmarshal(companyDataBytes, &request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println(request)
	err = validator.ValidateCompanyDataFile(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.UploadRawFile(claim, &request)

	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)

}
