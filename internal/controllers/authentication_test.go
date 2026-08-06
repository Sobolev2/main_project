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
		name   string
		body   string
		mockGetUser bool
		// GetUserByUsername
		usedUser    models.UserPublic
		usedUserErr error

		wantStatus int
	}{
		{
			name: "bad request",
			body: "",
			mockGetUser: false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty username or password",
			body: `{"user_name":"","first_name":"Test","last_name":"User","password":""}`,
			mockGetUser: false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "password too short",
			body: `{"user_name":"testuser","first_name":"Test","last_name":"User","password":"short"}`,
			mockGetUser: false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "username already exists",
			body: `{"user_name":"testuser","first_name":"Test","last_name":"User","password":"password123"}`,
			mockGetUser: true,
			usedUser: models.UserPublic{ID: 1, UserName: "testuser"},
			usedUserErr: nil,
			wantStatus: http.StatusConflict,
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
		name   string
		body   string
	
		// GetUserByUsername
		usedUserErr error


		// CreateUser
		createUserErr error

		wantStatus int
	}{
		{
			name: "internal error GetUserByUsername",
			body: `{"user_name":"username","first_name":"Test","last_name":"User","password":"password123"}`,
			usedUserErr: errors.New("internal error GetUserByUsername"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "internal error CreateUser",
			body: `{"user_name":"userName","first_name":"Test","last_name":"User","password":"password123"}`,
			usedUserErr: pgx.ErrNoRows,
			createUserErr: errors.New("internal error CreateUser"),
			wantStatus: http.StatusInternalServerError,
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