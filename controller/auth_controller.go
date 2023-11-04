package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naomigrain/go-docker/exception"
	"github.com/naomigrain/go-docker/model/web"
	"github.com/naomigrain/go-docker/service"
)

type AuthController interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	Refresh(c echo.Context) error
}

type authControllerImpl struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) *authControllerImpl {
	return &authControllerImpl{
		AuthService: authService,
	}
}

func (ct *authControllerImpl) Register(c echo.Context) error {
	registerReq := new(web.RegisterRequest)

	if errBind := c.Bind(registerReq); errBind != nil {
		return &exception.BadRequestError{Message: errBind.Error()}
	}

	errRegister := ct.AuthService.Register(*registerReq)
	if errRegister != nil {
		return errRegister
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
	}
	return c.JSON(res.Code, res)
}

func (ct *authControllerImpl) Login(c echo.Context) error {
	loginReq := new(web.LoginRequest)

	if errBind := c.Bind(loginReq); errBind != nil {
		return &exception.BadRequestError{Message: errBind.Error()}
	}

	loginRes, errLogin := ct.AuthService.Login(*loginReq)
	if errLogin != nil {
		return errLogin
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   loginRes,
	}
	return c.JSON(res.Code, res)
}

func (ct *authControllerImpl) Refresh(c echo.Context) error {
	refreshReq := new(web.RefreshRequest)

	if errBind := c.Bind(refreshReq); errBind != nil {
		return &exception.BadRequestError{Message: errBind.Error()}
	}

	refreshRes, errRefresh := ct.AuthService.Refresh(*refreshReq)
	if errRefresh != nil {
		return errRefresh
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   refreshRes,
	}
	return c.JSON(res.Code, res)
}
