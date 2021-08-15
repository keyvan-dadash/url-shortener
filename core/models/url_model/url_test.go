package url_model_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/sod-lol/url-shortener/core/models/url_model"
	"github.com/sod-lol/url-shortener/libs/encoder/base62"
	"github.com/stretchr/testify/assert"
)

func TestCreatingURL(t *testing.T) {
	assert := assert.New(t)

	id := rand.Uint64()
	urlObj := url_model.CreateURLObj(id, "https://google.com", time.Hour*2)

	assert.Equal(urlObj.ID, id)
	assert.Equal(urlObj.OriginalURL, "https://google.com")
	assert.Equal(urlObj.ShortURL, base62.Encode(id))
	assert.Equal(urlObj.ExpireTime, time.Hour*2)
	assert.Equal(urlObj.Clicked, 0)
}
