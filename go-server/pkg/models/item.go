package models

import "time"

type Item struct {
	ID                 int       `db:"id" json:"id"`
	UserID             int       `db:"user_id" json:"user_id"`
	PlaidAccessToken   string    `db:"plaid_access_token" json:"plaid_access_token"`
	PlaidItemID        string    `db:"plaid_item_id" json:"plaid_item_id"`
	PlaidInstitutionID string    `db:"plaid_institution_id" json:"plaid_institution_id"`
	Status             string    `db:"status" json:"status"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
	TransactionsCursor *string   `db:"transactions_cursor" json:"transactions_cursor"`
}
