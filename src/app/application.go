package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sial-soft/oauth-api/src/http"
	"github.com/sial-soft/oauth-api/src/repository/db"
	"github.com/sial-soft/oauth-api/src/repository/rest"
	access_token2 "github.com/sial-soft/oauth-api/src/services/access_token"
	"log"
)

var (
	router = gin.Default()
)

func StartApplication() {

	atHandler := http.NewHandler(access_token2.NewService(rest.NewRestUserRepository(), db.NewDbRepository()))
	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	log.Fatal(router.Run(":8080"))
}
