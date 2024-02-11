package service

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/n9mi/go-docker/exception"
	"github.com/n9mi/go-docker/helper"
	"github.com/n9mi/go-docker/model/domain"
	"github.com/n9mi/go-docker/model/web"
	"github.com/n9mi/go-docker/repository"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(registerReq web.RegisterRequest) error
	Login(loginReq web.LoginRequest) (web.LoginResponse, error)
	Refresh(refreshReq web.RefreshRequest) (web.RefreshResponse, error)
}

type authServiceImpl struct {
	DB             *gorm.DB
	Validate       *validator.Validate
	UserRepository repository.UserRepository
	RoleRepository repository.RoleRepository
}

func NewAuthService(db *gorm.DB, validate *validator.Validate,
	userRepository repository.UserRepository, roleRepository repository.RoleRepository) *authServiceImpl {
	return &authServiceImpl{
		DB:             db,
		Validate:       validate,
		UserRepository: userRepository,
		RoleRepository: roleRepository,
	}
}

func (s *authServiceImpl) Register(registerReq web.RegisterRequest) error {
	errValidate := s.Validate.Struct(registerReq)
	if errValidate != nil {
		return errValidate
	}

	tx := s.DB.Begin()
	userFound, _ := s.UserRepository.FindByEmail(tx, registerReq.Email)
	roleFound, errRole := s.RoleRepository.FindById(tx, registerReq.RoleID)

	// user can register with the same email but different role
	// if user already registered with the same role
	if userFound.ID != 0 && s.UserRepository.HasRole(tx, userFound, roleFound.Name) {
		return &exception.BadRequestError{Message: "user already exists"}
	}

	// if user register with unknown role
	// or register to the admin role
	if roleFound.ID == 0 || roleFound.Name == "admin" || errRole != nil {
		return &exception.BadRequestError{Message: "role doesn't exists"}
	}

	// if existing user register with new role
	if userFound.ID != 0 {
		errUpdateRole := s.UserRepository.AppendRoles(tx, userFound, []domain.Role{roleFound})
		if errUpdateRole != nil {
			if errRollback := tx.Rollback().Error; errRollback != nil {
				return errRollback
			}
			return errUpdateRole
		}
	} else {
		// or the user is completely new
		_, errCreate := s.UserRepository.Save(tx, domain.User{
			Name:     registerReq.Name,
			Email:    registerReq.Email,
			Password: helper.HashUserPassword(registerReq.Password),
		}, []*domain.Role{&roleFound})
		if errCreate != nil {
			if errRollback := tx.Rollback().Error; errRollback != nil {
				return errRollback
			}
			return errCreate
		}
	}

	if errCommit := tx.Commit().Error; errCommit != nil {
		return errCommit
	}

	return nil
}

func (s *authServiceImpl) Login(loginReq web.LoginRequest) (web.LoginResponse, error) {
	var loginRes web.LoginResponse

	errValidate := s.Validate.Struct(loginReq)
	if errValidate != nil {
		return loginRes, errValidate
	}

	tx := s.DB.Begin()
	userFound, errUser := s.UserRepository.FindByEmail(tx, loginReq.Email)

	if userFound.ID == 0 || errUser != nil {
		return loginRes, &exception.BadRequestError{Message: "user doesn't exists"}
	}

	if !helper.ComparePassword(loginReq.Password, userFound.Password) {
		return loginRes, &exception.BadRequestError{Message: "wrong password"}
	}

	// if role doesn't found
	roleFound, errRole := s.RoleRepository.FindById(tx, loginReq.RoleID)
	if roleFound.ID == 0 || errRole != nil {
		return loginRes, echo.ErrForbidden
	}

	// if user doesn't have the role
	if !s.UserRepository.HasRole(tx, userFound, roleFound.Name) {
		return loginRes, echo.ErrForbidden
	}

	token, errToken := helper.GenerateLoginToken(userFound.Email, roleFound.ID)
	if errToken != nil {
		return loginRes, errToken
	}

	loginRes = web.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return loginRes, nil
}

func (s *authServiceImpl) Refresh(refreshReq web.RefreshRequest) (web.RefreshResponse, error) {
	var refreshRes web.RefreshResponse

	if errVal := s.Validate.Struct(refreshReq); errVal != nil {
		return refreshRes, errVal
	}

	refreshToken, errParse := helper.ParseRefreshToken(refreshReq.RefreshToken)
	if errParse != nil {
		return refreshRes, errParse
	}

	if refreshToken.ExpiresAt < time.Now().Unix() {
		return refreshRes, &exception.BadRequestError{Message: "Expired refresh token"}
	}

	oldAccessTokenData, errParse := helper.ParseAccessToken(refreshReq.AccessToken)
	if errParse != nil {
		return refreshRes, errParse
	}

	newAccessToken, errGenerate := helper.GenerateAccessToken(
		oldAccessTokenData.Email,
		oldAccessTokenData.RoleID,
	)
	if errGenerate != nil {
		return refreshRes, errGenerate
	}

	refreshRes.AccessToken = newAccessToken

	return refreshRes, nil
}
