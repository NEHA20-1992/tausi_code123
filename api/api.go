package api

import (
	"github.com/NEHA20-1992/tausi_code/api/handler"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Initialize(DB *gorm.DB, router *mux.Router) {
	handler.InitializeRouters(DB, router)
}
