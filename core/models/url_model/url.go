package url_model

import (
	"time"

	"github.com/sod-lol/url-shortener/libs/encoder/base62"
)

//URL is strcut which represent all short url attributs
type URL struct {
	id          int64
	OriginalURL string
	ShortURL    string
	CreatedTime time.Time
	ExpireTime  time.Time
	Clicked     int64
}

func CreateURLObj(ID int64, OriginalURL string, ExpireTime time.Time) *URL {

	return &URL{
		id:          ID,
		OriginalURL: OriginalURL,
		ShortURL:    base62.Encode(ID),
		CreatedTime: time.Now(),
		ExpireTime:  ExpireTime,
		Clicked:     0,
	}
}
