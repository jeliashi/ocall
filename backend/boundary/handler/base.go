package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func GetId(c *gin.Context) (uuid.UUID, error) {
	if id, err := uuid.Parse(c.Param("id")); err != nil {
		return uuid.Nil, gin.Error{Err: errors.Wrap(err, "unable to parse id"), Type: gin.ErrorTypeBind}
	} else {
		return id, nil
	}
}
