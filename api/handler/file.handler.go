package handler

import (
	"encoding/json"
	"errors"

	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	// "io"
	// "github.com/go-sql-driver/mysql"
	//"github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"

	// "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/NEHA20-1992/tausi_code/api/auth"

	"github.com/NEHA20-1992/tausi_code/api/model"
	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var ErrModelMismatch = errors.New("no such model is available for this company ")

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
		companyService: service.GetCompanyService(db),
		modelService:   service.GetModelService(db),
		service:        service.GetFileService(db)}
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
	var request = model.FileFilterRequest{}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	result, result1, err := h.service.DownloadRawDataFile(claim, companyName, request.FileName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if request.FileName == "" {
		response.JSON(w, http.StatusOK, result1)
	} else {
		response.JSONDOWNLOAD(w, http.StatusOK, result)
	}

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
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]

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

	if claim.CompanyId != 1 {
		response.ERROR(w, http.StatusUnauthorized, ErrNotAdmin)
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
	var companyName, modelName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]
	aCompany, err := h.companyService.Get(claim, companyName)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	path := filepath.Join("./files/analyzed_data/", strconv.FormatInt(time.Now().Unix(), 10)+aFileHeader.Filename)
	file, _ := excelize.OpenReader(aFile)

	if err := file.SaveAs(path); err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	aModel, err := h.modelService.Get(claim, aCompany.Name, modelName)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.UploadProcessedFile(claim, aCompany.Name, aModel.Name, filepath.Base(path), file)

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

	var File multipart.File
	var FileHeader *multipart.FileHeader
	File, FileHeader, err = req.FormFile("rawFile")
	if err != nil || FileHeader == nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var companyName, modelName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]
	file, _ := excelize.OpenReader(File)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	path := filepath.Join("./files/raw_data/", strconv.FormatInt(time.Now().Unix(), 10)+FileHeader.Filename)

	if err := file.SaveAs(path); err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	aCompany, err := h.companyService.Get(claim, companyName)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, ErrCompanyMismatch)
		return
	}
	if aCompany.ID != claim.CompanyId {
		response.ERROR(w, http.StatusBadRequest, ErrCompanyMismatch)
		return
	}
	_, err = h.modelService.Get(claim, aCompany.Name, modelName)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, ErrModelMismatch)
		return

	}

	result, err := h.service.UploadRawFile(claim, aCompany.Name, modelName, filepath.Base(path), file)

	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
