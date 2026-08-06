package repository

import (
	"context"
	"semen_project/internal/models"
)
type MessageRepo interface {
	GetMessageById(messageId int) (models.Message, error)
	UpdateMarkReadToRead(messageId int) error
	GetAllChatMessages(chatId int) ([]models.Message, error)
	UpdateMessage(messageId int, content string) error
	DeleteMessage(messageId int) error
	SendMessage(chatID, senderID int, content string) error
	GetCountNotReadMessages(chatID, userID int) (int, error)
	GetMessageStatus(messageID int) (bool, error)
	GetChatById(chatId int) (models.Chat, error)
}
func (s *Store) GetMessageById(message_id int) (models.Message, error) {
	var message models.Message
	query := `
	SELECT id, chat_id, sender_id, content, mark_read
	FROM messages
	WHERE id = $1
	`
	err := s.db.QueryRow(context.Background(), query, message_id).Scan(&message.ID,
		 &message.ChatID, &message.SenderID, &message.Content, &message.MarkRead)
		if err != nil {
			return models.Message{}, err
		}
	return message, nil
}
func (s *Store) UpdateMarkReadToRead(message_id int) error {
	query := `
	UPDATE messages
	SET mark_read = true
	WHERE id = $1
	`
	_,err := s.db.Exec(context.Background(), query, message_id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetAllChatMessages(chat_id int) ([]models.Message, error) {
	var messages []models.Message
	query := `
	SELECT id, chat_id, sender_id, content, mark_read
	FROM messages
	WHERE chat_id = $1
	ORDER BY id DESC
	`
	rows, err := s.db.Query(context.Background(), query, chat_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Content, &message.MarkRead)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}
	return messages, nil
}
func (s *Store) DeleteMessage(message_id int) error {
	query := `
	DELETE FROM messages
	WHERE id = $1
	`
	_, err := s.db.Exec(context.Background(), query, message_id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) SendMessage(chatID int, senderID int, content string) error {
	query := `
	INSERT INTO messages (chat_id, sender_id, content)
	VALUES ($1, $2, $3)
	`
	_, err := s.db.Exec(context.Background(), query, chatID, senderID, content)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetCountNotReadMessages(chatID int, userID int) (int,error) {
	query := `
	SELECT COUNT(*)
	FROM messages
	WHERE chat_id = $1 AND sender_id != $2 AND mark_read = false
	`
	var count int
	err := s.db.QueryRow(context.Background(),query, chatID, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (s *Store) GetMessageStatus(MessageID int) (bool,error) {
	query := `
	SELECT mark_read
	FROM messages
	WHERE id = $1
	`
	var status bool
	err := s.db.QueryRow(context.Background(), query, MessageID).Scan(&status)
	if err != nil {
		return false, err
	}
	return status, nil
}
func (s *Store) UpdateMessage(messageID int,content string) error {
	query := `
	UPDATE messages
	SET content =$1
	WHERE id = $2
	`
	_,err := s.db.Exec(context.Background(), query, content, messageID)
	if err != nil {
		return err
	}
	return nil
}