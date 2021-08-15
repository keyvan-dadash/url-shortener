package url_c

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sod-lol/url-shortener/core/models/url_model"
	"github.com/sod-lol/url-shortener/core/repos/url_repo"
)

func HandlerRedirect(urlRepo url_repo.URLRepo) gin.HandlerFunc {
	return func(c *gin.Context) {

		shortUrl := c.Param("shortUrl")
		ctx := c.Request.Context()
		url, err := urlRepo.GetURLByShortURL(ctx, shortUrl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		url.Clicked += 1
		err = urlRepo.UpdateURL(ctx, url)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err1": err.Error(),
			})
			return
		}

		c.Redirect(http.StatusMovedPermanently, url.OriginalURL)

	}
}

type shortURLRequest struct {
	URL string `form:"url" json:"url" xml:"url"  binding:"required"`
}

func HandlerShortURLRequest(urlRepo url_repo.URLRepo) gin.HandlerFunc {
	return func(c *gin.Context) {

		var urlRequestJson shortURLRequest

		if err := c.ShouldBindJSON(&urlRequestJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		url := url_model.CreateURLObj(rand.Uint64(), urlRequestJson.URL, time.Hour*2)

		err := urlRepo.SaveURL(ctx, url)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"short_url": url.ShortURL,
		})

	}
}
