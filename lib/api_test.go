package tikwm

import (
	"context"
	"log/slog"
	"testing"
)

const (
	testVideoId  = "7002172928477367557"
	testURL      = "https://www.tiktok.com/@gioscottii/video/7002172928477367557"
	testUniqueId = "gioscottii"
)

func TestAPI(t *testing.T) {
	for _, url := range []string{testURL, testVideoId} {
		post, err := DefaultAPI().Post(context.Background(), url)
		if err != nil {
			t.Fatal(err)
		}
		if post.Id != testVideoId {
			t.Errorf("post.Id is %s, want %s", post.Id, testVideoId)
		}
	}
}

func TestFeed(t *testing.T) {
	counter := 0
	DefaultAPI().Logger = slog.Default()

	for post, err := range FeedSeq(testUniqueId) {
		if err != nil {
			t.Fatal(err)
		}

		slog.Debug("feed", "post_id", post.ID(), "urls", post.ContentUrls())
		counter++
		if counter > 70 {
			break
		}
	}
}
