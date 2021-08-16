package url_r

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sod-lol/url-shortener/controllers/url_c"
	"github.com/sod-lol/url-shortener/core/repos/url_repo"
)

//SetUpUrlRoutes for URL shortner app
func SetUpUrlRoutes(ctx context.Context, g *gin.RouterGroup) {

	urlRepo, hasUrlRepo := url_repo.GetURLRepoFromContex(ctx)
	if !hasUrlRepo {
		logrus.Fatalf("[FATAL] context does not have URL Repo")
	}

	g.POST("/submit", url_c.HandlerShortURLRequest(urlRepo))
	g.GET("/:shortUrl", url_c.HandlerRedirect(urlRepo))
}
