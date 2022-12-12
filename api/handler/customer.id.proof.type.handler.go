package handler

import (
	"net/http"

	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"

	"gorm.io/gorm"
)

type CustomerIdProofTypeHandler interface {
	GetAll(w http.ResponseWriter, req *http.Request)
}

type CustomerIdProofTypeHandlerImpl struct {
	service service.CustomerIdProofTypeService
}

func GetCustomerIdProofTypeHandlerInstance(db *gorm.DB) (handler CustomerIdProofTypeHandler) {
	return CustomerIdProofTypeHandlerImpl{service: service.GetCustomerIdProofTypeService(db)}
}

func (h CustomerIdProofTypeHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
	result, err := h.service.GetAll()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}
