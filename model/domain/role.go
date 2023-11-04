package domain

type Role struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Users []*User `gorm:"many2many:user_roles;"`
}
