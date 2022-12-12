package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/NEHA20-1992/tausi_code/api/auth"
	"github.com/NEHA20-1992/tausi_code/api/helper"
	"github.com/NEHA20-1992/tausi_code/api/model"
	"github.com/NEHA20-1992/tausi_code/pkg/config"
	"gorm.io/gorm"
)

var ErrAlreadyExists = errors.New("email already exist")
var ErrOldPasswordMismatch = errors.New("old password is mismatch")
var ErrAuthenticationInvalidEmail = errors.New("we can't seem to find your account")
var ErrUnableToSetResetCode = errors.New("unable to set reset code")

type UserService interface {
	Create(claim *auth.AuthenticatedClaim, companyName string, User *model.User) (*model.User, error)
	Update(claim *auth.AuthenticatedClaim, companyName string, existingEmailAddress string, User *model.User) (*model.User, error)
	ChangePassword(claim *auth.AuthenticatedClaim, ChangePasswordRequest *model.ChangePasswordRequest) (result bool, err error)
	GetByEmail(claim *auth.AuthenticatedClaim, emailAddress string) (result *model.User, err error)
	ResetPassword(claim *auth.AuthenticatedClaim, emailAddress string, password string) (result *model.User, err error)
	GetById(claim *auth.AuthenticatedClaim, userId uint32) (result *model.User, err error)
	GetUserNameById(claim *auth.AuthenticatedClaim, userId uint32) (userName string, err error)
	Get(claim *auth.AuthenticatedClaim, companyName string, emailAddress string) (*model.User, error)
	GetAll(claim *auth.AuthenticatedClaim, companyName string) ([]model.User, error)
	SetResetCode(claim *auth.AuthenticatedClaim, emailAddress string) (result *model.User, err error)
}

type UserServiceImpl struct {
	db              *gorm.DB
	companyService  CompanyService
	userRoleService UserRoleService
}

func GetUserService(db *gorm.DB) UserService {
	return UserServiceImpl{
		db:              db,
		companyService:  GetCompanyService(db),
		userRoleService: GetUserRoleService(db)}
}

func (m UserServiceImpl) Create(claim *auth.AuthenticatedClaim, companyName string, user *model.User) (result *model.User, err error) {
	var company *model.Company
	company, err = m.companyService.Get(claim, companyName)
	if err != nil {
		return
	}

	if user == nil {
		return
	}

	user.CompanyId = company.ID

	var existingUser model.User
	err = m.db.Model(&model.User{}).Select("user_id").Where("email = ?", user.Email).Find(&existingUser).Error
	if err != nil || existingUser.ID > 0 {
		err = ErrAlreadyExists
		return
	}

	err = m.create(claim, user)
	if err != nil {
		return
	}

	var createdUser model.User
	err = m.db.Model(&user).Where("user_id = ?", user.ID).First(&createdUser).Error
	if err == nil && createdUser.ID > 0 {
		err = m.loadUserRoles(&createdUser)
		if err != nil {
			return
		}

		var subject string = "Welcome to Tausi - Credit Scoring Engine"

		var htmlBody string = "<h3>Hi " + user.FirstName + "</h3>" +
			"<div style='border: 7px solid lightgray;margin-left: 100px;padding: 25px;width: 400px;text-align: center;'>" +

			"<div style='font-size: 18px;color: black;font-weight: bold;'>Welcome to Tausi</div>" +

			"<div style='margin: 20px 0px 20px 0px;'>" +
			"<div>Your Tausi account is created. You may click on the link below set your password,</div>" +
			"</div>" +

			"<div style='background-color: #000;padding: 10px 25px;text-align: center;display: inline-block;font-size: 13px;margin: 4px 2px;cursor: pointer;'><a style='text-decoration: none;color:#fff' href='" + config.ServerConfiguration.Amazonses.PasswordResetUrl +
			"?email=" + user.Email + "&resetCode=" + user.ResetCode + "'>Reset your password</a></div>" +
			"<div style='margin: 20px 0px 20px 0px;'>" +

			"<div>You can also change your password at any time from Change Password option after login with your current password.</div>" +
			"</div>" +
			"</div>"

		err = helper.SendEmailServiceSmtp(user.Email, user.FirstName, subject, htmlBody, user, 1)
		//err = helper.SendEmailService(user.Email, subject, htmlBody)
		if err != nil {
			return
		}

		createdUser.Password = ""
		err = m.updateMeta(claim, nil, &createdUser)
		if err != nil {
			return
		}
		result = &createdUser
		return
	}

	return
}

func (m UserServiceImpl) Update(claim *auth.AuthenticatedClaim, companyName string, existingEmailAddress string, user *model.User) (result *model.User, err error) {
	var company *model.Company
	company, err = m.companyService.Get(claim, companyName)
	if err != nil {
		return
	}
	user.CompanyId = company.ID

	if existingEmailAddress != user.Email {
		var existingUser2 model.User
		err = m.db.Model(&model.User{}).Where("email = ?", user.Email).Find(&existingUser2).Error
		if err != nil || existingUser2.ID > 0 {
			err = ErrAlreadyExists
			return
		}
	}

	var existingUser model.User
	err = m.db.Model(&model.User{}).Select("user_id").Where("email = ?", existingEmailAddress).First(&existingUser).Error
	if err != nil {
		return
	}
	user.ID = existingUser.ID

	var currentUserId uint32 = 1
	if claim != nil {
		currentUserId = claim.UserId
	}
	m.db.Debug().Model(&model.User{}).Where("user_id = ?", user.ID).Take(&model.User{}).UpdateColumns(
		map[string]interface{}{
			"first_name":         user.FirstName,
			"last_name":          user.LastName,
			"contact_number":     user.ContactNumber,
			"email":              user.Email,
			"email_verified":     user.EmailVerified,
			"active":             user.Active,
			"last_updated_by_id": currentUserId,
			"last_updated_at":    time.Now(),
		},
	)
	if m.db.Error != nil {
		return &model.User{}, m.db.Error
	}

	m.updateUserRoles(claim, user)

	// This is the display the updated user
	updatedRecord := model.User{}
	err = m.db.Debug().Model(&updatedRecord).Where("user_id = ?", user.ID).Take(&updatedRecord).Error
	if err != nil {
		return
	}
	m.loadUserRoles(&updatedRecord)
	updatedRecord.Password = ""
	err = m.updateMeta(claim, nil, &updatedRecord)
	if err != nil {
		return
	}
	result = &updatedRecord

	//updateUserRoles(user)

	return
}

func (m UserServiceImpl) ChangePassword(claim *auth.AuthenticatedClaim, changePasswordRequest *model.ChangePasswordRequest) (result bool, err error) {

	var existingUser model.User
	err = m.db.
		Model(&model.User{}).
		Where("user_id = ?", claim.UserId).
		First(&existingUser).
		Error
	if err != nil {
		return
	}

	err = existingUser.VerifyPassword(changePasswordRequest.Password)
	if err != nil {
		err = ErrOldPasswordMismatch
		return
	}

	var currentUserId uint32
	if claim != nil {
		currentUserId = claim.UserId
	}

	existingUser.Password = changePasswordRequest.NewPassword
	err = existingUser.Hash()
	if err != nil {
		panic(err)
	}

	m.db.Debug().Model(&model.User{}).Where("user_id = ?", claim.UserId).Take(&model.User{}).UpdateColumns(
		map[string]interface{}{
			"password":           existingUser.Password,
			"last_updated_by_id": currentUserId,
			"last_updated_at":    time.Now(),
		},
	)

	result = true

	return
}

func (m UserServiceImpl) GetByEmail(claim *auth.AuthenticatedClaim, emailAddress string) (result *model.User, err error) {
	record := model.User{}
	err = m.db.Where("email = ?", emailAddress).First(&record).Error
	if err != nil {
		return
	}
	result = &record
	err = m.updateMeta(claim, nil, &record)
	if err != nil {
		return
	}
	m.loadUserRoles(&record)

	return
}

func (m UserServiceImpl) ResetPassword(claim *auth.AuthenticatedClaim, emailAddress string, password string) (result *model.User, err error) {
	record := model.User{}
	err = m.db.Where("email = ?", emailAddress).First(&record).Error
	if err != nil {
		return
	}

	record.ResetCode = ""
	record.Password = password
	err = record.Hash()
	if err != nil {
		panic(err)
	}

	m.db.Debug().Model(&model.User{}).Where("email = ?", emailAddress).Take(&model.User{}).UpdateColumns(
		map[string]interface{}{
			"password":   record.Password,
			"reset_code": record.ResetCode,
		},
	)

	// m.db.Model(&record).UpdateColumn("password", record.Password)
	// err = m.db.Error

	result = &record

	return
}

func (m UserServiceImpl) SetResetCode(claim *auth.AuthenticatedClaim, emailAddress string) (result *model.User, err error) {
	record := model.User{}
	err = m.db.Where("email = ?", emailAddress).First(&record).Error
	if err != nil {
		err = ErrAuthenticationInvalidEmail
		return
	}

	record.ResetCode = model.RandStringBytes(6)
	err = m.db.Model(&record).UpdateColumn("reset_code", record.ResetCode).Error
	if err != nil {
		err = ErrUnableToSetResetCode
		return
	}
	result = &record
	var subject string = "Welcome to Tausi - Credit Scoring Engine"

	var htmlBody string = "<h3>Hi " + result.FirstName + "</h3>" +
		"<div style='border: 7px solid lightgray;margin-left: 100px;padding: 25px;width: 400px;text-align: center;'>" +

		"<div style='font-size: 18px;color: black;font-weight: bold;'>Forget your password?</div>" +

		"<div style='margin: 20px 0px 20px 0px;'>" +
		"<div>You have requested for a new password for the following account </div>" +
		"<div style='color=blue'>(" + result.Email + ")</div>" +
		"</div>" +

		"<div style='margin: 20px 0px 20px 0px;'>" +
		"<div>You may click on the link below reset/update your password:</div>" +
		"</div>" +

		"<div style='background-color: #000;padding: 10px 25px;text-align: center;display: inline-block;font-size: 13px;margin: 4px 2px;cursor: pointer;'><a style='text-decoration: none;color:#fff' href='" + config.ServerConfiguration.Amazonses.PasswordResetUrl +
		"?email=" + result.Email + "&resetCode=" + result.ResetCode + "'>Reset your password</a></div>" +
		"<div style='margin: 20px 0px 20px 0px;'>" +

		"<div>You can also change your password at any time from Change Password option after login with your current password.</div>" +
		"</div>" +
		"</div>"

	err = helper.SendEmailServiceSmtp(record.Email, record.FirstName, subject, htmlBody, result, 2)
	//err = helper.SendEmailService(record.Email, subject, htmlBody)
	return
}

func (m UserServiceImpl) GetById(claim *auth.AuthenticatedClaim, userId uint32) (result *model.User, err error) {
	record := model.User{}
	err = m.db.Where("user_id = ?", userId).First(&record).Error
	if err == nil {
		result = &record
	}
	record.Password = ""
	err = m.updateMeta(claim, nil, &record)
	if err != nil {
		return
	}

	m.loadUserRoles(&record)

	return
}

func (m UserServiceImpl) GetUserNameById(claim *auth.AuthenticatedClaim, userId uint32) (userName string, err error) {
	record := model.User{}
	err = m.db.Select("first_name", "last_name").Where("user_id = ?", userId).First(&record).Error
	if err == nil {
		userName = fmt.Sprintf("%s, %s", record.LastName, record.FirstName)
	}

	if err != nil {
		return
	}

	m.loadUserRoles(&record)

	return
}

func (m UserServiceImpl) Get(claim *auth.AuthenticatedClaim, companyName string, emailAddress string) (result *model.User, err error) {
	var company *model.Company
	company, err = m.companyService.Get(claim, companyName)
	if err != nil {
		return
	}

	record := model.User{}
	err = m.db.Where("company_id = ? AND email = ?", company.ID, emailAddress).First(&record).Error
	if err != nil {
		return
	}
	record.Password = ""
	result = &record

	err = m.updateMeta(claim, nil, &record)
	if err != nil {
		return
	}

	m.loadUserRoles(&record)

	return
}

func (m UserServiceImpl) GetAll(claim *auth.AuthenticatedClaim, companyName string) (result []model.User, err error) {
	var company *model.Company
	company, err = m.companyService.Get(claim, companyName)
	if err != nil {
		return
	}

	result = []model.User{}
	err = m.db.Model(&model.User{}).Limit(1000).Where("company_id = ?", company.ID).Find(&result).Error
	if err != nil {
		result = nil
	}

	var userMap map[int](*string) = make(map[int]*string)
	var resultList []model.User = make([]model.User, len(result))
	for indxValue, aRecord := range result {
		var newRecord = &aRecord

		err = m.updateMeta(claim, userMap, newRecord)
		if err != nil {
			return
		}

		newRecord.Password = ""
		m.loadUserRoles(newRecord)
		resultList[indxValue] = *newRecord
	}

	result = resultList

	return
}

func (m UserServiceImpl) updateMeta(claim *auth.AuthenticatedClaim, userMap map[int](*string), newRecord *model.User) (err error) {
	if userMap == nil {
		userMap = make(map[int]*string)
	}
	//config.ServerConfiguration.Application.PasswordResetUrl

	userName, err := getUserName(m.db, claim, userMap, newRecord.CreatedById)
	if err != nil {
		return
	}
	newRecord.CreatedBy = userName

	if newRecord.LastUpdatedById > 0 {
		userName, err = getUserName(m.db, claim, userMap, newRecord.LastUpdatedById)
		if err != nil {
			return
		}
		newRecord.LastUpdatedBy = userName
	}

	newRecord.Admin = (newRecord.CompanyId == 1)

	var company model.Company
	err = m.db.Model(&company).Where("company_id = ?", newRecord.CompanyId).First(&company).Error
	if err != nil {
		return
	}

	newRecord.CompanyName = company.Name

	return
}

func (m UserServiceImpl) create(claim *auth.AuthenticatedClaim, user *model.User) (err error) {
	var currentUserId uint32 = 1
	if claim != nil {
		currentUserId = claim.UserId
	}
	user.CreatedById = currentUserId
	passwordText := user.Password
	user.ResetCode = model.RandStringBytes(6)
	m.db.Model(&user).Omit("last_updated_by_id", "last_updated_at").Create(&user)
	err = m.db.Error
	if err != nil {
		return
	}

	if len(passwordText) > 0 {
		err = m.updatePassword(user, passwordText)
	}

	m.updateUserRoles(claim, user)

	return
}

func (m UserServiceImpl) updatePassword(user *model.User, passwordText string) (err error) {
	user.Password = passwordText
	err = user.Hash()
	if err != nil {
		return
	}

	m.db.Model(&user).UpdateColumn("password", user.Password)
	err = m.db.Error

	return
}

func (m UserServiceImpl) loadUserRoles(user *model.User) (err error) {
	var roles []model.UserRole
	m.db.Model(user).Association("UserRoles").Find(&roles)
	err = m.db.Error
	if err == nil {
		var roleNames []string = make([]string, len(roles))
		for indx, aRole := range roles {
			roleNames[indx] = aRole.Name
		}
		user.Roles = roleNames
	}

	return
}

func (m UserServiceImpl) updateUserRoles(claim *auth.AuthenticatedClaim, user *model.User) (err error) {
	m.db.Debug().Model(user).Association("UserRoles").Clear()
	if len(user.Roles) > 0 {
		var urList []model.UserRole = make([]model.UserRole, len(user.Roles))
		for indx, roleName := range user.Roles {
			var ur *model.UserRole
			ur, err = m.userRoleService.Get(roleName)
			if err != nil {
				return
			}
			urList[indx] = *ur
		}

		if len(urList) > 0 {
			m.db.Debug().Model(user).Association("UserRoles").Append(urList)
		}
	}

	return
}
