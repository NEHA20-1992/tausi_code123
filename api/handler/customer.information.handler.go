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

var ErrNoData = errors.New("no data found")
var ErrCompanyNotFound = errors.New("enter valid company")
var ErrModelNotFound = errors.New("enter valid model name")
var ErrCompanyMismatch = errors.New("company mismatch")

type CustomerInformationHandler interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetAllCompayCustomer(w http.ResponseWriter, req *http.Request)
	GetAllCustomerExcel(w http.ResponseWriter, req *http.Request)
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
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if aCompany.ID != claim.CompanyId {
		response.ERROR(w, http.StatusBadRequest, err)
	}

	result, err := h.service.Create(claim, companyName, modelName, &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h CustomerInformationHandlerImpl) GetAllCustomerExcel(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	// if claim.CompanyId != 1 {
	// 	response.ERROR(w, http.StatusUnauthorized, err)
	// 	return
	// }

	var request = model.CustomerFilterRequest{}
	// queryData := req.FormValue("query")
	// // if aCompanyData == "" {
	// // 	response.ERROR(w, http.StatusBadRequest, err)
	// // 	return
	// // }
	// queryData, err = url.PathUnescape(queryData)
	// if err != nil {
	// 	response.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	// var DataBytes []byte = []byte(queryData)
	// err = json.Unmarshal(DataBytes, &request)
	// if err != nil {
	// 	response.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	// request.Query = req.FormValue("Query")
	request.CompanyName = req.FormValue("CompanyName")
	request.ModelName = req.FormValue("ModelName")
	request.City = req.FormValue("City")
	MinGroupScore := req.FormValue("MinGroupScore")
	request.MinGroupScore, err = strconv.ParseFloat(MinGroupScore, 64)
	MaxGroupScore := req.FormValue("MaxGroupScore")
	request.MaxGroupScore, err = strconv.ParseFloat(MaxGroupScore, 64)
	MinPercentage := req.FormValue("MinPercentage")
	request.MinPercentage, err = strconv.ParseFloat(MinPercentage, 64)
	MaxPercentage := req.FormValue("MaxPercentage")
	request.MaxPercentage, err = strconv.ParseFloat(MaxPercentage, 64)
	request.Sort = req.FormValue("Sort")
	PageNumber, err := strconv.ParseUint(req.FormValue("PageNumber"), 10, 64)
	request.PageNumber = uint32(PageNumber)
	Size, err := strconv.ParseUint(req.FormValue("Size"), 10, 64)
	request.Size = uint32(Size)
	if request.CompanyName != "" {
		c, err := h.companyService.Get(claim, request.CompanyName)
		if err != nil {
			response.ERROR(w, http.StatusBadRequest, ErrCompanyNotFound)
			return
		}
		request.CID = c.ID
	}
	if request.ModelName != "" {
		m, err := h.modelService.Get(claim, request.CompanyName, request.ModelName)
		if err != nil {
			response.ERROR(w, http.StatusBadRequest, ErrModelNotFound)
			return
		}
		request.MID = m.ID
	}
	result, err := h.service.GetAllExcel(claim, &request)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if result == nil {
		response.ERROR(w, http.StatusOK, ErrNoData)

	} else {
		response.JSONDOWNLOAD(w, http.StatusOK, result)
	}
}

func (h CustomerInformationHandlerImpl) GetAllCustomer(w http.ResponseWriter, req *http.Request) {
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var request = model.CustomerFilterRequest{}
	// queryData := req.FormValue("query")

	// queryData, err = url.PathUnescape(queryData)
	// if err != nil {
	// 	response.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	// var DataBytes []byte = []byte(queryData)
	// err = json.Unmarshal(DataBytes, &request)
	// if err != nil {
	// 	response.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	// request.Query = req.FormValue("Query")
	request.CompanyName = req.FormValue("CompanyName")
	request.ModelName = req.FormValue("ModelName")
	request.City = req.FormValue("City")
	MinGroupScore := req.FormValue("MinGroupScore")
	request.MinGroupScore, err = strconv.ParseFloat(MinGroupScore, 64)
	MaxGroupScore := req.FormValue("MaxGroupScore")
	request.MaxGroupScore, err = strconv.ParseFloat(MaxGroupScore, 64)
	MinPercentage := req.FormValue("MinPercentage")
	request.MinPercentage, err = strconv.ParseFloat(MinPercentage, 64)
	MaxPercentage := req.FormValue("MaxPercentage")
	request.MaxPercentage, err = strconv.ParseFloat(MaxPercentage, 64)
	request.Sort = req.FormValue("Sort")
	PageNumber, err := strconv.ParseUint(req.FormValue("PageNumber"), 10, 64)
	request.PageNumber = uint32(PageNumber)
	Size, err := strconv.ParseUint(req.FormValue("Size"), 10, 64)
	request.Size = uint32(Size)
	if request.CompanyName != "" {
		aCompany, err := h.companyService.Get(claim, request.CompanyName)
		if err != nil {
			response.ERROR(w, http.StatusBadRequest, ErrCompanyNotFound)
			return
		}
		request.CID = aCompany.ID
	}
	if request.ModelName != "" {
		aModel, err := h.modelService.Get(claim, request.CompanyName, request.ModelName)
		if err != nil {
			response.ERROR(w, http.StatusBadRequest, ErrModelNotFound)
			return
		}
		request.MID = aModel.ID
	}
	result, err := h.service.GetAll(claim, &request)
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
		response.ERROR(w, http.StatusBadRequest, err)
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
