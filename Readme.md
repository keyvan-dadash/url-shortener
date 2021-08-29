
# Url Shortner

## Run on docker
for running on docker simply run this command on root of project
```
docker-compose up -d --build
```

## Migrate to another db
for migrate to other databases:

- #### Add databse config

    add your databse config to services folder

-  #### Implement interface

    implement this interface(which is reside at repos/url_repo/url_repo.go):
    ```go
    //URLRepo should be implemented for all databases.
    type URLRepo interface {

        //query API's
        IsValidShortURL(ctx context.Context, shortUrl string) bool
        GetURLByShortURL(ctx context.Context, shortUrl string) (*url_model.URL, error)
        IsValidID(ctx context.Context, id uint64) bool
        GetURLByID(ctx context.Context, id uint64) (*url_model.URL, error)

        //modify AIP's
        SaveURL(ctx context.Context, url *url_model.URL) error
        UpdateURL(ctx context.Context, url *url_model.URL) error
        DeleteURLByID(ctx context.Context, id uint64) error
        DeleteURLByShortURL(ctx context.Context, shortUrl string) error
    }
    ```
- #### Write Test!!
    simply write your test!!

## Urls
this project has 3 urls

- #### Submit url
    submit your url at:
    ```
    {base_url}/url/submit
    ```
    and with body of:
    ```json
    {
        "url": "https://www.example.com"
    }
    ```

- #### Redirect to url
    ```
    {base_url}/url/{short_url}
    ```

- #### Url metrics
    ```
    {base_url}/url/info/{short_url}
    ```
