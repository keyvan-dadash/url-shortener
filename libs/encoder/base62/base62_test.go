package base62

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase62Encode(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(Encode(uint64(473475984745793)), "2ARoEW4Fl")
}

func TestBase62Decode(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(Decode("2ARoEW4Fl"), uint64(473475984745793))
}
