package controllers_test

import (
	"semen_project/internal/controllers"
	"semen_project/internal/mocks"
	"semen_project/internal/models"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"errors"
	"strings"

	"semen_project/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// register
func TestRegister_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockAuthRepo)
	handler := controllers.NewAuthHandler(mockRepo, "test_secret")
	router := gin.New()
	body := `{"user_name":"testuser","first_name":"Test","last_name":"User","password":"password123"}`
	router.POST("/register", func(c *gin.Context) {
		handler.Register(c)
	})
	user := models.UserPublic{
		ID:        1,
		UserName:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}
	mockRepo.On("GetUserByUsername", "testuser").Return(models.UserPublic{}, pgx.ErrNoRows)
	mockRepo.On("CreateUser", "testuser", "Test", "User", mock.Anything).Return(&user, nil)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	assert.Contains(
		t,
		recorder.Body.String(),
		"testuser",
	)
	mockRepo.AssertExpectations(t)
}
func TestRegister_Errors(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		mockGetUser bool
		// GetUserByUsername
		usedUser    models.UserPublic
		usedUserErr error

		wantStatus int
	}{
		{
			name:        "bad request",
			body:        "",
			mockGetUser: false,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "empty username or password",
			body:        `{"user_name":"","first_name":"Test","last_name":"User","password":""}`,
			mockGetUser: false,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "password too short",
			body:        `{"user_name":"testuser","first_name":"Test","last_name":"User","password":"short"}`,
			mockGetUser: false,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "username already exists",
			body:        `{"user_name":"testuser","first_name":"Test","last_name":"User","password":"password123"}`,
			mockGetUser: true,
			usedUser:    models.UserPublic{ID: 1, UserName: "testuser"},
			usedUserErr: nil,
			wantStatus:  http.StatusConflict,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepo)
			handler := controllers.NewAuthHandler(mockRepo, "test_secret")

			if test.mockGetUser {
				mockRepo.On("GetUserByUsername", mock.Anything).Return(test.usedUser, test.usedUserErr)
			}
			router := gin.New()
			router.POST("/register", handler.Register)

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(
				t,
				test.wantStatus,
				recorder.Code,
			)
			if !test.mockGetUser {
				mockRepo.AssertNotCalled(
					t,
					"GetUserByUsername",
					mock.Anything,
				)
			}
			mockRepo.AssertExpectations(t)
		})

	}
}
func TestRegister_InternalError(t *testing.T) {
	tests := []struct {
		name string
		body string

		// GetUserByUsername
		usedUserErr error

		// CreateUser
		createUserErr error

		wantStatus int
	}{
		{
			name:        "internal error GetUserByUsername",
			body:        `{"user_name":"username","first_name":"Test","last_name":"User","password":"password123"}`,
			usedUserErr: errors.New("internal error GetUserByUsername"),
			wantStatus:  http.StatusInternalServerError,
		},
		{
			name:          "internal error CreateUser",
			body:          `{"user_name":"userName","first_name":"Test","last_name":"User","password":"password123"}`,
			usedUserErr:   pgx.ErrNoRows,
			createUserErr: errors.New("internal error CreateUser"),
			wantStatus:    http.StatusInternalServerError,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepo)
			handler := controllers.NewAuthHandler(mockRepo, "test_secret")
			mockRepo.On("GetUserByUsername", mock.Anything).Return(models.UserPublic{}, test.usedUserErr)
			if errors.Is(test.usedUserErr, pgx.ErrNoRows) {
				mockRepo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&models.UserPublic{}, test.createUserErr)
			}
			router := gin.New()
			router.POST("/register", handler.Register)

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(
				t,
				test.wantStatus,
				recorder.Code,
			)
			mockRepo.AssertExpectations(t)
		})
	}
}

// login
func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.MockAuthRepo)
	handler := controllers.NewAuthHandler(mockRepo, "test_secret")

	router := gin.New()
	router.POST("/login", handler.Login)

	// Настоящий bcrypt-хэш
	hash, err := repository.HashPassword("password")
	assert.NoError(t, err)

	user := models.UserPublic{
		ID:       1,
		UserName: "testuser",
	}

	mockRepo.On("GetUserByUsername", "testuser").
		Return(user, nil)

	mockRepo.On("GetPasswordById", 1).
		Return(hash, nil)

	mockRepo.On("CheckPassword", hash, "password").Return(nil)

	mockRepo.On("GenerateAccessToken", 1, "test_secret").
		Return("test-access-token", nil)

	mockRepo.On("SaveRefreshToken", 1, mock.AnythingOfType("string")).
		Return(nil)

	body := `{
		"user_name":"testuser",
		"password":"password"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	assert.Contains(t, recorder.Body.String(), "access_token")
	assert.Contains(t, recorder.Body.String(), "refresh_token")

	mockRepo.AssertExpectations(t)
}
func TestLogin_Errors(t *testing.T) {
	tests := []struct {
		name string
		body string

		mockGetUser bool

		user    models.UserPublic
		userErr error

		password    string
		passwordErr error

		wantStatus int
	}{
		{
			name:       "bad request",
			body:       "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty username",
			body:       `{"user_name":"","password":"password123"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty password",
			body:       `{"user_name":"testuser","password":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "user not found",
			body:        `{"user_name":"testuser","password":"password123"}`,
			mockGetUser: true,
			userErr:     pgx.ErrNoRows,
			wantStatus:  http.StatusNotFound,
		},
		{
			name:        "password not found",
			body:        `{"user_name":"testuser","password":"password123"}`,
			mockGetUser: true,
			user: models.UserPublic{
				ID:       1,
				UserName: "testuser",
			},
			passwordErr: pgx.ErrNoRows,
			wantStatus:  http.StatusNotFound,
		},
		{
			name:        "wrong password",
			body:        `{"user_name":"testuser","password":"password123"}`,
			mockGetUser: true,
			user: models.UserPublic{
				ID:       1,
				UserName: "testuser",
			},
			password:   "$2a$10$aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			wantStatus: http.StatusUnauthorized,
		},
	}

	gin.SetMode(gin.TestMode)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepo)

			handler := controllers.NewAuthHandler(mockRepo, "test_secret")

			if test.mockGetUser {
				mockRepo.
					On("GetUserByUsername", mock.Anything).
					Return(test.user, test.userErr)

				if test.userErr == nil {
					mockRepo.
						On("GetPasswordById", 1).
						Return(test.password, test.passwordErr)

					if test.passwordErr == nil {
						mockRepo.
							On("CheckPassword", test.password, "password123").
							Return(errors.New("wrong password"))
					}
				}
			}

			router := gin.New()
			router.POST("/login", handler.Login)

			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(test.body),
			)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, test.wantStatus, rec.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
func TestLogin_InternalError(t *testing.T) {
	tests := []struct {
		name string
		body string

		userErr      error
		passwordErr  error
		saveTokenErr error

		wantStatus int
	}{
		{
			name: "GetUserByUsername error",
			body: `{
				"user_name":"testuser",
				"password":"password123"
			}`,
			userErr:    errors.New("db error"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "GetPasswordById error",
			body: `{
				"user_name":"testuser",
				"password":"password123"
			}`,
			passwordErr: errors.New("db error"),
			wantStatus:  http.StatusInternalServerError,
		},
		{
			name: "SaveRefreshToken error",
			body: `{
				"user_name":"testuser",
				"password":"password123"
			}`,
			saveTokenErr: errors.New("db error"),
			wantStatus:   http.StatusInternalServerError,
		},
	}

	gin.SetMode(gin.TestMode)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			hash, err := repository.HashPassword("password123")
			assert.NoError(t, err)

			mockRepo := new(mocks.MockAuthRepo)

			handler := controllers.NewAuthHandler(
				mockRepo,
				"test_secret",
			)

			user := models.UserPublic{
				ID:       1,
				UserName: "testuser",
			}

			mockRepo.
				On("GetUserByUsername", "testuser").
				Return(user, test.userErr)

			if test.userErr == nil {

				mockRepo.
					On("GetPasswordById", 1).
					Return(hash, test.passwordErr)

				if test.passwordErr == nil {

					mockRepo.
						On("CheckPassword", hash, "password123").
						Return(nil)

					mockRepo.
						On("GenerateAccessToken", 1, "test_secret").
						Return("test-access-token", nil)

					mockRepo.
						On(
							"SaveRefreshToken",
							1,
							mock.AnythingOfType("string"),
						).
						Return(test.saveTokenErr)
				}
			}

			router := gin.New()
			router.POST("/login", handler.Login)

			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(test.body),
			)

			req.Header.Set(
				"Content-Type",
				"application/json",
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(
				t,
				test.wantStatus,
				rec.Code,
			)

			mockRepo.AssertExpectations(t)
		})
	}
}

// refresh
func TestRefresh_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockAuthRepo)
	handler := controllers.NewAuthHandler(mockRepo, "test_secret")

	router := gin.New()
	router.POST("/refresh", handler.Refresh)

	hash := repository.HashRefreshToken("refreshtoken")

	mockRepo.On("GetRefreshTokenByTokenHash", hash).Return(models.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(60 * time.Second),
		CreatedAt: time.Now(),
	}, nil)
	mockRepo.On("GenerateAccessToken", 1, "test_secret").Return("AccessToken", nil)

	body := `{
	"refresh_token": "refreshtoken"
	}`
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{
        "access_token": "AccessToken"
    }`, recorder.Body.String())
	mockRepo.AssertExpectations(t)
}
func TestRefresh_BadRequest(t *testing.T) {
	hash := repository.HashRefreshToken("token")
	tests := []struct {
		name string
		body string

		mockGetRefresh bool

		getRefreshErr error
		refreshToken  models.RefreshToken

		wantStatus int
	}{{
		name:       "Empry request",
		body:       "",
		wantStatus: http.StatusBadRequest,
	},
		{
			name:       "Invalid JSON",
			body:       `{"refresh_token": ""`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Token not found",
			body: `{"refresh_token": "token"}`,

			mockGetRefresh: true,
			getRefreshErr:  pgx.ErrNoRows,
			wantStatus:     http.StatusUnauthorized,
		},
		{
			name: "Token expired",
			body: `{"refresh_token": "token"}`,

			mockGetRefresh: true,
			getRefreshErr:  nil,
			refreshToken: models.RefreshToken{
				ID:        1,
				UserID:    1,
				TokenHash: hash,
				ExpiresAt: time.Now().Add(-time.Second),
				CreatedAt: time.Now(),
			},

			wantStatus: http.StatusUnauthorized,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepo)

			handler := controllers.NewAuthHandler(mockRepo, "test_secret")

			if test.mockGetRefresh {
				mockRepo.On(
					"GetRefreshTokenByTokenHash",
					hash,
				).Return(test.refreshToken, test.getRefreshErr)

				if test.getRefreshErr == nil &&
					test.refreshToken.ExpiresAt.Before(time.Now()) {

					mockRepo.On(
						"DeleteRefreshTokenByTokenHash",
						test.refreshToken.TokenHash,
					).Return(nil)
				}
			}
			router := gin.New()

			router.POST("/refresh", handler.Refresh)

			req := httptest.NewRequest(
				http.MethodPost,
				"/refresh",
				strings.NewReader(test.body),
			)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, test.wantStatus, recorder.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
func TestRefresh_InternalError(t *testing.T) {
	hash := repository.HashRefreshToken("token")
	tests := []struct {
		name string
		body string

		mockGetRefresh bool

		getRefreshErr error
		refreshToken  models.RefreshToken

		mockGenerateAccessToken bool

		accessToken            string
		generateAccessTokenErr error

		wantStatus int
	}{
		{
			name: "Internal error GetRefreshTokenByTokenHash",
			body: `{"refresh_token": "token"}`,

			mockGetRefresh: true,
			getRefreshErr:  errors.New("internal error"),
			wantStatus:     http.StatusInternalServerError,
		},
		{
			name: "Internal error GenerateAccessToken",
			body: `{"refresh_token": "token"}`,

			mockGetRefresh: true,
			getRefreshErr:  nil,
			refreshToken: models.RefreshToken{
				ID:        1,
				UserID:    1,
				TokenHash: hash,
				ExpiresAt: time.Now().Add(time.Second * 60),
				CreatedAt: time.Now(),
			},

			mockGenerateAccessToken: true,
			generateAccessTokenErr:  errors.New("internal error"),
			wantStatus:              http.StatusInternalServerError,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepo)

			handler := controllers.NewAuthHandler(mockRepo, "test_secret")

			if test.mockGetRefresh {
				mockRepo.On("GetRefreshTokenByTokenHash", hash).Return(test.refreshToken, test.getRefreshErr)
				if test.getRefreshErr == nil &&
					test.refreshToken.ExpiresAt.Before(time.Now()) {

					mockRepo.On(
						"DeleteRefreshTokenByTokenHash",
						test.refreshToken.TokenHash,
					).Return(nil)
				}
				if test.mockGenerateAccessToken {
					mockRepo.On("GenerateAccessToken", test.refreshToken.UserID, "test_secret").Return(mock.Anything, test.generateAccessTokenErr)
				}
			}
			router := gin.New()

			router.POST("/refresh", handler.Refresh)

			req := httptest.NewRequest(
				http.MethodPost,
				"/refresh",
				strings.NewReader(test.body),
			)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, test.wantStatus, recorder.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

// logout
func TestLogout_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockAuthRepo)
	handler := controllers.NewAuthHandler(mockRepo, "test_secret")

	router := gin.New()
	router.POST("/logout", handler.Logout)

	hash := repository.HashRefreshToken("token")

	mockRepo.On("GetRefreshTokenByTokenHash", hash).Return(models.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(60 * time.Second),
		CreatedAt: time.Now(),
	}, nil)
	mockRepo.On("DeleteRefreshTokenByTokenHash", hash).Return(nil)

	body := `{
	"refresh_token": "token"
	}`
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	mockRepo.AssertExpectations(t)
}
func TestLogout_BadRequest(t *testing.T) {
	hash := repository.HashRefreshToken("token")

	tests := []struct {
		name string
		body string

		mockGetRefreshToken bool
		getRefreshTokenErr  error
		storedToken         models.RefreshToken

		expectDeleteRefreshToken bool

		wantStatus int
	}{
		{
			name:       "Bad request",
			body:       "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Token not found",
			body: `{"refresh_token": "token"}`,

			mockGetRefreshToken: true,
			getRefreshTokenErr:  pgx.ErrNoRows,

			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Token Expired",
			body: `{"refresh_token": "token"}`,

			mockGetRefreshToken: true,
			getRefreshTokenErr:  nil,
			storedToken: models.RefreshToken{
				ID:        1,
				UserID:    1,
				TokenHash: hash,
				ExpiresAt: time.Now().Add(-time.Second),
				CreatedAt: time.Now(),
			},
			expectDeleteRefreshToken: true,
			wantStatus:               http.StatusUnauthorized,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepo)
			handler := controllers.NewAuthHandler(mockRepo, "test_secret")
			router := gin.New()
			if test.mockGetRefreshToken {
				mockRepo.On("GetRefreshTokenByTokenHash", hash).Return(test.storedToken, test.getRefreshTokenErr)
				if test.expectDeleteRefreshToken {
					mockRepo.On("DeleteRefreshTokenByTokenHash", hash).Return(nil)
				}
			}
			router.POST("/logout", handler.Logout)

			req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, test.wantStatus, recorder.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
func TestLogout_InternalError(t *testing.T) {
	hash := repository.HashRefreshToken("token")

	tests := []struct {
		name string
		body string

		mockGetRefreshToken bool
		getRefreshTokenErr  error
		storedToken         models.RefreshToken

		expiredDeleteToken bool

		mockDeleteRefreshToken bool
		deleteRefreshTokenErr  error

		wantStatus int
	}{
		{
			name: "Internal error GetRefreshTokenByTokenHash",
			body: `{"refresh_token": "token"}`,

			mockGetRefreshToken: true,
			getRefreshTokenErr:  errors.New("Internal error GetRefreshTokenByTokenHash"),
			expiredDeleteToken: false,
			wantStatus:          http.StatusInternalServerError,
		},
		{
			name: "internal error DeleteRefreshTokenByTokenHash",
			body: `{"refresh_token": "token"}`,

			mockGetRefreshToken: true,
			getRefreshTokenErr:  nil,
			storedToken: models.RefreshToken{
				ID:        1,
				UserID:    1,
				TokenHash: hash,
				ExpiresAt: time.Now().Add(time.Minute),
				CreatedAt: time.Now(),
			},
			expiredDeleteToken: false,
			mockDeleteRefreshToken: true,
			deleteRefreshTokenErr:  errors.New("internal error DeleteRefreshTokenByTokenHash"),
			wantStatus:             http.StatusInternalServerError,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepo)
			handler := controllers.NewAuthHandler(mockRepo, "test_secret")

			if test.mockGetRefreshToken {
				mockRepo.On("GetRefreshTokenByTokenHash", hash).Return(test.storedToken, test.getRefreshTokenErr)
				if test.expiredDeleteToken{
					mockRepo.On("DeleteRefreshTokenByTokenHash", hash).Return(nil)
				}
				if test.mockDeleteRefreshToken {
					mockRepo.On("DeleteRefreshTokenByTokenHash", hash).Return(test.deleteRefreshTokenErr)
				}
			}
			router := gin.New()
			router.POST("/logout", handler.Logout)

			req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, test.wantStatus, recorder.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
// LogoutAllDevicesExceptThis
func TestLogoutAllDevicesExceptThis_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockAuthRepo)
	handler := controllers.NewAuthHandler(mockRepo, "test_secret")

	router := gin.New()
	router.POST("/logout/all-devices", func(c *gin.Context) {
		c.Set("userID", 1)
		handler.LogoutAllDevicesExceptThis(c)
	})

	hash := repository.HashRefreshToken("token")

	mockRepo.On("GetRefreshTokenByTokenHash", hash).Return(models.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(60 * time.Second),
		CreatedAt: time.Now(),
	}, nil)
	mockRepo.On("DeleteAllUserRefreshTokensExceptThis", 1 ,hash).Return(nil)

	body := `{
	"refresh_token": "token"
	}`
	req := httptest.NewRequest(http.MethodPost, "/logout/all-devices", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	mockRepo.AssertExpectations(t)
}
func TestLogoutAllDevicesExceptThis_BadRequest(t *testing.T) {
	hash := repository.HashRefreshToken("token")
	tests := []struct{
		name string
		userID int
		body string

		mockGetToken bool
		GetTokenErr error
		storedToken models.RefreshToken

		deleteExpiredToken bool

		wantStatus int
	}{
		{
			name: "userId not found in context",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Bad request body",
			userID: 1,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Token not found",
			userID: 1,
			body: `{"refresh_token": "token"}`,

			mockGetToken: true,
			GetTokenErr: pgx.ErrNoRows,
			
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "User not own this token",
			userID: 2,
			body: `{"refresh_token": "token"}`,

			mockGetToken: true,
			GetTokenErr: nil,
			storedToken: models.RefreshToken{
				ID: 1,
				UserID: 1,
				TokenHash: hash,
				ExpiresAt: time.Now().Add(time.Minute),
				CreatedAt: time.Now(),
			},
			
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "token expired",
			userID: 1,
			body: `{"refresh_token": "token"}`,

			mockGetToken: true,
			GetTokenErr: nil,
			storedToken: models.RefreshToken{
				ID: 1,
				UserID: 1,
				TokenHash: hash,
				ExpiresAt: time.Now().Add(-time.Second),
				CreatedAt: time.Now(),
			},

			deleteExpiredToken: true,

			wantStatus: http.StatusUnauthorized,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepo)
			handler := controllers.NewAuthHandler(mockRepo, "test_secret")

			if test.mockGetToken {
				mockRepo.On("GetRefreshTokenByTokenHash", hash).Return(test.storedToken, test.GetTokenErr)
				if test.deleteExpiredToken {
					mockRepo.On("DeleteRefreshTokenByTokenHash", hash).Return(nil)
				}
			}

			router := gin.New()
			router.POST("/logout/all-devices", func(c *gin.Context){
				if test.userID > 0 {
					c.Set("userID", test.userID)
				}
				handler.LogoutAllDevicesExceptThis(c)
			})
			
			req := httptest.NewRequest(http.MethodPost, "/logout/all-devices", strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, test.wantStatus, recorder.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
func TestLogoutAllDevicesExceptThis_InternalError(t *testing.T) {
	hash := repository.HashRefreshToken("token")
	tests := []struct{
		name string
		userID int
		body string

		mockGetToken bool
		GetTokenErr error
		storedToken models.RefreshToken

		deleteExpiredToken bool

		mockDeleteAllTokens bool
		DeleteAllTokenErr error

		wantStatus int
	}{
		{
			name: "Internal error GetRefreshTokenBYTokenHash",
			userID: 1,
			body: `{"refresh_token": "token"}`,

			mockGetToken: true,
			GetTokenErr: errors.New("internal error GetRefreshTokenBYTokenHash"),

			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal error DeleteAllUserRefreshTokensExceptThis",
			userID: 1,
			body: `{"refresh_token": "token"}`,		
			
			mockGetToken: true,
			GetTokenErr: nil,
			storedToken: models.RefreshToken{
				ID: 1,
				UserID: 1,
				TokenHash: hash,
				ExpiresAt: time.Now().Add(time.Minute),
				CreatedAt: time.Now(),
			},
			
			deleteExpiredToken: false,
			
			mockDeleteAllTokens: true,
			DeleteAllTokenErr: errors.New("Internal error DeleteAllUserRefreshTokensExceptThis"),

			wantStatus: http.StatusInternalServerError,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepo)
			handler := controllers.NewAuthHandler(mockRepo, "test_secret")

			if test.mockGetToken {
				mockRepo.On("GetRefreshTokenByTokenHash", hash).Return(test.storedToken, test.GetTokenErr)
				if test.deleteExpiredToken {
					mockRepo.On("DeleteRefreshTokenByTokenHash", hash).Return(nil)
				}
				if test.mockDeleteAllTokens {
					mockRepo.On("DeleteAllUserRefreshTokensExceptThis", test.storedToken.UserID, hash).Return(test.DeleteAllTokenErr)
				}
			}

			router := gin.New()
			router.POST("/logout/all-devices", func(c *gin.Context){
				if test.userID > 0 {
					c.Set("userID", test.userID)
				}
				handler.LogoutAllDevicesExceptThis(c)
			})
			
			req := httptest.NewRequest(http.MethodPost, "/logout/all-devices", strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, test.wantStatus, recorder.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}