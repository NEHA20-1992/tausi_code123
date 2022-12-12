package handler

import (
	"net/http"

	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"
	"gorm.io/gorm"
)

type CityHandler interface {
	GetAll(w http.ResponseWriter, req *http.Request)
}

type CityHandlerImpl struct {
	service service.CityService
}

func GetCityHandlerInstance(db *gorm.DB) (handler CityHandler) {
	return CityHandlerImpl{service: service.GetCityService(db)}
}

func (h CityHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
	result, err := h.service.GetAll()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
