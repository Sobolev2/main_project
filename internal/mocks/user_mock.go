package mocks

import (
	"semen_project/internal/models"

	"semen_project/internal/repository"

	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(username string, firstName string, lastName string, password string) (*models.UserPublic, error) {
	args := m.Called(username, firstName, lastName, password)
	user, _ := args.Get(0).(*models.UserPublic)
	return user, args.Error(1)
}
func (m *MockUserRepo) GetUserById(id int) (models.UserPublic, error) {
	args := m.Called(id)
	user, _ := args.Get(0).(models.UserPublic)
	return user, args.Error(1)
}
func (m *MockUserRepo) GetPasswordById(id int) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}
func (m *MockUserRepo) GetAllUsersExceptId(id int) ([]models.UserPublic, error) {
	args := m.Called(id)
	users, _ := args.Get(0).([]models.UserPublic)
	return users, args.Error(1)
}
func (m *MockUserRepo) GetUserByUsernameExceptId(id int, username string) (models.UserPublic, error) {
	args := m.Called(id, username)
	user, _ := args.Get(0).(models.UserPublic)
	return user, args.Error(1)
}
func (m *MockUserRepo) GetUserByUsername(username string) (models.UserPublic, error) {
	args := m.Called(username)
	user, _ := args.Get(0).(models.UserPublic)
	return user, args.Error(1)
}

func (m *MockUserRepo) GetAllUsers() ([]models.UserPublic, error) {
	args := m.Called()
	users, _ := args.Get(0).([]models.UserPublic)
	return users, args.Error(1)
}
func (m *MockUserRepo) UpdateUser(userName, firstName, lastName string, id int) error {
	args := m.Called(userName, firstName, lastName, id)
	return args.Error(0)
}
func (m *MockUserRepo) UpdatePassword(id int, password string) error {
	args := m.Called(id, password)
	return args.Error(0)
}
func (m *MockUserRepo) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockUserRepo) SearchUsers(query string) ([]models.UserPublic, error) {
	args := m.Called(query)
	users, _ := args.Get(0).([]models.UserPublic)
	return users, args.Error(1)
}

var _ repository.UserRepo = (*MockUserRepo)(nil)
