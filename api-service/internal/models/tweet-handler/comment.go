package tweethandlermodels

type CreateCommentRequest struct {
	TweetId string `json:"tweet_id"`
	Content string `json:"content"`
}


type DeleteCommentRequest struct{
	TweetId string `json:"tweet_id"`
	CommentId string `json:"comment_id"`
}