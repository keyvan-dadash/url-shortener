package base62_test

import (
	"testing"

	"github.com/sod-lol/url-shortener/libs/encoder/base62"
	"github.com/stretchr/testify/assert"
)

func TestBase62Encode(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(base62.Encode(uint64(473475984745793)), "2ARoEW4Fl")
}

func TestBase62Decode(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(base62.Decode("2ARoEW4Fl"), uint64(473475984745793))
}
