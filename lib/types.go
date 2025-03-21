package tikwm

type UserPost struct {
	Id          string  `json:"id"`
	VideoId     string  `json:"video_id"`
	Region      string  `json:"region"`
	Title       string  `json:"title"`
	Cover       string  `json:"cover"`
	OriginCover string  `json:"origin_cover"`
	Duration    float64 `json:"duration"`
	Play        string  `json:"play"`
	Wmplay      string  `json:"wmplay"`
	Hdplay      string  `json:"hdplay"`
	Size        int64   `json:"size"`
	WmSize      int64   `json:"wm_size"`
	HdSize      int64   `json:"hd_size"`
	Music       string  `json:"music"`
	MusicInfo   struct {
		Id       string `json:"id"`
		Title    string `json:"title"`
		Play     string `json:"play"`
		Cover    string `json:"cover"`
		Author   string `json:"author"`
		Original bool   `json:"original"`
		// Duration for some videos is string.
		Duration interface{} `json:"duration"`
		Album    string      `json:"album"`
	} `json:"music_info"`
	PlayCount     int         `json:"play_count"`
	DiggCount     int         `json:"digg_count"`
	CommentCount  int         `json:"comment_count"`
	ShareCount    int         `json:"share_count"`
	DownloadCount int         `json:"download_count"`
	CollectCount  int         `json:"collect_count"`
	CreateTime    int64       `json:"create_time"`
	Anchors       interface{} `json:"anchors"`
	AnchorsExtras string      `json:"anchors_extras"`
	IsAd          bool        `json:"is_ad"`
	CommerceInfo  struct {
		AdvPromotable          bool `json:"adv_promotable"`
		AuctionAdInvited       bool `json:"auction_ad_invited"`
		BrandedContentType     int  `json:"branded_content_type"`
		WithCommentFilterWords bool `json:"with_comment_filter_words"`
	} `json:"commerce_info"`
	CommercialVideoInfo string `json:"commercial_video_info"`
	Author              struct {
		Id       string `json:"id"`
		UniqueId string `json:"unique_id"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	} `json:"author"`
	Images []string `json:"images"`
}

func (post *UserPost) IsAlbum() bool {
	return len(post.Images) != 0
}

func (post *UserPost) IsVideo() bool {
	return !post.IsAlbum()
}

func (post *UserPost) ContentUrls(hd ...bool) []string {
	HD := at(hd, 0, true)

	urls := post.Images
	if post.IsVideo() {
		if post.Hdplay != "" && HD {
			urls = []string{post.Hdplay}
		} else if post.Play != "" {
			urls = []string{post.Play}
		} else {
			urls = []string{post.Wmplay}
		}
	}
	return urls
}

// ID is the simplest way to get video's id
func (post *UserPost) ID() string {
	if post.Id != "" {
		return post.Id
	}
	return post.VideoId
}

type UserFeed struct {
	Posts   []*UserPost `json:"videos"`
	Cursor  string      `json:"cursor"`
	HasMore bool        `json:"hasMore"`
}

type UserDetail struct {
	User struct {
		Id                  string      `json:"id"`
		UniqueId            string      `json:"uniqueId"`
		Nickname            string      `json:"nickname"`
		AvatarThumb         string      `json:"avatarThumb"`
		AvatarMedium        string      `json:"avatarMedium"`
		AvatarLarger        string      `json:"avatarLarger"`
		Signature           string      `json:"signature"`
		Verified            bool        `json:"verified"`
		SecUid              string      `json:"secUid"`
		Secret              bool        `json:"secret"`
		Ftc                 bool        `json:"ftc"`
		Relation            int         `json:"relation"`
		OpenFavorite        bool        `json:"openFavorite"`
		CommentSetting      interface{} `json:"commentSetting"`
		DuetSetting         interface{} `json:"duetSetting"`
		StitchSetting       interface{} `json:"stitchSetting"`
		PrivateAccount      bool        `json:"privateAccount"`
		IsADVirtual         bool        `json:"isADVirtual"`
		IsUnderAge18        bool        `json:"isUnderAge18"`
		InsId               string      `json:"ins_id"`
		TwitterId           string      `json:"twitter_id"`
		YoutubeChannelTitle string      `json:"youtube_channel_title"`
		YoutubeChannelId    string      `json:"youtube_channel_id"`
	} `json:"user"`
	Stats struct {
		FollowingCount int `json:"followingCount"`
		FollowerCount  int `json:"followerCount"`
		HeartCount     int `json:"heartCount"`
		VideoCount     int `json:"videoCount"`
		DiggCount      int `json:"diggCount"`
		Heart          int `json:"heart"`
	} `json:"stats"`
}
