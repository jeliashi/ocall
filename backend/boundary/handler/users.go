package handler

import (
	"backend/boundary/middleware"
	"backend/boundary/presenter"
	"github.com/gin-gonic/gin"
	"net/http"

	"backend/models"
	"backend/usecase/users"
)

type UserController struct {
	userService users.Service
}

// @Summary Create a new profile
// @Description Create a new profile
// @Tags profiles
// @Accept  json
// @Produce  json
// @Security BearerToken
// @Param profile body models.Profile true "Profile object to be created"
// @Success 201 {object} presenter.IdResponse
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /profiles [post]
func (u *UserController) createProfile(c *gin.Context) {
	var profile models.Profile
	if err := c.Bind(&profile); err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if firebaseId, exists := c.Get(models.FirebaseContextKey); exists {
		profile.UserIDs = []models.UserID{{FirebaseId: firebaseId.(string), Permissions: models.Admin}}
	}

	if id, err := u.userService.CreateProfile(c, profile); err != nil {
		presenter.HandleErr(c, err)
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"id": id.String()})
	}
}

// @Summary Get a profile by ID
// @Description Returns the profile with the specified ID
// @ID get-profile-by-id
// @Tags profiles
// @Produce json
// @Security BearerToken
// @Param id path string true "Profile ID"
// @Success 200 {object} models.Profile
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /profiles/{id} [get]
func (u *UserController) getProfile(c *gin.Context) {
	id, err := GetId(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if profile, exists := c.Get("profile"); exists {
		c.JSON(http.StatusOK, profile)
		return
	}
	if profile, err := u.userService.GetProfileByID(c, id); err != nil {
		presenter.HandleErr(c, err)
		return
	} else {
		c.JSON(http.StatusOK, profile)
	}
}

// Update a profile by ID
// PATCH /profiles/:id
// @Summary Update a profile by ID
// @Description Update a profile by ID
// @Tags Profiles
// @Accept json
// @Produce json
// @Security BearerToken
// @Param id path string true "Profile ID"
// @Param application body models.Profile true "Profile object"
// @Success 200 {object} models.Profile
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 401 {object} presenter.ErrorResponse
// @Failure 403 {object} presenter.ErrorResponse
// @Failure 404 {object} presenter.ErrorResponse
// @Router /events/{id} [patch]
func (u *UserController) updateProfile(c *gin.Context) {
	var profile models.Profile
	if err := c.Bind(&profile); err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if profile, err := u.userService.UpdateProfile(c, profile); err != nil {
		presenter.HandleErr(c, err)
		return
	} else {
		c.JSON(http.StatusOK, profile)
	}
}

// Delete a profile by ID
// DELETE /profiles/:id
// @Summary Delete a profile by ID
// @Description Delete a profile by ID
// @Tags Profiles
// @Param id path string true "Profile ID"
// @Security BearerToken
// @Success 204 "No Content"
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /applications/{id} [delete]
func (u *UserController) deleteProfile(c *gin.Context) {
	id, err := GetId(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if err := u.userService.DeleteProfile(c, id); err != nil {
		presenter.HandleErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func CreateUserController(
	service users.Service,
	router *gin.Engine,
	firebase middleware.FirebaseMiddleware,
	permissionsMiddleware middleware.PermissionsMiddleware,
) UserController {
	handler := UserController{userService: service}
	router.Use(firebase.AuthMiddleware)
	router.POST("/profiles", handler.createProfile)
	router.GET("/profile/:id", handler.getProfile)
	router.PATCH("/profiles/:id", permissionsMiddleware.ProfileModifier, handler.updateProfile)
	router.DELETE("/profiles/:id", permissionsMiddleware.ProfileModifier, handler.deleteProfile)
	return handler
}
