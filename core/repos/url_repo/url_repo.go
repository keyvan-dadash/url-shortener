package url_repo

import (
	"context"

	"github.com/sod-lol/url-shortener/core/models/url_model"
)

//URLRepo is interface which is every database storage for
// storing url obj should implement this
type URLRepo interface {

	//query API's
	IsValidShortURL(ctx context.Context, shortUrl string) bool
	GetURLByShortURL(ctx context.Context, shortUrl string) (*url_model.URL, error)
	IsValidID(ctx context.Context, id uint64) bool
	GetURLByID(ctx context.Context, id uint64) (*url_model.URL, error)

	//modify AIP's
	SaveURL(ctx context.Context, url *url_model.URL) error
	DeleteURLByID(ctx context.Context, id uint64) error
	DeleteURLByShortURL(ctx context.Context, shortUrl string) error
}
