package domain

type User struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string
	Password string
	Roles    []*Role `gorm:"many2many:user_roles;"`
}
