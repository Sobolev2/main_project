package repository

import (
	"context"
		"semen_project/internal/models"

)
type FollowRepo interface {
	FollowUser(followerId, userId int) error
	UnFollowUser(followerId, userId int) error
	GetAllUserFollowers(userId int) ([]models.UserPublic, error)
	GetAllUserFollowing(userId int) ([]models.UserPublic, error)
	GetFollowStatus(userID, followerID int) (bool, error)
	GetCountFollowing(userID int) (int, error)
	GetCountFollowers(userID int) (int, error)
	GetUserById(userId int) (models.UserPublic, error)
	CreateFriendship(userFirstId, userSecondId int) error
	GetFriendship(userFirstID, userSecondID int) (bool, error)
	DeleteFriend(userFirstId, userSecondId int) error
}
func (s *Store) FollowUser(follower_id int, user_id int) error{
	query := `
	INSERT INTO followers (follower_id, user_id)
	VALUES ($1, $2)
	`
	_, err := s.db.Exec(context.Background(), query, follower_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) UnFollowUser(follower_id int, user_id int) error {
	query := `
	DELETE FROM followers WHERE follower_id = $1 AND user_id = $2`
	_, err := s.db.Exec(context.Background(), query, follower_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetAllUserFollowers(user_id int) ([]models.UserPublic, error) {
	var followers []models.UserPublic
	query := `
	SELECT u.id, u.user_name, u.first_name, u.last_name
	FROM followers f
	JOIN users u ON f.follower_id = u.id
	WHERE f.user_id = $1
	`
	rows, err := s.db.Query(context.Background(), query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var follower models.UserPublic
		err := rows.Scan(&follower.ID, &follower.UserName, &follower.FirstName, &follower.LastName)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}	
	return followers, nil
}
func (s *Store) GetAllUserFollowing(user_id int) ([]models.UserPublic, error) {
	var following_users []models.UserPublic
	query := `
	SELECT u.id, u.user_name, u.first_name, u.last_name
	FROM followers f
	JOIN users u ON f.user_id = u.id
	WHERE f.follower_id = $1
	`
	rows, err := s.db.Query(context.Background(), query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var follower models.UserPublic
		err := rows.Scan(&follower.ID, &follower.UserName, &follower.FirstName, &follower.LastName)
		if err != nil {
			return nil, err
		}
		following_users = append(following_users, follower)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}	
	return following_users, nil
}
func (s *Store) GetFollowStatus(userID int, followerID int) (bool, error) {
	query := `
	SELECT 1
	FROM followers
	WHERE user_id = $1 AND follower_id = $2
	`
	var exists int
	err := s.db.QueryRow(context.Background(), query, userID, followerID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s *Store) GetCountFollowing(userID int) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM followers
	WHERE follower_id = $1
	`
	var count int
	err := s.db.QueryRow(context.Background(), query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (s *Store) GetCountFollowers(userID int) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM followers
	WHERE user_id = $1
	`
	var count int
	err := s.db.QueryRow(context.Background(), query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
