package service

import (
	"context"
	"log"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/internal/repository"
)

// Structure representing shop
type ShopService struct {
	userRepo        repository.UserRepositoryInt
	transactionRepo repository.TransactionRepositoryInt
	inventoryRepo   repository.InventoryRepositoryInt
}

// Constructor for the shop
func NewShopService(
	uRepo repository.UserRepositoryInt,
	tRepo repository.TransactionRepositoryInt,
	iRepo repository.InventoryRepositoryInt,
) *ShopService {
	return &ShopService{
		userRepo:        uRepo,
		transactionRepo: tRepo,
		inventoryRepo:   iRepo,
	}
}

// Function that buys item itemName for user with userID during transaction, returns error
func (s *ShopService) BuyItem(ctx context.Context, userID string, itemName string) error {
	item, exists := model.Items[itemName]
	if !exists {
		return model.ErrItemNotFound
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return err
	}

	if user.Coins < item.Price {
		return model.ErrInsufficientFunds
	}

	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		log.Printf("Transaction error: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.userRepo.UpdateUserCoinsTx(ctx, tx, userID, -item.Price); err != nil {
		log.Printf("Updating user coins error: %v", err)
		return err
	}

	if err := s.inventoryRepo.AddToInventoryTx(ctx, tx, userID, itemName, 1); err != nil {
		log.Printf("Adding item to invetory error: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Transaction commit error: %v", err)
		return err
	}

	return nil
}
