package http

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/sial-soft/oauth-api/src/domain/access_token"
	access_token2 "github.com/sial-soft/oauth-api/src/services/access_token"
	"github.com/sial-soft/oauth-api/src/utils/rest_errors"
	"net/http"
)

type AccessTokenHandler interface {
	GetById(c *gin.Context)
	Create(c *gin.Context)
}

type accessTokenHandler struct {
	service access_token2.Service
}

func NewHandler(service access_token2.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetById(c *gin.Context) {
	accessToken, err := handler.service.GetById(c.Param("access_token_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var at access_token.AccessTokenRequest
	if err := c.ShouldBindJSON(&at); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	if _, err := handler.service.Create(at); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, at)
}
