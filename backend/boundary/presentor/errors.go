package presentor

import (
	"backend/models"
	"github.com/pkg/errors"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func HandleErr(err error) (int, ErrorResponse) {
	if errors.Is(err, &models.NotFoundErr{}) {
		return http.StatusNotFound, ErrorResponse{Error: err.Error()}
	}
	return http.StatusInternalServerError, ErrorResponse{Error: err.Error()}
}
