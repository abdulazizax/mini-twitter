package userhandler

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	pb "github.com/abdulazizax/mini-twitter/api-service/genproto/user"
	"github.com/abdulazizax/mini-twitter/api-service/internal/items/middleware"
	userhandlermodles "github.com/abdulazizax/mini-twitter/api-service/internal/models/user-handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body pb.RegisterUserRequest true "User registration information"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /auth/sign-up [post]
func (h *UserHandler) RegisterUserHandler(c *gin.Context) {
	h.logger.Info("RegisterUserHandler called")

	var req pb.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	_, err := h.user.GetUserByEmail(c.Request.Context(), &pb.GetUserByEmailRequest{Email: req.Email})
	if err != nil {
		h.logger.Error("Failed to get user by email", "error", err)
		c.IndentedJSON(400, gin.H{"error": "User with this email already exists"})
		return
	}

	_, err = h.user.GetUserByUsername(c.Request.Context(), &pb.GetUserByUsernameRequest{Username: req.Username})
	if err != nil {
		h.logger.Error("Failed to get user by username", "error", err)
		c.IndentedJSON(400, gin.H{"error": "User with this username already exists"})
		return
	}

	body, err := protojson.Marshal(&req)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.registered", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	h.logger.Info("User registered successfully", "email", req.Email)
	c.IndentedJSON(201, gin.H{"message": "User registered successfully"})
}

// @Summary Get user information
// @Security     BearerAuth
// @Description Get the information of the authenticated user
// @Tags Users
// @Produce json
// @Success 200 {object} pb.User
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user [get]
func (h *UserHandler) GetUserHandler(c *gin.Context) {
	h.logger.Info("GetUserHandler called")

	userId := middleware.GetUserId(c, h.config)
	if userId == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.user.GetUser(c.Request.Context(), &pb.GetUserRequest{UserId: userId})
	if err != nil {
		h.logger.Error("Failed to get user", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to get user"})
		return
	}

	h.logger.Info("User information retrieved successfully", "userId", userId)
	c.IndentedJSON(200, user)
}

// @Summary Update user information
// @Security     BearerAuth
// @Description Update the information of the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body userhandlermodles.UpdateUserRequest true "User update information"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user [put]
func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	h.logger.Info("UpdateUserHandler called")

	userId := middleware.GetUserId(c, h.config)
	if userId == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	username := middleware.GetUsername(c, h.config)
	if username == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	userEmail := middleware.GetUserEmail(c, h.config)
	if userEmail == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req userhandlermodles.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	request := pb.UpdateUserRequest{
		UserId:            userId,
		Username:          username,
		Email:             userEmail,
		PhoneNumber:       req.PhoneNumber,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Bio:               req.Bio,
		ProfilePictureUrl: req.ProfilePictureUrl,
	}

	body, err := protojson.Marshal(&request)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.updated", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	h.logger.Info("User updated successfully", "email", userEmail)
	c.IndentedJSON(200, gin.H{"message": "User updated successfully"})
}

// @Summary Delete user
// @Security     BearerAuth
// @Description Delete the authenticated user's account
// @Tags Users
// @Produce json
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user [post]
func (h *UserHandler) DeleteUserHandler(c *gin.Context) {
	h.logger.Info("DeleteUserHandler called")

	userId := middleware.GetUserId(c, h.config)
	if userId == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	body, err := protojson.Marshal(&pb.DeleteUserRequest{UserId: userId})
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.deleted", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	h.logger.Info("User deleted successfully", "userId", userId)
	c.IndentedJSON(200, gin.H{"message": "User deleted successfully"})
}

// @Summary User login
// @Description Authenticate a user and return login information
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body pb.LoginRequest true "User login credentials"
// @Success 200 {object} pb.LoginResponse
// @Failure 400 {object} gin.H
// @Router /auth/sign-in [post]
func (h *UserHandler) LoginHandler(c *gin.Context) {
	h.logger.Info("LoginHandler called")

	var req pb.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := h.user.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to login", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid credentials"})
		return
	}

	h.logger.Info("User logged in successfully", "email", req.Email)
	c.IndentedJSON(200, user)
}

// @Summary User logout
// @Description Log out the authenticated user
// @Tags Auth
// @Produce json
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /auth/sign-out [post]
func (h *UserHandler) LogoutHandler(c *gin.Context) {
	h.logger.Info("LogoutHandler called")

	userId := middleware.GetUserId(c, h.config)
	if userId == "" {
		h.logger.Warn("User not authenticated")
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	_, err := h.user.Logout(c.Request.Context(), &pb.LogoutRequest{UserId: userId})
	if err != nil {
		h.logger.Error("Failed to logout", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to logout"})
		return
	}

	h.logger.Info("User logged out successfully", "userId", userId)
	c.IndentedJSON(200, gin.H{"message": "Logged out successfully"})
}

// @Summary Send verification code
// @Description Send a verification code to the user's email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body pb.SendVerificationCodeRequest true "Send verification code request"
// @Success 200 {object} pb.RawResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /auth/send-verification-email [post]
func (h *UserHandler) SendVerificationCodeHandler(c *gin.Context) {
	h.logger.Info("SendVerificationCodeHandler called")

	var req pb.SendVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	resp, err := h.user.SendVerificationCode(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to send verification code", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to send verification code"})
		return
	}

	h.logger.Info("Verification code sent successfully", "email", req.Email)
	c.IndentedJSON(200, resp)
}

// @Summary Update user password
// @Description Update the user's password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body pb.UpdateUserPasswordRequest true "Update user password request"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /auth/update-password [post]
func (h *UserHandler) UpdateUserPasswordHandler(c *gin.Context) {
	h.logger.Info("UpdatePasswordHandler called")

	var req pb.UpdateUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", "error", err)
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	body, err := protojson.Marshal(&req)
	if err != nil {
		h.logger.Error("Failed to marshal request", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = h.producer.Send("user.password.updated", body)
	if err != nil {
		h.logger.Error("Failed to send message to producer", "error", err)
		c.IndentedJSON(500, gin.H{"error": "Failed to update user password"})
		return
	}

	h.logger.Info("User password updated successfully", "email", req.Email)
	c.IndentedJSON(200, gin.H{"message": "User password updated successfully"})
}

// @Summary uploadFile
// @Security     BearerAuth
// @Description Upload a media file
// @Tags Users
// @Accept multipart/form-data
// @Param file formData file true "UploadMediaForm"
// @Success 201 {object} gin.H
// @Failure 400 {object} error
// @Router /user/uploadmedia [post]
func (h *UserHandler) UploadMediaHandler(c *gin.Context) {
	userId := middleware.GetUserId(c, h.config)
	if userId == "" {
		h.logger.Warn("User not authenticated", "userId", userId)
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var file userhandlermodles.File
	err := c.ShouldBind(&file)
	if err != nil {
		h.logger.Error("Failed to bind file", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	// Upload file to temporary location and send to MinIO
	fileExt := filepath.Ext(file.File.Filename)
	newFile := uuid.NewString() + fileExt

	// Tekshir bucket mavjudligini, agar yo'q bo'lsa, yarat
	exists, errExists := h.minio.BucketExists(context.Background(), "photos")
	if errExists != nil {
		h.logger.Error("Error checking if bucket exists", "error", errExists)
		c.AbortWithError(500, errExists)
		return
	}
	if !exists {
		err = h.minio.MakeBucket(context.Background(), "photos", minio.MakeBucketOptions{})
		if err != nil {
			h.logger.Error("Failed to create MinIO bucket", "error", err)
			c.AbortWithError(500, err)
			return
		}

		// Set MinIO bucket policy
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"AWS": ["*"]
					},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, "photos")

		err = h.minio.SetBucketPolicy(c.Request.Context(), "photos", policy)
		if err != nil {
			h.logger.Error("Failed to set MinIO bucket policy", "error", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		h.logger.Info("Bucket already exists", "bucket", "photos")
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		h.logger.Error("Failed to create temporary file", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer os.Remove(tempFile.Name()) // Remove file later

	err = c.SaveUploadedFile(&file.File, tempFile.Name())
	if err != nil {
		h.logger.Error("Failed to save uploaded file", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Faylning Content-Type ni dinamik aniqlash
	fileType := mime.TypeByExtension(fileExt)
	if fileType == "" {
		fileType = "application/octet-stream" // Agar Content-Type aniqlanmasa, default
	}

	// Upload to MinIO
	info, err := h.minio.FPutObject(c.Request.Context(), "photos", newFile, tempFile.Name(), minio.PutObjectOptions{
		ContentType: fileType, // Fayl turi dinamik
	})
	if err != nil {
		h.logger.Error("Failed to upload file to MinIO", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	h.logger.Info("File uploaded successfully", "bucket", info.Bucket, "file", newFile)

	// Create file URL
	madeURL := fmt.Sprintf("http://localhost:9000/photos/%s", newFile)

	h.logger.Info("File URL created", "url", madeURL)
	c.JSON(201, gin.H{"url": madeURL})
}
