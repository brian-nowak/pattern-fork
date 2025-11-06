package models

import "time"

type Transaction struct {
	ID                     int       `db:"id" json:"id"`
	AccountID              int       `db:"account_id" json:"account_id"`
	PlaidTransactionID     string    `db:"plaid_transaction_id" json:"plaid_transaction_id"`
	PlaidCategoryID        *string   `db:"plaid_category_id" json:"plaid_category_id"`
	Category               *string   `db:"category" json:"category"`
	Type                   string    `db:"type" json:"type"`
	Name                   string    `db:"name" json:"name"`
	Amount                 float64   `db:"amount" json:"amount"`
	IsoCurrencyCode        *string   `db:"iso_currency_code" json:"iso_currency_code"`
	UnofficialCurrencyCode *string   `db:"unofficial_currency_code" json:"unofficial_currency_code"`
	Date                   time.Time `db:"date" json:"date"`
	Pending                bool      `db:"pending" json:"pending"`
	AccountOwner           *string   `db:"account_owner" json:"account_owner"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
}
