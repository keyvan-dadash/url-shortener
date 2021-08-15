package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sod-lol/url-shortener/routers/url_r"
)

func InitRoutes(ctx context.Context, g *gin.RouterGroup) {

	url_r.SetUpUrlRoutes(ctx, g.Group("/url"))
}
