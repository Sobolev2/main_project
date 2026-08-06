package mocks

import (
	"semen_project/internal/models"

	"semen_project/internal/repository"

	"github.com/stretchr/testify/mock"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) GetUserByUsername(username string) (models.UserPublic, error) {
	args := m.Called(username)
	user, _ := args.Get(0).(models.UserPublic)
	return user, args.Error(1)
}
func (m *MockAuthRepo) GetPasswordById(id int) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}
func (m *MockAuthRepo) CreateUser(username, firstName, lastName, password string) (*models.UserPublic, error) {
	args := m.Called(username, firstName, lastName, password)
	user, _ := args.Get(0).(*models.UserPublic)
	return user, args.Error(1)
}
func (m *MockAuthRepo) DeleteAllUserRefreshTokensExceptThis(userID int, tokenHash string) error {
	args := m.Called(userID, tokenHash)
	return args.Error(0)
}
func (m *MockAuthRepo) SaveRefreshToken(userID int, tokenHash string) error {
	args := m.Called(userID, tokenHash)
	return args.Error(0)
}
func (m *MockAuthRepo) GetRefreshTokenByTokenHash(hash string) (models.RefreshToken, error) {
	args := m.Called(hash)
	token, _ := args.Get(0).(models.RefreshToken)
	return token, args.Error(1)
}
func (m *MockAuthRepo) DeleteRefreshTokenByTokenHash(tokenHash string) error {
	args := m.Called(tokenHash)
	return args.Error(0)
}

var _ repository.AuthRepo = (*MockAuthRepo)(nil)
