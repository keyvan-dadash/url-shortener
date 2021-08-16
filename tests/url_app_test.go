package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sod-lol/url-shortener/core/repos/url_repo"
	"github.com/sod-lol/url-shortener/routers"
	"github.com/sod-lol/url-shortener/services/redis"
	"github.com/stretchr/testify/assert"
)

func setUpRouters() (*gin.Engine, *url_repo.URLRedisStorage) {
	router := gin.New()

	root := context.Background()

	redisClient := redis.CreateRedisClient("127.0.0.1:10332", "", 0)
	urlRepo := url_repo.URLRedisStorage{
		Client: redisClient,
	}

	ctxWithRepo := url_repo.SetURLRepoInContext(root, &urlRepo)

	routers.InitRoutes(ctxWithRepo, &router.RouterGroup)

	router.Run(":8080")

	return router, &urlRepo
}

func shutdownRouter(redisClient *url_repo.URLRedisStorage) {
	redisClient.Close()
}

func TestSubmitURL(t *testing.T) {
	assert := assert.New(t)
	routerSetup, urlRepo := setUpRouters()
	defer shutdownRouter(urlRepo)

	w := httptest.NewRecorder()

	jsonBody := `{"url": "https://www.google.com"}`
	req, _ := http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)
}

func TestSubmitMultiplyURL(t *testing.T) {
	assert := assert.New(t)
	routerSetup, urlRepo := setUpRouters()
	defer shutdownRouter(urlRepo)

	//1 test
	w := httptest.NewRecorder()

	jsonBody := `{"url": "https://www.google.com"}`
	req, _ := http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	//2 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.facebook.com"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	//3 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.twitch.tv"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	//4 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.youtube.com"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)
}

func TestSubmitMultiplySameURL(t *testing.T) {
	assert := assert.New(t)
	routerSetup, urlRepo := setUpRouters()
	defer shutdownRouter(urlRepo)

	//1 test
	w := httptest.NewRecorder()

	jsonBody := `{"url": "https://www.google.com"}`
	req, _ := http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	//2 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.google.com"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	//3 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.google.com"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	//4 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.google.com"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)
}

type RespTestBody struct {
	ShortURL string `json:"short_url"`
}

func TestRedirect(t *testing.T) {
	assert := assert.New(t)
	routerSetup, urlRepo := setUpRouters()
	defer shutdownRouter(urlRepo)

	//test router
	testRouter := gin.New()
	testRouter.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"test_string": "Good!",
		})
	})

	//run on other goroutine because it's blocking
	go func() {
		testRouter.Run(":10923")
	}()

	w := httptest.NewRecorder()

	//set up test url
	jsonBody := `{"url": "http://127.0.0.1:10923/test"}`
	req, _ := http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	resp := &RespTestBody{}
	err := json.Unmarshal(w.Body.Bytes(), resp)

	assert.Equal(err, nil)

	//begin testing redirect
	w = httptest.NewRecorder()
	completeURL := string("/url/") + resp.ShortURL
	req, _ = http.NewRequest("GET", completeURL, nil)
	routerSetup.ServeHTTP(w, req)

	assert.Equal("http://127.0.0.1:10923/test", w.Header().Get("Location"))
	assert.Equal(http.StatusMovedPermanently, w.Code)
}
