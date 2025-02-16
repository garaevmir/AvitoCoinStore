package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/garaevmir/avitocoinstore/internal/model"
)

type InventoryRepositoryInt interface {
	GetUserInventory(ctx context.Context, userID string) ([]model.InventoryItem, error)
	AddToInventoryTx(ctx context.Context, tx pgx.Tx, userID, item string, quantity int) error
}

type InventoryRepository struct {
	pool DB
}

func NewInventoryRepository(db DB) *InventoryRepository {
	return &InventoryRepository{pool: db}
}

func (r InventoryRepository) GetUserInventory(ctx context.Context, userID string) ([]model.InventoryItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT item_name, quantity 
         FROM inventory 
         WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.InventoryItem
	for rows.Next() {
		var item model.InventoryItem
		if err := rows.Scan(&item.Name, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r InventoryRepository) AddToInventoryTx(ctx context.Context, tx pgx.Tx, userID, item string, quantity int) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO inventory (user_id, item_name, quantity)
         VALUES ($1, $2, $3)
         ON CONFLICT (user_id, item_name) DO UPDATE
         SET quantity = inventory.quantity + excluded.quantity`,
		userID, item, quantity,
	)
	return err
}
