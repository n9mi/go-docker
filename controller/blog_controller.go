package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/n9mi/go-docker/exception"
	"github.com/n9mi/go-docker/model/web"
	"github.com/n9mi/go-docker/service"
)

type BlogController interface {
	GetAll(c echo.Context) error
	GetById(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type blogControllerImpl struct {
	BlogService service.BlogService
}

func NewBlogController(blogService service.BlogService) *blogControllerImpl {
	return &blogControllerImpl{
		BlogService: blogService,
	}
}

func (ct *blogControllerImpl) GetAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))

	blogs, errFind := ct.BlogService.FindAll(page, pageSize)
	if errFind != nil {
		return errFind
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   blogs,
	}
	return c.JSON(res.Code, res)
}

func (ct *blogControllerImpl) GetById(c echo.Context) error {
	id := c.Param("id")
	idInt, errConv := strconv.Atoi(id)
	if errConv != nil {
		return &exception.NotFoundError{Entity: "blog"}
	}

	blog, errFind := ct.BlogService.FindById(uint(idInt))
	if errFind != nil {
		return errFind
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   blog,
	}
	return c.JSON(res.Code, res)
}

func (ct *blogControllerImpl) Create(c echo.Context) error {
	blogReq := new(web.BlogRequest)

	if errBind := c.Bind(blogReq); errBind != nil {
		return &exception.BadRequestError{Message: errBind.Error()}
	}

	blogRes, errCreate := ct.BlogService.Create(*blogReq)
	if errCreate != nil {
		return errCreate
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   blogRes,
	}
	return c.JSON(res.Code, res)
}

func (ct *blogControllerImpl) Update(c echo.Context) error {
	id := c.Param("id")
	idInt, errConv := strconv.Atoi(id)
	if errConv != nil {
		return &exception.NotFoundError{Entity: "blog"}
	}

	blogReq := new(web.BlogRequest)
	if errBind := c.Bind(blogReq); errBind != nil {
		return &exception.BadRequestError{Message: errBind.Error()}
	}

	blogReq.ID = uint(idInt)
	blogRes, errUpdate := ct.BlogService.Update(*blogReq)
	if errUpdate != nil {
		return errUpdate
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   blogRes,
	}
	return c.JSON(res.Code, res)
}

func (ct *blogControllerImpl) Delete(c echo.Context) error {
	id := c.Param("id")
	idInt, errConv := strconv.Atoi(id)
	if errConv != nil {
		return &exception.NotFoundError{Entity: "blog"}
	}

	errDel := ct.BlogService.Delete(uint(idInt))
	if errDel != nil {
		return errDel
	}

	res := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   nil,
	}
	return c.JSON(res.Code, res)
}
