package controllers

import (
	"net/http"
	"semen_project/internal/dto/authentication_requests"
	"semen_project/internal/repository"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)
type AuthHandler struct {
	Repo repository.AuthRepo
	Secret string
}
func NewAuthHandler(repo repository.AuthRepo, secret string) *AuthHandler {
	return &AuthHandler{
		Repo: repo,
		Secret: secret,
	}
}
func (h *AuthHandler) Register(ctx *gin.Context) {
	var user dto.RegisterRequest
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в request body"})
		return
	}
	if strings.TrimSpace(user.UserName) == "" || strings.TrimSpace(user.Password) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "UserName и Password не могут быть пустыми"})
		return
	}
	if len(user.Password) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "Пароль должен содержать минимум 8 символов"})
		return
	}
	checkUserName, err := h.Repo.GetUserByUsername(user.UserName)
	if err != nil && err != pgx.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных при проверке UserName"})
		return		
	}
	if checkUserName.ID != 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким UserName уже существует"})
		return		
	}
	user.Password, err = repository.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}
	createdUser, err := h.Repo.CreateUser(user.UserName, user.FirstName, user.LastName, user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не получилось создать Пользователя"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": createdUser})
}
func (h *AuthHandler) Login(ctx *gin.Context) {

	// проверка на пустые поля и валидация данных

	var loginData dto.LoginRequest
	err := ctx.BindJSON(&loginData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в чтении запроса"})
		return	
	}
	if strings.TrimSpace(loginData.UserName) == "" || strings.TrimSpace(loginData.Password) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "UserName и Password не могут быть пустыми"})
		return
	}
	usercheck, err := h.Repo.GetUserByUsername(loginData.UserName)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return		
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных при проверке UserName"})
		return			
	}
	userCheckPassword, err	:= h.Repo.GetPasswordById(usercheck.ID)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return		
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных при проверке пароля"})
		return			
	}
	err = repository.CheckPassword(userCheckPassword, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный UserName или пароль"})
		return	
	}

	// создание токенов

	accessToken, err := repository.GenerateAccessToken(usercheck.ID, h.Secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации access токена"})
		return
	}
	refreshToken, err := repository.GenerateRefreshToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации refresh токена"})
		return
	}
	refreshTokenHash := repository.HashRefreshToken(refreshToken)

	err = h.Repo.SaveRefreshToken(usercheck.ID, refreshTokenHash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении refresh токена"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}
func (h *AuthHandler) Refresh(ctx *gin.Context) {
	var refreshToken dto.RefreshRequest
	err := ctx.BindJSON(&refreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в request body"})
		return		
	}

	hashToken := repository.HashRefreshToken(refreshToken.RefreshToken)

	RealToken, err := h.Repo.GetRefreshTokenByTokenHash(hashToken)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не найден"})
		return		
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке refresh токена"})
		return			
	}
	if RealToken.ExpiresAt.Before(time.Now()) {
    _ = h.Repo.DeleteRefreshTokenByTokenHash(RealToken.TokenHash)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh токен истек"})
		return
	}

	NewAccessToken, err := repository.GenerateAccessToken(RealToken.UserID, h.Secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации access токена"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": NewAccessToken})
}
func (h *AuthHandler) Logout(ctx *gin.Context) {
	var refreshToken dto.LogoutRequest
	err := ctx.BindJSON(&refreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в request body"})
		return
	}

	tokenHash := repository.HashRefreshToken(refreshToken.RefreshToken)

	tokenToCompare, err := h.Repo.GetRefreshTokenByTokenHash(tokenHash)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не найден"})
		return		
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке refresh токена"})
		return			
	}
	if tokenToCompare.ExpiresAt.Before(time.Now()) {
		_ = h.Repo.DeleteRefreshTokenByTokenHash(tokenHash)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh токен истек"})
		return
	}
	
	err = h.Repo.DeleteRefreshTokenByTokenHash(tokenHash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении refresh токена"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Вы успешно вышли из аккаунта"})
	// нужно ли как-то забрать у пользователя access токен??
}
func (h *AuthHandler) LogoutAllDevicesExceptThis(ctx *gin.Context) {
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при чтении токена"})
		return	
	}
	id := value.(int)

	var thisToken dto.LogoutAllDevicesExceptThisRequest
	err := ctx.BindJSON(&thisToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка в request body"})
		return			
	}
	hashThisToken := repository.HashRefreshToken(thisToken.RefreshToken)

	tokenToCompare, err := h.Repo.GetRefreshTokenByTokenHash(hashThisToken)
	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Токен не найден"})
		return		
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке refresh токена"})
		return			
	}
    if tokenToCompare.UserID != id {
    ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не принадлежит пользователю"})
    return
    }
	if tokenToCompare.ExpiresAt.Before(time.Now()) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh токен истек"})
		return
	}

	err = h.Repo.DeleteAllUserRefreshTokensExceptThis(id, hashThisToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выходе с аккаунта на всех устройствах"})
		return	
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Вы успешно вышли с аккаунта на всех устройствах"})
}
