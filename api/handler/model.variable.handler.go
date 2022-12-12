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

type ModelVariableHandler interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetAll(w http.ResponseWriter, req *http.Request)
	Get(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
}

type ModelVariableHandlerImpl struct {
	service service.ModelVariableService
}

func GetModelVariableHandlerInstance(db *gorm.DB) (handler ModelVariableHandler) {
	return ModelVariableHandlerImpl{
		service: service.GetModelVariableService(db)}
}

func (h ModelVariableHandlerImpl) Create(w http.ResponseWriter, req *http.Request) {
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
	var modelName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]

	var request = model.ModelVariable{}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.Create(claim, companyName, modelName, &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h ModelVariableHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// if claim.CompanyId != 1 {
	// 	response.ERROR(w, http.StatusUnauthorized, ErrNotAdmin)
	// 	return
	// }

	var companyName, modelName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]

	result, err := h.service.GetAll(claim, companyName, modelName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h ModelVariableHandlerImpl) Get(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if claim.CompanyId != 1 {
		response.ERROR(w, http.StatusUnauthorized, ErrNotAdmin)
		return
	}

	var companyName, modelName, variableName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]
	variableName = vars["variableName"]

	result, err := h.service.Get(claim, companyName, modelName, variableName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)

}

func (h ModelVariableHandlerImpl) Update(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if claim.CompanyId != 1 {
		response.ERROR(w, http.StatusUnauthorized, ErrNotAdmin)
		return
	}

	var companyName, modelName, variableName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]
	variableName = vars["variableName"]

	var request = model.ModelVariable{}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//request.Name = variableName

	result, err := h.service.Update(claim, companyName, modelName, variableName, &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h ModelVariableHandlerImpl) Delete(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if claim.CompanyId != 1 {
		response.ERROR(w, http.StatusUnauthorized, ErrNotAdmin)
		return
	}

	var companyName, modelName, variableName string
	var vars = mux.Vars(req)
	companyName = vars["companyName"]
	modelName = vars["modelName"]
	variableName = vars["variableName"]

	result, err := h.service.Delete(claim, companyName, modelName, variableName)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
