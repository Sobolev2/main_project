package repository


import (
	"context"
	"semen_project/internal/models"
)
type RefreshTokenRepo interface {

	SaveRefreshToken(userID int, tokenHash string) error
	DeleteRefreshTokenById(tokenID int) error
	DeleteAllUserRefreshTokensExceptThis(userID int, tokenHash string) error
	GetRefreshTokenById(tokenID int) (models.RefreshToken, error)
	UpdateRefreshToken(r models.RefreshToken) error
	GetRefreshTokenByTokenHash(hashToken string) (models.RefreshToken, error)
	DeleteRefreshTokenByTokenHash(tokenHash string) error
}

func (s *Store) SaveRefreshToken(userID int, tokenHash string) error {
	query := `
	INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at)
	VALUES ($1, $2, NOW() + INTERVAL '30 days', NOW())
	`
	_, err := s.db.Exec(context.Background(), query, userID, tokenHash)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) DeleteRefreshTokenById(token_id int) error {
	query := `DELETE FROM refresh_tokens WHERE id = $1`
	_, err := s.db.Exec(context.Background(), query, token_id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) DeleteAllUserRefreshTokensExceptThis(user_id int, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1 AND token_hash != $2`
	_, err := s.db.Exec(context.Background(), query, user_id, tokenHash)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetRefreshTokenById(token_id int) (models.RefreshToken, error) {
	query := `
	SELECT id, user_id, token_hash, expires_at, created_at
	FROM refresh_tokens
	WHERE id = $1
	`
	var r models.RefreshToken
	err := s.db.QueryRow(context.Background(), query, token_id).Scan(
		&r.ID,
		&r.UserID,
		&r.TokenHash,
		&r.ExpiresAt,
		&r.CreatedAt,
	)
	if err != nil {
		return models.RefreshToken{}, err
	}
	return r, nil
}
func (s *Store) UpdateRefreshToken(r models.RefreshToken) error {
	query := `
	UPDATE refresh_tokens
	SET token_hash = $1, expires_at = $2
	WHERE id = $3
	`
	_, err := s.db.Exec(context.Background(), query, r.TokenHash, r.ExpiresAt, r.ID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetRefreshTokenByTokenHash(hashToken string) (models.RefreshToken, error) {
	query := `
	SELECT id, user_id, token_hash, expires_at, created_at
	FROM refresh_tokens
	WHERE token_hash = $1
	`
	var r models.RefreshToken
	err := s.db.QueryRow(context.Background(), query, hashToken).Scan(
		&r.ID,
		&r.UserID,
		&r.TokenHash,
		&r.ExpiresAt,
		&r.CreatedAt,
	)
	if err != nil {
		return models.RefreshToken{}, err
	}
	return r, nil
}
func (s *Store) DeleteRefreshTokenByTokenHash(token_hash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := s.db.Exec(context.Background(), query, token_hash)
	if err != nil {
		return err
	}
	return nil
}
