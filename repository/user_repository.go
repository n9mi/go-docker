package repository

import (
	"github.com/naomigrain/go-docker/exception"
	"github.com/naomigrain/go-docker/model/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindById(tx *gorm.DB, id uint) (domain.User, error)
	FindByEmail(tx *gorm.DB, email string) (domain.User, error)
	IsUserExistsByEmail(tx *gorm.DB, email string) bool
	Save(tx *gorm.DB, user domain.User, role []*domain.Role) (domain.User, error)
	HasRole(tx *gorm.DB, user domain.User, roleName string) bool
	Roles(tx *gorm.DB, user domain.User) ([]domain.Role, error)
	AppendRoles(tx *gorm.DB, user domain.User, roles []domain.Role) error
}

type userRepositoryImpl struct {
}

func NewUserRepository() *userRepositoryImpl {
	return &userRepositoryImpl{}
}

func (r *userRepositoryImpl) FindById(tx *gorm.DB, id uint) (domain.User, error) {
	var user domain.User
	if err := tx.Find(&user, id).Error; err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, &exception.NotFoundError{Entity: "user"}
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByEmail(tx *gorm.DB, email string) (domain.User, error) {
	var user domain.User
	if err := tx.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, &exception.NotFoundError{Entity: "user"}
	}

	return user, nil
}

func (r *userRepositoryImpl) IsUserExistsByEmail(tx *gorm.DB, email string) bool {
	var count int64
	if tx.Model(&domain.User{}).Where("email = ?", email).Count(&count); count <= 0 {
		return false
	}

	return true
}

func (r *userRepositoryImpl) Save(tx *gorm.DB, user domain.User, roles []*domain.Role) (domain.User, error) {
	if len(roles) > 0 {
		user.Roles = roles
	}

	if err := tx.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepositoryImpl) HasRole(tx *gorm.DB, user domain.User, roleName string) bool {
	var role domain.Role
	if tx.Model(&user).Where("name = ?", roleName).Association("Roles").
		Find(&role); role.ID != 0 {
		return true
	}

	return false
}

func (r *userRepositoryImpl) Roles(tx *gorm.DB, user domain.User) ([]domain.Role, error) {
	var roles []domain.Role
	if err := tx.Model(&user).Association("Roles").Find(&roles); err != nil {
		return roles, err
	}

	return roles, nil
}

func (r *userRepositoryImpl) AppendRoles(tx *gorm.DB, user domain.User, roles []domain.Role) error {
	if err := tx.Model(&user).Association("Roles").Append(roles); err != nil {
		return err
	}

	return nil
}
