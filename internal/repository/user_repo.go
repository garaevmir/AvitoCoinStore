package repository

import (
	"context"
	"log"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInt interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateUserCoinsTx(ctx context.Context, tx pgx.Tx, userID string, delta int) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (username, password_hash, coins) 
         VALUES ($1, $2, $3)
		 RETURNING id`,
		user.Username, user.PasswordHash, user.Coins,
	).Scan(&user.ID)
	if err != nil {
		log.Printf("Error creting user: %v", err)
		return err
	}
	return nil
}

func (r UserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, coins 
         FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Printf("Database error: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r UserRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, coins 
         FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins)
	return &user, err
}

func (r UserRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (r UserRepository) UpdateUserCoinsTx(ctx context.Context, tx pgx.Tx, userID string, delta int) error {
	_, err := tx.Exec(ctx,
		"UPDATE users SET coins = coins + $1 WHERE id = $2",
		delta, userID,
	)
	return err
}
