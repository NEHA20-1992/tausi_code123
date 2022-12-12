package handler

import (
	"net/http"

	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"
	"gorm.io/gorm"
)

type CompanyTypeHandler interface {
	GetAll(w http.ResponseWriter, req *http.Request)
}

type CompanyTypeHandlerImpl struct {
	service service.CompanyTypeService
}

func GetCompanyTypeHandlerInstance(db *gorm.DB) (handler CompanyTypeHandler) {
	return CompanyTypeHandlerImpl{service: service.GetCompanyTypeService(db)}
}

func (h CompanyTypeHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
	result, err := h.service.GetAll()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}
