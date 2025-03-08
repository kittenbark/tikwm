package tikwm

import (
	"context"
	"iter"
)

// PostContext corresponds to /
// url examples: "7002172928477367557", "https://www.tiktok.com/@gioscottii/video/7002172928477367557"
func PostContext(ctx context.Context, url string, hd ...bool) (*UserPost, error) {
	return DefaultAPI().Post(ctx, url, hd...)
}

// Post corresponds to /
// url examples: "7002172928477367557", "https://www.tiktok.com/@gioscottii/video/7002172928477367557"
func Post(url string, hd ...bool) (*UserPost, error) {
	return PostContext(context.Background(), url, hd...)
}

// FeedContext corresponds to user/posts
func FeedContext(ctx context.Context, uniqueID string, count int, cursor string) (*UserFeed, error) {
	return DefaultAPI().Feed(ctx, uniqueID, count, cursor)
}

// Feed corresponds to user/posts
func Feed(uniqueID string, count int, cursor string) (*UserFeed, error) {
	return DefaultAPI().Feed(context.Background(), uniqueID, count, cursor)
}

func FeedContextSeq(ctx context.Context, uniqueID string, hd ...bool) iter.Seq2[*UserPost, error] {
	return DefaultAPI().FeedSeq(ctx, uniqueID, hd...)
}

func FeedSeq(uniqueID string, hd ...bool) iter.Seq2[*UserPost, error] {
	return DefaultAPI().FeedSeq(context.Background(), uniqueID, hd...)
}

// DetailsContext corresponds to user/info
func DetailsContext(ctx context.Context, uniqueID string) (*UserDetail, error) {
	return DefaultAPI().Details(ctx, uniqueID)
}

// Details corresponds to user/info
func Details(uniqueID string) (*UserDetail, error) {
	return DefaultAPI().Details(context.Background(), uniqueID)
}
