package controllers

import (
	"net/http"
	commentrequests "semen_project/internal/dto/comment_requests"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"semen_project/internal/repository"
)
type CommentHandler struct {
	Repo repository.CommentRepo
}
func NewCommentHandler(repo repository.CommentRepo) *CommentHandler {
	return &CommentHandler{
		Repo: repo,
	}
}
func (h *CommentHandler) CreateComment(ctx *gin.Context) {
	var content commentrequests.CreateComment
	err := ctx.BindJSON(&content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении request body"})
		return
	}
	if len(content.Content) > 10000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Превышен лимит символов - 10000"})
		return
	}
	if strings.TrimSpace(content.Content) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Комментарий не может быть пустым"})
		return
	}
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при чтении id из токена"})
		return
	}
	userID := value.(int)
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id"})
		return
	}
	idparam := ctx.Param("id")
	postID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id поста"})
		return
	}
	if postID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id поста"})
		return
	}
	_, err = h.Repo.GetPostById(postID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Поста с таким id не существует"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении поста"})
		return
	}
	err = h.Repo.CreateComment(postID, userID, content.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании комментария"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Комментарий успешно отправлен"})
}
func (h *CommentHandler) DeleteComment(ctx *gin.Context) {
	idparam := ctx.Param("id")
	commentID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id комментария"})
		return
	}
	if commentID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id комментария"})
		return
	}
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id из токена"})
		return
	}
	userID := value.(int)
	comment, err := h.Repo.GetCommentById(commentID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Комментарий с таким id не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении комментария"})
		return
	}
	if userID != comment.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете удалить чужой комментарий"})
		return
	}
	err = h.Repo.DeleteComment(commentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении комментария"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Комментарий успешно удален"})
}
func (h *CommentHandler) UpdateComment(ctx *gin.Context) {
	var content commentrequests.UpdateComment
	err := ctx.BindJSON(&content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении request body"})
		return
	}
	if len(content.Content) > 10000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Превышен лимит символов - 10000"})
		return
	}
	if strings.TrimSpace(content.Content) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Комментарий не может быть пустым"})
		return
	}
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при чтении id из токена"})
		return
	}
	userID := value.(int)
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id"})
		return
	}
	idparam := ctx.Param("id")
	commentID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id поста"})
		return
	}
	if commentID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id поста"})
		return
	}
	comment, err := h.Repo.GetCommentById(commentID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Комментарий не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении комментария"})
		return
	}
	if userID != comment.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете обновить чужой комментарий"})
		return
	}
	err = h.Repo.UpdateComment(commentID, content.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении комментария"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Комментарий успешно обновлен"})
}
func (h *CommentHandler) GetAllPostComments(ctx *gin.Context) {
	idparam := ctx.Param("id")
	postID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id поста"})
		return
	}
	if postID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id поста"})
		return
	}
	_, err = h.Repo.GetPostById(postID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении поста "})
		return
	}
	comments, err := h.Repo.GetAllPostComments(postID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении всех комментариев"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"comments": comments})
	
}
func (h *CommentHandler) GetCountComments(ctx *gin.Context) {
	idparam := ctx.Param("id")
	postID, err := strconv.Atoi(idparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении id поста"})
		return
	}
	if postID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id поста"})
		return
	}
	_, err = h.Repo.GetPostById(postID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении поста"})
		return
	}
	count, err := h.Repo.GetCountComments(postID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении числа комментариев"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"count": count})
}
