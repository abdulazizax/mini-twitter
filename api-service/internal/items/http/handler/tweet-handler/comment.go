package tweethandler

import (
	"log/slog"

	pb "github.com/abdulazizax/mini-twitter/api-service/genproto/comment"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/middleware"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/msgbroker"
	tweethandlermodels "github.com/abdulazizax/mini-twitter/api-service/internal/models/tweet-handler"
	"github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

type CommentHandler struct {
	comment  pb.CommentServiceClient
	logger   *slog.Logger
	config   *config.Config
	producer *msgbroker.Producer
}

func NewCommentHandler(comment pb.CommentServiceClient, logger *slog.Logger, config *config.Config, producer *msgbroker.Producer) *CommentHandler {
	return &CommentHandler{
		comment:  comment,
		logger:   logger,
		config:   config,
		producer: producer,
	}
}

// @Summary Create a new comment
// @Security     BearerAuth
// @Description Create a new comment for a tweet
// @Tags Comments
// @Accept json
// @Produce json
// @Param request body tweethandlermodels.CreateCommentRequest true "Comment creation request"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet/comment [post]
func (h *CommentHandler) CreateCommentHandler(c *gin.Context) {
	h.logger.Info("CreateCommentHandler called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req tweethandlermodels.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", slog.String("error", err.Error()))
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	request := pb.CreateCommentRequest{
		Username: username,
		TweetId:  req.TweetId,
		Content:  req.Content,
	}

	body, err := protojson.Marshal(&request)
	if err != nil {
		h.logger.Error("Failed to marshal request", slog.String("error", err.Error()))
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.comment.created", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", slog.String("error", err.Error()))
		c.IndentedJSON(500, gin.H{"error": "Failed to create comment"})
		return
	}

	h.logger.Info("Comment created successfully", slog.String("username", username), slog.String("tweet_id", req.TweetId))
	c.IndentedJSON(201, gin.H{"message": "Comment created successfully"})
}

// @Summary Delete a comment
// @Security     BearerAuth
// @Description Delete a comment from a tweet
// @Tags Comments
// @Accept json
// @Produce json
// @Param request body tweethandlermodels.DeleteCommentRequest true "Comment deletion request"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet/comment/delete [post]
func (h *CommentHandler) DeleteCommentHandler(c *gin.Context) {
	h.logger.Info("DeleteCommentHandler called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req tweethandlermodels.DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", slog.String("error", err.Error()))
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	request := pb.DeleteCommentRequest{
		Username:  username,
		TweetId:   req.TweetId,
		CommentId: req.CommentId,
	}

	body, err := protojson.Marshal(&request)
	if err != nil {
		h.logger.Error("Failed to marshal request", slog.String("error", err.Error()))
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.comment.deleted", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", slog.String("error", err.Error()))
		c.IndentedJSON(500, gin.H{"error": "Failed to delete comment"})
		return
	}

	h.logger.Info("Comment deleted successfully", slog.String("username", username), slog.String("tweet_id", req.TweetId), slog.String("comment_id", req.CommentId))
	c.IndentedJSON(200, gin.H{"message": "Comment deleted successfully"})
}

// @Summary Get comments for a tweet
// @Security     BearerAuth
// @Description Retrieve all comments for a specific tweet
// @Tags Comments
// @Produce json
// @Param tweet_id path string true "Tweet ID"
// @Success 200 {object} pb.GetCommentsForTweetResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet/comment/{tweet_id} [get]
func (h *CommentHandler) GetCommentsForTweetHandler(c *gin.Context) {
	h.logger.Info("GetCommentsForTweetHandler called")

	tweetID := c.Param("tweet_id")
	if tweetID == "" {
		h.logger.Warn("Tweet ID is missing")
		c.IndentedJSON(400, gin.H{"error": "Tweet ID is required"})
		return
	}

	request := &pb.GetCommentsForTweetRequest{
		TweetId: tweetID,
	}

	response, err := h.comment.GetCommentsForTweet(c.Request.Context(), request)
	if err != nil {
		h.logger.Error("Failed to get comments for tweet", slog.String("error", err.Error()), slog.String("tweet_id", tweetID))
		c.IndentedJSON(500, gin.H{"error": "Failed to get comments for tweet"})
		return
	}

	h.logger.Info("Comments retrieved successfully", slog.String("tweet_id", tweetID), slog.Int("comment_count", len(response.Comments)))
	c.IndentedJSON(200, response)
}
