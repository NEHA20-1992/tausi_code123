package handler

import (
	"net/http"

	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"
	"gorm.io/gorm"
)

type CountryHandler interface {
	GetAll(w http.ResponseWriter, req *http.Request)
}

type CountryHandlerImpl struct {
	service service.CountryService
}

func GetCountryHandlerInstance(db *gorm.DB) (handler CountryHandler) {
	return CountryHandlerImpl{service: service.GetCountryService(db)}
}

func (h CountryHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
	result, err := h.service.GetAll()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
