package controllers

import (
	"semen_project/internal/kafka"
	"semen_project/internal/repository"

	pb "semen_project/proto/gen/notificationpb"
)

type Handlers struct {
	AuthHandler         *AuthHandler
	ChatHandler         *ChatHandler
	CommentHandler      *CommentHandler
	FollowHandler       *FollowHandler
	FriendHandler       *FriendHandler
	LikeHandler         *LikeHandler
	MessageHandler      *MessageHandler
	PostHandler         *PostHandler
	UserHandler         *UserHandler
	NotificationHandler *NotificationHandler
}

func NewHandlers(
	store *repository.Store,
	secret string,
	producer *kafka.MessageProducer,
	notificationClient pb.NotificationServiceClient,
) *Handlers {
	return &Handlers{
		AuthHandler:         NewAuthHandler(store, secret),
		UserHandler:         NewUserHandler(store),
		PostHandler:         NewPostHandler(store),
		CommentHandler:      NewCommentHandler(store),
		LikeHandler:         NewLikeHandler(store),
		FriendHandler:       NewFriendHandler(store),
		FollowHandler:       NewFollowHandler(store),
		MessageHandler:      NewMessageHandler(store, producer),
		ChatHandler:         NewChatHandler(store),
		NotificationHandler: NewNotificationHandler(notificationClient),
	}
}
