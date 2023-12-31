package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zenorachi/todo-service/internal/entity"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, entity.ErrEmptyAuthHeader.Error())
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, entity.ErrInvalidAuthHeader.Error())
		return
	}

	userId, err := h.tokenManager.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}
