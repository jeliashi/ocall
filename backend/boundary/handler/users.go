package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ocall/backend/boundary/middleware"

	"ocall/backend/boundary/presentor"
	"ocall/backend/models"
	"ocall/backend/usecase/users"
)

type UserController struct {
	UserService users.Service
}

// @description Create a new Profile type associated with firebase user
// @accept json
// @produce json
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Bearer
// @Param data body models.Profile true "input profile model"
// @Success 201 object presenter.IdResponse "Successful Profile Created"
// @Failure 400 object presenter.ErrorResponse "StatusBadRequest"
// @Failure 422 object presenter.ErrorResponse "StatusUnprocessableEntity"
func (u *UserController) createProfile(c *gin.Context) {
	pType, ok := models.ParseProfile(c.Param("pType"))
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	var input models.Profile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, presentor.ErrorResponse{Error: err.Error()})
		return
	}
	input.ProfileType = pType

	firebaseID := models.GetFirebaseIDFromContext(c)
	if firebaseID != "" {
		input.UserIDs = []models.UserID{{FirebaseId: firebaseID, Permissions: models.Admin}}
	}

	id, err := u.UserService.CreateProfile(c, input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, presentor.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, presentor.IdResponse{Id: id.String()})
}

// @description Get a Profile type by id
// @produce json
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Bearer
// @Success 200 object models.Profile
// @Failure 400 object presenter.ErrorResponse "StatusBadRequest"
// @Failure 404 object presenter.ErrorResponse "StatusNotFound"
// @Failure 422 object presenter.ErrorResponse "StatusUnprocessableEntity"
func (u *UserController) getProfile(c *gin.Context) {
	profileID := c.Value("id").(uuid.UUID)

	profile, err := u.UserService.GetProfileByID(c, profileID)
	if err != nil {
		code, response := presentor.HandleErr(err)
		c.JSON(code, response)
		return
	}
	c.JSON(http.StatusOK, profile)
	return
}

// @description Update a Profile type by id
// @accept json
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Bearer
// @Success 202
// @Failure 400 object presenter.ErrorResponse "StatusBadRequest"
// @Failure 404 object presenter.ErrorResponse "StatusNotFound"
// @Failure 422 object presenter.ErrorResponse "StatusUnprocessableEntity"
func (u *UserController) patchProfile(c *gin.Context) {
	var input models.Profile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, presentor.ErrorResponse{Error: err.Error()})
		return
	}
	input.ID = c.Value("id").(uuid.UUID)

	if err := u.UserService.UpdateProfile(c, input); err != nil {
		code, message := presentor.HandleErr(err)
		c.JSON(code, message)
		return
	}
	c.Status(http.StatusAccepted)
	return
}

// @description Delete a Profile type by id
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Bearer
// @Success 200
// @Failure 400 object presenter.ErrorResponse "StatusBadRequest"
// @Failure 404 object presenter.ErrorResponse "StatusNotFound"
// @Failure 422 object presenter.ErrorResponse "StatusUnprocessableEntity"
func (u *UserController) deleteProfile(c *gin.Context) {
	if err := u.UserService.DeleteProfile(c, c.Value("id").(uuid.UUID)); err != nil {
		code, message := presentor.HandleErr(err)
		c.JSON(code, message)
		return
	}
	c.Status(http.StatusOK)
	return
}

func (u *UserController) RegisterProfileEndpoints(
	router *gin.Engine,
	f *middleware.FirebaseMiddleware,
	p *middleware.PermissionsMiddleware,
) {
	router.Use(f.AuthMiddleware)
	router.POST("/:pType", u.createProfile)
	router.GET("/:id/", p.AllowedUser, u.getProfile)
	router.PATCH("/:id/patch", p.AllowedUser, u.patchProfile)
	router.DELETE("/:id/delete", p.AllowedUser, u.deleteProfile)
}
