# tikwm library for Go

This library streamlines working with https://tikwm.com (the best 3rd party API for TikTok, thx for the service btw).

Ensure every request is timed out properly, no 429.
```go
// Get single video:
post, err := tikwm.Post("7002172928477367557")

// Get user's feed iteratively:
for post, err := range tikwm.FeedSeq("gioscottii") {
    ...
}
```

## Env vars to run as proxy (sync local requests)
Check `cmd/proxy` for more info.
- `KITTENBARK_TIKWM_PROXY_URL` 
- `KITTENBARK_TIKWM_PROXY_PORT`
- `KITTENBARK_TIKWM_PROXY_PREFIX`