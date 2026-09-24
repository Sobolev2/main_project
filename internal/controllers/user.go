package controllers

import (
	"net/http"
	dto "semen_project/internal/dto/user_requests"
	"semen_project/internal/repository"
	"strconv"
	"strings"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	Repo repository.UserRepo
}

func NewUserHandler(repo repository.UserRepo) *UserHandler {
	return &UserHandler{
		Repo: repo,
	}
}

func (h *UserHandler) GetMe(ctx *gin.Context) {
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID не найден в контексте"})
		return
	}
	id := value.(int)
	user, err := h.Repo.GetUserById(id)
	if errors.Is(err, pgx.ErrNoRows) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка в получении пользователя по id"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var updateUser dto.UpdateUserRequest
	err := ctx.BindJSON(&updateUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в request body"})
		return
	}
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "userID не найден в контексте"})
		return
	}
	id := value.(int)
	LastUser, err := h.Repo.GetUserById(id)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении прошлых данных"})
		return
	}

	if strings.TrimSpace(updateUser.UserName) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "UserName не может быть пустым"})
		return
	}
	if strings.TrimSpace(updateUser.UserName) == strings.TrimSpace(LastUser.UserName) &&
		strings.TrimSpace(updateUser.FirstName) == strings.TrimSpace(LastUser.FirstName) &&
		strings.TrimSpace(updateUser.LastName) == strings.TrimSpace(LastUser.LastName) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка, все новые данные совпадают с прошлыми"})
		return
	}
	usedUserName, err := h.Repo.GetUserByUsernameExceptId(id, updateUser.UserName)
	if err != nil && err != pgx.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке уникальности username"})
		return
	}
	if usedUserName.ID != 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким username уже существует"})
		return
	}
	err = h.Repo.UpdateUser(updateUser.UserName, updateUser.FirstName, updateUser.LastName, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении данных пользователя"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно обновлен"})
}
func (h *UserHandler) GetAllUsers(ctx *gin.Context) {
	users, err := h.Repo.GetAllUsers()
	if err != nil && err != pgx.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении всех пользователей"})
		return
	}
	if len(users) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Список пользователей пуст",
			"users":   users,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Список пользователей найден", "users": users})
}
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	var password dto.DeleteUserRequest
	err := ctx.BindJSON(&password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в request body"})
		return
	}
	if strings.TrimSpace(password.Password) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password не может быть пустым"})
		return
	}
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": " userID не найден в контексте"})
		return
	}
	id := value.(int)

	lastPassword, err := h.Repo.GetPasswordById(id)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получении пароля по id"})
		return
	}
	err = h.Repo.CheckPassword(lastPassword, password.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}
	err = h.Repo.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении пользователя"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": " Пользователь успешно удален"})
}
func (h *UserHandler) UpdatePassword(ctx *gin.Context) {
	var passwords dto.UpdatePasswordRequest
	err := ctx.BindJSON(&passwords)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в request body"})
		return
	}
	if strings.TrimSpace(passwords.LastPassword) == "" || strings.TrimSpace(passwords.NewPassword) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "LastPassword и NewPassword не могут быть пустыми"})
		return
	}
	if len(passwords.NewPassword) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Пароль должен содержать минимум 8 символов",
		})
		return
	}
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": " userID не найден в контексте"})
		return
	}
	id := value.(int)

	dbPassword, err := h.Repo.GetPasswordById(id)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить пароль из базы данных"})
		return
	}

	err = h.Repo.CheckPassword(dbPassword, passwords.LastPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}
	err = h.Repo.CheckPassword(dbPassword, passwords.NewPassword)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Новый пароль совпадает с текущим"})
		return
	}
	HashNewPassword, err := repository.HashPassword(passwords.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось захешировать новый пароль"})
		return
	}

	err = h.Repo.UpdatePassword(id, HashNewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось изменить пароль"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Пароль успешно изменен"})
}
func (h *UserHandler) GetUserById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}
	if userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id"})
		return
	}
	user, err := h.Repo.GetUserById(userID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователя"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Пользователь найден", "user": user})
}
func (h *UserHandler) SearchUsers(ctx *gin.Context) {
	query := strings.TrimSpace(ctx.Query("query"))
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Параметр query не может быть пустым"})
		return
	}
	if len(query) < 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Минимум 2 символа для поиска"})
		return
	}
	users, err := h.Repo.SearchUsers(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске пользователей"})
		return
	}
	if len(users) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Пользователи не найдены",
			"users":   users,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Пользователи найдены",
		"users":   users,
	})
}
func (h *UserHandler) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Параметр username не может быть пустым"})
		return
	}
	user, err := h.Repo.GetUserByUsername(username)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователя"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Пользователь найден", "user": user})
}
