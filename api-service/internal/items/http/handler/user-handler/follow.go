package userhandler

import (
	pb "github.com/abdulazizax/mini-twitter/api-service/genproto/user"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/middleware"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary Follow a user
// @Security     BearerAuth
// @Description Follow a user by their user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path string true "User ID to follow"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/follow/{user_id} [post]
func (h *UserHandler) FollowUserHandler(c *gin.Context) {
	h.logger.Info("FollowUserHandler called")

	userId := middleware.GetUserId(c, h.config)
	if userId == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	user_id := c.Param("user_id")

	request := pb.FollowUserRequest{
		FollowerId: userId,
		FollowedId: user_id,
	}

	body, err := protojson.Marshal(&request)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.followed", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to follow user"})
		return
	}

	h.logger.Info("User followed successfully", "follower_id", userId, "followed_id", user_id)
	c.IndentedJSON(200, gin.H{"message": "User followed successfully"})
}

// @Summary Unfollow a user
// @Security     BearerAuth
// @Description Unfollow a user by their user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path string true "User ID to unfollow"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/unfollow/{user_id} [post]
func (h *UserHandler) UnfollowUserHandler(c *gin.Context) {
	h.logger.Info("UnfollowUserHandler called")

	userId := middleware.GetUserId(c, h.config)
	if userId == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	user_id := c.Param("user_id")

	request := pb.UnfollowUserRequest{
		FollowerId: userId,
		FollowedId: user_id,
	}

	body, err := protojson.Marshal(&request)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.unfollowed", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to unfollow user"})
		return
	}

	h.logger.Info("User unfollowed successfully", "follower_id", userId, "followed_id", user_id)
	c.IndentedJSON(200, gin.H{"message": "User unfollowed successfully"})
}

// @Summary Get followers
// @Security     BearerAuth
// @Description Get the list of followers for the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} pb.GetFollowersResponse
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/followers [get]
func (h *UserHandler) GetFollowersHandler(c *gin.Context) {
	h.logger.Info("GetFollowersHandler called")

	userId := middleware.GetUserId(c, h.config)
	if userId == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	request := &pb.GetFollowersRequest{
		UserId: userId,
	}

	resp, err := h.user.GetFollowers(c.Request.Context(), request)
	if err != nil {
		h.logger.Error("Failed to get followers", "error", err, "user_id", userId)
		c.IndentedJSON(500, gin.H{"error": "Failed to get followers"})
		return
	}

	h.logger.Info("Followers retrieved successfully", "user_id", userId, "followers_count", len(resp.Followers))
	c.IndentedJSON(200, resp)
}

// @Summary Get following
// @Security     BearerAuth
// @Description Get the list of users the authenticated user is following
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} pb.GetFollowingResponse
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/following [get]
func (h *UserHandler) GetFollowingHandler(c *gin.Context) {
	h.logger.Info("GetFollowingHandler called")

	userId := middleware.GetUserId(c, h.config)
	if userId == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	request := &pb.GetFollowingRequest{
		UserId: userId,
	}

	resp, err := h.user.GetFollowing(c.Request.Context(), request)
	if err != nil {
		h.logger.Error("Failed to get following", "error", err, "user_id", userId)
		c.IndentedJSON(500, gin.H{"error": "Failed to get following"})
		return
	}

	h.logger.Info("Following retrieved successfully", "user_id", userId, "following_count", len(resp.Following))
	c.IndentedJSON(200, resp)
}
