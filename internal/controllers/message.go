package controllers

import (
	"net/http"
	"semen_project/internal/dto/message_requests"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"semen_project/internal/repository"
)
type MessageHandler struct {
	Repo repository.MessageRepo
}
func NewMessageHandler(repo repository.MessageRepo) *MessageHandler {
	return &MessageHandler{
		Repo: repo,
	}
}
func (h *MessageHandler) SendMessage(ctx *gin.Context) {
	idparam := ctx.Param("id")
	chatID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id чата"})
		return
	}
	if chatID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id чата"})
		return
	}

	var content messagerequests.SendMessageRequest
	err = ctx.ShouldBindJSON(&content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при чтении json body"})
		return
	}
	if strings.TrimSpace(content.Content) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Сообщение не может быть пустым"})
		return
	}
	if len(content.Content) > 10000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Превышен лимит колличества символом(10000)"})
		return
	}

	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id пользователя из токена"})
		return
	}
	userID, ok := value.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id пользователя в токене"})
		return
	}
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id пользователя"})
		return
	}

	chat, err := h.Repo.GetChatById(chatID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Такого чата не существует"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении чата"})
		return
	}

	if userID != chat.UserFirst && userID != chat.UserSecond {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете отправлять сообщения в чужом чате"})
		return
	}

	err = h.Repo.SendMessage(chatID, userID, content.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при отправке сообщения"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Сообщение отправлено"})
}
func (h *MessageHandler) UpdateMessage(ctx *gin.Context) {
	idparam := ctx.Param("id")
	messageID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id сообщения"})
		return
	}
	if messageID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id сообщения"})
		return
	}

	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id пользователя из токена"})
		return
	}
	userID, ok := value.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id пользователя в токене"})
		return
	}
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id пользователя"})
		return
	}

	var content messagerequests.UpdateMessageRequest
	err = ctx.ShouldBindJSON(&content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при чтении json body"})
		return
	}
	if strings.TrimSpace(content.Content) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Сообщение не может быть пустым"})
		return
	}
	if len(content.Content) > 10000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Превышен лимит колличества символом(10000)"})
		return
	}

	message, err := h.Repo.GetMessageById(messageID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Сообщение с таким id не найдено"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении сообщения"})
		return
	}

	if userID != message.SenderID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете редактировать чужое сообщение"})
		return
	}

	err = h.Repo.UpdateMessage(messageID, content.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении сообщения"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Сообщение успешно отредактировано"})
}
func (h *MessageHandler) DeleteMessage(ctx *gin.Context) {
	idparam := ctx.Param("id")
	messageID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id сообщения"})
		return
	}
	if messageID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id сообщения"})
		return
	}

	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id пользователя из токена"})
		return
	}
	userID, ok := value.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id пользователя в токене"})
		return
	}
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id пользователя"})
		return
	}

	message, err := h.Repo.GetMessageById(messageID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cообщение с таким id не найдено"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении сообщения"})
		return
	}

	if userID != message.SenderID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете удалить чужое сообщение"})
		return
	}

	err = h.Repo.DeleteMessage(messageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении сообщения"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Сообщение успешно удалено"})
}
func (h *MessageHandler) GetAllChatMessages(ctx *gin.Context) {
	idparam := ctx.Param("id")
	chatID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id чата"})
		return
	}
	if chatID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id чата"})
		return
	}

	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id пользователя из токена"})
		return
	}
	userID, ok := value.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id пользователя в токене"})
		return
	}
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id пользователя"})
		return
	}	
	chat, err := h.Repo.GetChatById(chatID)

	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Чат с таким id не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении чата"})
		return
	}

	if userID != chat.UserFirst && userID != chat.UserSecond {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете просматривать чат, в котором не состоите"})
		return
	}

	messages, err := h.Repo.GetAllChatMessages(chatID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении всех сообщений чата"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"messages": messages})
}
func (h *MessageHandler) GetCountNotReadMessages(ctx *gin.Context) {
	idparam := ctx.Param("id")
	chatID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id чата"})
		return
	}
	if chatID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id чата"})
		return
	}
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id пользователя из токена"})
		return
	}
	userID, ok := value.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id пользователя в токене"})
		return
	}
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id пользователя"})
		return
	}	

	chat, err := h.Repo.GetChatById(chatID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "такого чата не существует"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении чата"})
		return
	}

	if userID != chat.UserFirst && userID != chat.UserSecond {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете просматривать чужой чат"})
		return
	}
	// получить кол-во непрочитанных сообщений где юзер айди равен другому пользователю из бд
	count, err := h.Repo.GetCountNotReadMessages(chatID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении числа непрочитанных сообщений пользователя"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"count": count})
}
func (h *MessageHandler) UpdateMarkReadToRead(ctx *gin.Context) {
	idparam := ctx.Param("id")
	messageID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id сообщения"})
		return
	}
	if messageID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id сообщения"})
		return
	}

	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id пользователя из токена"})
		return
	}
	userID, ok := value.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id пользователя в токене"})
		return
	}
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id пользователя"})
		return
	}	
	message, err := h.Repo.GetMessageById(messageID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Сообщениe с таким id не существует"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении сообщения"})
		return
	}
	chat, err := h.Repo.GetChatById(message.ChatID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Чат с таким id не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении чата"})
		return
	}
	if userID != chat.UserFirst && userID != chat.UserSecond {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете просматривать чужой чат"})
		return
	}
	if message.SenderID == userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете изменить статус просмотра своего сообщения"})
		return
	}
	status, err := h.Repo.GetMessageStatus(messageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении статуса сообщения"})
		return
	}
	if status {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Сообщение уже прочитано"})
		return
	}
	err = h.Repo.UpdateMarkReadToRead(messageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновалении статуса сообщения"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Статус сообщения успешно изменен на прочитано"})
}