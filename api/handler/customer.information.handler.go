package handler

import (
	"encoding/json"
	"errors"

	"net/http"
	"strconv"

	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"
	"github.com/NEHA20-1992/tausi_code/api/validator"

	// "github.com/gofiber/fiber/v2"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var ErrCompanyMismatch = errors.New("Enter valid company,company name not found")
var ErrEnterValidCompany = errors.New("company name can't be empty")

type CustomerInformationHandler interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetAllCompayCustomer(w http.ResponseWriter, req *http.Request)
	GetAllCustomer(w http.ResponseWriter, req *http.Request)
	// GetAllCustomerx(w http.ResponseWriter, req *http.Request)

	Get(w http.ResponseWriter, req *http.Request)
	GetCreditScore(w http.ResponseWriter, req *http.Request)
}

type CustomerInformationHandlerImpl struct {
	service        service.CustomerInformationService
	companyService service.CompanyService
	modelService   service.ModelService
}

func GetCustomerInformationHandlerInstance(db *gorm.DB) (handler CustomerInformationHandler) {
	return CustomerInformationHandlerImpl{
		service:        service.GetCustomerInformationService(db),
		companyService: service.GetCompanyService(db),
		modelService:   service.GetModelService(db)}
}

func (h CustomerInformationHandlerImpl) Create(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var companyName string
	var modelName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]

	var request = model.CustomerInformation{}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = validator.ValidateCustomer(&request)
	if err != nil {
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
	}

	result, err := h.service.Create(claim, companyName, modelName, &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h CustomerInformationHandlerImpl) GetAllCustomer(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var request = model.CustomerFilterRequest{}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	result, err := h.service.GetAll(claim, request.CompanyName, request.ModelName, request.PageNumber, request.Size, request.Cid, request.Search, request.CreditScore)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)

}

func (h CustomerInformationHandlerImpl) GetAllCompayCustomer(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
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
	if aCompany.ID != claim.CompanyId {
		response.ERROR(w, http.StatusBadRequest, ErrCompanyMismatch)
	}

	result, err := h.service.GetAllCompayCustomer(claim, companyName, modelName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h CustomerInformationHandlerImpl) Get(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var companyName, modelName string
	var customerInformationId uint64
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]
	customerInformationId, err = strconv.ParseUint(vars["customerInformationId"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	aCompany, err := h.companyService.Get(claim, companyName)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if aCompany.ID != claim.CompanyId {
		response.ERROR(w, http.StatusBadRequest, ErrCompanyMismatch)
	}

	result, err := h.service.Get(claim, companyName, modelName, uint32(customerInformationId))
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h CustomerInformationHandlerImpl) GetCreditScore(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var companyName, modelName string
	var customerInformationId uint64
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]
	customerInformationId, err = strconv.ParseUint(vars["customerInformationId"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	aCompany, err := h.companyService.Get(claim, companyName)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if aCompany.ID != claim.CompanyId {
		response.ERROR(w, http.StatusBadRequest, ErrCompanyMismatch)
	}

	result, err := h.service.GetCreditScore(claim, companyName, modelName, uint32(customerInformationId))
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)

}
