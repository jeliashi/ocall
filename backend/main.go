package main

import (
	"backend/boundary/handler"
	"backend/data/repository"
	"backend/models"
	"backend/usecase/agenda"
	"backend/usecase/users"
	_ "github.com/lib/pq"

	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"backend/boundary/middleware"
	"unsafe"
)

func main() {
	viper.AutomaticEnv()
	_ = viper.BindEnv("superUser", "OCALL_SUPERUSER")
	_ = viper.BindEnv("superPw", "OCALL_SUPERPW")
	_ = viper.BindEnv("dbUri", "OCALL_DB_URI")
	user := viper.GetString("superUser")
	pw := viper.GetString("superPw")
	uri := viper.GetString("dbUri")
	userBase := fmt.Sprintf("%s:%s", user, pw)

	firebaseMiddleware := middleware.FirebaseMiddleware{
		SuperUserEncoded: base64.StdEncoding.EncodeToString(stringToBytes(userBase)),
		Client:           nil,
	}

	db, err := sql.Open("postgres", uri)
	if err != nil {
		fmt.Print(errors.Wrapf(err, "Unable to connect to db %s", uri).Error())
		return
	}
	orm, err := gorm.Open(
		postgres.New(postgres.Config{Conn: db}),
	)
	err = orm.AutoMigrate(&models.Profile{}, &models.Tag{}, &models.Event{}, &models.Application{})
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
	handler.CreateUserController(uService, router, firebaseMiddleware, permissionMiddleWare)
	handler.CreateAgendaHanlder(aService, router, firebaseMiddleware, permissionMiddleWare)

	if _err := router.Run(":8080"); _err != nil {
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
