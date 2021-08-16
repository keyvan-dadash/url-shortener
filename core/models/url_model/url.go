package url_model

import (
	"time"

	"github.com/sod-lol/url-shortener/libs/encoder/base62"
)

//URL is strcut which represent all short url attributs
type URL struct {
	ID          uint64
	OriginalURL string
	ShortURL    string
	CreatedTime time.Time
	ExpireTime  time.Duration
	Clicked     uint64
}

//CreateURLObj is function which give ID, OriginalURL and ExpireTime then give you whole
// URL obj which is auto generated short url
func CreateURLObj(ID uint64, OriginalURL string, ExpireTime time.Duration) *URL {

	return &URL{
		ID:          ID,
		OriginalURL: OriginalURL,
		ShortURL:    base62.Encode(ID),
		CreatedTime: time.Now(),
		ExpireTime:  ExpireTime,
		Clicked:     0,
	}
}
