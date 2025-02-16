package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/garaevmir/avitocoinstore/internal/model"
)

// Interface for inventory repository, needed for testing
type InventoryRepositoryInt interface {
	GetUserInventory(ctx context.Context, userID string) ([]model.InventoryItem, error)
	AddToInventoryTx(ctx context.Context, tx pgx.Tx, userID, item string, quantity int) error
}

// Inventory repository for inventory manipulations
type InventoryRepository struct {
	pool DB
}

// Constructor for inventory repository
func NewInventoryRepository(db DB) *InventoryRepository {
	return &InventoryRepository{pool: db}
}

// Function to extract user inventory by userID from database, returns slice of InventoryItem and error
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

// Function to add item to users inventory by userID on database during a transaction, returns error
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
