package models

import "time"

type Account struct {
	ID                      int       `db:"id" json:"id"`
	ItemID                  int       `db:"item_id" json:"item_id"`
	PlaidAccountID          string    `db:"plaid_account_id" json:"plaid_account_id"`
	Name                    string    `db:"name" json:"name"`
	Mask                    string    `db:"mask" json:"mask"`
	OfficialName            *string   `db:"official_name" json:"official_name"`
	CurrentBalance          *float64  `db:"current_balance" json:"current_balance"`
	AvailableBalance        *float64  `db:"available_balance" json:"available_balance"`
	IsoCurrencyCode         *string   `db:"iso_currency_code" json:"iso_currency_code"`
	UnofficialCurrencyCode  *string   `db:"unofficial_currency_code" json:"unofficial_currency_code"`
	Type                    string    `db:"type" json:"type"`
	Subtype                 string    `db:"subtype" json:"subtype"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at"`
}
