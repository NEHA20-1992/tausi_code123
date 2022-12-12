package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ModelHandler interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetAll(w http.ResponseWriter, req *http.Request)
	Get(w http.ResponseWriter, req *http.Request)
	GetDetails(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	GetCustomerCount(w http.ResponseWriter, req *http.Request)
}

type ModelHandlerImpl struct {
	service        service.ModelService
	companyService service.CompanyService
}

func GetModelHandlerInstance(db *gorm.DB) (handler ModelHandler) {
	return ModelHandlerImpl{
		service:        service.GetModelService(db),
		companyService: service.GetCompanyService(db)}
}

func (h ModelHandlerImpl) Create(w http.ResponseWriter, req *http.Request) {
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

	var request = model.Model{}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.Create(claim, companyName, &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h ModelHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var companyName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]

	aCompany, err := h.companyService.Get(claim, companyName)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if aCompany.ID != claim.CompanyId && claim.CompanyId != 1 {
		response.ERROR(w, http.StatusBadRequest, ErrCompanyMismatch)
	}

	result, err := h.service.GetAll(claim, companyName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h ModelHandlerImpl) Get(w http.ResponseWriter, req *http.Request) {
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

	if aCompany.ID != claim.CompanyId && claim.CompanyId != 1 {
		response.ERROR(w, http.StatusBadRequest, ErrCompanyMismatch)
	}

	result, err := h.service.Get(claim, companyName, modelName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h ModelHandlerImpl) GetCustomerCount(w http.ResponseWriter, req *http.Request) {
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

	if aCompany.ID != claim.CompanyId && claim.CompanyId != 1 {
		response.ERROR(w, http.StatusBadRequest, ErrCompanyMismatch)
	}

	result, err := h.service.GetCustomerCount(claim, companyName, modelName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h ModelHandlerImpl) GetDetails(w http.ResponseWriter, req *http.Request) {
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

	if aCompany.ID != claim.CompanyId && claim.CompanyId != 1 {
		response.ERROR(w, http.StatusBadRequest, ErrCompanyMismatch)
	}

	result, err := h.service.GetDetails(claim, companyName, modelName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)

}

func (h ModelHandlerImpl) Update(w http.ResponseWriter, req *http.Request) {
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

	var request = model.Model{}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.Update(claim, companyName, modelName, &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
