package controllers

import (
	"semen_project/internal/repository"
	"semen_project/internal/kafka"
)
type Handlers struct {
	AuthHandler    *AuthHandler
	ChatHandler    *ChatHandler
	CommentHandler *CommentHandler
	FollowHandler  *FollowHandler
	FriendHandler  *FriendHandler
	LikeHandler    *LikeHandler
	MessageHandler *MessageHandler
	PostHandler    *PostHandler
	UserHandler    *UserHandler
}
func NewHandlers(
	store *repository.Store,
	secret string,
	producer *kafka.Producer,
) *Handlers {
	return &Handlers{
		AuthHandler:    NewAuthHandler(store, secret),
		UserHandler:    NewUserHandler(store),
		PostHandler:    NewPostHandler(store),
		CommentHandler: NewCommentHandler(store),
		LikeHandler:    NewLikeHandler(store, producer),
		FriendHandler:  NewFriendHandler(store),
		FollowHandler:  NewFollowHandler(store),
		MessageHandler: NewMessageHandler(store),
		ChatHandler:    NewChatHandler(store),
	}
}