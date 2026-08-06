package repository

import (
	"context"
	"semen_project/internal/models"

)
type CommentRepo interface {
	CreateComment(postID, userID int, content string) error
	DeleteComment(commentId int) error
	UpdateComment(commentId int, content string) error
	GetAllPostComments(postId int) ([]models.Comment, error)
	GetCommentById(commentId int) (models.Comment, error)
	GetAllUserCommentsUnderCurrentPost(postID, userID int) ([]models.Comment, error)
	GetCountComments(postID int) (int, error)
	GetPostById(postId int) (models.Post, error)
}
func (s *Store) CreateComment(postID int, userID int, content string) error{
	query := `
	INSERT INTO comments (post_id, user_id, content)
	VALUES ($1, $2, $3)
	`
	_, err := s.db.Exec(context.Background(),query, postID, userID, content)
	if err != nil {
		return  err
	}
	return nil
}
func (s *Store) DeleteComment(comment_id int) error {	
	query := `
	DELETE FROM comments WHERE id = $1
	`
	_, err := s.db.Exec(context.Background(),query, comment_id)
	if err != nil {
		return  err
	}
	return nil
}
func (s *Store) UpdateComment(comment_id int,content string) error {
	query := `
	UPDATE comments
	SET content = $2
	WHERE id = $1
	`
	_, err := s.db.Exec(context.Background(),query,comment_id, content)
	if err != nil {
		return  err
	}
	return nil
}
func (s *Store) GetAllPostComments(post_id int) ([]models.Comment, error) {
	var comments []models.Comment
	query := `
	SELECT id, post_id, user_id, content
	FROM comments
	WHERE post_id = $1
	ORDER By id DESC
	`
	rows, err := s.db.Query(context.Background(),query, post_id)
	if err != nil {
		return  nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}	
	return comments, nil
}
func (s *Store) GetCommentById(comment_id int) (models.Comment, error) {
	query := `
	SELECT id, post_id, user_id, content
	FROM comments
	WHERE id = $1
	`
	var comment models.Comment
	err := s.db.QueryRow(context.Background(), query, comment_id).Scan( 
	&comment.ID,
    &comment.PostID,
    &comment.UserID,
    &comment.Content,)
	if err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}
func (s *Store) GetAllUserCommentsUnderCurrentPost(postID int, userID int) ([]models.Comment, error) {
	var comments []models.Comment
	query := `
	SELECT id, post_id, user_id, content
	FROM comments
	WHERE post_id = $1 and user_id = $2
	ORDER BY id DESC
	`
	rows, err := s.db.Query(context.Background(),query, postID, userID)
	if err != nil {
		return  nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}	
	return comments, nil
}
func (s *Store) GetCountComments(postID int) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM comments
	WHERE post_id = $1
	`
	var count int
	err := s.db.QueryRow(context.Background(), query, postID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
