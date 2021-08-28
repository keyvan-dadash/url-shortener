package url_repo

import "context"

type repoKey string

var (
	URLRepoKey = repoKey("url-repo")
)

func SetURLRepoInContext(parentCtx context.Context, urlRepo URLRepo) context.Context {
	return context.WithValue(parentCtx, URLRepoKey, urlRepo)
}

func GetURLRepoFromContex(ctx context.Context) (URLRepo, bool) {

	urlRepo := ctx.Value(URLRepoKey).(URLRepo)

	if urlRepo == nil {
		return nil, false
	}

	return urlRepo, true
}
