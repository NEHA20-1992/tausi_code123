package handler

import (
	"net/http"

	"github.com/NEHA20-1992/tausi_code/api/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InitializeRouters(db *gorm.DB, router *mux.Router) {
	registerAuthenticationHandlers(db, router)
	registerReferenceDataHandlers(db, router)
	registerCompanyHandlers(db, router)
	registerUserHandlers(db, router)
	registerModelHandlers(db, router)
	registerModelVariableHandlers(db, router)
	registerCustomerInformationHandlers(db, router)
	fileHandlers(db, router)
}

func registerAuthenticationHandlers(db *gorm.DB, router *mux.Router) {
	router.HandleFunc("/auth/login",
		middleware.JSONResponder(
			GetAuthenticationHandlerInstance(db).GenerateToken)).
		Methods(http.MethodPost)

	router.HandleFunc("/auth/getResetCode/{userEmail}",
		middleware.JSONResponder(
			GetAuthenticationHandlerInstance(db).GetResetCode)).
		Methods(http.MethodGet)

	router.HandleFunc("/auth/resetPassword",
		middleware.JSONResponder(
			GetAuthenticationHandlerInstance(db).ResetPassword)).
		Methods(http.MethodPost)

	router.HandleFunc("/auth/refresh",
		middleware.JSONResponder(
			GetAuthenticationHandlerInstance(db).RefreshToken)).
		Methods(http.MethodPost)
	router.HandleFunc("/auth/user",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetAuthenticationHandlerInstance(db).GetCurrentUser))).
		Methods(http.MethodGet)

	router.HandleFunc("/auth/changePassword",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetUserHandlerInstance(db).ChangePassword))).
		Methods(http.MethodPost)
}

func registerReferenceDataHandlers(db *gorm.DB, router *mux.Router) {
	router.HandleFunc("/countries",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetCountryHandlerInstance(db).GetAll))).
		Methods(http.MethodGet)

	router.HandleFunc("/cities",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetCityHandlerInstance(db).GetAll))).
		Methods(http.MethodGet)

	router.HandleFunc("/userroles",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetUserRoleHandlerInstance(db).GetAll))).
		Methods(http.MethodGet)

	router.HandleFunc("/companytypes",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetCompanyTypeHandlerInstance(db).GetAll))).
		Methods(http.MethodGet)

	router.HandleFunc("/customeridprooftypes",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetCustomerIdProofTypeHandlerInstance(db).GetAll))).
		Methods(http.MethodGet)
}

func registerCompanyHandlers(db *gorm.DB, router *mux.Router) {
	router.HandleFunc("/companies", middleware.JSONResponder(middleware.Authenticate(GetCompanyHandlerInstance(db).GetAll))).Methods(http.MethodGet)
	router.HandleFunc("/companies/{companyName}", middleware.JSONResponder(middleware.Authenticate(GetCompanyHandlerInstance(db).Get))).Methods(http.MethodGet)
	router.HandleFunc("/companies", middleware.JSONResponder(middleware.Authenticate(GetCompanyHandlerInstance(db).Create))).Methods(http.MethodPost)
	router.HandleFunc("/companies/{companyName}", middleware.JSONResponder(middleware.Authenticate(GetCompanyHandlerInstance(db).Update))).Methods(http.MethodPut)

	router.HandleFunc("/companies/getCount",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetCompanyHandlerInstance(db).GetCount))).
		Methods(http.MethodGet)

	router.HandleFunc("/companies/{companyName}/getCraCount",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetCompanyHandlerInstance(db).GetCraCount))).
		Methods(http.MethodGet)

	router.HandleFunc("/companies/{companyName}/logo",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetCompanyHandlerInstance(db).GetCompanyLogo))).
		Methods(http.MethodGet)

}

func registerUserHandlers(db *gorm.DB, router *mux.Router) {
	router.HandleFunc("/companies/{companyName}/users",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetUserHandlerInstance(db).GetAll))).
		Methods(http.MethodGet)

	router.HandleFunc("/companies/{companyName}/users/{emailAddress}",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetUserHandlerInstance(db).Get))).
		Methods(http.MethodGet)

	router.HandleFunc("/companies/{companyName}/users",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetUserHandlerInstance(db).Create))).
		Methods(http.MethodPost)

	router.HandleFunc("/companies/{companyName}/users/{emailAddress}",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetUserHandlerInstance(db).Update))).
		Methods(http.MethodPut)
}

func registerModelHandlers(db *gorm.DB, router *mux.Router) {
	router.HandleFunc("/companies/{companyName}/models", middleware.JSONResponder(middleware.Authenticate(GetModelHandlerInstance(db).GetAll))).Methods(http.MethodGet)
	router.HandleFunc("/companies/{companyName}/models/{modelName}", middleware.JSONResponder(middleware.Authenticate(GetModelHandlerInstance(db).Get))).Methods(http.MethodGet)
	router.HandleFunc("/companies/{companyName}/models", middleware.JSONResponder(middleware.Authenticate(GetModelHandlerInstance(db).Create))).Methods(http.MethodPost)
	router.HandleFunc("/companies/{companyName}/models/{modelName}", middleware.JSONResponder(middleware.Authenticate(GetModelHandlerInstance(db).Update))).Methods(http.MethodPut)

	router.HandleFunc("/companies/{companyName}/models/{modelName}/details", middleware.JSONResponder(middleware.Authenticate(GetModelHandlerInstance(db).GetDetails))).Methods(http.MethodGet)

	router.HandleFunc("/companies/{companyName}/models/{modelName}/customerCount",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetModelHandlerInstance(db).GetCustomerCount))).
		Methods(http.MethodGet)

}

func registerModelVariableHandlers(db *gorm.DB, router *mux.Router) {
	router.HandleFunc("/companies/{companyName}/models/{modelName}/variables", middleware.JSONResponder(middleware.Authenticate(GetModelVariableHandlerInstance(db).GetAll))).Methods(http.MethodGet)
	router.HandleFunc("/companies/{companyName}/models/{modelName}/variables/{variableName}", middleware.JSONResponder(middleware.Authenticate(GetModelVariableHandlerInstance(db).Get))).Methods(http.MethodGet)
	router.HandleFunc("/companies/{companyName}/models/{modelName}/variables", middleware.JSONResponder(middleware.Authenticate(GetModelVariableHandlerInstance(db).Create))).Methods(http.MethodPost)
	router.HandleFunc("/companies/{companyName}/models/{modelName}/variables/{variableName}", middleware.JSONResponder(middleware.Authenticate(GetModelVariableHandlerInstance(db).Update))).Methods(http.MethodPut)
	router.HandleFunc("/companies/{companyName}/models/{modelName}/variables/{variableName}", middleware.JSONResponder(middleware.Authenticate(GetModelVariableHandlerInstance(db).Delete))).Methods(http.MethodDelete)
}

func registerCustomerInformationHandlers(db *gorm.DB, router *mux.Router) {
	// router.HandleFunc("/companies/customers",func(c *fiber.Ctx)error{
	// 	info =[]model.CustomerInformation

	// 	sql := "SELECT * FROM info"
	// 	if s := c.Query("s");s!=""{
	// 		sql =fmt.Sprintf("%s where first_name like '%%%s%%' or first_name like '%%%s%%'",sql,s,s)
	// 	}
	// 	if sort :=c.Query("sort");sort != {
	// 		sql := fmt.Sprintf("%s ORDER BY income %s",sql,sort)
	// 	}
	// 	page ,_ :=strconv.Atoi(C.Query("page",1))
	// 	perPage :=9
	// 	db.Raw(sql).Scan(&info)
	// 	return c.JSON(info)
	// })
	router.HandleFunc("/companies/customers", middleware.JSONResponder(middleware.Authenticate(GetCustomerInformationHandlerInstance(db).GetAllCustomer))).Methods(http.MethodPost)
	router.HandleFunc("/companies/{companyName}/models/{modelName}/data", middleware.JSONResponder(middleware.Authenticate(GetCustomerInformationHandlerInstance(db).GetAllCompayCustomer))).Methods(http.MethodGet)
	router.HandleFunc("/companies/{companyName}/models/{modelName}/data/{customerInformationId}", middleware.JSONResponder(middleware.Authenticate(GetCustomerInformationHandlerInstance(db).Get))).Methods(http.MethodGet)
	router.HandleFunc("/companies/{companyName}/models/{modelName}/getCreditScore/{customerInformationId}", middleware.JSONResponder(middleware.Authenticate(GetCustomerInformationHandlerInstance(db).GetCreditScore))).Methods(http.MethodGet)
	router.HandleFunc("/companies/{companyName}/models/{modelName}/data", middleware.JSONResponder(middleware.Authenticate(GetCustomerInformationHandlerInstance(db).Create))).Methods(http.MethodPost)
}

func fileHandlers(db *gorm.DB, router *mux.Router) {
	router.HandleFunc("/companies/{companyName}/rawFiles",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetFileHandlerInstance(db).GetAllRawFile))).
		Methods(http.MethodPut)

	router.HandleFunc("/companies/{companyName}/models/{modelName}/processedFiles",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetFileHandlerInstance(db).GetAllProcessedFile))).
		Methods(http.MethodPut)

	router.HandleFunc("/companies/{companyName}/uploadRawFiles",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetFileHandlerInstance(db).UploadRawFile))).
		Methods(http.MethodPost)

	router.HandleFunc("/companies/{companyName}/models/{modelName}/uploadProcessedFiles",
		middleware.JSONResponder(
			middleware.Authenticate(
				GetFileHandlerInstance(db).UploadProcessedFile))).
		Methods(http.MethodPost)

}
