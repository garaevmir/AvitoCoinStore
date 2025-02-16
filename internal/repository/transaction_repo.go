package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/garaevmir/avitocoinstore/internal/model"
)

type TransactionRepositoryInt interface {
	TransferCoins(ctx context.Context, fromUserID, toUserID string, amount int) error
	GetTransactionHistory(ctx context.Context, userID string) (*model.TransactionHistory, error)
}

type TransactionRepository struct {
	pool DB
}

func NewTransactionRepository(db DB) *TransactionRepository {
	return &TransactionRepository{pool: db}
}

func (r TransactionRepository) TransferCoins(ctx context.Context, fromUserID, toUserID string, amount int) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Transaction error: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	var balance int
	err = tx.QueryRow(ctx, "SELECT coins FROM users WHERE id = $1 FOR UPDATE", fromUserID).Scan(&balance)
	if err != nil {
		log.Printf("Database error: %v", err)
		return err
	}
	if balance < amount {
		log.Printf("Database error: %v", err)
		return model.ErrInsufficientFunds
	}

	batch := &pgx.Batch{}
	batch.Queue("UPDATE users SET coins = coins - $1 WHERE id = $2", amount, fromUserID)
	batch.Queue("UPDATE users SET coins = coins + $1 WHERE id = $2", amount, toUserID)
	batch.Queue("INSERT INTO transactions (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)", fromUserID, toUserID, amount)

	br := tx.SendBatch(ctx, batch)

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			log.Printf("Database error: %v", err)
			return err
		}
	}

	br.Close()

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Transaction commit error: %v", err)
		return err
	}

	return nil
}

func (r TransactionRepository) GetTransactionHistory(ctx context.Context, userID string) (*model.TransactionHistory, error) {
	history := &model.TransactionHistory{
		Received: make([]model.ReceivedTransaction, 0),
		Sent:     make([]model.SentTransaction, 0),
	}

	rows, err := r.pool.Query(ctx,
		`SELECT u.username, t.amount, t.created_at 
         FROM transactions t
         JOIN users u ON t.from_user_id = u.id
         WHERE t.to_user_id = $1`,
		userID,
	)
	if err != nil {
		log.Printf("Database error: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.ReceivedTransaction
		if err := rows.Scan(&t.FromUser, &t.Amount, &t.Timestamp); err != nil {
			log.Printf("Database error: %v", err)
			return nil, err
		}
		history.Received = append(history.Received, t)
	}

	rows, err = r.pool.Query(ctx,
		`SELECT u.username, t.amount, t.created_at 
         FROM transactions t
         JOIN users u ON t.to_user_id = u.id
         WHERE t.from_user_id = $1`,
		userID,
	)
	if err != nil {
		log.Printf("Database error: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.SentTransaction
		if err := rows.Scan(&t.ToUser, &t.Amount, &t.Timestamp); err != nil {
			log.Printf("Database error: %v", err)
			return nil, err
		}
		history.Sent = append(history.Sent, t)
	}

	return history, nil
}
