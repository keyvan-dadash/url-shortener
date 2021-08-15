package url

import (
	"time"
)

//URL is strcut which represent all short url attributs
type URL struct {
	id          int64
	OriginalURL string
	ShortURL    string
	CreatedTime time.Duration
	ExpireTime  time.Duration
	Clicked     int64
}

func CreateURLObj(OriginalURL string, ExpireTime time.Duration) (*URL, error) {

}
