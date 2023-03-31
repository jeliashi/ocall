package middleware

import (
	"backend/models"
	"backend/usecase/agenda"
	"backend/usecase/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

type PermissionsMiddleware struct {
	UService users.Service
	AService agenda.Service
}

func (m *PermissionsMiddleware) setID(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, errors.Wrap(err, "unable to parse id").Error())
		return uuid.Nil, errors.New("finished call")
	}
	c.Set("id", id)
	return id, nil
}

func (m *PermissionsMiddleware) AllowedUser(c *gin.Context) {
	id, err := m.setID(c)
	if err != nil {
		return
	}
	_, ok := c.Get(models.FirebaseContextKey)
	profiles, err := m.UService.GetProfilesByUser(c)
	var match = !ok
	for _, p := range profiles {
		if id == p.ID {
			match = true
			break
		}
	}
	if !match {
		c.Status(http.StatusForbidden)
		return
	}
	c.Next()
}

func (m *PermissionsMiddleware) ApplicationViewer(c *gin.Context) {
	id, err := m.setID(c)
	if err != nil {
		return
	}
	app, err := m.AService.GetApplication(c, id)
	if err != nil {
		c.Next()
	}
	event, err := m.AService.GetEvent(c, app.EventRef)
	if err != nil {
		c.Next()
	}
	_, ok := c.Get(models.FirebaseContextKey)
	profiles, err := m.UService.GetProfilesByUser(c)
	var match = !ok
	for _, p := range profiles {
		if (app.Performer.ID == p.ID) || (event.Producer.ID == p.ID) {
			match = true
			break
		}
	}
	if !match {
		c.Status(http.StatusForbidden)
		return
	}
	c.Next()
}

func (m *PermissionsMiddleware) EventModifier(c *gin.Context) {
	id, err := m.setID(c)
	if err != nil {
		return
	}
	event, err := m.AService.GetEvent(c, id)
	_, ok := c.Get(models.FirebaseContextKey)
	profiles, err := m.UService.GetProfilesByUser(c)
	var match = !ok
	for _, p := range profiles {
		if event.Producer.ID == p.ID {
			match = true
			break
		}
	}
	if !match {
		c.Status(http.StatusForbidden)
		return
	}
	c.Next()
}
