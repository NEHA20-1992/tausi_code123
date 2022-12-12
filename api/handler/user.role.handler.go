package handler

import (
	"net/http"

	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"
	"gorm.io/gorm"
)

type UserRoleHandler interface {
	GetAll(w http.ResponseWriter, req *http.Request)
}

type UserRoleHandlerImpl struct {
	service service.UserRoleService
}

func GetUserRoleHandlerInstance(db *gorm.DB) (handler UserRoleHandler) {
	return UserRoleHandlerImpl{service: service.GetUserRoleService(db)}
}

func (h UserRoleHandlerImpl) GetAll(w http.ResponseWriter, req *http.Request) {
	result, err := h.service.GetAll()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}
