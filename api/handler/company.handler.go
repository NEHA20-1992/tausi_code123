package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"
	"github.com/NEHA20-1992/tausi_code/api/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var ErrTheLogoTooLarge = errors.New("The logo is too large. logo must be smaller than 1mb.")

type CompanyHandler interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetAll(w http.ResponseWriter, req *http.Request)
	GetCount(w http.ResponseWriter, req *http.Request)
	GetCraCount(w http.ResponseWriter, req *http.Request)
	Get(w http.ResponseWriter, req *http.Request)
	GetById(w http.ResponseWriter, req *http.Request)
	GetCompanyLogo(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
}

type CompanyHandlerImpl struct {
	service     service.CompanyService
	userService service.UserService
}

func GetCompanyHandlerInstance(db *gorm.DB) (handler CompanyHandler) {
	return CompanyHandlerImpl{
		service:     service.GetCompanyService(db),
		userService: service.GetUserService(db)}
}

func (h CompanyHandlerImpl) Create(w http.ResponseWriter, req *http.Request) {
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
	aFile, aFileHeader, err = req.FormFile("logo")
	if err != nil || aFileHeader == nil {
		// response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var logoContents []byte
	if aFile != nil && aFileHeader != nil {
		defer aFile.Close()
		logoContents, err = ioutil.ReadAll(aFile)
		if err != nil || logoContents == nil {
			response.ERROR(w, http.StatusBadRequest, err)
			return
		}
	}

	aCompanyData := req.FormValue("data")
	if aCompanyData == "" {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	aCompanyData, err = url.PathUnescape(aCompanyData)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var companyDataBytes []byte = []byte(aCompanyData)
	var request = model.Company{}
	err = json.Unmarshal(companyDataBytes, &request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = validator.ValidateCompany(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.Create(claim, &request, logoContents)

	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h CompanyHandlerImpl) GetCount(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	result, err := h.service.GetEntityCount(claim)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h CompanyHandlerImpl) GetCraCount(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var companyName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]

	result, err := h.service.GetCraEntityCount(claim, companyName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h CompanyHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	result, err := h.service.GetAll(claim)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h CompanyHandlerImpl) GetCompanyLogo(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var companyName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]

	result, err := h.service.GetCompanyLogo(claim, companyName)

	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
func (h CompanyHandlerImpl) GetById(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var companyID uint32
	var vars = mux.Vars(req)
	u64, err := strconv.ParseUint(vars["companyID"], 10, 32)
	companyID = uint32(u64)

	result, err := h.service.GetById(claim, companyID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
func (h CompanyHandlerImpl) Get(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	var companyName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]

	result, err := h.service.Get(claim, companyName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)

}

func (h CompanyHandlerImpl) Update(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// var request = model.Company{}
	// err = json.NewDecoder(req.Body).Decode(&request)
	// if err != nil {
	// 	response.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	err = req.ParseMultipartForm(1024 * 1024 * 10)
	if err != nil {
		panic(err)
	}

	var aFile multipart.File
	var aFileHeader *multipart.FileHeader
	aFile, aFileHeader, err = req.FormFile("logo")
	if err != nil || aFileHeader == nil {
		return
	}

	var logoContents []byte
	if aFile != nil && aFileHeader != nil {
		defer aFile.Close()
		logoContents, err = ioutil.ReadAll(aFile)
		if err != nil || logoContents == nil {
			response.ERROR(w, http.StatusBadRequest, err)
			return
		}
	}

	aCompanyData := req.FormValue("data")
	if aCompanyData == "" {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	aCompanyData, err = url.PathUnescape(aCompanyData)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var companyDataBytes []byte = []byte(aCompanyData)
	var request = model.Company{}
	err = json.Unmarshal(companyDataBytes, &request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = validator.ValidateCompany(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var existingCompanyName string
	var vars = mux.Vars(req)
	existingCompanyName = vars["companyName"]
	// request.Name = companyName

	result, err := h.service.Update(claim, existingCompanyName, &request, logoContents)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
