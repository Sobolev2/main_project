package controllers_test

import (
	"semen_project/internal/controllers"
	"semen_project/internal/mocks"
	"semen_project/internal/models"
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
		usedUser    models.UserPublic
		usedUserErr error

		wantStatus int
	}{
		{
			name:       "bad request body",
			body:       "",
			checkUsername: false,
			wantStatus: http.StatusBadRequest,
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
			getUserErr: pgx.ErrNoRows,
			checkUsername: false,
			wantStatus: http.StatusNotFound,
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
			wantStatus: http.StatusBadRequest,
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
			wantStatus: http.StatusBadRequest,
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
