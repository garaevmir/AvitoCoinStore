package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/garaevmir/avitocoinstore/internal/model"
)

// Interface for user repository, needed for testing
type UserRepositoryInt interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateUserCoinsTx(ctx context.Context, tx pgx.Tx, userID string, delta int) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

// User repository for user manipulations
type UserRepository struct {
	pool DB
}

// Constructor for user repository
func NewUserRepository(db DB) *UserRepository {
	return &UserRepository{pool: db}
}

// Function that writes user to database and assigns userID, returns error
func (r UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (username, password_hash, coins) 
         VALUES ($1, $2, $3)
		 RETURNING id`,
		user.Username, user.PasswordHash, user.Coins,
	).Scan(&user.ID)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return err
	}
	return nil
}

// Extracts user by given username if there exists such a user returns it's data otherwise returns nil,
// return user and error
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

// Extracts user by given userID, return user and error
// Unlike GetUserByUsername this function presumes that user exists, due to the fact that there exists userID
func (r UserRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, coins 
         FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins)
	return &user, err
}

// Function that starts transaction
func (r UserRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.BeginTx(ctx, pgx.TxOptions{})
}

// Function that updates amount of money available for user with userID
func (r UserRepository) UpdateUserCoinsTx(ctx context.Context, tx pgx.Tx, userID string, delta int) error {
	_, err := tx.Exec(ctx,
		"UPDATE users SET coins = coins + $1 WHERE id = $2",
		delta, userID,
	)
	return err
}
