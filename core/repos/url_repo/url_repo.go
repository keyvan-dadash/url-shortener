package url_repo

import (
	"context"

	"github.com/sod-lol/url-shortener/core/models/url_model"
)

//URLRepo is interface which is every database storage for
// storing url obj should implement this
type URLRepo interface {

	//query API's
	IsValidShortURL(ctx context.Context, shortUrl string)
	GetURLByShortURL(ctx context.Context, shortUrl string)
	IsValidID(ctx context.Context, id int64)
	GetURLByID(ctx context.Context, id int64)

	//modify AIP's
	SaveURL(ctx context.Context, url url_model.URL)
	DeleteURLByID(ctx context.Context, id int64)
	DeleteURLByShortURL(ctx context.Context, shortUrl string)
}
