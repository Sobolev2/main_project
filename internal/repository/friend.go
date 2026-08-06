package repository

import (
	"context"
	"semen_project/internal/models"

)
type FriendRepo interface {
	CreateFriendship(userFirstId, userSecondId int) error
	DeleteFriend(userFirstId, userSecondId int) error
	GetAllUserFriends(userId int) ([]models.UserPublic, error)
	GetCountFriends(userID int) (int, error)
	GetFriendship(userFirstID, userSecondID int) (bool, error)
}

func (s *Store) CreateFriendship(user_first_id int, user_second_id int) error {
	if user_first_id > user_second_id {
		user_first_id, user_second_id = user_second_id, user_first_id
	}
	query := `
	INSERT INTO friends (user_first, user_second)
	VALUES ($1, $2)
	`
	_, err := s.db.Exec(context.Background(), query, user_first_id, user_second_id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) DeleteFriend(user_first_id int, user_second_id int) error {
	if user_first_id > user_second_id {
		user_first_id, user_second_id = user_second_id, user_first_id
	}
	query := `
	DELETE FROM friends WHERE (user_first = $1 AND user_second = $2)`
	_, err := s.db.Exec(context.Background(), query, user_first_id, user_second_id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetAllUserFriends(user_id int) ([]models.UserPublic, error) {
	var users []models.UserPublic
	query := `
		SELECT
	    u.id,
	    u.user_name,
	    u.first_name,
	    u.last_name
	FROM friends f
	JOIN users u
	ON u.id =
	CASE
	    WHEN f.user_first = $1 THEN f.user_second
	    ELSE f.user_first
	END
	WHERE f.user_first = $1
	   OR f.user_second = $1;
	`
	rows, err := s.db.Query(context.Background(),query, user_id)
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
func (s *Store) GetCountFriends(userID int) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM friends
	WHERE user_first = $1 OR user_second = $1;
	`
	var count int
	err := s.db.QueryRow(context.Background(), query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (s *Store) GetFriendship(userFirstID int, userSecondID int) (bool, error) {
	if userFirstID > userSecondID {
		userFirstID, userSecondID = userSecondID, userFirstID
	}
	query := `
	SELECT 1
	FROM friends
	WHERE	user_first = $1 AND user_second = $2
	`
	var exists int
	err := s.db.QueryRow(context.Background(), query, userFirstID, userSecondID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return true, nil
}