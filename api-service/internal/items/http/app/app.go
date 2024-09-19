// Package api API.
//
// @title # MiniTwitter
// @version 1.03.67.83.145
//
// @description API Endpoints for MiniTwitter
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package app

import (
	"log/slog"

	_ "github.com/abdulazizax/mini-twitter/api-service/internal/items/http/app/docs"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/middleware"

	casbin "github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"

	"github.com/abdulazizax/mini-twitter/api-service/internal/items/http/handler"
	"github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Run initializes and starts the HTTP server for the MiniTwitter API.
// It sets up routing, middleware, and Swagger documentation.
//
// Parameters:
// - handler: Pointer to the Handler struct containing all route handlers
// - logger: Structured logger for logging
// - config: Application configuration
// - enforcer: Casbin enforcer for authorization
//
// Returns:
// - error: Any error that occurs during server startup
func Run(handler *handler.Handler, logger *slog.Logger, config *config.Config, enforcer *casbin.Enforcer) error {
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	// Swagger documentation setup
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url, ginSwagger.PersistAuthorization(true)))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Authentication routes
	auth := router.Group("auth")
	{
		auth.POST("/sign-up", handler.UserHandler.RegisterUserHandler)
		auth.POST("/sign-in", handler.UserHandler.LoginHandler)
		auth.POST("/sign-out", handler.UserHandler.LogoutHandler)
		auth.POST("/send-verification-email", handler.UserHandler.SendVerificationCodeHandler)
		auth.POST("/update-password", handler.UserHandler.UpdateUserPasswordHandler)
	}

	// User routes with authorization middleware
	user := router.Group("user")
	user.Use(middleware.AuthzMiddleware("/user", enforcer, config))
	{
		user.GET("/", handler.UserHandler.GetUserHandler)
		user.PUT("/", handler.UserHandler.UpdateUserHandler)
		user.POST("/", handler.UserHandler.DeleteUserHandler)
		user.POST("/follow/:user_id", handler.UserHandler.FollowUserHandler)
		user.POST("/unfollow/:user_id", handler.UserHandler.UnfollowUserHandler)
		user.GET("/followers", handler.UserHandler.GetFollowersHandler)
		user.GET("/following", handler.UserHandler.GetFollowingHandler)
		user.POST("/uploadmedia", handler.UserHandler.UploadMediaHandler)

		// Tweet-related routes
		tweet := user.Group("tweet")
		{
			tweet.POST("/", handler.TwitterHandler.TweetHandler.CreateTweetHandler)
			tweet.GET("/:tweet_serial", handler.TwitterHandler.TweetHandler.GetTweetHandler)
			tweet.PUT("/", handler.TwitterHandler.TweetHandler.UpdateTweetHandler)
			tweet.POST("/:tweet_serial/delete", handler.TwitterHandler.TweetHandler.DeleteTweetHandler)
			tweet.GET("/", handler.TwitterHandler.TweetHandler.GetAllTweets)

			tweet.POST("/views", handler.TwitterHandler.TweetHandler.IncreaseViewsCountHandler)
			tweet.POST("/reposts", handler.TwitterHandler.TweetHandler.IncreaseRepostCountHandler)
			tweet.POST("/shares", handler.TwitterHandler.TweetHandler.IncreaseSharesCountHandler)

			// Like-related routes for tweets
			like := tweet.Group("like")
			{
				like.POST("/:tweet_id", handler.TwitterHandler.LikeHandler.LikeTweetHandler)
				like.POST("/:tweet_id/unlike", handler.TwitterHandler.LikeHandler.UnLikeTweetHandler)
				like.GET("/:tweet_id", handler.TwitterHandler.LikeHandler.GetLikesTweetHandler)
			}

			// Comment-related routes
			comment := tweet.Group("comment")
			{
				comment.POST("/", handler.TwitterHandler.CommentHandler.CreateCommentHandler)
				comment.POST("/delete", handler.TwitterHandler.CommentHandler.DeleteCommentHandler)
				comment.GET("/:tweet_id", handler.TwitterHandler.CommentHandler.GetCommentsForTweetHandler)

				// Like-related routes for comments
				like := comment.Group("like")
				{
					like.POST("/:comment_id", handler.TwitterHandler.LikeHandler.LikeCommentHandler)
					like.POST("/:comment_id/unlike", handler.TwitterHandler.LikeHandler.UnLikeCommentHandler)
					like.GET("/:comment_id", handler.TwitterHandler.LikeHandler.GetLikesCommentHandler)
				}
			}
		}
	}

	// Start the server
	return router.Run(config.Server.ServerPort)
}
