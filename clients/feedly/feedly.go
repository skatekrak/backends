package feedly

type FeedlyClient struct {
	refreshToken string
	accessToken  string
}

func New(refreshToken string) *FeedlyClient {
	return &FeedlyClient{refreshToken: refreshToken, accessToken: ""}
}
