package tweethandlermodels

type CreateTweetRequest struct {
	Content string   `json:"content"`
	Media   []string `json:"media"`
}

type UpdateTweetRequest struct {
	TweetSerial int32    `json:"tweet_serial"`
	Content     string   `json:"content"`
	Media       []string `json:"media"`
}

