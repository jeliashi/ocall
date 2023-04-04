package handler

import (
	"backend/boundary/middleware"
	"backend/boundary/presenter"
	"backend/models"
	"backend/usecase/agenda"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nferruzzi/gormGIS"
	"net/http"
	"strconv"
	"time"
)

type AgendaController struct {
	agendaService agenda.Service
}

// @Summary Create a new event
// @Description Create a new event
// @Tags Events
// @Accept  json
// @Produce  json
// @Security BearerToken
// @Param event body models.Event true "Event object to be created"
// @Success 201 {object} presenter.IdResponse
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 401 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /events [post]
func (a *AgendaController) createEvent(c *gin.Context) {
	var event models.Event
	if err := c.Bind(&event); err != nil {
		presenter.HandleErr(c, err)
		return
	}

	if id, err := a.agendaService.CreateEvent(c, event); err != nil {
		presenter.HandleErr(c, err)
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"id": id.String()})
	}
}

func parseTime(c *gin.Context, key string, result *time.Time) error {
	if out, err := time.Parse(time.RFC3339, c.Query(key)); err != nil {
		return gin.Error{Err: err, Type: gin.ErrorTypeBind}
	} else {
		result = &out
		return nil
	}
}
func parseFloat(c *gin.Context, key string, result *float64) error {
	if out, err := strconv.ParseFloat(c.Query(key), 64); err != nil {
		return gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		}
	} else {
		result = &out
		return nil
	}
}

// GetAllEvents returns all events within a specified time range and distance from a center point.
// @Summary Get all events
// @Description Returns all events within a specified time range and distance from a center point.
// @Tags events
// @Accept  json
// @Produce  json
// @Param start_time query string true "Start time of the range (RFC3339)"
// @Param end_time query string true "End time of the range (RFC3339)"
// @Param lat query number true "latitude of search point"
// @Param lon query number true "longitude of search point"
// @Param distance_km query number true "Distance from the center point in kilometers"
// @Success 200 {array} models.Event
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /events [get]
func (a *AgendaController) getEventsByFilter(c *gin.Context) {
	var startTime, endTime time.Time
	if err := parseTime(c, "start_time", &startTime); err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if err := parseTime(c, "end_time", &endTime); err != nil {
		presenter.HandleErr(c, err)
		return
	}

	var lat, lon, distance float64
	if err := parseFloat(c, "lat", &lat); err != nil {
		return
	}
	if err := parseFloat(c, "lon", &lon); err != nil {
		return
	}
	centerPoint := gormGIS.GeoPoint{
		Lng: lon,
		Lat: lat,
	}
	if err := parseFloat(c, "distance_km", &distance); err != nil {
		return
	}
	events, err := a.agendaService.GetAllEvents(c, startTime, endTime, centerPoint, distance)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, events)
}

// @Summary Get an event by ID
// @Description Returns the event with the specified ID
// @ID get-event-by-id
// @Tags Events
// @Produce json
// @Security BearerToken
// @Param id path string true "Event ID"
// @Success 200 {object} models.Event
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 404 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /events/{id} [get]
func (a *AgendaController) getEvent(c *gin.Context) {
	id, err := GetId(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if profile, exists := c.Get("profile"); exists {
		c.JSON(http.StatusOK, profile)
		return
	}
	if profile, err := a.agendaService.GetEvent(c, id); err != nil {
		presenter.HandleErr(c, err)
		return
	} else {
		c.JSON(http.StatusOK, profile)
	}
}

// Update an event by ID
// PATCH /events/:id
// @Summary Update an event by ID
// @Description Update an event by ID
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerToken
// @Param id path string true "Event ID"
// @Param event body models.Event true "Event object"
// @Success 200 {object} models.Event
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 401 {object} presenter.ErrorResponse
// @Failure 403 {object} presenter.ErrorResponse
// @Failure 404 {object} presenter.ErrorResponse
// @Router /events/{id} [patch]
func (a *AgendaController) updateEvent(c *gin.Context) {
	var event models.Event
	if err := c.Bind(&event); err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if profile, err := a.agendaService.UpdateEvent(c, event); err != nil {
		presenter.HandleErr(c, err)
		return
	} else {
		c.JSON(http.StatusOK, profile)
	}
}

// Delete an event by ID
// DELETE /events/:id
// @Summary Delete an event by ID
// @Description Delete an event by ID
// @Tags Events
// @Param id path string true "Event ID"
// @Security BearerToken
// @Success 204 "No Content"
// @Failure 401 {object} presenter.ErrorResponse
// @Failure 403 {object} presenter.ErrorResponse
// @Failure 404 {object} presenter.ErrorResponse
// @Router /events/{id} [delete]
func (a *AgendaController) deleteEvent(c *gin.Context) {
	id, err := GetId(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if err := a.agendaService.DeleteEvent(c, id); err != nil {
		presenter.HandleErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Create a new application
// @Description Create a new application
// @Tags Applications
// @Accept  json
// @Produce  json
// @Security BearerToken
// @Param application body models.Application true "Application object to be created"
// @Success 201 {object} presenter.IdResponse
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 401 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /applications [post]
func (a *AgendaController) createApplication(c *gin.Context) {
	var application models.Application
	if err := c.Bind(&application); err != nil {
		presenter.HandleErr(c, err)
		return
	}

	if id, err := a.agendaService.CreateApplication(c, application); err != nil {
		presenter.HandleErr(c, err)
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"id": id.String()})
	}
}

// @Summary Get an Application by ID
// @Description Returns the application with the specified ID
// @ID get-application-by-id
// @Tags Applications
// @Produce json
// @Security BearerToken
// @Param id path string true "Application ID"
// @Success 200 {object} models.Application
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 404 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /applications/{id} [get]
func (a *AgendaController) getApplication(c *gin.Context) {
	id, err := GetId(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if application, err := a.agendaService.GetApplication(c, id); err != nil {
		presenter.HandleErr(c, err)
		return
	} else {
		c.JSON(http.StatusOK, application)
	}
}

// Update an application by ID
// PATCH /applications/:id
// @Summary Update an event by ID
// @Description Update an event by ID
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerToken
// @Param id path string true "Application ID"
// @Param application body models.Application true "Application object"
// @Success 200 {object} models.Application
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 401 {object} presenter.ErrorResponse
// @Failure 403 {object} presenter.ErrorResponse
// @Failure 404 {object} presenter.ErrorResponse
// @Router /applications/{id} [patch]
func (a *AgendaController) updateApplication(c *gin.Context) {
	var application models.Application
	if err := c.Bind(&application); err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if profile, err := a.agendaService.UpdateApplication(c, application); err != nil {
		presenter.HandleErr(c, err)
		return
	} else {
		c.JSON(http.StatusOK, profile)
	}
}

// Delete an application by ID
// DELETE /applications/:id
// @Summary Delete an application by ID
// @Description Delete an application by ID
// @Tags Applications
// @Param id path string true "Application ID"
// @Security BearerToken
// @Success 204 "No Content"
// @Failure 401 {object} presenter.ErrorResponse
// @Failure 403 {object} presenter.ErrorResponse
// @Failure 404 {object} presenter.ErrorResponse
// @Router /applications/{id} [delete]
func (a *AgendaController) deleteApplication(c *gin.Context) {
	id, err := GetId(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	if err := a.agendaService.DeleteApplication(c, id); err != nil {
		presenter.HandleErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Get Events by Producer ID
// @Description Returns the events belonging to producer
// @ID get-events-by-producer-id
// @Tags events
// @Produce json
// @Security BearerToken
// @Param id path string true "Producer ID"
// @Success 200 {object} []models.Event
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 404 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /producer/{id}/events [get]
func (a *AgendaController) getEventsByProducer(c *gin.Context) {
	id, err := GetId(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}

	events, err := a.agendaService.GetEventsByProducer(c, id)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, events)
	return
}

// @Summary Get Applications by Event ID
// @Description Returns the applications submitted to an event
// @ID get-applications-by-event-id
// @Tags applications
// @Produce json
// @Security BearerToken
// @Param id path string true "Event ID"
// @Success 200 {object} []models.Application
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 404 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /event/{id}/applications [get]
func (a *AgendaController) getApplicationsByEvent(c *gin.Context) {
	id, err := GetId(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	applications, err := a.agendaService.GetApplicationsByEvent(c, id)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, applications)
}

// @Summary Get Applications by Performer ID
// @Description Returns the applications submitted to an event
// @ID get-applications-by-performer-id
// @Tags applications
// @Produce json
// @Security BearerToken
// @Param id path string true "Performer ID"
// @Success 200 {object} []models.Application
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 404 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /performer/{id}/applications [get]
func (a *AgendaController) getApplicationsByPerformer(c *gin.Context) {
	id, err := GetId(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	applications, err := a.agendaService.GetApplicationsByPerformer(c, id)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, applications)
}

// @Summary Create a new tag
// @Description Create a new tag
// @Tags Tags
// @Produce  json
// @Security BasicAuth
// @Success 201 {object} presenter.IdResponse
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 401 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /tag/{name} [post]
func (a *AgendaController) createTag(c *gin.Context) {
	name := c.Param("name")
	tag := models.Tag{
		Name: name,
	}
	id, err := a.agendaService.CreateTag(c, tag)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, presenter.IdResponse{Id: fmt.Sprintf("%d", id)})
}

// @Summary Delete a tag
// @Description Delete a tag
// @Tags Tags
// @Produce  json
// @Security BasicAuth
// @Param event body models.Event true "Event object to be created"
// @Success 201 {object} presenter.IdResponse
// @Failure 400 {object} presenter.ErrorResponse
// @Failure 401 {object} presenter.ErrorResponse
// @Failure 500 {object} presenter.ErrorResponse
// @Router /tag/{name} [delete]
func (a *AgendaController) deleteTag(c *gin.Context) {
	name := c.Param("name")
	tag := models.Tag{
		Name: name,
	}
	err := a.agendaService.DeleteTag(c, tag)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func CreateAgendaHanlder(
	service agenda.Service,
	router *gin.Engine,
	firebaseMiddleware middleware.FirebaseMiddleware,
	permissionsMiddleware middleware.PermissionsMiddleware,
) AgendaController {
	handler := AgendaController{service}
	router.Use(firebaseMiddleware.AuthMiddleware)
	router.POST("/events", handler.createEvent)
	router.GET("/events/:id", handler.getEvent)
	router.GET("/events", handler.getEventsByFilter)
	router.GET("/producer/:id/events", handler.getEventsByProducer)
	router.PATCH("/events/:id", permissionsMiddleware.EventModifier, handler.updateEvent)
	router.DELETE("/events/:id", permissionsMiddleware.EventModifier, handler.deleteEvent)
	router.POST("/applications", handler.createApplication)
	router.GET("/applications/:id", permissionsMiddleware.ApplicationViewer, handler.getApplication)
	router.GET("/events/:id/applications", permissionsMiddleware.EventModifier, handler.getApplicationsByEvent)
	router.GET("/performer/:id/applications", permissionsMiddleware.ApplicationModifier, handler.getApplicationsByPerformer)
	router.PATCH("/applications/:id", permissionsMiddleware.ApplicationModifier, handler.updateApplication)
	router.DELETE("/applications/:id", permissionsMiddleware.ApplicationModifier, handler.deleteApplication)
	router.POST("/tag/:name", permissionsMiddleware.Admin, handler.createTag)
	router.POST("/tag/:name", permissionsMiddleware.Admin, handler.deleteTag)
	return handler
}
