package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naomigrain/go-docker/controller"
	"github.com/naomigrain/go-docker/repository"
	"github.com/naomigrain/go-docker/service"
	"gorm.io/gorm"
)

func AuthRouter(e *echo.Echo, mainUrl string, db *gorm.DB, validate *validator.Validate) {
	userRepository := repository.NewUserRepository()
	roleRepository := repository.NewRoleRepository()
	authService := service.NewAuthService(db, validate, userRepository, roleRepository)
	userController := controller.NewAuthController(authService)

	g := e.Group(mainUrl + "/auth")
	g.POST("/login", userController.Login)
	g.POST("/register", userController.Register)
	g.POST("/refresh", userController.Refresh)
}
