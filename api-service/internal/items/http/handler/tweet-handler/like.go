package tweethandler

import (
	"log/slog"

	pb "github.com/abdulazizax/mini-twitter/api-service/genproto/like"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/middleware"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/msgbroker"
	"github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

type LikeHandler struct {
	like     pb.LikeServiceClient
	logger   *slog.Logger
	config   *config.Config
	producer *msgbroker.Producer
}

func NewLikeHandler(like pb.LikeServiceClient, logger *slog.Logger, config *config.Config, producer *msgbroker.Producer) *LikeHandler {
	return &LikeHandler{
		like:     like,
		logger:   logger,
		config:   config,
		producer: producer,
	}
}

// @Summary Like a tweet
// @Security     BearerAuth
// @Description Adds a like to a tweet for the authenticated user
// @Tags Likes
// @Accept json
// @Produce json
// @Param tweet_id path string true "Tweet ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/tweet/like/{tweet_id} [post]
func (h *LikeHandler) LikeTweetHandler(c *gin.Context) {
	h.logger.Info("LikeTweetHandler called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	tweetId := c.Param("tweet_id")
	if tweetId == "" {
		h.logger.Warn("Tweet ID is missing")
		c.JSON(400, gin.H{"error": "Tweet ID is missing"})
		return
	}

	req := pb.LikeRequest{
		Username:   username,
		TargetId:   tweetId,
		TargetType: "tweet",
	}

	body, err := protojson.Marshal(&req)
	if err != nil {
		h.logger.Error("Failed to marshal request", slog.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.like", body)
	if err != nil {
		h.logger.Error("Failed to send message to Kafka", slog.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("Like request sent to Kafka", slog.String("username", username), slog.String("tweet_id", tweetId))
	c.JSON(200, gin.H{"message": "Like request sent to Kafka"})
}

// @Summary Unlike a tweet
// @Security     BearerAuth
// @Description Removes a like from a tweet for the authenticated user
// @Tags Likes
// @Accept json
// @Produce json
// @Param tweet_id path string true "Tweet ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/tweet/like/{tweet_id}/unlike [post]
func (h *LikeHandler) UnLikeTweetHandler(c *gin.Context) {
	h.logger.Info("UnLikeTweetHandler called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	tweetId := c.Param("tweet_id")
	if tweetId == "" {
		h.logger.Warn("Tweet ID is missing")
		c.JSON(400, gin.H{"error": "Tweet ID is missing"})
		return
	}

	req := pb.UnlikeRequest{
		Username:   username,
		TargetId:   tweetId,
		TargetType: "tweet",
	}

	body, err := protojson.Marshal(&req)
	if err != nil {
		h.logger.Error("Failed to marshal request", slog.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.unlike", body)
	if err != nil {
		h.logger.Error("Failed to send message to Kafka", slog.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("UnLike request sent to Kafka", slog.String("username", username), slog.String("tweet_id", tweetId))
	c.JSON(200, gin.H{"message": "UnLike request sent to Kafka"})
}

// @Summary Get likes for a tweet
// @Security     BearerAuth
// @Description Retrieves the likes for a specific tweet
// @Tags Likes
// @Accept json
// @Produce json
// @Param tweet_id path string true "Tweet ID"
// @Success 200 {object} pb.GetLikesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/tweet/like/{tweet_id} [get]
func (h *LikeHandler) GetLikesTweetHandler(c *gin.Context) {
	h.logger.Info("GetLikesTweetHandler called")

	tweetId := c.Param("tweet_id")
	if tweetId == "" {
		h.logger.Warn("Tweet ID is missing")
		c.JSON(400, gin.H{"error": "Tweet ID is missing"})
		return
	}

	req := pb.GetLikesRequest{
		TargetId:   tweetId,
		TargetType: "tweet",
	}

	resp, err := h.like.GetLikes(c, &req)
	if err != nil {
		h.logger.Error("Failed to get likes", slog.String("error", err.Error()), slog.String("tweet_id", tweetId))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("Likes retrieved successfully", slog.String("tweet_id", tweetId))
	c.JSON(200, resp)
}

// @Summary Like a comment
// @Security     BearerAuth
// @Description Adds a like to a comment for the authenticated user
// @Tags Likes
// @Accept json
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/tweet/comment/like/{comment_id} [post]
func (h *LikeHandler) LikeCommentHandler(c *gin.Context) {
	h.logger.Info("LikeCommentHandler called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	commentId := c.Param("comment_id")
	if commentId == "" {
		h.logger.Warn("Comment ID is missing")
		c.JSON(400, gin.H{"error": "Comment ID is missing"})
		return
	}

	req := pb.LikeRequest{
		Username:   username,
		TargetId:   commentId,
		TargetType: "comment",
	}

	body, err := protojson.Marshal(&req)
	if err != nil {
		h.logger.Error("Failed to marshal request", slog.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.comment.like", body)
	if err != nil {
		h.logger.Error("Failed to send message to Kafka", slog.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("Like request sent to Kafka", slog.String("username", username), slog.String("comment_id", commentId))
	c.JSON(200, gin.H{"message": "Like request sent to Kafka"})
}

// @Summary Unlike a comment
// @Security     BearerAuth
// @Description Removes a like from a comment for the authenticated user
// @Tags Likes
// @Accept json
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/tweet/comment/like/{comment_id}/unlike [post]
func (h *LikeHandler) UnLikeCommentHandler(c *gin.Context) {
	h.logger.Info("UnLikeCommentHandler called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	commentId := c.Param("comment_id")
	if commentId == "" {
		h.logger.Warn("Comment ID is missing")
		c.JSON(400, gin.H{"error": "Comment ID is missing"})
		return
	}

	req := pb.UnlikeRequest{
		Username:   username,
		TargetId:   commentId,
		TargetType: "comment",
	}

	body, err := protojson.Marshal(&req)
	if err != nil {
		h.logger.Error("Failed to marshal request", slog.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.comment.unlike", body)
	if err != nil {
		h.logger.Error("Failed to send message to Kafka", slog.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("UnLike request sent to Kafka", slog.String("username", username), slog.String("comment_id", commentId))
	c.JSON(200, gin.H{"message": "UnLike request sent to Kafka"})
}

// @Summary Get likes for a comment
// @Security     BearerAuth
// @Description Retrieves the likes for a specific comment
// @Tags Likes
// @Accept json
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Success 200 {object} pb.GetLikesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/tweet/comment/like/{comment_id} [get]
func (h *LikeHandler) GetLikesCommentHandler(c *gin.Context) {
	h.logger.Info("GetLikesCommentHandler called")

	commentId := c.Param("comment_id")
	if commentId == "" {
		h.logger.Warn("Comment ID is missing")
		c.JSON(400, gin.H{"error": "Comment ID is missing"})
		return
	}

	req := pb.GetLikesRequest{
		TargetId:   commentId,
		TargetType: "comment",
	}

	resp, err := h.like.GetLikes(c, &req)
	if err != nil {
		h.logger.Error("Failed to get likes", slog.String("error", err.Error()), slog.String("comment_id", commentId))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("Likes retrieved successfully", slog.String("comment_id", commentId))
	c.JSON(200, resp)
}
