package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"semen_project/internal/dto/post_requests"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"semen_project/internal/repository"
)
type PostHandler struct {
	Repo repository.PostRepo
}
func NewPostHandler(repo repository.PostRepo) *PostHandler {
	return &PostHandler{
		Repo: repo,
	}
}
func (h *PostHandler) CreatePost(ctx *gin.Context) {
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID не найден в контексте"})
		return			
	}
	userID := value.(int)
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id"})
		return
	}
	var post postrequests.CreatePostRequest
	err := ctx.BindJSON(&post)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в request body"})
		return		
	}
	if strings.TrimSpace(post.Content) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Content не может быть пустым"})
		return
	}
	if len(post.Content) > 10000 {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "Слишком длинный пост, лимит 10000 символов"})
    return
	}
	createdPost, err := h.Repo.CreatePost(userID, post.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании поста"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Пост успешно создан", "post": createdPost})
}
func (h *PostHandler) GetPostById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	postID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return		
	}
	if postID <= 0 {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id"})
	return
	}
	post, err := h.Repo.GetPostById(postID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении поста"})
		return
	}
	if post.ID <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Пост найден", "post": post})
}
func (h *PostHandler) GetFeed(ctx *gin.Context) {
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID не найден в контексте"})
		return			
	}
	userID := value.(int)
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id"})
		return
	}
	firstPartPosts, err := h.Repo.GetAllFriendsPosts(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении постов друзей"})
		return
	}
	secondPartPosts, err := h.Repo.GetAllNotFriendsPosts(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении постов пользователей"})
		return
	}
	feed := append(firstPartPosts, secondPartPosts...)
	if len(feed) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Лента постов пуста", "feed": feed})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Лента постов", "feed": feed})
}
func (h *PostHandler) GetAllPosts(ctx *gin.Context) {
	posts, err := h.Repo.GetAllPosts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении всех постов"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Все посты", "posts": posts})
}
func (h *PostHandler) GetAllUserPosts(ctx *gin.Context) {
	IdParam := ctx.Param("id")
	userID, err := strconv.Atoi(IdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return		
	}
	if userID <= 0 {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id"})
	return
	}
	posts, err := h.Repo.GetAllUserPosts(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении постов пользователя"})
		return
	}
	if len(posts) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Список постов пользователя пуст", "posts": posts})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Посты пользователя", "posts": posts})
}
func (h *PostHandler) DeletePost(ctx *gin.Context) {
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID не найден в контексте"})
		return			
	}
	userID := value.(int)
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id поста"})
		return
	}
	idParam := ctx.Param("id")
	postID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при чтении id поста"})
		return		
	}
	post, err := h.Repo.GetPostById(postID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return		
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении поста"})
		return
	}
	if post.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете удалить чужой пост"})
		return
	}
	err = h.Repo.DeletePost(post.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении поста"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Пост успешно удален"})
}
func (h *PostHandler) UpdatePost(ctx *gin.Context) {
	idParam := ctx.Param("id")
	postID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при чтении id поста"})
		return		
	}
	if postID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id поста"})
		return
	}
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID не найден в контексте"})
		return			
	}
	userID := value.(int)
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id пользователя"})
		return
	}
	var newContent postrequests.UpdatePostRequest
	err = ctx.BindJSON(&newContent)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в request body"})
		return		
	}
	post, err := h.Repo.GetPostById(postID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return		
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении поста"})
		return
	}
	if post.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Вы не можете изменить чужой пост"})
		return
	}
	if strings.TrimSpace(newContent.Content) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Content не может быть пустым"})
		return
	}
	if len(newContent.Content) > 10000 {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Слишком длинный пост, лимит 10000 символов"})
	return
	}
	updatedPost, err := h.Repo.UpdatePost(postID, newContent.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении поста"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Пост успешно обновлен", "post": updatedPost})
}