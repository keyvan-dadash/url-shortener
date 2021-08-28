package url_c

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sod-lol/url-shortener/core/models/url_model"
	"github.com/sod-lol/url-shortener/core/repos/url_repo"
)

//HandlerRedirect is gin func handler for when user enter short url and we must redirect
func HandlerRedirect(urlRepo url_repo.URLRepo) gin.HandlerFunc {
	return func(c *gin.Context) {

		shortUrl := c.Param("shortUrl")
		ctx := c.Request.Context()
		url, err := urlRepo.GetURLByShortURL(ctx, shortUrl)

		if err != nil && errors.Is(err, url_model.ErrShortURLDoesNotExists) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "try again later",
			})
			return
		}

		url.Clicked += 1
		err = urlRepo.UpdateURL(ctx, url)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "try again later",
			})
			return
		}

		c.Redirect(http.StatusMovedPermanently, url.OriginalURL)

	}
}

type shortURLRequest struct {
	URL string `form:"url" json:"url" xml:"url"  binding:"required"`
}

//HandlerShortURLRequest is gin func handler which is for submit url and take short url
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
				"error": "try again later",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"short_url": url.ShortURL,
		})

	}
}

//GET
// /HandlerShortURLInformation if function which is return information about given short url
func HandlerShortURLInformation(urlRepo url_repo.URLRepo) gin.HandlerFunc {
	return func(c *gin.Context) {

		shortUrl := c.Param("shortUrl")

		ctx := c.Request.Context()

		url_struct, err := urlRepo.GetURLByShortURL(ctx, shortUrl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "try again later",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"url_information": url_struct,
		})

	}
}
