package handlers

import (
	"compound/go-server/internal/db"
	plaidpkg "compound/go-server/internal/plaid"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExchangeTokenRequest represents the request body for exchanging a public token
type ExchangeTokenRequest struct {
	PublicToken string `json:"publicToken" binding:"required"`
	UserID      int    `json:"userId" binding:"required"`
}

// ExchangeToken handles POST /api/items
// Exchanges a public token from Plaid Link for an access token and stores the item
//
// Request body:
// {
//   "publicToken": "public-sandbox-...",
//   "userId": 123
// }
//
// Response:
// {
//   "item_id": "plaid-item-id",
//   "access_token": "access-sandbox-...",
//   "accounts": [
//     {
//       "id": "account-id",
//       "name": "Checking Account",
//       "mask": "1234",
//       "type": "depository",
//       "subtype": "checking"
//     }
//   ]
// }
func ExchangeToken(c *gin.Context) {
	var req ExchangeTokenRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "publicToken and userId are required",
		})
		return
	}

	// Exchange public token for access token with Plaid
	accessToken, itemID, err := plaidpkg.ExchangePublicToken(context.Background(), req.PublicToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to exchange token: " + err.Error(),
		})
		return
	}

	// Get accounts for the linked item
	accounts, err := plaidpkg.GetAccounts(context.Background(), accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get accounts: " + err.Error(),
		})
		return
	}

	// Get item details to retrieve institution ID
	item, err := plaidpkg.GetItem(context.Background(), accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get item details: " + err.Error(),
		})
		return
	}

	// Get institution details
	var institutionID string
	var institutionName string
	if item.InstitutionId.IsSet() {
		instIDPtr := item.InstitutionId.Get()
		if instIDPtr != nil {
			institutionID = *instIDPtr
			institution, err := plaidpkg.InstitutionsGetByID(context.Background(), *instIDPtr)
			if err == nil {
				institutionName = institution.GetName()
			}
		}
	}

	// Store item in database
	dbItem, err := db.CreateItem(context.Background(), req.UserID, accessToken, itemID, institutionID, "linked")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to store item: " + err.Error(),
		})
		return
	}

	// Store accounts in database and prepare response
	var responseAccounts []gin.H
	for _, account := range accounts {
		_, err := db.CreateOrUpdateAccount(
			context.Background(),
			dbItem.ID,
			account.GetAccountId(),
			account.GetName(),
			account.GetMask(),
			string(account.GetType()),
			string(account.GetSubtype()),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to store account: " + err.Error(),
			})
			return
		}

		responseAccounts = append(responseAccounts, gin.H{
			"id":      account.GetAccountId(),
			"name":    account.GetName(),
			"mask":    account.GetMask(),
			"type":    account.GetType(),
			"subtype": account.GetSubtype(),
		})
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"item_id":              dbItem.ID,
		"plaid_item_id":        itemID,
		"institution_name":     institutionName,
		"access_token":         accessToken,
		"accounts":             responseAccounts,
	})
}

// GetItemAccounts handles GET /api/items/:id/accounts
// Retrieves accounts for a specific item
func GetItemAccounts(c *gin.Context) {
	_ = c.Param("id") // itemIDStr - not yet implemented

	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not yet implemented",
	})
}
