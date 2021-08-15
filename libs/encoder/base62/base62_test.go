package base62_test

import (
	"testing"

	"github.com/sod-lol/url-shortener/libs/encoder/base62"
	"github.com/stretchr/testify/assert"
)

func TestBase62Encode(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("eM0cQu", base62.Encode(uint64(4238435464)))
}

func TestBase62Decode(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint64(4238435464), base62.Decode("eM0cQu"))
}
