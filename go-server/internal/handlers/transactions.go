package handlers

import (
	"compound/go-server/internal/db"
	"context"
	"net/http"
	"strconv"

	plaidpkg "compound/go-server/internal/plaid"

	"github.com/gin-gonic/gin"
)

// called inside SyncTransactionsForItem
// accepts each transaction object from the plaid data
// func CreateOrUpdateTransaction(c *gin.Context) {

// }

// SyncTransactionsForItem handles POST /api/transactions/sync
// accepts itemID
func SyncTransactionsForItem(c *gin.Context) {

	// parse itemID from incoming URL
	// Get user ID from URL parameter
	idStr := c.Param("itemID")
	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	// retrieve item from DB to get access token
	item, err := db.GetItemByID(context.Background(), itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get item: " + err.Error(),
		})
		return
	}

	// use access token + cursor to call plaid.NewTransactionsSyncRequest
	// returns plaid objects
	result, err := plaidpkg.SyncTransactions(context.Background(), item.PlaidAccessToken, item.TransactionsCursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to sync transactions: " + err.Error(),
		})
		return
	}

	// move this to CreateOrUpdateTransaction later
	// combine adds and updates
	allTransactions := append(result.Added, result.Modified...)

	// loop thru each
	for _, plaidTx := range allTransactions {
		// get our DB account ID from plaid account ID
		account, err := db.GetAccountByPlaidAccountID(context.Background(), plaidTx.GetAccountId())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get DB Account ID for: " + err.Error(),
			})
			return
		}

		// map category data
		categoryData := map[string]interface{}{
			"legacy":                    plaidTx.GetCategory(),
			"personal_finance_category": plaidTx.GetPersonalFinanceCategory(),
		}

		// Convert account owner string to *string
		accountOwner := plaidTx.GetAccountOwner()
		var ownerPtr *string
		if accountOwner != "" {
			ownerPtr = &accountOwner
		}

		_, err = db.CreateOrUpdateTransaction(
			context.Background(),
			account.ID,
			plaidTx.GetTransactionId(),
			categoryData,
			plaidTx.GetTransactionType(),
			plaidTx.GetName(),
			plaidTx.GetAmount(),
			plaidTx.GetIsoCurrencyCode(),
			plaidTx.GetUnofficialCurrencyCode(),
			plaidTx.GetDate(),
			plaidTx.GetPending(),
			ownerPtr,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to store transaction: " + err.Error(),
			})
			return
		}
	}

	// Update the cursor for the next sync
	err = db.UpdateItemTransactionsCursor(context.Background(), itemID, result.NextCursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update cursor: " + err.Error(),
		})
		return
	}

	// Return summary of what was synced
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"addedCount":    len(result.Added),
		"modifiedCount": len(result.Modified),
		"removedCount":  len(result.Removed),
	})
}

// handles GET /api/users/:user_id/transactions
// Returns all transactions for a specific user
func GetUserTransactions(c *gin.Context) {
	// parse user ID from URL parameter
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	// get all transactions for the user
	transactions, err := db.GetTransactionByUserID(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get transactions: " + err.Error(),
		})
		return
	}

	// return transactions
	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
	})
}
