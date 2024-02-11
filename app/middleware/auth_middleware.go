package middleware

import (
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/n9mi/go-docker/exception"
	"github.com/n9mi/go-docker/helper"
	"github.com/n9mi/go-docker/repository"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB, enforcer *casbin.Enforcer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var casbinSubject string

			authHeader := c.Request().Header.Get("Authorization")
			if !strings.Contains(authHeader, "Bearer") {
				casbinSubject = "guest"
			} else {
				tokenStr := strings.Replace(authHeader, "Bearer ", "", -1)
				parsedAccessToken, errAccess := helper.ParseAccessToken(tokenStr)
				if errAccess != nil {
					return errAccess
				}

				userRepository := repository.NewUserRepository()
				roleRepository := repository.NewRoleRepository()

				tx := db.Begin()

				// check if user exists by email
				userFound, errUser := userRepository.FindByEmail(tx, parsedAccessToken.Email)
				if errUser != nil {
					return &exception.TokenError{}
				}

				// check if role exists by role_id
				roleFound, errRole := roleRepository.FindById(tx, parsedAccessToken.RoleID)
				if errRole != nil {
					return &exception.TokenError{}
				}

				// check if user has the role
				if !userRepository.HasRole(tx, userFound, roleFound.Name) {
					return &exception.TokenError{}
				}

				// assign role name as subject to casbin policy
				casbinSubject = roleFound.Name
			}

			// apply casbin enforcer
			enforceRes, errEnforce := enforcer.Enforce(casbinSubject, c.Request().URL.String(), c.Request().Method)

			if errEnforce != nil {
				return errEnforce
			}

			if enforceRes {
				return next(c)
			} else {
				return echo.ErrForbidden
			}
		}
	}
}
