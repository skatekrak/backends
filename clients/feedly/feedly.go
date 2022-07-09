package feedly

type FeedlyClient struct {
	RefreshToken string
	AccessToken  string
}

func New(refreshToken string) *FeedlyClient {
	return &FeedlyClient{RefreshToken: refreshToken}
}

func (f *FeedlyClient) HasAccessToken() bool {
	return f.AccessToken != ""
}
