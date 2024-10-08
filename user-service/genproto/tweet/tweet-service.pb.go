// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: tweet-service/tweet-service.proto

package tweet

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Tweet model
type Tweet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	TweetSerial   int32                  `protobuf:"varint,3,opt,name=tweet_serial,json=tweetSerial,proto3" json:"tweet_serial,omitempty"`
	Content       string                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Media         []string               `protobuf:"bytes,5,rep,name=media,proto3" json:"media,omitempty"`
	CommentsCount int32                  `protobuf:"varint,7,opt,name=comments_count,json=commentsCount,proto3" json:"comments_count,omitempty"`
	ViewsCount    int32                  `protobuf:"varint,8,opt,name=views_count,json=viewsCount,proto3" json:"views_count,omitempty"`
	RepostCount   int32                  `protobuf:"varint,9,opt,name=repost_count,json=repostCount,proto3" json:"repost_count,omitempty"`
	SharesCount   int32                  `protobuf:"varint,10,opt,name=shares_count,json=sharesCount,proto3" json:"shares_count,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Tweet) Reset() {
	*x = Tweet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tweet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tweet) ProtoMessage() {}

func (x *Tweet) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tweet.ProtoReflect.Descriptor instead.
func (*Tweet) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{0}
}

func (x *Tweet) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Tweet) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Tweet) GetTweetSerial() int32 {
	if x != nil {
		return x.TweetSerial
	}
	return 0
}

func (x *Tweet) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Tweet) GetMedia() []string {
	if x != nil {
		return x.Media
	}
	return nil
}

func (x *Tweet) GetCommentsCount() int32 {
	if x != nil {
		return x.CommentsCount
	}
	return 0
}

func (x *Tweet) GetViewsCount() int32 {
	if x != nil {
		return x.ViewsCount
	}
	return 0
}

func (x *Tweet) GetRepostCount() int32 {
	if x != nil {
		return x.RepostCount
	}
	return 0
}

func (x *Tweet) GetSharesCount() int32 {
	if x != nil {
		return x.SharesCount
	}
	return 0
}

func (x *Tweet) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Tweet) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// CreateTweet
type CreateTweetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Content  string   `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	Media    []string `protobuf:"bytes,3,rep,name=media,proto3" json:"media,omitempty"`
}

func (x *CreateTweetRequest) Reset() {
	*x = CreateTweetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTweetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTweetRequest) ProtoMessage() {}

func (x *CreateTweetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTweetRequest.ProtoReflect.Descriptor instead.
func (*CreateTweetRequest) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateTweetRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *CreateTweetRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *CreateTweetRequest) GetMedia() []string {
	if x != nil {
		return x.Media
	}
	return nil
}

type CreateTweetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tweet *Tweet `protobuf:"bytes,1,opt,name=tweet,proto3" json:"tweet,omitempty"`
}

func (x *CreateTweetResponse) Reset() {
	*x = CreateTweetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTweetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTweetResponse) ProtoMessage() {}

func (x *CreateTweetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTweetResponse.ProtoReflect.Descriptor instead.
func (*CreateTweetResponse) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{2}
}

func (x *CreateTweetResponse) GetTweet() *Tweet {
	if x != nil {
		return x.Tweet
	}
	return nil
}

// GetTweet
type GetTweetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username    string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	TweetSerial int32  `protobuf:"varint,2,opt,name=tweet_serial,json=tweetSerial,proto3" json:"tweet_serial,omitempty"`
}

func (x *GetTweetRequest) Reset() {
	*x = GetTweetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTweetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTweetRequest) ProtoMessage() {}

func (x *GetTweetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTweetRequest.ProtoReflect.Descriptor instead.
func (*GetTweetRequest) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetTweetRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *GetTweetRequest) GetTweetSerial() int32 {
	if x != nil {
		return x.TweetSerial
	}
	return 0
}

type GetTweetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tweet *Tweet `protobuf:"bytes,1,opt,name=tweet,proto3" json:"tweet,omitempty"`
}

func (x *GetTweetResponse) Reset() {
	*x = GetTweetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTweetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTweetResponse) ProtoMessage() {}

func (x *GetTweetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTweetResponse.ProtoReflect.Descriptor instead.
func (*GetTweetResponse) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetTweetResponse) GetTweet() *Tweet {
	if x != nil {
		return x.Tweet
	}
	return nil
}

// UpdateTweet
type UpdateTweetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username    string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	TweetSerial int32    `protobuf:"varint,2,opt,name=tweet_serial,json=tweetSerial,proto3" json:"tweet_serial,omitempty"`
	Content     string   `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	Media       []string `protobuf:"bytes,4,rep,name=media,proto3" json:"media,omitempty"`
}

func (x *UpdateTweetRequest) Reset() {
	*x = UpdateTweetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTweetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTweetRequest) ProtoMessage() {}

func (x *UpdateTweetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTweetRequest.ProtoReflect.Descriptor instead.
func (*UpdateTweetRequest) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateTweetRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UpdateTweetRequest) GetTweetSerial() int32 {
	if x != nil {
		return x.TweetSerial
	}
	return 0
}

func (x *UpdateTweetRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *UpdateTweetRequest) GetMedia() []string {
	if x != nil {
		return x.Media
	}
	return nil
}

type UpdateTweetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tweet *Tweet `protobuf:"bytes,1,opt,name=tweet,proto3" json:"tweet,omitempty"`
}

func (x *UpdateTweetResponse) Reset() {
	*x = UpdateTweetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTweetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTweetResponse) ProtoMessage() {}

func (x *UpdateTweetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTweetResponse.ProtoReflect.Descriptor instead.
func (*UpdateTweetResponse) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateTweetResponse) GetTweet() *Tweet {
	if x != nil {
		return x.Tweet
	}
	return nil
}

// DeleteTweet
type DeleteTweetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username    string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	TweetSerial int32  `protobuf:"varint,2,opt,name=tweet_serial,json=tweetSerial,proto3" json:"tweet_serial,omitempty"`
}

func (x *DeleteTweetRequest) Reset() {
	*x = DeleteTweetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTweetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTweetRequest) ProtoMessage() {}

func (x *DeleteTweetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTweetRequest.ProtoReflect.Descriptor instead.
func (*DeleteTweetRequest) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteTweetRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *DeleteTweetRequest) GetTweetSerial() int32 {
	if x != nil {
		return x.TweetSerial
	}
	return 0
}

// GetAllTweets
type GetAllTweetsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *GetAllTweetsRequest) Reset() {
	*x = GetAllTweetsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllTweetsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllTweetsRequest) ProtoMessage() {}

func (x *GetAllTweetsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllTweetsRequest.ProtoReflect.Descriptor instead.
func (*GetAllTweetsRequest) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{8}
}

func (x *GetAllTweetsRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type GetAllTweetsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tweets []*Tweet `protobuf:"bytes,1,rep,name=tweets,proto3" json:"tweets,omitempty"`
}

func (x *GetAllTweetsResponse) Reset() {
	*x = GetAllTweetsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllTweetsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllTweetsResponse) ProtoMessage() {}

func (x *GetAllTweetsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllTweetsResponse.ProtoReflect.Descriptor instead.
func (*GetAllTweetsResponse) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{9}
}

func (x *GetAllTweetsResponse) GetTweets() []*Tweet {
	if x != nil {
		return x.Tweets
	}
	return nil
}

// Status message
type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{10}
}

func (x *Status) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// General purpose ID message for increasing counts
type Id struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username    string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	TweetSerial int32  `protobuf:"varint,2,opt,name=tweet_serial,json=tweetSerial,proto3" json:"tweet_serial,omitempty"`
}

func (x *Id) Reset() {
	*x = Id{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tweet_service_tweet_service_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Id) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Id) ProtoMessage() {}

func (x *Id) ProtoReflect() protoreflect.Message {
	mi := &file_tweet_service_tweet_service_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Id.ProtoReflect.Descriptor instead.
func (*Id) Descriptor() ([]byte, []int) {
	return file_tweet_service_tweet_service_proto_rawDescGZIP(), []int{11}
}

func (x *Id) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Id) GetTweetSerial() int32 {
	if x != nil {
		return x.TweetSerial
	}
	return 0
}

var File_tweet_service_tweet_service_proto protoreflect.FileDescriptor

var file_tweet_service_tweet_service_proto_rawDesc = []byte{
	0x0a, 0x21, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x74, 0x77, 0x65, 0x65, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x74, 0x77, 0x65, 0x65, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8a, 0x03, 0x0a, 0x05,
	0x54, 0x77, 0x65, 0x65, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x77, 0x65, 0x65, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x74, 0x77, 0x65, 0x65, 0x74, 0x53, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6d,
	0x65, 0x64, 0x69, 0x61, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x63, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x76,
	0x69, 0x65, 0x77, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x76, 0x69, 0x65, 0x77, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c,
	0x72, 0x65, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0b, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x21, 0x0a, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x73, 0x68, 0x61, 0x72, 0x65, 0x73, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a,
	0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x60, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x22, 0x39, 0x0a, 0x13, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x22, 0x0a, 0x05, 0x74, 0x77, 0x65, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52, 0x05,
	0x74, 0x77, 0x65, 0x65, 0x74, 0x22, 0x50, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x54, 0x77, 0x65, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x77, 0x65, 0x65, 0x74, 0x5f, 0x73, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x74, 0x77, 0x65, 0x65,
	0x74, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x22, 0x36, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x54, 0x77,
	0x65, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x05, 0x74,
	0x77, 0x65, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x77, 0x65,
	0x65, 0x74, 0x2e, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52, 0x05, 0x74, 0x77, 0x65, 0x65, 0x74, 0x22,
	0x83, 0x01, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x77, 0x65, 0x65, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x74, 0x77, 0x65, 0x65, 0x74, 0x53,
	0x65, 0x72, 0x69, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x22, 0x39, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54,
	0x77, 0x65, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x05,
	0x74, 0x77, 0x65, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x77,
	0x65, 0x65, 0x74, 0x2e, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52, 0x05, 0x74, 0x77, 0x65, 0x65, 0x74,
	0x22, 0x53, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x77, 0x65, 0x65, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x74, 0x77, 0x65, 0x65, 0x74, 0x53,
	0x65, 0x72, 0x69, 0x61, 0x6c, 0x22, 0x31, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x54,
	0x77, 0x65, 0x65, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3c, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x41,
	0x6c, 0x6c, 0x54, 0x77, 0x65, 0x65, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x24, 0x0a, 0x06, 0x74, 0x77, 0x65, 0x65, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52, 0x06,
	0x74, 0x77, 0x65, 0x65, 0x74, 0x73, 0x22, 0x22, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x43, 0x0a, 0x02, 0x49, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x74, 0x77, 0x65, 0x65, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0b, 0x74, 0x77, 0x65, 0x65, 0x74, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x32,
	0xeb, 0x03, 0x0a, 0x0c, 0x54, 0x77, 0x65, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x44, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x77, 0x65, 0x65, 0x74, 0x12,
	0x19, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x77,
	0x65, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x74, 0x77, 0x65,
	0x65, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x54, 0x77, 0x65,
	0x65, 0x74, 0x12, 0x16, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x77,
	0x65, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x74, 0x77, 0x65,
	0x65, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x77, 0x65,
	0x65, 0x74, 0x12, 0x19, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e,
	0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x77, 0x65, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x0b, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x54, 0x77, 0x65, 0x65, 0x74, 0x12, 0x19, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x77, 0x65, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x47, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x54, 0x77, 0x65, 0x65,
	0x74, 0x73, 0x12, 0x1a, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c,
	0x6c, 0x54, 0x77, 0x65, 0x65, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b,
	0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x54, 0x77, 0x65,
	0x65, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x12, 0x49,
	0x6e, 0x63, 0x72, 0x65, 0x61, 0x73, 0x65, 0x56, 0x69, 0x65, 0x77, 0x73, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x09, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x49, 0x64, 0x1a, 0x0d, 0x2e, 0x74,
	0x77, 0x65, 0x65, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x13, 0x49,
	0x6e, 0x63, 0x72, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x09, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x49, 0x64, 0x1a, 0x0d, 0x2e,
	0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x13,
	0x49, 0x6e, 0x63, 0x72, 0x65, 0x61, 0x73, 0x65, 0x53, 0x68, 0x61, 0x72, 0x65, 0x73, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x09, 0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x49, 0x64, 0x1a, 0x0d,
	0x2e, 0x74, 0x77, 0x65, 0x65, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x10, 0x5a,
	0x0e, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x77, 0x65, 0x65, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tweet_service_tweet_service_proto_rawDescOnce sync.Once
	file_tweet_service_tweet_service_proto_rawDescData = file_tweet_service_tweet_service_proto_rawDesc
)

func file_tweet_service_tweet_service_proto_rawDescGZIP() []byte {
	file_tweet_service_tweet_service_proto_rawDescOnce.Do(func() {
		file_tweet_service_tweet_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_tweet_service_tweet_service_proto_rawDescData)
	})
	return file_tweet_service_tweet_service_proto_rawDescData
}

var file_tweet_service_tweet_service_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_tweet_service_tweet_service_proto_goTypes = []any{
	(*Tweet)(nil),                 // 0: tweet.Tweet
	(*CreateTweetRequest)(nil),    // 1: tweet.CreateTweetRequest
	(*CreateTweetResponse)(nil),   // 2: tweet.CreateTweetResponse
	(*GetTweetRequest)(nil),       // 3: tweet.GetTweetRequest
	(*GetTweetResponse)(nil),      // 4: tweet.GetTweetResponse
	(*UpdateTweetRequest)(nil),    // 5: tweet.UpdateTweetRequest
	(*UpdateTweetResponse)(nil),   // 6: tweet.UpdateTweetResponse
	(*DeleteTweetRequest)(nil),    // 7: tweet.DeleteTweetRequest
	(*GetAllTweetsRequest)(nil),   // 8: tweet.GetAllTweetsRequest
	(*GetAllTweetsResponse)(nil),  // 9: tweet.GetAllTweetsResponse
	(*Status)(nil),                // 10: tweet.Status
	(*Id)(nil),                    // 11: tweet.Id
	(*timestamppb.Timestamp)(nil), // 12: google.protobuf.Timestamp
}
var file_tweet_service_tweet_service_proto_depIdxs = []int32{
	12, // 0: tweet.Tweet.created_at:type_name -> google.protobuf.Timestamp
	12, // 1: tweet.Tweet.updated_at:type_name -> google.protobuf.Timestamp
	0,  // 2: tweet.CreateTweetResponse.tweet:type_name -> tweet.Tweet
	0,  // 3: tweet.GetTweetResponse.tweet:type_name -> tweet.Tweet
	0,  // 4: tweet.UpdateTweetResponse.tweet:type_name -> tweet.Tweet
	0,  // 5: tweet.GetAllTweetsResponse.tweets:type_name -> tweet.Tweet
	1,  // 6: tweet.TweetService.CreateTweet:input_type -> tweet.CreateTweetRequest
	3,  // 7: tweet.TweetService.GetTweet:input_type -> tweet.GetTweetRequest
	5,  // 8: tweet.TweetService.UpdateTweet:input_type -> tweet.UpdateTweetRequest
	7,  // 9: tweet.TweetService.DeleteTweet:input_type -> tweet.DeleteTweetRequest
	8,  // 10: tweet.TweetService.GetAllTweets:input_type -> tweet.GetAllTweetsRequest
	11, // 11: tweet.TweetService.IncreaseViewsCount:input_type -> tweet.Id
	11, // 12: tweet.TweetService.IncreaseRepostCount:input_type -> tweet.Id
	11, // 13: tweet.TweetService.IncreaseSharesCount:input_type -> tweet.Id
	2,  // 14: tweet.TweetService.CreateTweet:output_type -> tweet.CreateTweetResponse
	4,  // 15: tweet.TweetService.GetTweet:output_type -> tweet.GetTweetResponse
	6,  // 16: tweet.TweetService.UpdateTweet:output_type -> tweet.UpdateTweetResponse
	10, // 17: tweet.TweetService.DeleteTweet:output_type -> tweet.Status
	9,  // 18: tweet.TweetService.GetAllTweets:output_type -> tweet.GetAllTweetsResponse
	10, // 19: tweet.TweetService.IncreaseViewsCount:output_type -> tweet.Status
	10, // 20: tweet.TweetService.IncreaseRepostCount:output_type -> tweet.Status
	10, // 21: tweet.TweetService.IncreaseSharesCount:output_type -> tweet.Status
	14, // [14:22] is the sub-list for method output_type
	6,  // [6:14] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_tweet_service_tweet_service_proto_init() }
func file_tweet_service_tweet_service_proto_init() {
	if File_tweet_service_tweet_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tweet_service_tweet_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Tweet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateTweetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*CreateTweetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GetTweetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*GetTweetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateTweetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateTweetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteTweetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*GetAllTweetsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*GetAllTweetsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*Status); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tweet_service_tweet_service_proto_msgTypes[11].Exporter = func(v any, i int) any {
			switch v := v.(*Id); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tweet_service_tweet_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tweet_service_tweet_service_proto_goTypes,
		DependencyIndexes: file_tweet_service_tweet_service_proto_depIdxs,
		MessageInfos:      file_tweet_service_tweet_service_proto_msgTypes,
	}.Build()
	File_tweet_service_tweet_service_proto = out.File
	file_tweet_service_tweet_service_proto_rawDesc = nil
	file_tweet_service_tweet_service_proto_goTypes = nil
	file_tweet_service_tweet_service_proto_depIdxs = nil
}
