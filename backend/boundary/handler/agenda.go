package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ocall/backend/boundary/middleware"
	"ocall/backend/boundary/presentor"
	"ocall/backend/models"
	"ocall/backend/usecase/agenda"
)

type AgendaController struct {
	Service agenda.Service
}

// createEvent creates an event.
// @Summary Create an Event
// @Description Creates an event based on the request body.
// @Accept json
// @Produce json
// @Param application body models.Event true "Event object"
// @Success 201 {object} presentor.IdResponse
// @Failure 400 {object} presentor.ErrorResponse
// @Failure 500 {object} presentor.ErrorResponse
// @Router /event [post]
func (a *AgendaController) createEvent(c *gin.Context) {
	var input models.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		_, response := presentor.HandleErr(err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	id, err := a.Service.CreateEvent(c, input)
	if err != nil {
		code, response := presentor.HandleErr(err)
		c.JSON(code, response)
		return
	}
	c.JSON(http.StatusCreated, presentor.IdResponse{Id: id.String()})
}

// getEvent retrieve an event.
// @Summary Retrieves an Event
// @Description Creates an event based on the request body.
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} models.Event
// @Failure 400 {object} presentor.ErrorResponse
// @Failure 500 {object} presentor.ErrorResponse
// @Router /event/{id} [get]
func (a *AgendaController) getEvent(c *gin.Context) {
	id := c.Value("id").(uuid.UUID)

	event, err := a.Service.GetEvent(c, id)
	if err != nil {
		code, response := presentor.HandleErr(err)
		c.JSON(code, response)
		return
	}
	c.JSON(http.StatusOK, event)
}

// updateEvent retrieve an event.
// @Summary Retrieves an Event
// @Description Creates an event based on the request body.
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param event body models.Event true "Event object"
// @Success 200 {object} models.Event
// @Failure 400 {object} presentor.ErrorResponse
// @Failure 500 {object} presentor.ErrorResponse
// @Router /event/{id} [patch]
func (a *AgendaController) updateEvent(c *gin.Context) {
	var input models.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, presentor.ErrorResponse{Error: err.Error()})
		return
	}
	input.ID = c.Value("id").(uuid.UUID)
	if err := a.Service.UpdateEvent(c, input); err != nil {
		code, message := presentor.HandleErr(err)
		c.JSON(code, message)
		return
	}
	c.Status(http.StatusAccepted)
}

// deleteEvent deletes an event.
// @Summary Deletes an Event
// @Description Deletes an event based on the request body.
// @Param id path string true "Event ID"
// @Success 200 {object} models.Event
// @Failure 400 {object} presentor.ErrorResponse
// @Failure 500 {object} presentor.ErrorResponse
// @Router /event/{id} [delete]
func (a *AgendaController) deleteEvent(c *gin.Context) {
	if err := a.Service.DeleteEvent(c, c.Value("id").(uuid.UUID)); err != nil {
		code, message := presentor.HandleErr(err)
		c.JSON(code, message)
		return
	}
	c.Status(http.StatusOK)
}

// createApplication creates an application.
// @Summary Create an Application
// @Description Creates an application based on the request body.
// @Accept json
// @Produce json
// @Param application body models.Application true "Application object"
// @Success 201 {object} presentor.IdResponse
// @Failure 400 {object} presentor.ErrorResponse
// @Failure 500 {object} presentor.ErrorResponse
// @Router /event [post]
func (a *AgendaController) createApplication(c *gin.Context) {
	var input models.Application
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, presentor.ErrorResponse{Error: err.Error()})
		return
	}
	id, err := a.Service.CreateApplication(c, input)
	if err != nil {
		code, message := presentor.HandleErr(err)
		c.JSON(code, message)
		return
	}
	c.JSON(http.StatusCreated, presentor.IdResponse{Id: id.String()})
}

// getApplication retrieve an application.
// @Summary Retrieves an Application
// @Description Retrieves an application based on id in route
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} models.Application
// @Failure 400 {object} presentor.ErrorResponse
// @Failure 500 {object} presentor.ErrorResponse
// @Router /event/{id} [get]
func (a *AgendaController) getApplication(c *gin.Context) {
	app, err := a.Service.GetApplication(c, c.Value("id").(uuid.UUID))
	if err != nil {
		code, message := presentor.HandleErr(err)
		c.JSON(code, message)
		return
	}
	c.JSON(http.StatusOK, app)
}

// updateApplication updates an application.
// @Summary Updates an Application
// @Description Updates an application based on the request body.
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param application body models.Application true "Application object"
// @Success 200 {object} models.Application
// @Failure 400 {object} presentor.ErrorResponse
// @Failure 500 {object} presentor.ErrorResponse
// @Router /application/{id} [patch]
func (a *AgendaController) updateApplication(c *gin.Context) {
	var input models.Application
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, presentor.ErrorResponse{Error: err.Error()})
		return
	}
	input.ID = c.Value("id").(uuid.UUID)
	if err := a.Service.UpdateApplication(c, input); err != nil {
		code, message := presentor.HandleErr(err)
		c.JSON(code, message)
		return
	}
	c.Status(http.StatusAccepted)
}

// deleteApplication deletes an application.
// @Summary Deletes an Application
// @Description Deletes an application based on the request body.
// @Param id path string true "Application ID"
// @Success 200 {object} models.Application
// @Failure 400 {object} presentor.ErrorResponse
// @Failure 500 {object} presentor.ErrorResponse
// @Router /application/{id} [delete]
func (a *AgendaController) deleteApplication(c *gin.Context) {
	if err := a.Service.DeleteApplication(c, c.Value("id").(uuid.UUID)); err != nil {
		code, message := presentor.HandleErr(err)
		c.JSON(code, message)
		return
	}
	c.Status(http.StatusOK)
}

func (a *AgendaController) RegisterAgendaEndpoints(
	router *gin.Engine,
	f *middleware.FirebaseMiddleware,
	p *middleware.PermissionsMiddleware,
) {
	router.Use(f.AuthMiddleware)
	router.POST("/event", a.createEvent)
	router.GET("/event/:id", a.getEvent)
	router.PATCH("/event/:id", p.EventModifier, a.updateEvent)
	router.DELETE("/event/:id", p.EventModifier, a.deleteEvent)
	router.POST("/application", a.createApplication)
	router.GET("/application/:id", p.ApplicationViewer, a.getApplication)
	router.PATCH("/application/:id", p.ApplicationViewer, a.updateApplication)
	router.DELETE("/application/:id", p.ApplicationViewer, a.deleteApplication)
}
