package exception

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/n9mi/go-docker/model/web"
)

func CustomErrorHandler(err error, c echo.Context) {
	var res web.ErrorReponse

	if castedErr, ok := err.(*NotFoundError); ok {
		res.Code = http.StatusNotFound
		res.Status = "NOT FOUND"
		res.Message = castedErr.Error()
	} else if ok := errors.Is(err, echo.ErrNotFound); ok {
		res.Code = http.StatusNotFound
		res.Status = "NOT FOUND"
		res.Message = err.Error()
	} else if castedErr, ok := err.(validator.ValidationErrors); ok {
		res.Code = http.StatusBadRequest
		res.Status = "BAD REQUEST"
		res.Message = castedErr.Error()
	} else if castedErr, ok := err.(*BadRequestError); ok {
		res.Code = http.StatusBadRequest
		res.Status = "BAD REQUEST"
		res.Message = castedErr.Error()
	} else if castedErr, ok := err.(*NotFoundValidate); ok {
		res.Code = http.StatusBadRequest
		res.Status = "BAD REQUEST"
		res.Message = castedErr.Error()
	} else if castedErr, ok := err.(*TokenError); ok {
		res.Code = http.StatusForbidden
		res.Status = "FORBIDDEN"
		res.Message = castedErr.Error()
	} else if ok := errors.Is(err, echo.ErrForbidden); ok {
		res.Code = http.StatusForbidden
		res.Status = "FORBIDDEN"
	} else {
		res.Code = http.StatusInternalServerError
		res.Status = "FAIL"
		res.Message = err.Error()
	}

	c.Logger().Error(err)
	c.JSON(res.Code, res)
}
