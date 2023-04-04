package middleware

import (
	"backend/boundary/presenter"
	"backend/models"
	"backend/usecase/agenda"
	"backend/usecase/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

const ParamIdContextKey string = "contextId"

type PermissionsMiddleware struct {
	uService users.Service
	aService agenda.Service
}

func NewPermissionsMiddleware(uService users.Service, aService agenda.Service) PermissionsMiddleware {
	return PermissionsMiddleware{uService: uService, aService: aService}
}

func (m *PermissionsMiddleware) setID(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.Wrap(err, "unable to parse id").Error()})
		return uuid.Nil, errors.New("finished call")
	}
	c.Set(ParamIdContextKey, id)
	return id, nil
}

func (m *PermissionsMiddleware) ProfileModifier(c *gin.Context) {
	profileId, err := m.setID(c)
	if err != nil {
		return
	}
	firebaseId, exists := c.Get(models.FirebaseContextKey)
	if !exists {
		c.Next()
		return
	}

	_users, err := m.uService.GetUsersByProfileId(c, profileId)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	for _, user := range _users {
		if user.FirebaseId == firebaseId {
			c.Next()
			return
		}
	}

	c.Status(http.StatusForbidden)
	return

}

func (m *PermissionsMiddleware) ApplicationViewer(c *gin.Context) {
	firebaseID, exists := c.Get(models.FirebaseContextKey)
	if !exists {
		c.Next()
		return
	}
	applicationId, err := m.setID(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	app, err := m.aService.GetApplication(c, applicationId)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	event, err := m.aService.GetEvent(c, app.EventRef)
	if err != nil {
		c.Next()
	}
	appUsers, err := m.uService.GetUsersByProfileId(c, app.PerformerID)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	eventUsers, err := m.uService.GetUsersByProfileId(c, event.ProducerID)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	profileUsers := append(appUsers, eventUsers...)
	for _, p := range profileUsers {
		if p.FirebaseId == firebaseID {
			c.Next()
			return
		}
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"how": "did we get here?"})
}

func (m *PermissionsMiddleware) EventModifier(c *gin.Context) {
	firebaseId, ok := c.Get(models.FirebaseContextKey)
	if !ok {
		c.Next()
		return
	}
	eventId, err := m.setID(c)
	if err != nil {
		return
	}
	event, err := m.aService.GetEvent(c, eventId)
	profileUsers, err := m.uService.GetUsersByProfileId(c, event.ProducerID)
	for _, p := range profileUsers {
		if p.FirebaseId == firebaseId {
			c.Next()
			return
		}
	}
	c.Status(http.StatusForbidden)
	return

}

func (m *PermissionsMiddleware) ApplicationModifier(c *gin.Context) {
	firebaseId, exists := c.Get(models.FirebaseContextKey)
	if !exists {
		c.Next()
		return
	}
	id, err := m.setID(c)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	app, err := m.aService.GetApplication(c, id)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	profileUsers, err := m.uService.GetUsersByProfileId(c, app.Performer.ID)
	if err != nil {
		presenter.HandleErr(c, err)
		return
	}
	for _, u := range profileUsers {
		if u.FirebaseId == firebaseId {
			c.Next()
			return
		}
	}
	c.Status(http.StatusForbidden)
}

func (m *PermissionsMiddleware) Admin(c *gin.Context) {
	if _, exists := c.Get(models.FirebaseContextKey); exists {
		c.Status(http.StatusUnauthorized)
		return
	}
	c.Next()
}
