package repository

import (
	"github.com/naomigrain/go-docker/exception"
	"github.com/naomigrain/go-docker/model/domain"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Save(tx *gorm.DB, role domain.Role) (domain.Role, error)
	FindById(tx *gorm.DB, id uint) (domain.Role, error)
	FindByName(tx *gorm.DB, name string) (domain.Role, error)
}

type roleRepositoryImpl struct {
}

func NewRoleRepository() *roleRepositoryImpl {
	return &roleRepositoryImpl{}
}

func (r *roleRepositoryImpl) Save(tx *gorm.DB, role domain.Role) (domain.Role, error) {
	if err := tx.Save(&role).Error; err != nil {
		return role, err
	}

	return role, nil
}

func (r *roleRepositoryImpl) FindById(tx *gorm.DB, id uint) (domain.Role, error) {
	var role domain.Role
	if err := tx.Find(&role, id).Error; err != nil {
		return role, &exception.NotFoundError{Entity: "role"}
	}

	return role, nil
}

func (r *roleRepositoryImpl) FindByName(tx *gorm.DB, name string) (domain.Role, error) {
	var role domain.Role
	if err := tx.Where("name = ?", name).First(&role).Error; err != nil {
		return role, &exception.NotFoundError{Entity: "role"}
	}

	return role, nil
}
