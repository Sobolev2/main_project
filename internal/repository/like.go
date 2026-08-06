package repository

import (
	"context"
	"semen_project/internal/models"

)
type LikeRepo interface {
	CreateLike(postId, userId int) error
	DeleteLike(postID, userID int) error
	GetAllUserLikes(userId int) ([]models.Like, error)
	GetAllPostLikes(postId int) ([]models.Like, error)
	GetLikeStatus(postID, userID int) (bool, error)
	GetCountLikes(postID int) (int, error)
	GetPostById(postId int) (models.Post, error)
	GetUserById(userId int) (models.UserPublic, error)
}
func (s *Store) CreateLike(post_id int, user_id int) error {
	query := `
	INSERT INTO likes (post_id, user_id)
	VALUES ($1, $2)
	`
	_, err := s.db.Exec(context.Background(), query, post_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) DeleteLike(postID int, userID int) error {
	query := `
	DELETE FROM likes
	WHERE post_id = $1 AND user_id = $2
	`

	_, err := s.db.Exec(context.Background(), query, postID, userID)
	return err
}
func (s *Store) GetAllUserLikes(user_id int) ([]models.Like, error) {
	var likes []models.Like
	query := `
	SELECT id, post_id, user_id
	FROM likes
	WHERE user_id = $1
	`
	rows, err := s.db.Query(context.Background(),query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var like models.Like
		err := rows.Scan(&like.ID, &like.PostID, &like.UserID)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}	
	return likes, nil
}
func (s *Store) GetAllPostLikes(post_id int) ([]models.Like, error)  {
	var likes []models.Like
	query := `
	SELECT id, post_id, user_id
	FROM likes
	WHERE post_id = $1
	`
	rows, err := s.db.Query(context.Background(),query, post_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var like models.Like
		err := rows.Scan(&like.ID, &like.PostID, &like.UserID)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}	
	return likes, nil
}
func (s *Store) GetLikeStatus(postID int, userID int) (bool, error) {
	query := `
	SELECT 1
	FROM likes
	WHERE post_id = $1 AND user_id = $2
	`

	var exists int
	err := s.db.QueryRow(context.Background(), query, postID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return true, nil
}
func (s *Store) GetCountLikes(postID int) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM likes
	WHERE post_id = $1
	`

	var count int
	err := s.db.QueryRow(context.Background(), query, postID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}