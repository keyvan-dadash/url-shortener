package url_repo

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/sod-lol/url-shortener/core/models/url_model"
)

type URLRedisStorage struct {
	*redis.Client
}

//Queries

func (urs *URLRedisStorage) IsValidShortURL(ctx context.Context, shortUrl string) bool {
	if _, err := urs.Client.Get(ctx, shortUrl).Result(); err != nil {
		if err != redis.Nil {
			logrus.Errorf("[ERROR] redis get key %v faild with err %v\n", shortUrl, err)
		}
		return false
	}

	return true
}

func (urs *URLRedisStorage) GetURLByShortURL(ctx context.Context, shortUrl string) (*url_model.URL, error) {
	val, err := urs.Client.Get(ctx, shortUrl).Result()
	if err != nil {
		if err != redis.Nil {
			logrus.Errorf("[ERROR] redis get key %v faild with err %v\n", shortUrl, err)
		}
		return nil, err
	}

	url := new(url_model.URL)

	if err = json.Unmarshal([]byte(val), url); err != nil {
		logrus.Errorf("[ERROR] unmarshal json at key %v faild with err %v\n", shortUrl, err)
		return nil, err
	}

	return url, nil
}

func (urs *URLRedisStorage) IsValidID(ctx context.Context, id uint64) bool {
	if _, err := urs.Client.Get(ctx, strconv.FormatUint(id, 10)).Result(); err != nil {
		if err != redis.Nil {
			logrus.Errorf("[ERROR] redis get key %v faild with err %v\n", id, err)
		}
		return false
	}

	return true
}

func (urs *URLRedisStorage) GetURLByID(ctx context.Context, id uint64) (*url_model.URL, error) {
	val, err := urs.Client.Get(ctx, strconv.FormatUint(id, 10)).Result()
	if err != nil {
		if err != redis.Nil {
			logrus.Errorf("[ERROR] redis get key %v faild with err %v\n", id, err)
		}
		return nil, err
	}

	url := new(url_model.URL)

	if err = json.Unmarshal([]byte(val), url); err != nil {
		logrus.Errorf("[ERROR] unmarshal json at key %v faild with err %v\n", id, err)
		return nil, err
	}

	return url, nil
}

//Modifiers

func (urs *URLRedisStorage) SaveURL(ctx context.Context, url *url_model.URL) error {
	if urs.IsValidID(ctx, url.ID) {
		for {
			id := rand.Uint64()
			if !urs.IsValidID(ctx, id) {
				url.ID = id
				break
			}
		}
	}

	bytes, err := json.Marshal(url)

	if err != nil {
		logrus.Errorf("[ERROR] marshal json at key %v faild with err %v\n", url.ShortURL, err)
		return err
	}

	urlStr := string(bytes)

	err = urs.Client.Set(ctx, strconv.FormatUint(url.ID, 10), urlStr, url.ExpireTime).Err()

	if err != nil {
		logrus.Errorf("[ERROR] redis set key %v faild with err %v\n", url.ID, err)
		return err
	}

	err = urs.Client.Set(ctx, string(url.ShortURL), urlStr, url.ExpireTime).Err()

	if err != nil {
		logrus.Errorf("[ERROR] redis set key %v faild with err %v\n", url.ID, err)
		return err
	}

	return err
}

func (urs *URLRedisStorage) UpdateURL(ctx context.Context, url *url_model.URL) error {

	bytes, err := json.Marshal(url)

	if err != nil {
		logrus.Errorf("[ERROR] marshal json at key %v faild with err %v\n", url.ShortURL, err)
		return err
	}

	urlStr := string(bytes)

	err = urs.Client.Set(ctx, strconv.FormatUint(url.ID, 10), urlStr, url.ExpireTime).Err()

	if err != nil {
		logrus.Errorf("[ERROR] redis set key %v faild with err %v\n", url.ID, err)
		return err
	}

	err = urs.Client.Set(ctx, string(url.ShortURL), urlStr, url.ExpireTime).Err()

	if err != nil {
		logrus.Errorf("[ERROR] redis set key %v faild with err %v\n", url.ID, err)
		return err
	}

	return err
}

func (urs *URLRedisStorage) DeleteURLByID(ctx context.Context, id uint64) error {

	url, err := urs.GetURLByID(ctx, id)

	if err != nil {
		logrus.Errorf("[ERROR] get url by key %v faild with err %v\n", id, err)
		return err
	}

	err = urs.Client.Del(ctx, strconv.FormatUint(url.ID, 10), url.ShortURL).Err()

	if err != nil {
		logrus.Errorf("[ERROR] delete key %v faild with err %v\n", url.ID, err)
	}

	return err
}

func (urs *URLRedisStorage) DeleteURLByShortURL(ctx context.Context, shortUrl string) error {

	url, err := urs.GetURLByShortURL(ctx, shortUrl)

	if err != nil {
		logrus.Errorf("[ERROR] get url by key %v faild with err %v\n", shortUrl, err)
		return err
	}

	err = urs.Client.Del(ctx, strconv.FormatUint(url.ID, 10), url.ShortURL).Err()

	if err != nil {
		logrus.Errorf("[ERROR] delete key %v faild with err %v\n", url.ID, err)
	}

	return err
}
