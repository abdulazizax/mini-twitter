package userhandlermodles

import "mime/multipart"

type UpdateUserRequest struct {
	PhoneNumber       string `json:"phone_number"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Bio               string `json:"bio"`
	ProfilePictureUrl string `json:"profile_picture_url"`
}

type File struct {
	File multipart.FileHeader `form:"file" binding:"required"`
}

type UploadMediaRequest struct {
	Media  string `json:"media"`
	UserId string `json:"user_id"`
}
