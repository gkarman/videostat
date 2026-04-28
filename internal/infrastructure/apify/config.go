package apify

type Config struct {
	Token  string
	Host   string
	Limits Limits
}

type Limits struct {
	YouTubeLimits   YouTubeLimits
	TikTokLimits    TikTokLimits
	InstagramLimits InstagramLimits
}

type YouTubeLimits struct {
	MaxVideos int
	Days      int
}

type TikTokLimits struct {
	MaxVideos int
	Days      int
}

type InstagramLimits struct {
	MaxVideos int
	Days      int
}
