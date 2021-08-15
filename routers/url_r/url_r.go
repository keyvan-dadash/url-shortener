package url_r

import "github.com/gin-gonic/gin"

func SetUpUrlRoutes(g *gin.RouterGroup) {

	g.POST("/submit/")
	g.GET("/{shortUrl}/")
}
