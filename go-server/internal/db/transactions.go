package db

import (
	"compound/go-server/pkg/models"
	"context"
	"fmt"
)

// CreateOrUpdateTransaction creates or updates a transaction in the database
func CreateOrUpdateTransaction(ctx context.Context, accountID int, plaidTransactionID, plaidCategoryID, category, txType, name string, amount float64, isoCurrencyCode, unofficialCurrencyCode string, date string, pending bool, accountOwner *string) (*models.Transaction, error) {
	query := `INSERT INTO transactions_table (account_id, plaid_transaction_id, plaid_category_id, category, type, name, amount, iso_currency_code, unofficial_currency_code, date, pending, account_owner, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
	          ON CONFLICT (plaid_transaction_id) DO UPDATE SET
	            category = EXCLUDED.category,
	            plaid_category_id = EXCLUDED.plaid_category_id,
	            type = EXCLUDED.type,
	            name = EXCLUDED.name,
	            amount = EXCLUDED.amount,
	            iso_currency_code = EXCLUDED.iso_currency_code,
	            unofficial_currency_code = EXCLUDED.unofficial_currency_code,
	            pending = EXCLUDED.pending,
	            account_owner = EXCLUDED.account_owner,
	            updated_at = NOW()
	          RETURNING id, account_id, plaid_transaction_id, plaid_category_id, category, type, name, amount, iso_currency_code, unofficial_currency_code, date, pending, account_owner, created_at, updated_at`

	transaction := &models.Transaction{}
	err := conn.QueryRow(ctx, query, accountID, plaidTransactionID, plaidCategoryID, category, txType, name, amount, isoCurrencyCode, unofficialCurrencyCode, date, pending, accountOwner).Scan(
		&transaction.ID,
		&transaction.AccountID,
		&transaction.PlaidTransactionID,
		&transaction.PlaidCategoryID,
		&transaction.Category,
		&transaction.Type,
		&transaction.Name,
		&transaction.Amount,
		&transaction.IsoCurrencyCode,
		&transaction.UnofficialCurrencyCode,
		&transaction.Date,
		&transaction.Pending,
		&transaction.AccountOwner,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return transaction, nil
}

// GetTransactionsByAccountID retrieves all transactions for a specific account
func GetTransactionsByAccountID(ctx context.Context, accountID int) ([]*models.Transaction, error) {
	query := `SELECT id, account_id, plaid_transaction_id, plaid_category_id, category, type, name, amount, iso_currency_code, unofficial_currency_code, date, pending, account_owner, created_at, updated_at
	          FROM transactions_table WHERE account_id=$1`

	rows, err := conn.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.PlaidTransactionID,
			&transaction.PlaidCategoryID,
			&transaction.Category,
			&transaction.Type,
			&transaction.Name,
			&transaction.Amount,
			&transaction.IsoCurrencyCode,
			&transaction.UnofficialCurrencyCode,
			&transaction.Date,
			&transaction.Pending,
			&transaction.AccountOwner,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return transactions, nil
}

// GetTransactionByID retrieves a single transaction by ID
func GetTransactionByID(ctx context.Context, transactionID int) (*models.Transaction, error) {
	query := `SELECT id, account_id, plaid_transaction_id, plaid_category_id, category, type, name, amount, iso_currency_code, unofficial_currency_code, date, pending, account_owner, created_at, updated_at
	          FROM transactions_table WHERE id=$1`

	transaction := &models.Transaction{}
	err := conn.QueryRow(ctx, query, transactionID).Scan(
		&transaction.ID,
		&transaction.AccountID,
		&transaction.PlaidTransactionID,
		&transaction.PlaidCategoryID,
		&transaction.Category,
		&transaction.Type,
		&transaction.Name,
		&transaction.Amount,
		&transaction.IsoCurrencyCode,
		&transaction.UnofficialCurrencyCode,
		&transaction.Date,
		&transaction.Pending,
		&transaction.AccountOwner,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return transaction, nil
}

// DeleteTransaction deletes a transaction from the database
func DeleteTransaction(ctx context.Context, transactionID int) error {
	query := `DELETE FROM transactions_table WHERE id=$1`

	result, err := conn.Exec(ctx, query, transactionID)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}
