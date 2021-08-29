package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	gin.DefaultWriter = ioutil.Discard

	router := gin.New()

	root := context.Background()

	redisClient := redis.CreateRedisClient("127.0.0.1:10332", "", 0)
	urlRepo := url_repo.URLRedisStorage{
		Client: redisClient,
	}

	ctxWithRepo := url_repo.SetURLRepoInContext(root, &urlRepo)

	routers.InitRoutes(ctxWithRepo, &router.RouterGroup)

	return router, &urlRepo
}

func shutdownRouter(redisClient *url_repo.URLRedisStorage) {
	redisClient.Close()
}

type RespTestBody struct {
	ShortURL string `json:"short_url"`
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

	//clean up
	resp := &RespTestBody{}
	err := json.Unmarshal(w.Body.Bytes(), resp)

	assert.Equal(err, nil)

	root := context.Background()
	err = urlRepo.DeleteURLByShortURL(root, resp.ShortURL)

	assert.Equal(err, nil)

}

type GetURLInfoRespTestBody struct {
	ID          uint64 `json:"ID"`
	OriginalURL string `json:"OriginalURL"`
	ShortURL    string `json:"ShortURL"`
	Clicked     uint64 `json:"Clicked"`
}

func TestGetURLInfo(t *testing.T) {
	assert := assert.New(t)
	routerSetup, urlRepo := setUpRouters()
	defer shutdownRouter(urlRepo)

	w := httptest.NewRecorder()

	//submit url
	jsonBody := `{"url": "https://www.google.com"}`
	req, _ := http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)
	assert.Equal(http.StatusCreated, w.Code)

	respGoogle := &RespTestBody{}
	err := json.Unmarshal(w.Body.Bytes(), respGoogle)

	assert.Equal(err, nil)

	//get url info test
	w = httptest.NewRecorder()

	completeURL := string("/url/info/") + respGoogle.ShortURL
	req, _ = http.NewRequest("GET", completeURL, nil)
	routerSetup.ServeHTTP(w, req)

	respInfo := w.Result()
	body, _ := io.ReadAll(respInfo.Body)

	type Response struct {
		URLInformation GetURLInfoRespTestBody `json:"url_information"`
	}

	var wholeResp Response
	err = json.Unmarshal(body, &wholeResp)
	respGoogleInfo := wholeResp.URLInformation
	assert.Equal(nil, err)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("https://www.google.com", respGoogleInfo.OriginalURL)
	assert.Equal(respGoogle.ShortURL, respGoogleInfo.ShortURL)
	assert.Equal(uint64(0), respGoogleInfo.Clicked)

	//clean up
	resp := &RespTestBody{}
	err = json.Unmarshal(w.Body.Bytes(), resp)

	assert.Equal(err, nil)
	fmt.Println(respGoogleInfo.ShortURL)
	root := context.Background()
	err = urlRepo.DeleteURLByShortURL(root, respGoogleInfo.ShortURL)

	assert.Equal(err, nil)

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

	respGoogle := &RespTestBody{}
	err := json.Unmarshal(w.Body.Bytes(), respGoogle)

	assert.Equal(err, nil)

	//2 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.facebook.com"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	respFacebook := &RespTestBody{}
	err = json.Unmarshal(w.Body.Bytes(), respFacebook)

	assert.Equal(err, nil)

	//3 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.twitch.tv"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	respTwitch := &RespTestBody{}
	err = json.Unmarshal(w.Body.Bytes(), respTwitch)

	assert.Equal(err, nil)

	//4 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.youtube.com"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	respYoutube := &RespTestBody{}
	err = json.Unmarshal(w.Body.Bytes(), respYoutube)

	assert.Equal(err, nil)

	//Clean up all test
	root := context.Background()

	//clean up google
	err = urlRepo.DeleteURLByShortURL(root, respGoogle.ShortURL)
	assert.Equal(err, nil)

	//clean up facebook
	err = urlRepo.DeleteURLByShortURL(root, respFacebook.ShortURL)
	assert.Equal(err, nil)

	//clean up twitch
	err = urlRepo.DeleteURLByShortURL(root, respTwitch.ShortURL)
	assert.Equal(err, nil)

	//clean up youtube
	err = urlRepo.DeleteURLByShortURL(root, respYoutube.ShortURL)
	assert.Equal(err, nil)
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

	respGoogle := &RespTestBody{}
	err := json.Unmarshal(w.Body.Bytes(), respGoogle)

	assert.Equal(err, nil)

	//2 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.google.com"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	respGoogle1 := &RespTestBody{}
	err = json.Unmarshal(w.Body.Bytes(), respGoogle1)

	assert.Equal(err, nil)

	//3 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.google.com"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	respGoogle2 := &RespTestBody{}
	err = json.Unmarshal(w.Body.Bytes(), respGoogle2)

	assert.Equal(err, nil)

	//4 test
	w = httptest.NewRecorder()

	jsonBody = `{"url": "https://www.google.com"}`
	req, _ = http.NewRequest("POST", "/url/submit", strings.NewReader(string(jsonBody)))
	routerSetup.ServeHTTP(w, req)

	assert.Equal(http.StatusCreated, w.Code)

	respGoogle3 := &RespTestBody{}
	err = json.Unmarshal(w.Body.Bytes(), respGoogle3)

	assert.Equal(err, nil)

	//Clean up all test
	root := context.Background()

	//clean up google1
	err = urlRepo.DeleteURLByShortURL(root, respGoogle.ShortURL)
	assert.Equal(err, nil)

	//clean up google2
	err = urlRepo.DeleteURLByShortURL(root, respGoogle1.ShortURL)
	assert.Equal(err, nil)

	//clean up google3
	err = urlRepo.DeleteURLByShortURL(root, respGoogle2.ShortURL)
	assert.Equal(err, nil)

	//clean up google4
	err = urlRepo.DeleteURLByShortURL(root, respGoogle3.ShortURL)
	assert.Equal(err, nil)
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

	//clean up
	root := context.Background()
	err = urlRepo.DeleteURLByShortURL(root, resp.ShortURL)
	assert.Equal(err, nil)
}
