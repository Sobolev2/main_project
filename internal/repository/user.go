package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"semen_project/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(username, firstName, lastName, password string) (*models.UserPublic, error)
	GetUserById(id int) (models.UserPublic, error)
	GetPasswordById(id int) (string, error)
	GetAllUsersExceptId(id int) ([]models.UserPublic, error)
	GetUserByUsernameExceptId(id int, username string) (models.UserPublic, error)
	GetUserByUsername(username string) (models.UserPublic, error)
	GetAllUsers() ([]models.UserPublic, error)
	UpdateUser(userName, firstName, lastName string, id int) error
	UpdatePassword(id int, password string) error
	DeleteUser(id int) error
	SearchUsers(query string) ([]models.UserPublic, error)
	CheckPassword(hash string, password string) error
}

type Store struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewStore(db *pgxpool.Pool, redisClient *redis.Client) *Store {
	return &Store{
		db:    db,
		redis: redisClient,
	}
}
func (s *Store) GetUserById(id int) (models.UserPublic, error) {
	ctx := context.Background()

	key := fmt.Sprintf("user:%d", id)

	cachedUser, err := s.redis.Get(ctx, key).Result()

	if err == nil {
		var user models.UserPublic

		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			return models.UserPublic{}, err
		}

		return user, nil
	}

	if err != redis.Nil {
		return models.UserPublic{}, err
	}

	var user models.UserPublic

	query := `
		SELECT id, user_name, first_name, last_name
		FROM users
		WHERE id = $1
	`

	err = s.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
	)

	if err != nil {
		return models.UserPublic{}, err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return models.UserPublic{}, err
	}

	err = s.redis.Set(
		ctx,
		key,
		data,
		5*time.Minute,
	).Err()

	if err != nil {
		return models.UserPublic{}, err
	}

	return user, nil
}
func (s *Store) CreateUser(username string, firstName string, lastName string, password string) (*models.UserPublic, error) {
	query := `
	INSERT INTO users (user_name, first_name, last_name, password)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_name, first_name, last_name;
	`
	createdUser := &models.UserPublic{}
	err := s.db.QueryRow(
		context.Background(),
		query,
		username, firstName, lastName, password,
	).Scan(
		&createdUser.ID,
		&createdUser.UserName,
		&createdUser.FirstName,
		&createdUser.LastName,
	)
	if err != nil {
		return &models.UserPublic{}, err
	}
	return createdUser, nil
}
func (s *Store) GetPasswordById(id int) (string, error) {
	var password string
	query := `SELECT password FROM users WHERE id = $1`
	err := s.db.QueryRow(context.Background(), query, id).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}
func (s *Store) GetAllUsersExceptId(id int) ([]models.UserPublic, error) {
	var users []models.UserPublic
	query := `SELECT id, user_name, first_name, last_name FROM users WHERE id != $1`
	rows, err := s.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.UserPublic
		err := rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
func (s *Store) GetUserByUsernameExceptId(id int, username string) (models.UserPublic, error) {
	var user models.UserPublic
	query := `SELECT id, user_name, first_name, last_name
	FROM users
	WHERE user_name = $1 AND id != $2`
	err := s.db.QueryRow(context.Background(), query, username, id).Scan(
		&user.ID,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return models.UserPublic{}, err
	}
	return user, nil
}
func (s *Store) GetUserByUsername(username string) (models.UserPublic, error) {
	var user models.UserPublic

	query := `
	SELECT id, user_name, first_name, last_name
	FROM users
	WHERE user_name = $1
	`

	err := s.db.QueryRow(context.Background(), query, username).Scan(
		&user.ID,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return models.UserPublic{}, err
	}

	return user, nil
}
func (s *Store) GetAllUsers() ([]models.UserPublic, error) {
	var users []models.UserPublic
	query := `SELECT id, user_name, first_name, last_name FROM users`
	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.UserPublic
		err := rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
func (s *Store) UpdateUser(userName string, firstName string, lastName string, id int) error {
	query := `UPDATE users 
	SET user_name = $2, first_name = $3, last_name = $4
	WHERE id = $1 
	`
	_, err := s.db.Exec(context.Background(), query, id, userName, firstName, lastName)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) UpdatePassword(id int, password string) error {
	query := `
	UPDATE users
	SET password = $2
	WHERE id = $1
	`
	_, err := s.db.Exec(context.Background(), query, id, password)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) SearchUsers(query string) ([]models.UserPublic, error) {
	var users []models.UserPublic

	querytext := `
	SELECT id, user_name, first_name, last_name
	FROM users
	WHERE
		user_name ILIKE '%' || $1 || '%'
		OR first_name ILIKE '%' || $1 || '%'
		OR last_name ILIKE '%' || $1 || '%'
		OR (first_name || ' ' || last_name) ILIKE '%' || $1 || '%'
	ORDER BY id
	LIMIT 50;
	`

	rows, err := s.db.Query(context.Background(), querytext, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserPublic

		err := rows.Scan(
			&user.ID,
			&user.UserName,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func (s *Store) CheckPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
