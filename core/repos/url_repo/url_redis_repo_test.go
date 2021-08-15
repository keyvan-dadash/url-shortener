package url_repo_test

import (
	"context"
	"errors"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	"github.com/sod-lol/url-shortener/core/models/url_model"
	"github.com/sod-lol/url-shortener/core/repos/url_repo"
	"github.com/sod-lol/url-shortener/services/redis"
	"github.com/stretchr/testify/assert"
)

func getRedisRepo() url_repo.URLRepo {
	redisClient := redis.CreateRedisClient("redis-storage:6379", "", 0)
	urlRepo := url_repo.URLRedisStorage{
		Client: redisClient,
	}

	return &urlRepo
}

func TestSaveURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	//save url obj
	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	//clean up
	redisRepo.(*url_repo.URLRedisStorage).Del(root, strconv.FormatUint(urlObj.ID, 10), urlObj.ShortURL)
}

func TestUpdateURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	//save url obj
	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	//begin to test update url
	urlObj.Clicked += 1

	err = redisRepo.UpdateURL(root, urlObj)

	assert.Equal(err, nil)
	assert.Equal(urlObj.Clicked, 1)

	//clean up
	redisRepo.(*url_repo.URLRedisStorage).Del(root, strconv.FormatUint(urlObj.ID, 10), urlObj.ShortURL)
}

func TestDeleteURLByID(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	//save url obj
	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	//begin to test delete url by id (also clean up)
	err = redisRepo.DeleteURLByID(root, urlObj.ID)

	assert.Equal(err, nil)

	_, err = redisRepo.(*url_repo.URLRedisStorage).Get(root, strconv.FormatUint(urlObj.ID, 10)).Result()

	assert.Equal(err, redisv8.Nil)

	_, err = redisRepo.(*url_repo.URLRedisStorage).Get(root, urlObj.ShortURL).Result()

	assert.Equal(err, redisv8.Nil)

}

func TestDeleteURLByShortURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	//save url obj
	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	//begin to test delete url by url (also clean up)
	err = redisRepo.DeleteURLByShortURL(root, urlObj.ShortURL)

	assert.Equal(err, nil)

	_, err = redisRepo.(*url_repo.URLRedisStorage).Get(root, urlObj.ShortURL).Result()

	assert.Equal(err, redisv8.Nil)

	_, err = redisRepo.(*url_repo.URLRedisStorage).Get(root, strconv.FormatUint(urlObj.ID, 10)).Result()

	assert.Equal(err, redisv8.Nil)

}

func TestIsValidShortURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	//save url obj
	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	//begin to test is valid short url
	assert.Equal(redisRepo.IsValidShortURL(root, urlObj.ShortURL), true)

	//clean up
	err = redisRepo.DeleteURLByShortURL(root, urlObj.ShortURL)

	assert.Equal(err, nil)

	//check not exist after clean up
	assert.Equal(redisRepo.IsValidShortURL(root, urlObj.ShortURL), false)
}

func TestIsValidID(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	//save url obj
	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	//begin to test is valid id
	assert.Equal(redisRepo.IsValidID(root, urlObj.ID), true)

	err = redisRepo.DeleteURLByID(root, urlObj.ID)

	assert.Equal(err, nil)

	//check not exist after clean up
	assert.Equal(redisRepo.IsValidID(root, urlObj.ID), false)
}

func TestGetURLByShortURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	//save url obj
	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	//begin to test get url by short url
	urlObjFromRepo, err := redisRepo.GetURLByShortURL(root, urlObj.ShortURL)

	assert.Equal(err, nil)

	assert.Equal(reflect.DeepEqual(urlObj, urlObjFromRepo), true)

	//clean up
	redisRepo.(*url_repo.URLRedisStorage).Del(root, strconv.FormatUint(urlObj.ID, 10), urlObj.ShortURL)
}

func TestGetURLByID(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	//save url obj
	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	//begin to test get url by id
	urlObjFromRepo, err := redisRepo.GetURLByID(root, urlObj.ID)

	assert.Equal(err, nil)

	assert.Equal(reflect.DeepEqual(urlObj, urlObjFromRepo), true)

	//clean up
	redisRepo.(*url_repo.URLRedisStorage).Del(root, strconv.FormatUint(urlObj.ID, 10), urlObj.ShortURL)
}

func TestTimeoutAndDeletedURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	//save url obj
	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Second*1)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	//clean up with time out
	time.Sleep(time.Second * 1)

	_, err = redisRepo.(*url_repo.URLRedisStorage).Get(root, urlObj.ShortURL).Result()

	assert.Equal(err, redisv8.Nil)

	_, err = redisRepo.(*url_repo.URLRedisStorage).Get(root, strconv.FormatUint(urlObj.ID, 10)).Result()

	assert.Equal(err, redisv8.Nil)
}

func TestNotFoundGetURLByShortURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Second*1)

	_, err := redisRepo.GetURLByShortURL(root, urlObj.ShortURL)

	assert.Equal(errors.Is(err, url_model.ErrShortURLDoesNotExists), true)

}

func TestNotFoundGetURLByID(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*url_repo.URLRedisStorage).Close()

	id := rand.Uint64()

	_, err := redisRepo.GetURLByID(root, id)

	assert.Equal(errors.Is(err, url_model.ErrShortURLDoesNotExists), true)

}
