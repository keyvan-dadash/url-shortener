package url_repo

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
	"github.com/sod-lol/url-shortener/services/redis"
	"github.com/stretchr/testify/assert"
)

func getRedisRepo() URLRepo {
	redisClient := redis.CreateRedisClient("redis-storage:6379", "", 0)
	urlRepo := URLRedisStorage{
		Client: redisClient,
	}

	return &urlRepo
}

func TestSaveURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)
}

func TestUpdateURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	urlObj.Clicked += 1

	err = redisRepo.UpdateURL(root, urlObj)

	assert.Equal(err, nil)
	assert.Equal(urlObj.Clicked, 1)
}

func TestDeleteURLByID(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	err = redisRepo.DeleteURLByID(root, urlObj.ID)

	assert.Equal(err, nil)

	_, err = redisRepo.(*URLRedisStorage).Get(root, strconv.FormatUint(urlObj.ID, 10)).Result()

	assert.Equal(err, redisv8.Nil)

	_, err = redisRepo.(*URLRedisStorage).Get(root, urlObj.ShortURL).Result()

	assert.Equal(err, redisv8.Nil)

}

func TestDeleteURLByShortURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	err = redisRepo.DeleteURLByShortURL(root, urlObj.ShortURL)

	assert.Equal(err, nil)

	_, err = redisRepo.(*URLRedisStorage).Get(root, urlObj.ShortURL).Result()

	assert.Equal(err, redisv8.Nil)

	_, err = redisRepo.(*URLRedisStorage).Get(root, strconv.FormatUint(urlObj.ID, 10)).Result()

	assert.Equal(err, redisv8.Nil)

}

func TestIsValidShortURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	assert.Equal(redisRepo.IsValidShortURL(root, urlObj.ShortURL), true)

	err = redisRepo.DeleteURLByShortURL(root, urlObj.ShortURL)

	assert.Equal(err, nil)

	assert.Equal(redisRepo.IsValidShortURL(root, urlObj.ShortURL), false)
}

func TestIsValidID(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	assert.Equal(redisRepo.IsValidID(root, urlObj.ID), true)

	err = redisRepo.DeleteURLByID(root, urlObj.ID)

	assert.Equal(err, nil)

	assert.Equal(redisRepo.IsValidID(root, urlObj.ID), false)
}

func TestGetURLByShortURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	urlObjFromRepo, err := redisRepo.GetURLByShortURL(root, urlObj.ShortURL)

	assert.Equal(err, nil)

	assert.Equal(reflect.DeepEqual(urlObj, urlObjFromRepo), true)
}

func TestGetURLByID(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	urlObjFromRepo, err := redisRepo.GetURLByID(root, urlObj.ID)

	assert.Equal(err, nil)

	assert.Equal(reflect.DeepEqual(urlObj, urlObjFromRepo), true)
}

func TestTimeoutAndDeletedURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Second*1)

	err := redisRepo.SaveURL(root, urlObj)

	assert.Equal(err, nil)

	time.Sleep(time.Second * 1)

	_, err = redisRepo.(*URLRedisStorage).Get(root, urlObj.ShortURL).Result()

	assert.Equal(err, redisv8.Nil)

	_, err = redisRepo.(*URLRedisStorage).Get(root, strconv.FormatUint(urlObj.ID, 10)).Result()

	assert.Equal(err, redisv8.Nil)
}

func TestNotFoundGetURLByShortURL(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Second*1)

	_, err := redisRepo.GetURLByShortURL(root, urlObj.ShortURL)

	assert.Equal(errors.Is(err, url_model.ErrShortURLDoesNotExists), true)

}

func TestNotFoundGetURLByID(t *testing.T) {
	assert := assert.New(t)
	root := context.Background()

	redisRepo := getRedisRepo()
	defer redisRepo.(*URLRedisStorage).Close()

	id := rand.Uint64()

	_, err := redisRepo.GetURLByID(root, id)

	assert.Equal(errors.Is(err, url_model.ErrShortURLDoesNotExists), true)

}
