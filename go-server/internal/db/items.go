package db

import (
	"compound/go-server/pkg/models"
	"context"
	"fmt"
)

// GetItemByID retrieves an item from the database by ID
func GetItemByID(ctx context.Context, id int) (*models.Item, error) {
	query := `SELECT id, user_id, plaid_access_token, plaid_item_id, plaid_institution_id,
	                 status, created_at, updated_at, transactions_cursor
	          FROM items WHERE id=$1`

	item := &models.Item{}
	err := conn.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.UserID,
		&item.PlaidAccessToken,
		&item.PlaidItemID,
		&item.PlaidInstitutionID,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.TransactionsCursor,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return item, nil
}

// GetItemsByUserID retrieves all items for a user
func GetItemsByUserID(ctx context.Context, userID int) ([]*models.Item, error) {
	query := `SELECT id, user_id, plaid_access_token, plaid_item_id, plaid_institution_id,
	                 status, created_at, updated_at, transactions_cursor
	          FROM items WHERE user_id=$1`

	rows, err := conn.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.PlaidAccessToken,
			&item.PlaidItemID,
			&item.PlaidInstitutionID,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.TransactionsCursor,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return items, nil
}

// CreateItem creates a new item in the database
func CreateItem(ctx context.Context, userID int, plaidAccessToken, plaidItemID, plaidInstitutionID, status string) (*models.Item, error) {
	query := `INSERT INTO items (user_id, plaid_access_token, plaid_item_id, plaid_institution_id, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	          RETURNING id, user_id, plaid_access_token, plaid_item_id, plaid_institution_id, status, created_at, updated_at, transactions_cursor`

	item := &models.Item{}
	err := conn.QueryRow(ctx, query, userID, plaidAccessToken, plaidItemID, plaidInstitutionID, status).Scan(
		&item.ID,
		&item.UserID,
		&item.PlaidAccessToken,
		&item.PlaidItemID,
		&item.PlaidInstitutionID,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.TransactionsCursor,
	)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return item, nil
}

// UpdateItemTransactionsCursor updates the transactions cursor for an item
func UpdateItemTransactionsCursor(ctx context.Context, itemID int, cursor string) error {
	query := `UPDATE items SET transactions_cursor=$1, updated_at=NOW() WHERE id=$2`

	result, err := conn.Exec(ctx, query, cursor, itemID)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

// UpdateItemStatus updates the status of an item
func UpdateItemStatus(ctx context.Context, itemID int, status string) error {
	query := `UPDATE items SET status=$1, updated_at=NOW() WHERE id=$2`

	result, err := conn.Exec(ctx, query, status, itemID)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

// DeleteItem deletes an item from the database
func DeleteItem(ctx context.Context, itemID int) error {
	query := `DELETE FROM items WHERE id=$1`

	result, err := conn.Exec(ctx, query, itemID)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}
