package main

import (
	"encoding/json"
	"github.com/kittenbark/tikwm/lib"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var kGatewayPrefix = LookupOr("KITTENBARK_TIKWM_PROXY_PREFIX", "/v1/tikwm")
var kGatewayPort = LookupOr("KITTENBARK_TIKWM_PROXY_PORT", ":42011")

func main() {
	tikwm.DefaultAPI().Urls = []string{tikwm.DefaultURL}
	mux := http.NewServeMux()

	mux.HandleFunc(kGatewayPrefix+"/", HandlePost)
	mux.HandleFunc(kGatewayPrefix+"/user/posts", HandleUserPosts)
	mux.HandleFunc(kGatewayPrefix+"/user/info", HandleUserInfo)

	slog.Info("tikwm#proxy_start", "port", kGatewayPort, "api_prefix", kGatewayPrefix)
	if err := http.ListenAndServe(kGatewayPort, mux); err != nil {
		panic(err)
	}
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	post, err := tikwm.Post(query.Get("url"), query.Get("hd") == "1")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(NewError(http.StatusBadRequest, err))
		return
	}

	data, _ := json.Marshal(post)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func HandleUserPosts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	feed, err := tikwm.Feed(
		query.Get("unique_id"),
		QueryInt(query, "count", tikwm.DefaultUserFeedSize),
		query.Get("cursor"),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(NewError(http.StatusBadRequest, err))
		return
	}

	data, _ := json.Marshal(feed)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	details, err := tikwm.Details(query.Get("unique_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(NewError(http.StatusBadRequest, err))
		return
	}

	data, _ := json.Marshal(details)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func NewError(status int, err error) []byte {
	type Error struct {
		Status  int    `json:"code"`
		Message string `json:"msg"`
	}
	data, _ := json.Marshal(&Error{status, err.Error()})
	return data
}

func QueryInt(query url.Values, key string, otherwise int) int {
	val := query.Get(key)
	result, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return otherwise
	}
	return int(result)
}

func LookupOr(env string, otherwise string) string {
	result, ok := os.LookupEnv(env)
	if !ok {
		return otherwise
	}
	return result
}
