package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/n9mi/go-docker/controller"
	"github.com/n9mi/go-docker/repository"
	"github.com/n9mi/go-docker/service"
	"gorm.io/gorm"
)

func BlogRouter(e *echo.Echo, mainUrl string, db *gorm.DB, validate *validator.Validate) {
	blogRepository := repository.NewBlogRepository()
	blogService := service.NewBlogService(db, validate, blogRepository)
	blogController := controller.NewBlogController(blogService)

	g := e.Group(mainUrl + "/blogs")
	g.GET("", blogController.GetAll)
	g.GET("/:id", blogController.GetById)
	g.POST("", blogController.Create)
	g.PUT("/:id", blogController.Update)
	g.DELETE("/:id", blogController.Delete)
}
