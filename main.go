package main

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/n9mi/go-docker/app/database"
	"github.com/n9mi/go-docker/app/middleware"
	"github.com/n9mi/go-docker/app/router"
	"github.com/n9mi/go-docker/config"
)

func main() {
	isUsingDotEnv := true
	if isUsingDotEnv {
		godotenv.Load()
	}

	e := router.InitializeEcho()

	enforcer, errEnforcer := casbin.NewEnforcer("./casbin/model.conf", "./casbin/policy.csv")
	if errEnforcer != nil {
		e.Logger.Fatal(errEnforcer.Error())
	}

	dbConfig := config.NewDatabaseConfig(isUsingDotEnv)
	db, errDBConn := database.NewDB(dbConfig)
	defer database.Drop(db)

	if errDBConn != nil {
		e.Logger.Fatal(errDBConn.Error())
	}

	database.Drop(db)
	database.Migrate(db)
	database.Seed(db)

	e.Use(middleware.AuthMiddleware(db, enforcer))
	router.AssignRouter(e, db, validator.New())

	appConfig := config.NewAppConfig(isUsingDotEnv)
	e.Logger.Fatal(e.Start(":" + appConfig.AppPort))
}
