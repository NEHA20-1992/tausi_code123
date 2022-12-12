package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"github.com/NEHA20-1992/tausi_code/api/response"
	"github.com/NEHA20-1992/tausi_code/api/service"
	"github.com/NEHA20-1992/tausi_code/api/validator"
	"github.com/NEHA20-1992/tausi_code/pkg/logger"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var ErrAuthenticationMismatch = errors.New("Email/Password mismatch")
var ErrAuthenticationInvalidEmail = errors.New("We can't seem to find your account")
var ErrTokenPurposeMismatch = errors.New("valid referesh token is required")
var ErrAccessTokenRequired = errors.New("valid access token is required")
var ErrAuthenticationResetCodeMismatch = errors.New("reset code is mismatch")
var ErrAuthenticationNewAndConfirmPasswordMismatch = errors.New("New/Confirm password mismatch")

type AuthenticationHandler interface {
	GenerateToken(w http.ResponseWriter, req *http.Request)
	RefreshToken(w http.ResponseWriter, req *http.Request)
	GetCurrentUser(w http.ResponseWriter, req *http.Request)
	ResetPassword(w http.ResponseWriter, req *http.Request)
	GetResetCode(w http.ResponseWriter, req *http.Request)
}

type AuthenticationHandlerImpl struct {
	service service.UserService
}

func GetAuthenticationHandlerInstance(db *gorm.DB) (handler AuthenticationHandler) {
	return AuthenticationHandlerImpl{service: service.GetUserService(db)}
}

func (h AuthenticationHandlerImpl) GenerateToken(w http.ResponseWriter, req *http.Request) {
	var authRequest = model.AuthenticationRequest{}
	err := json.NewDecoder(req.Body).Decode(&authRequest)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.service.GetByEmail(nil, authRequest.Email)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, ErrAuthenticationMismatch)
		return
	}

	err = u.VerifyPassword(authRequest.Password)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, ErrAuthenticationMismatch)
		return
	}

	authResponse, err := auth.GenerateAuthenticationResponse(u, time.Hour*1, time.Hour*24)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, authResponse)
}

func (h AuthenticationHandlerImpl) RefreshToken(w http.ResponseWriter, req *http.Request) {
	var refreshAuthRequest = model.RefreshAuthenticationRequest{}
	err := json.NewDecoder(req.Body).Decode(&refreshAuthRequest)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var claim *auth.AuthenticatedClaim
	claim, err = auth.ValidateTokenString(refreshAuthRequest.RefreshToken)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	logger.AccessLogger.Println(response.ToJSON(claim))
	if claim.Purpose != "refresh" {
		response.ERROR(w, http.StatusBadRequest, ErrTokenPurposeMismatch)
		return
	}

	var claimUser *model.User
	claimUser, err = h.service.GetById(claim, claim.UserId)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	authResponse, err := auth.GenerateAuthenticationResponse(claimUser, time.Hour*1, time.Hour*24)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	response.JSON(w, http.StatusOK, authResponse)
}

func (h AuthenticationHandlerImpl) GetCurrentUser(w http.ResponseWriter, req *http.Request) {
	var claim *auth.AuthenticatedClaim
	claim, err := auth.ValidateToken(req)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	logger.AccessLogger.Println(response.ToJSON(claim))
	if claim.Purpose != "access" {
		response.ERROR(w, http.StatusBadRequest, ErrAccessTokenRequired)
		return
	}

	var claimUser *model.User
	claimUser, err = h.service.GetById(claim, claim.UserId)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	response.JSON(w, http.StatusOK, claimUser)
}

func (h AuthenticationHandlerImpl) GetResetCode(w http.ResponseWriter, req *http.Request) {
	var userEmail string
	var vars = mux.Vars(req)
	userEmail = vars["userEmail"]

	user, err := h.service.SetResetCode(nil, userEmail)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, ErrAuthenticationMismatch)
		return
	}
	var resp string
	if user != nil {
		resp = "Reset code successfully send to your email."
	}
	response.JSON(w, http.StatusOK, resp)
}

func (h AuthenticationHandlerImpl) ResetPassword(w http.ResponseWriter, req *http.Request) {

	var resetPasswordRequest = model.ResetPasswordRequest{}
	err := json.NewDecoder(req.Body).Decode(&resetPasswordRequest)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = validator.ValidateResetPassword(&resetPasswordRequest)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.service.GetByEmail(nil, resetPasswordRequest.Email)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, ErrAuthenticationInvalidEmail)
		return
	}

	if u.ResetCode != resetPasswordRequest.ResetCode {
		response.ERROR(w, http.StatusBadRequest, ErrAuthenticationResetCodeMismatch)
		return
	}

	if resetPasswordRequest.NewPassword != resetPasswordRequest.ConfirmPassword {
		response.ERROR(w, http.StatusBadRequest, ErrAuthenticationNewAndConfirmPasswordMismatch)
		return
	}

	u, err = h.service.ResetPassword(nil, resetPasswordRequest.Email, resetPasswordRequest.NewPassword)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, ErrAuthenticationMismatch)
		return
	}

	authResponse, err := auth.GenerateAuthenticationResponse(u, time.Hour*1, time.Hour*24)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, authResponse)
}
