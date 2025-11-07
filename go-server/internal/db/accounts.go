package db

import (
	"compound/go-server/pkg/models"
	"context"
	"fmt"
)

// CreateOrUpdateAccount creates or updates an account in the database
func CreateOrUpdateAccount(ctx context.Context, itemID int, plaidAccountID, name, mask, accountType, subtype string) (*models.Account, error) {
	query := `INSERT INTO accounts_table (item_id, plaid_account_id, name, mask, type, subtype, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	          ON CONFLICT (plaid_account_id) DO UPDATE SET
	            name = EXCLUDED.name,
	            mask = EXCLUDED.mask,
	            type = EXCLUDED.type,
	            subtype = EXCLUDED.subtype,
	            updated_at = NOW()
	          RETURNING id, item_id, plaid_account_id, name, mask, type, subtype, created_at, updated_at`

	account := &models.Account{}
	err := conn.QueryRow(ctx, query, itemID, plaidAccountID, name, mask, accountType, subtype).Scan(
		&account.ID,
		&account.ItemID,
		&account.PlaidAccountID,
		&account.Name,
		&account.Mask,
		&account.Type,
		&account.Subtype,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return account, nil
}

// GetAccountsByItemID retrieves all accounts for a specific item
func GetAccountsByItemID(ctx context.Context, itemID int) ([]*models.Account, error) {
	query := `SELECT id, item_id, plaid_account_id, name, mask, type, subtype, created_at, updated_at
	          FROM accounts_table WHERE item_id=$1`

	rows, err := conn.Query(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var accounts []*models.Account
	for rows.Next() {
		account := &models.Account{}
		err := rows.Scan(
			&account.ID,
			&account.ItemID,
			&account.PlaidAccountID,
			&account.Name,
			&account.Mask,
			&account.Type,
			&account.Subtype,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return accounts, nil
}

// GetAccountByID retrieves a single account by ID
func GetAccountByID(ctx context.Context, accountID int) (*models.Account, error) {
	query := `SELECT id, item_id, plaid_account_id, name, mask, type, subtype, created_at, updated_at
	          FROM accounts_table WHERE id=$1`

	account := &models.Account{}
	err := conn.QueryRow(ctx, query, accountID).Scan(
		&account.ID,
		&account.ItemID,
		&account.PlaidAccountID,
		&account.Name,
		&account.Mask,
		&account.Type,
		&account.Subtype,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return account, nil
}

// DeleteAccount deletes an account from the database
func DeleteAccount(ctx context.Context, accountID int) error {
	query := `DELETE FROM accounts_table WHERE id=$1`

	result, err := conn.Exec(ctx, query, accountID)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// Query accounts_table by plaidAccountID, return Account object
func GetAccountByPlaidAccountID(ctx context.Context, plaidAccountID string) (*models.Account, error) {
	query := `SELECT id, item_id, plaid_account_id, name, mask, type, subtype, created_at, updated_at
            FROM accounts_table WHERE plaid_account_id=$1`

	account := &models.Account{}

	err := conn.QueryRow(ctx, query, plaidAccountID).Scan(
		&account.ID,
		&account.ItemID,
		&account.PlaidAccountID,
		&account.Name,
		&account.Mask,
		&account.Type,
		&account.Subtype,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return account, nil
}
