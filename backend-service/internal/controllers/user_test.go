package controllers_test

import (
	"semen_project/internal/controllers"
	"semen_project/internal/mocks"
	"semen_project/internal/models"
	"semen_project/internal/repository"
	"testing"

	"net/http"
	"net/http/httptest"

	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// GetMe
func TestGetMe(t *testing.T) {
	tests := []struct {
		name       string
		userID     bool
		mockUser   models.UserPublic
		mockError  error
		wantStatus int
	}{
		{
			name:       "Success",
			userID:     true,
			mockUser:   models.UserPublic{ID: 1, UserName: "test"},
			mockError:  nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "user not found",
			userID:     true,
			mockUser:   models.UserPublic{},
			mockError:  pgx.ErrNoRows,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "no user id",
			userID:     false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "database error",
			userID:     true,
			mockUser:   models.UserPublic{},
			mockError:  errors.New("database error"),
			wantStatus: http.StatusInternalServerError,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepo)
			handler := controllers.NewUserHandler(mockRepo)
			if test.userID {
				mockRepo.On("GetUserById", 1).Return(test.mockUser, test.mockError)
			}
			router := gin.New()
			router.GET("/me", func(ctx *gin.Context) {
				if test.userID {
					ctx.Set("userID", 1)
				}
				handler.GetMe(ctx)
			})
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(
				t,
				test.wantStatus,
				recorder.Code,
			)
			mockRepo.AssertExpectations(t)
			if !test.userID {
				mockRepo.AssertNotCalled(t, "GetUserById", mock.Anything)
			}
		})
	}
}

// UpdateUser
func TestUpdateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	lastUser := models.UserPublic{
		ID:        1,
		UserName:  "olduser",
		FirstName: "Old",
		LastName:  "User",
	}

	mockRepo := new(mocks.MockUserRepo)
	router := gin.New()
	handler := controllers.NewUserHandler(mockRepo)

	mockRepo.On("GetUserById", 1).Return(lastUser, nil)
	mockRepo.On("GetUserByUsernameExceptId", 1, "newuser").Return(models.UserPublic{}, pgx.ErrNoRows)
	mockRepo.On("UpdateUser", "newuser", "New", "User", 1).Return(nil)

	router.PUT("/update", func(c *gin.Context) {
		c.Set("userID", 1)
		handler.UpdateUser(c)
	})

	body := `{
		"user_name": "newuser",
		"first_name": "New",
		"last_name": "User"
	}`
	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)

}
func TestUpdateUser_BadRequest(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		userID bool

		// GetUserById
		lastUser   models.UserPublic
		getUserErr error

		// GetUserByUsernameExceptId
		checkUsername bool
		usedUser      models.UserPublic
		usedUserErr   error

		wantStatus int
	}{
		{
			name:          "bad request body",
			body:          "",
			checkUsername: false,
			wantStatus:    http.StatusBadRequest,
		},
		{
			name: "not found user",
			body: `{
				"user_name": "newuser",
				"first_name": "New",
				"last_name": "User"}`,
			userID: true,
			lastUser: models.UserPublic{
				ID:        1,
				UserName:  "lastUser",
				FirstName: "last",
				LastName:  "User"},
			getUserErr:    pgx.ErrNoRows,
			checkUsername: false,
			wantStatus:    http.StatusNotFound,
		},
		{
			name: "empty new username",
			body: `{
				"user_name": "",
				"first_name": "New",
				"last_name": "User"}`,
			userID: true,
			lastUser: models.UserPublic{
				ID:        1,
				UserName:  "lastUser",
				FirstName: "last",
				LastName:  "User"},
			checkUsername: false,
			wantStatus:    http.StatusBadRequest,
		},
		{
			name: "data not changed",
			body: `{
				"user_name": "lastUser",
				"first_name": "last",
				"last_name": "User"}`,
			userID: true,
			lastUser: models.UserPublic{
				ID:        1,
				UserName:  "lastUser",
				FirstName: "last",
				LastName:  "User"},
			checkUsername: false,
			wantStatus:    http.StatusBadRequest,
		},
		{
			name: "username already exists",
			body: `{
				"user_name": "newUser",
				"first_name": "new",
				"last_name": "User"}`,
			userID: true,
			lastUser: models.UserPublic{
				ID:        1,
				UserName:  "lastUser",
				FirstName: "last",
				LastName:  "User"},
			checkUsername: true,
			usedUser: models.UserPublic{
				ID:        2,
				UserName:  "someUser",
				FirstName: "some",
				LastName:  "User"},
			usedUserErr: nil,
			wantStatus:  http.StatusConflict,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepo)
			handlers := controllers.NewUserHandler(mockRepo)
			if test.userID {
				mockRepo.On("GetUserById", 1).Return(test.lastUser, test.getUserErr)
				if test.getUserErr == nil && test.checkUsername {
					mockRepo.On("GetUserByUsernameExceptId", 1, mock.Anything).Return(test.usedUser, test.usedUserErr)
				}
			}
			router := gin.New()
			router.PUT("/update", func(c *gin.Context) {
				if test.userID {
					c.Set("userID", 1)
				}
				handlers.UpdateUser(c)
			})
			req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(
				t,
				test.wantStatus,
				recorder.Code,
			)
			mockRepo.AssertExpectations(t)
			if !test.userID {
				mockRepo.AssertNotCalled(t, "GetUserById", 1)
			}
		})
	}
}
func TestUpdateUser_InternalServerError(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		userID bool
		// GetUserById
		lastUser   models.UserPublic
		getUserErr error

		// GetUserByUsernameExceptId
		usedUser    models.UserPublic
		usedUserErr error

		// UpdateUser
		updateUserErr error

		wantStatus int
	}{
		{
			name: "userID not found",
			body: `{
				"user_name": "newuser",
				"first_name": "New",
				"last_name": "User"}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "internal error GetUserByID",
			body: `{
				"user_name": "newuser",
				"first_name": "New",
				"last_name": "User"}`,
			userID:     true,
			lastUser:   models.UserPublic{},
			getUserErr: errors.New("database error"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "internal error GetUserByUsernameExceptId",
			body: `{
				"user_name": "newuser",
				"first_name": "New",
				"last_name": "User"}`,
			userID: true,
			lastUser: models.UserPublic{
				ID:        1,
				UserName:  "lastUser",
				FirstName: "last",
				LastName:  "User"},
			getUserErr:  nil,
			usedUser:    models.UserPublic{},
			usedUserErr: errors.New("database error"),
			wantStatus:  http.StatusInternalServerError,
		},
		{
			name: "internal error UpdateUser",
			body: `{
				"user_name": "newuser",
				"first_name": "New",
				"last_name": "User"}`,
			userID: true,
			lastUser: models.UserPublic{
				ID:        1,
				UserName:  "lastUser",
				FirstName: "last",
				LastName:  "User"},
			getUserErr:    nil,
			usedUser:      models.UserPublic{},
			usedUserErr:   pgx.ErrNoRows,
			updateUserErr: errors.New("database error"),
			wantStatus:    http.StatusInternalServerError,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepo)
			handlers := controllers.NewUserHandler(mockRepo)
			if test.userID {
				mockRepo.On("GetUserById", 1).Return(test.lastUser, test.getUserErr)
				if test.getUserErr == nil {
					mockRepo.On("GetUserByUsernameExceptId", 1, mock.Anything).Return(test.usedUser, test.usedUserErr)
					if test.usedUserErr == nil || test.usedUserErr == pgx.ErrNoRows {
						mockRepo.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything, 1).Return(test.updateUserErr)
					}
				}
			}
			router := gin.New()
			router.PUT("/update", func(c *gin.Context) {
				if test.userID {
					c.Set("userID", 1)
				}
				handlers.UpdateUser(c)
			})
			req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(test.body))
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

// GetAllUsers
func TestGetAllUSers(t *testing.T) {
	tests := []struct {
		name           string
		users          []models.UserPublic
		getAllUsersErr error
		wantStatus     int
	}{
		{
			name:           "Internal error GetAllUsers",
			getAllUsersErr: errors.New("Internal error GetAllUsers"),
			wantStatus:     http.StatusInternalServerError,
		},
		{
			name:       "success, empty users list",
			users:      []models.UserPublic{},
			wantStatus: http.StatusOK,
		},
		{
			name: "Success, users found",
			users: []models.UserPublic{
				{
					ID:        1,
					UserName:  "first_user_name",
					FirstName: "first_name",
					LastName:  "last_name",
				},
				{
					ID:        2,
					UserName:  "second_user_name",
					FirstName: "first_name",
					LastName:  "last_name",
				},
			},
			wantStatus: http.StatusOK,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepo)
			handler := controllers.NewUserHandler(mockRepo)

			mockRepo.On("GetAllUsers").Return([]models.UserPublic{{}}, test.getAllUsersErr)

			router := gin.New()

			router.GET("/users", handler.GetAllUsers)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
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

// DeleteUser
func TestDeleteUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.MockUserRepo)
	handler := controllers.NewUserHandler(mockRepo)

	hash, _ := repository.HashPassword("password")

	mockRepo.On("GetPasswordById", 1).Return(hash, nil)
	mockRepo.On("CheckPassword", hash, "password").Return(nil)
	mockRepo.On("DeleteUser", 1).Return(nil)

	router := gin.New()
	body := `{
	"password": "password"
	}`
	router.DELETE("/user", func(c *gin.Context) {
		c.Set("userID", 1)
		handler.DeleteUser(c)
	})
	req := httptest.NewRequest(http.MethodDelete, "/user", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
}
func TestDeleteUser_BadRequest(t *testing.T) {
	hash, _ := repository.HashPassword("password")
	tests := []struct {
		name   string
		body   string
		userID bool

		mockGetPassword bool
		GetPasswordErr  error
		lastPassword    string

		mockCheckPassword bool
		CheckPasswordErr  error

		wantStatus int
	}{
		{
			name:       "Bad request",
			body:       "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Bad request, empty password",
			body:       `{"password": ""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:            "User not found",
			body:            `{"password": "password"}`,
			mockGetPassword: true,
			GetPasswordErr:  pgx.ErrNoRows,
			wantStatus:      http.StatusNotFound,
		},
		{
			name:              "Password is not confirm",
			body:              `{"password": "password"}`,
			mockGetPassword:   true,
			GetPasswordErr:    nil,
			lastPassword:      hash,
			mockCheckPassword: true,
			CheckPasswordErr:  errors.New("Password is not confirm"),
			wantStatus:        http.StatusUnauthorized,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepo)
			handler := controllers.NewUserHandler(mockRepo)

			if test.mockGetPassword {
				mockRepo.On("GetPasswordById", 1).Return(test.lastPassword, test.GetPasswordErr)
				if test.mockCheckPassword {
					mockRepo.On("CheckPassword", test.lastPassword, "password").Return(test.CheckPasswordErr)
				}
			}

			router := gin.New()
			router.DELETE("/user", func(c *gin.Context) {
				c.Set("userID", 1)
				handler.DeleteUser(c)
			})

			req := httptest.NewRequest(http.MethodDelete, "/user", strings.NewReader(test.body))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.wantStatus, recorder.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestDeleteUser_InternalError(t *testing.T) {
	hash, _ := repository.HashPassword("password")
	tests := []struct {
		name   string
		body   string
		userID bool

		mockGetPassword bool
		GetPasswordErr  error
		lastPassword    string

		mockCheckPassword bool

		mockDeleteUser bool
		DeleteUserErr  error

		wantStatus int
	}{
		{
			name: "UserID не найден в контексте",
			body: `{"password": "password"}`,

			wantStatus: http.StatusInternalServerError,
		},
		{
			name:   "Internal error GetPasswordById",
			body:   `{"password": "password"}`,
			userID: true,

			mockGetPassword: true,
			GetPasswordErr:  errors.New("Database error"),

			wantStatus: http.StatusInternalServerError,
		},
		{
			name:   "Internal error GetPasswordById",
			body:   `{"password": "password"}`,
			userID: true,

			mockGetPassword: true,
			lastPassword:    hash,

			mockCheckPassword: true,

			mockDeleteUser: true,
			DeleteUserErr:  errors.New("database error"),

			wantStatus: http.StatusInternalServerError,
		},
	}
	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepo)
			handler := controllers.NewUserHandler(mockRepo)

			if test.mockGetPassword {
				mockRepo.On("GetPasswordById", 1).Return(test.lastPassword, test.GetPasswordErr)
				if test.mockCheckPassword {
					mockRepo.On("CheckPassword", hash, "password").Return(nil)
					if test.mockDeleteUser {
						mockRepo.On("DeleteUser", 1).Return(test.DeleteUserErr)
					}
				}
			}

			router := gin.New()
			router.DELETE("/user", func(c *gin.Context) {
				if test.userID {
					c.Set("userID", 1)
				}
				handler.DeleteUser(c)
			})

			req := httptest.NewRequest(http.MethodDelete, "/user", strings.NewReader(test.body))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.wantStatus, recorder.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
