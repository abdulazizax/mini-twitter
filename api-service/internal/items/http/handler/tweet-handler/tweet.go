package tweethandler

import (
	"log/slog"
	"strconv"

	pb "github.com/abdulazizax/mini-twitter/api-service/genproto/tweet"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/middleware"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/msgbroker"
	tweethandlermodels "github.com/abdulazizax/mini-twitter/api-service/internal/models/tweet-handler"
	"github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

type TweetHandler struct {
	tweet    pb.TweetServiceClient
	logger   *slog.Logger
	config   *config.Config
	producer *msgbroker.Producer
}

func NewTweetHandler(tweet pb.TweetServiceClient, logger *slog.Logger, config *config.Config, producer *msgbroker.Producer) *TweetHandler {
	return &TweetHandler{
		tweet:    tweet,
		logger:   logger,
		config:   config,
		producer: producer,
	}
}

// @Summary Create a new tweet
// @Security     BearerAuth
// @Description Create a new tweet for the authenticated user
// @Tags Tweets
// @Accept json
// @Produce json
// @Param request body tweethandlermodels.CreateTweetRequest true "Tweet details"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet [post]
func (h *TweetHandler) CreateTweetHandler(c *gin.Context) {
	h.logger.Info("CreateTweetHandler called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req tweethandlermodels.CreateTweetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request"})
		return
	}

	request := &pb.CreateTweetRequest{
		Username: username,
		Content:  req.Content,
		Media:    req.Media,
	}

	body, err := protojson.Marshal(request)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.created", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to create tweet"})
		return
	}

	h.logger.Info("Tweet created successfully", "username", username)
	c.IndentedJSON(201, gin.H{"message": "Tweet created successfully"})
}

// @Summary Get a specific tweet
// @Security     BearerAuth
// @Description Get a specific tweet by username and tweet serial
// @Tags Tweets
// @Produce json
// @Param tweet_serial path int true "Tweet serial number"
// @Success 200 {object} pb.GetTweetResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet/{tweet_serial} [get]
func (h *TweetHandler) GetTweetHandler(c *gin.Context) {
	h.logger.Info("GetTweetHandler called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	tweet_serial := c.Param("tweet_serial")
	if tweet_serial == "" {
		h.logger.Warn("Tweet serial is empty")
		c.IndentedJSON(400, gin.H{"error": "Invalid tweet serial"})
		return
	}

	int64TweetSerial, err := strconv.Atoi(tweet_serial)
	if err != nil {
		h.logger.Error("Failed to convert tweet_serial to int", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid tweet_serial"})
		return
	}

	request := &pb.GetTweetRequest{
		Username:    username,
		TweetSerial: int32(int64TweetSerial),
	}

	resp, err := h.tweet.GetTweet(c.Request.Context(), request)
	if err != nil {
		h.logger.Error("Failed to get tweet", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to get tweet"})
		return
	}

	h.logger.Info("Tweet retrieved successfully", "username", username)
	c.IndentedJSON(200, resp)
}

// @Summary Update a tweet
// @Security     BearerAuth
// @Description Update a tweet for the authenticated user
// @Tags Tweets
// @Accept json
// @Produce json
// @Param request body tweethandlermodels.UpdateTweetRequest true "Updated tweet details"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet [put]
func (h *TweetHandler) UpdateTweetHandler(c *gin.Context) {
	h.logger.Info("UpdateTweetHandler called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req tweethandlermodels.UpdateTweetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request"})
		return
	}

	request := &pb.UpdateTweetRequest{
		Username:    username,
		TweetSerial: req.TweetSerial,
		Content:     req.Content,
		Media:       req.Media,
	}

	body, err := protojson.Marshal(request)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.updated", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to update tweet"})
		return
	}

	h.logger.Info("Tweet updated successfully", "username", username)
	c.IndentedJSON(200, gin.H{"message": "Tweet updated successfully"})
}

// @Summary Delete a tweet
// @Security     BearerAuth
// @Description Delete a tweet for the authenticated user
// @Tags Tweets
// @Produce json
// @Param tweet_serial path int true "Tweet serial number"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet/{tweet_serial}/delete [post]
func (h *TweetHandler) DeleteTweetHandler(c *gin.Context) {
	h.logger.Info("DeleteTweetHandler called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	tweet_serial := c.Param("tweet_serial")
	if tweet_serial == "" {
		h.logger.Warn("Tweet serial is empty")
		c.IndentedJSON(400, gin.H{"error": "Invalid tweet serial"})
		return
	}

	int64TweetSerial, err := strconv.Atoi(tweet_serial)
	if err != nil {
		h.logger.Error("Failed to convert tweet_serial to int", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid tweet_serial"})
		return
	}

	request := &pb.DeleteTweetRequest{
		Username:    username,
		TweetSerial: int32(int64TweetSerial),
	}

	body, err := protojson.Marshal(request)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.deleted", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to delete tweet"})
		return
	}

	h.logger.Info("Tweet deleted successfully", "username", username)
	c.IndentedJSON(200, gin.H{"message": "Tweet deleted successfully"})
}

// @Summary Get all tweets for a user
// @Security     BearerAuth
// @Description Get all tweets for the authenticated user
// @Tags Tweets
// @Produce json
// @Success 200 {object} pb.GetAllTweetsResponse
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet [get]
func (h *TweetHandler) GetAllTweets(c *gin.Context) {
	h.logger.Info("GetAllTweets called")

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	request := &pb.GetAllTweetsRequest{
		Username: username,
	}

	resp, err := h.tweet.GetAllTweets(c.Request.Context(), request)
	if err != nil {
		h.logger.Error("Failed to get all tweets", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to get all tweets"})
		return
	}

	h.logger.Info("All tweets retrieved successfully", "username", username)
	c.IndentedJSON(200, resp)
}

// @Summary Increase views count for a tweet
// @Security BearerAuth
// @Description Increase the views count for a specific tweet
// @Tags Tweets
// @Accept json
// @Produce json
// @Param id body pb.Id true "Tweet ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet/views [post]
func (h *TweetHandler) IncreaseViewsCountHandler(c *gin.Context) {
	h.logger.Info("IncreaseViewsCountHandler called")

	var request pb.Id
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request"})
		return
	}

	body, err := protojson.Marshal(&request)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.views.increased", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to increase views count"})
		return
	}

	h.logger.Info("Views count increased successfully", "username", request.Username)
	c.IndentedJSON(200, gin.H{"message": "Views count increased successfully"})
}

// @Summary Increase repost count for a tweet
// @Security BearerAuth
// @Description Increase the repost count for a specific tweet
// @Tags Tweets
// @Accept json
// @Produce json
// @Param id body pb.Id true "Tweet ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet/reposts [post]
func (h *TweetHandler) IncreaseRepostCountHandler(c *gin.Context) {
	h.logger.Info("IncreaseRepostCountHandler called")

	var request pb.Id
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request"})
		return
	}

	body, err := protojson.Marshal(&request)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.repost.increased", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to increase repost count"})
		return
	}

	h.logger.Info("Repost count increased successfully", "username", request.Username)
	c.IndentedJSON(200, gin.H{"message": "Repost count increased successfully"})
}

// @Summary Increase shares count for a tweet
// @Security BearerAuth
// @Description Increase the shares count for a specific tweet
// @Tags Tweets
// @Accept json
// @Produce json
// @Param id body pb.Id true "Tweet ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/tweet/shares [post]
func (h *TweetHandler) IncreaseSharesCountHandler(c *gin.Context) {
	h.logger.Info("IncreaseSharesCountHandler called")

	var request pb.Id
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request"})
		return
	}

	body, err := protojson.Marshal(&request)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.tweet.shares.increased", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to increase shares count"})
		return
	}

	h.logger.Info("Shares count increased successfully", "username", request.Username)
	c.IndentedJSON(200, gin.H{"message": "Shares count increased successfully"})
}
