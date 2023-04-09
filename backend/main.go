package main

import (
	"backend/boundary/handler"
	"backend/data/repository"
	"backend/docs"
	"backend/models"
	"backend/usecase/agenda"
	"backend/usecase/users"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"backend/boundary/middleware"
	"unsafe"
)

// @title MyService API
// @description This is a sample API for MyService
// @version 1
// @host localhost:8080
// @BasePath /
func main() {
	viper.AutomaticEnv()
	_ = viper.BindEnv("superUser", "OCALL_SUPERUSER")
	_ = viper.BindEnv("superPw", "OCALL_SUPERPW")
	_ = viper.BindEnv("dbUri", "OCALL_DB_URI")
	user := viper.GetString("superUser")
	pw := viper.GetString("superPw")
	uri := viper.GetString("dbUri")
	userBase := fmt.Sprintf("%s:%s", user, pw)

	firebaseMiddleware := middleware.NewFirebaseMiddleware(base64.StdEncoding.EncodeToString(stringToBytes(userBase)), nil)

	db := postgres.New(postgres.Config{DSN: uri, PreferSimpleProtocol: true})

	orm, err := gorm.Open(db, &gorm.Config{})
	if err != nil {
		fmt.Print(errors.Wrapf(err, "Unable to connect to db %s", uri).Error())
		return
	}
	err = AutoMigrate(orm)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	if err != nil {
		fmt.Printf("unable to open with gorm")
		return
	}
	uRepo := repository.NewUserRepo(orm)
	aRepo := repository.NewAgendaRepo(orm)
	uService := users.NewService(&uRepo)
	aService := agenda.NewService(&aRepo)

	permissionMiddleWare := middleware.NewPermissionsMiddleware(uService, aService)
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	handler.RegisterUserController(uService, v1, firebaseMiddleware, permissionMiddleWare)
	handler.RegisterAgendaHanlder(aService, v1, firebaseMiddleware, permissionMiddleWare)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if _err := router.Run(); _err != nil {
		log.Printf("you failed")
	}

}

func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func AutoMigrate(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	enums := []func(db *gorm.DB) error{
		models.AutoMigratePermission,
		models.AutoMigrateProfileType,
		models.AutoMigrateApplicationStatus,
		models.AutoMigrateEventApplicationStatus,
	}
	for _, f := range enums {
		if err := f(db); err != nil {
			return err
		}
	}
	err := db.AutoMigrate(&models.Profile{}, &models.UserID{}, &models.Tag{}, &models.Event{}, &models.Application{})
	if err != nil {
		return err
	}
	return nil
}
