package repository

import (
	"github.com/n9mi/go-docker/exception"
	"github.com/n9mi/go-docker/helper"
	"github.com/n9mi/go-docker/model/domain"
	"gorm.io/gorm"
)

type BlogRepository interface {
	FindAll(tx *gorm.DB, page int, pageSize int) ([]domain.ScanBlogs, error)
	FindById(tx *gorm.DB, id uint) (domain.ScanBlog, error)
	IsExists(tx *gorm.DB, id uint) bool
	BelongsToUser(tx *gorm.DB, blog domain.Blog) (domain.User, error)
	Save(tx *gorm.DB, blog domain.Blog) (domain.Blog, error)
	Delete(tx *gorm.DB, id uint) error
}

type blogRepositoryImpl struct {
}

func NewBlogRepository() *blogRepositoryImpl {
	return &blogRepositoryImpl{}
}

func (r *blogRepositoryImpl) FindAll(tx *gorm.DB, page int, pageSize int) ([]domain.ScanBlogs, error) {
	var blogs []domain.ScanBlogs
	if page > 0 && pageSize > 0 {
		tx = tx.Scopes(helper.Paginate(page, pageSize))
	}
	if err := tx.Model(&domain.Blog{}).
		Order("blogs.id asc").
		Select("blogs.id, blogs.title, blogs.summary, users.name as created_by").
		Joins("inner join users on users.id = blogs.user_id").
		Scan(&blogs).Error; err != nil {
		return blogs, err
	}

	return blogs, nil
}

func (r *blogRepositoryImpl) FindById(tx *gorm.DB, id uint) (domain.ScanBlog, error) {
	var blog domain.ScanBlog
	if err := tx.Model(&domain.Blog{}).
		Select("blogs.id, blogs.title, blogs.summary, blogs.content, users.name as created_by").
		Joins("inner join users on users.id = blogs.user_id").
		Where("blogs.id = ?", id).
		Scan(&blog).Error; err != nil {
		return blog, err
	}

	if blog.ID == 0 {
		return blog, &exception.NotFoundError{Entity: "blog"}
	}

	return blog, nil
}

func (r *blogRepositoryImpl) IsExists(tx *gorm.DB, id uint) bool {
	var count int64
	if tx.Model(&domain.Blog{}).Where("id = ?", id).Count(&count); count <= 0 {
		return false
	}

	return true
}

func (r *blogRepositoryImpl) BelongsToUser(tx *gorm.DB, blog domain.Blog) (domain.User, error) {
	var userBelongs domain.User
	if err := tx.Model(&blog).Association("User").Find(&userBelongs); err != nil {
		return userBelongs, err
	}

	return userBelongs, nil
}

func (r *blogRepositoryImpl) Save(tx *gorm.DB, blog domain.Blog) (domain.Blog, error) {
	if err := tx.Save(&blog).Error; err != nil {
		return blog, err
	}

	return blog, nil
}

func (r *blogRepositoryImpl) Delete(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&domain.Blog{}, id).Error; err != nil {
		return err
	}

	return nil
}
