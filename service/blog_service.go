package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/naomigrain/go-docker/exception"
	"github.com/naomigrain/go-docker/model/domain"
	"github.com/naomigrain/go-docker/model/web"
	"github.com/naomigrain/go-docker/repository"
	"gorm.io/gorm"
)

type BlogService interface {
	FindAll(page int, pageSize int) ([]web.BlogListResponse, error)
	FindById(id uint) (web.BlogResponse, error)
	Create(blogReq web.BlogRequest) (web.BlogResponse, error)
	Update(blogReq web.BlogRequest) (web.BlogResponse, error)
	Delete(id uint) error
}

type blogServiceImpl struct {
	DB             *gorm.DB
	Validate       *validator.Validate
	BlogRepository repository.BlogRepository
}

func NewBlogService(db *gorm.DB, validate *validator.Validate,
	blogRepository repository.BlogRepository) *blogServiceImpl {
	return &blogServiceImpl{
		DB:             db,
		Validate:       validate,
		BlogRepository: blogRepository,
	}
}

func (s *blogServiceImpl) FindAll(page int, pageSize int) ([]web.BlogListResponse, error) {
	var blogsListRes []web.BlogListResponse

	tx := s.DB.Begin()
	blogsDom, errFind := s.BlogRepository.FindAll(tx, page, pageSize)
	if errFind != nil {
		return blogsListRes, errFind
	}

	for _, b := range blogsDom {
		blogsListRes = append(blogsListRes, web.BlogListResponse{
			ID:        b.ID,
			Title:     b.Title,
			Summary:   b.Summary,
			CreatedBy: b.CreatedBy,
		})
	}

	return blogsListRes, nil
}

func (s *blogServiceImpl) FindById(id uint) (web.BlogResponse, error) {
	var blogRes web.BlogResponse

	tx := s.DB.Begin()
	blogDom, errFind := s.BlogRepository.FindById(tx, id)
	if errFind != nil {
		return blogRes, errFind
	}

	blogRes = web.BlogResponse{
		ID:        blogDom.ID,
		Title:     blogDom.Title,
		Summary:   blogDom.Summary,
		Content:   blogDom.Content,
		CreatedBy: blogDom.CreatedBy,
	}

	return blogRes, nil
}

func (s *blogServiceImpl) Create(blogReq web.BlogRequest) (web.BlogResponse, error) {
	var blogRes web.BlogResponse

	if errValidate := s.Validate.Struct(blogReq); errValidate != nil {
		return blogRes, errValidate
	}

	tx := s.DB.Begin()

	userRepository := repository.NewUserRepository()
	userBelongs, errFind := userRepository.FindById(tx, blogReq.UserID)
	if errFind != nil {
		return blogRes, &exception.NotFoundValidate{Entity: "user"}
	}

	blogDom, errSave := s.BlogRepository.Save(tx, domain.Blog{
		Title:   blogReq.Title,
		Summary: blogReq.Summary,
		Content: blogReq.Content,
		UserID:  userBelongs.ID,
	})
	if errSave != nil {
		if errRollback := tx.Rollback().Error; errRollback != nil {
			return blogRes, errRollback
		}
		return blogRes, errSave
	}
	if errCommit := tx.Commit().Error; errCommit != nil {
		return blogRes, errCommit
	}

	blogRes = web.BlogResponse{
		ID:        blogDom.ID,
		Title:     blogDom.Title,
		Summary:   blogDom.Summary,
		Content:   blogDom.Content,
		CreatedBy: userBelongs.Name,
	}
	return blogRes, nil
}

func (s *blogServiceImpl) Update(blogReq web.BlogRequest) (web.BlogResponse, error) {
	var blogRes web.BlogResponse

	if errValidate := s.Validate.Struct(blogReq); errValidate != nil {
		return blogRes, errValidate
	}

	tx := s.DB.Begin()
	if !s.BlogRepository.IsExists(tx, blogReq.ID) {
		return blogRes, &exception.NotFoundError{Entity: "blog"}
	}

	userRepository := repository.NewUserRepository()
	userBelongs, errFind := userRepository.FindById(tx, blogReq.UserID)
	if errFind != nil {
		return blogRes, &exception.NotFoundValidate{Entity: "user"}
	}

	blogDom, errSave := s.BlogRepository.Save(tx, domain.Blog{
		ID:      blogReq.ID,
		Title:   blogReq.Title,
		Summary: blogReq.Summary,
		Content: blogReq.Content,
		UserID:  userBelongs.ID,
	})
	if errSave != nil {
		if errRollback := tx.Rollback().Error; errRollback != nil {
			return blogRes, errRollback
		}
		return blogRes, errSave
	}
	if errCommit := tx.Commit().Error; errCommit != nil {
		return blogRes, errCommit
	}

	blogRes = web.BlogResponse{
		ID:        blogDom.ID,
		Title:     blogDom.Title,
		Summary:   blogDom.Summary,
		Content:   blogDom.Content,
		CreatedBy: userBelongs.Name,
	}
	return blogRes, nil
}

func (s *blogServiceImpl) Delete(id uint) error {
	tx := s.DB.Begin()
	if !s.BlogRepository.IsExists(tx, id) {
		return &exception.NotFoundError{Entity: "blog"}
	}

	errDel := s.BlogRepository.Delete(tx, id)
	if errDel != nil {
		if errRollback := tx.Rollback().Error; errRollback != nil {
			return errRollback
		}

		return errDel
	}
	if errCommit := tx.Commit().Error; errCommit != nil {
		return errCommit
	}

	return nil
}
