package url_r

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sod-lol/url-shortener/controllers/url_c"
	"github.com/sod-lol/url-shortener/core/repos/url_repo"
)

func SetUpUrlRoutes(ctx context.Context, g *gin.RouterGroup) {

	userRepo := ctx.Value("user-repo").(url_repo.URLRepo)
	g.POST("/submit", url_c.HandlerShortURLRequest(userRepo))
	g.GET("/:shortUrl/", url_c.HandlerRedirect(userRepo))
}
