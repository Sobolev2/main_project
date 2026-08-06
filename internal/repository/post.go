package repository

import (
	"context"
	"semen_project/internal/models"

)
type PostRepo interface {
	CreatePost(userId int, content string) (*models.Post, error)
	GetPostById(postId int) (models.Post, error)
	DeletePost(postId int) error
	GetAllPosts() ([]models.Post, error)
	GetAllUserPosts(userId int) ([]models.Post, error)
	UpdatePost(postId int, content string) (*models.Post, error)
	GetAllFriendsPosts(userId int) ([]models.Post, error)
	GetAllNotFriendsPosts(userId int) ([]models.Post, error)
}
func (s *Store) CreatePost(userId int, content string) (*models.Post, error)  {
	post := &models.Post{}
	query := ` INSERT INTO posts (user_id, content)
	VALUES ($1, $2)
	RETURNING id, user_id, content;
	`
	err := s.db.QueryRow(context.Background(), query, userId, content).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
	)
	if err != nil {
		return nil, err
	}
	return post, nil 
}
func (s *Store) GetPostById(post_id int) (models.Post, error)  {
	var post models.Post
	query := `SELECT id, user_id, content FROM posts WHERE id = $1 `
	err := s.db.QueryRow(context.Background(), query, post_id).Scan(&post.ID, &post.UserID, &post.Content)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}
func (s *Store) DeletePost(post_id int) error {
	query := `DELETE FROM posts WHERE id = $1 `
	_, err := s.db.Exec(context.Background(), query, post_id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetAllPosts() ([]models.Post, error)  {
	var posts []models.Post
	query := ` SELECT Id, user_id, content From posts ORDER BY id DESC;`
	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content)
		if err != nil {
			return nil, err
		}
	posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}	
	return posts, nil
}
func (s *Store) GetAllUserPosts(user_id int) ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT id, user_id, content FROM posts WHERE user_id = $1 ORDER BY id DESC;`
	rows, err := s.db.Query(context.Background(), query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}	
	return posts, nil
}
func (s *Store) UpdatePost(post_id int, content string) (*models.Post, error)  {
	updatedpost := &models.Post{}
	query := `
	UPDATE posts
	SET content = $2
	WHERE id = $1
	RETURNING id, user_id, content
	`
	err := s.db.QueryRow(context.Background(), query, post_id, content).Scan(
		&updatedpost.ID,
		&updatedpost.UserID,
		&updatedpost.Content,
	)
	if err != nil {
		return nil, err
	}
	return updatedpost, nil
}
// for feed
func (s *Store) GetAllFriendsPosts(user_id int) ([]models.Post, error) {
	var posts []models.Post
	query := `
	SELECT p.id, p.user_id, p.content
	FROM posts p
	JOIN friends f
	ON (
 	   (f.user_first = $1 AND p.user_id = f.user_second)
  	  OR
  	  (f.user_second = $1 AND p.user_id = f.user_first)
	)
	ORDER BY p.id DESC;	
	`
	rows, err := s.db.Query(context.Background(), query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}	
	return posts, nil
}
func (s *Store) GetAllNotFriendsPosts(user_id int) ([]models.Post, error) {
	var posts []models.Post
	query := `
	SELECT p.id, p.user_id, p.content
	FROM posts p
	WHERE p.user_id <> $1
	AND NOT EXISTS (
	    SELECT 1
	    FROM friends f
	    WHERE
	        (f.user_first = $1 AND f.user_second = p.user_id)
	        OR
	        (f.user_second = $1 AND f.user_first = p.user_id)
	)
	ORDER BY p.id DESC;
	`
	rows, err := s.db.Query(context.Background(), query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
    return nil, err
	}	
	return posts, nil
}