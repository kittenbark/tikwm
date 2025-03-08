package tikwm

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	DefaultURL          = "https://tikwm.com/api"
	DefaultTimeout      = time.Millisecond * 1_100
	DefaultUserFeedSize = 33

	EnvProxyUrl = "KITTENBARK_TIKWM_PROXY_URL"
)

type API struct {
	Urls    []string
	Timeout time.Duration
	Http    *http.Client
	Sync    *sync.Mutex
	Logger  *slog.Logger
}

var DefaultAPI = sync.OnceValue(func() *API {
	return &API{
		Urls:    apiUrls(),
		Timeout: DefaultTimeout,
		Http:    http.DefaultClient,
		Sync:    &sync.Mutex{},
		Logger:  slog.New(slog.DiscardHandler),
	}
})

func (api *API) Post(ctx context.Context, url string, hd ...bool) (*UserPost, error) {
	query := map[string]string{"url": url}
	if len(hd) == 0 || hd[0] {
		query["hd"] = "1"
	}
	return rawMult[UserPost](api, ctx, "", query)
}

func (api *API) Feed(ctx context.Context, uniqueID string, count int, cursor string) (*UserFeed, error) {
	query := map[string]string{"unique_id": uniqueID, "count": strconv.Itoa(count), "cursor": cursor}
	if _, err := strconv.ParseInt(uniqueID, 10, 64); err == nil {
		query = map[string]string{"user_id": uniqueID, "count": strconv.Itoa(count), "cursor": cursor}
	}
	return rawMult[UserFeed](api, ctx, "user/posts", query)
}

func (api *API) Details(ctx context.Context, uniqueID string) (*UserDetail, error) {
	return rawMult[UserDetail](api, ctx, "user/info", map[string]string{"unique_id": uniqueID})
}

func (api *API) FeedSeq(ctx context.Context, uniqueID string, hd ...bool) iter.Seq2[*UserPost, error] {
	HD := at(hd, 0, true)
	cursor := "0"
	return func(yield func(*UserPost, error) bool) {
		for {
			feed, err := api.Feed(ctx, uniqueID, DefaultUserFeedSize, cursor)
			if err != nil {
				yield(nil, err)
				return
			}
			for _, post := range feed.Posts {
				if HD {
					post, err = api.Post(ctx, post.ID(), HD)
				}
				if !yield(post, err) {
					return
				}
			}
			if !feed.HasMore {
				break
			}
			cursor = feed.Cursor
		}
	}
}

func Raw[T any](api *API, ctx context.Context, apiURL string, method string, query map[string]string) (*T, error) {
	api.Logger.InfoContext(ctx, "tikwm#raw", "apiURL", apiURL, "method", method, "query", query)
	if api.Sync != nil {
		api.Sync.Lock()
		defer time.AfterFunc(api.Timeout, api.Sync.Unlock)
	}

	requestURL, err := url.JoinPath(apiURL, method)
	if err != nil {
		api.Logger.ErrorContext(ctx, "tikwm#url", "url", apiURL, "method", method, "err", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, nil)
	if err != nil {
		api.Logger.ErrorContext(ctx, "tikwm#new_req", "url", apiURL, "method", method, "err", err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	reqQuery := req.URL.Query()
	for key, val := range query {
		reqQuery.Add(key, url.QueryEscape(val))
	}
	req.URL.RawQuery = reqQuery.Encode()

	type ResponseSchema struct {
		Code          int     `json:"code"`
		Msg           string  `json:"msg"`
		ProcessedTime float64 `json:"processed_time"`
		Data          *T      `json:"data"`
	}
	httpResponse, err := api.Http.Do(req)
	if err != nil {
		api.Logger.ErrorContext(ctx, "tikwm#request", "url", apiURL, "method", method, "err", err)
		return nil, err
	}

	var resp ResponseSchema
	if err = json.NewDecoder(httpResponse.Body).Decode(&resp); err != nil {
		api.Logger.ErrorContext(ctx, "tikwm#decode", "url", apiURL, "method", method, "err", err)
		return nil, err
	}
	if resp.Code != 0 {
		queryStr := "???"
		if buf, err := json.Marshal(query); err == nil {
			queryStr = string(buf)
		}
		err = fmt.Errorf("tikwm error: %s (%d) [%s, query: %s]", resp.Msg, resp.Code, method, queryStr)
		api.Logger.ErrorContext(ctx, "tikwm#api_error", "url", apiURL, "method", method, "err", err)
		return nil, err
	}

	return resp.Data, nil
}

func rawMult[T any](api *API, ctx context.Context, method string, query map[string]string) (res *T, err error) {
	for _, apiURL := range api.Urls {
		if res, err = Raw[T](api, ctx, apiURL, method, query); err == nil {
			return res, nil
		}
	}
	return
}

func at[T any](list []T, pos int, otherwise T) T {
	if pos < len(list) {
		return list[pos]
	}
	return otherwise
}

func apiUrls() []string {
	result := []string{}
	if proxyUrl := os.Getenv(EnvProxyUrl); proxyUrl != "" {
		result = append(result, proxyUrl)
	}
	result = append(result, DefaultURL)
	return result
}
