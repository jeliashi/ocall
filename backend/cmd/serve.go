package cmd

import (
	"backend/boundary/handler"
	"backend/data/repository"
	"backend/usecase/agenda"
	"backend/usecase/users"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend/boundary/middleware"
	"unsafe"
)

func main() {
	viper.AutomaticEnv()
	viper.BindEnv("superUser", "OCALL_SUPERUSER")
	viper.BindEnv("superPw", "OCALL_SUPERPW")
	viper.BindEnv("dbUri", "OCALL_DB_URI")
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
		fmt.Printf("Unable to connect to db")
		return
	}
	orm, err := gorm.Open(
		postgres.New(postgres.Config{Conn: db}),
	)
	if err != nil {
		fmt.Printf("unable to open with gorm")
		return
	}
	uRepo := repository.NewUserRepo(orm)
	aRepo := repository.NewAgendaRepo(orm)
	uService := users.NewService(&uRepo)
	aService := agenda.NewService(&aRepo)
	permissionMiddleWare := middleware.PermissionsMiddleware{
		uService, aService,
	}
	uHandler := handler.UserController{uService}
	aHandler := handler.AgendaController{aService}

	router := gin.Default()
	uHandler.RegisterProfileEndpoints(router, &firebaseMiddleware, &permissionMiddleWare)
	aHandler.RegisterAgendaEndpoints(router, &firebaseMiddleware, &permissionMiddleWare)

	router.Run(":8080")

}

func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
