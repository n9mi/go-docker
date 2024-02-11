package database

import (
	"fmt"
	"strconv"

	"github.com/n9mi/go-docker/config"
	"github.com/n9mi/go-docker/helper"
	"github.com/n9mi/go-docker/model/domain"
	"github.com/n9mi/go-docker/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newDatabaseConnString(dbConfig config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbConfig.DBHost,
		dbConfig.DBUser,
		dbConfig.DBPassword,
		dbConfig.DBName,
		dbConfig.DBPort,
	)
}

func NewDB(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	dsn := newDatabaseConnString(dbConfig)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&domain.Role{})
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Blog{})
}

func Seed(db *gorm.DB) {
	roles := SeedRole(db)
	users := SeedUser(db, roles)
	SeedBlog(db, users)
}

func SeedRole(db *gorm.DB) []domain.Role {
	roleNames := []string{
		"admin",
		"creator",
		"subscriber",
	}
	var roles []domain.Role

	roleRepository := repository.NewRoleRepository()
	for _, rN := range roleNames {
		tx := db.Begin()
		result, _ := roleRepository.Save(tx, domain.Role{
			Name: rN,
		})
		tx.Commit()
		roles = append(roles, result)
	}

	return roles
}

func SeedUser(db *gorm.DB, roles []domain.Role) []domain.User {
	var users []domain.User

	userRepository := repository.NewUserRepository()
	for _, r := range roles {
		for i := 1; i <= 2; i++ {
			tx := db.Begin()
			userName := r.Name + " " + strconv.Itoa(i)
			result, _ := userRepository.Save(tx, domain.User{
				Name:     userName,
				Email:    userName + "@gmail.com",
				Password: helper.HashUserPassword(r.Name + strconv.Itoa(i)),
			}, []*domain.Role{&r})
			tx.Commit()

			users = append(users, result)
		}
	}

	return users
}

func SeedBlog(db *gorm.DB, users []domain.User) []domain.Blog {
	userRepository := repository.NewUserRepository()
	blogRepoistory := repository.NewBlogRepository()
	var blogs []domain.Blog

	for _, u := range users {
		tx := db.Begin()
		if userRepository.HasRole(tx, u, "creator") {
			result, _ := blogRepoistory.Save(tx, domain.Blog{
				Title:   "Title " + helper.GenerateRandomString(10),
				Summary: "Summary " + helper.GenerateRandomString(50),
				Content: "Content " + helper.GenerateRandomString(100),
				UserID:  u.ID,
			})
			tx.Commit()
			blogs = append(blogs, result)
		}
	}

	return blogs
}

func Drop(db *gorm.DB) {
	db.Migrator().DropTable(&domain.Blog{})
	db.Migrator().DropTable("user_roles")
	db.Migrator().DropTable(&domain.User{})
	db.Migrator().DropTable(&domain.Role{})
}
