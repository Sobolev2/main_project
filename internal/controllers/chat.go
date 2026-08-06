package controllers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"semen_project/internal/repository"
)
type ChatHandler struct {
	Repo repository.ChatRepo
}
func NewChatHandler(repo repository.ChatRepo) *ChatHandler {
	return &ChatHandler{
		Repo: repo,
	}
}
func (h *ChatHandler) CreateChat(ctx *gin.Context) {
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID не найден в контексте"})
		return
	}

	userID := value.(int)
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный userID"})
		return
	}

	idParam := ctx.Param("id")
	userSecondID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id пользователя"})
		return
	}

	if userSecondID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный UserSecondID"})
		return
	}

	if userID == userSecondID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "нельзя создать чат с самим собой"})
		return
	}

	_, err = h.Repo.GetUserById(userSecondID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователя"})
		return
	}

	// Проверяем существование чата
	chat, err := h.Repo.GetChatByUsersID(userID, userSecondID)

	if err != nil && err != pgx.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке существования чата"})
		return
	}

	if chat.ID > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Чат с этим пользователем уже существует"})
		return
	}

	// Создаем чат
	err = h.Repo.CreateChat(userID, userSecondID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании чата"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Чат успешно создан"})
}
func (h *ChatHandler) GetAllUserChats(ctx *gin.Context) {
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID не найден в контексте"})
		return			
	}
	userID := value.(int)
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный userID"})
		return			
	}
	chats, err := h.Repo.GetAllUserChats(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении чатов пользователя"})
		return			
	}
	ctx.JSON(http.StatusOK, gin.H{"chats": chats})
}
func (h *ChatHandler) DeleteChat(ctx *gin.Context) {
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID не найден в контексте"})
		return			
	}
	userID := value.(int)
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный userID"})
		return			
	}
	idParam := ctx.Param("id")
	chatID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return		
	}
	if chatID <= 0 {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id"})
	return
	}
	chat, err := h.Repo.GetChatById(chatID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Чат не найден"})
		return			
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении чата"})
		return			
	}
	if chat.UserFirst != userID && chat.UserSecond != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете удалить чужой чат"})
		return			
	}
	err = h.Repo.DeleteChat(chatID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении чата"})
		return			
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Чат успешно удален"})
}
func (h *ChatHandler) GetChatByUserID(ctx *gin.Context) {
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID не найден в контексте"})
		return			
	}
	userID := value.(int)
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный userID"})
		return			
	}
	idParam := ctx.Param("id")
	otherUserID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return		
	}
	if otherUserID <= 0 {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id"})
	return
	}
	chat, err := h.Repo.GetChatByUsersID(userID, otherUserID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Чат не найден"})
		return			
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении чата"})
		return			
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Чат найден", "chat": chat})
}

