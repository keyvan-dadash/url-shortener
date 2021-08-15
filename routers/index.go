package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sod-lol/url-shortener/routers/url_r"
)

func InitRoutes(g *gin.RouterGroup) {

	url_r.SetUpUrlRoutes(g.Group("/url"))
}
