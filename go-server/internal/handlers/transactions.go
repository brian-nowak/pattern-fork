package handlers

import (
	"compound/go-server/internal/db"
	"context"
	"net/http"
	"strconv"

	plaidpkg "compound/go-server/internal/plaid"

	"github.com/gin-gonic/gin"
)

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
	// returns :
	// result.Added = resp.GetAdded()
	// result.Modified = resp.GetModified()
	// result.Removed = resp.GetRemoved()
	// result.NextCursor = resp.GetNextCursor()
	// result.HasMore = resp.GetHasMore()
	result, err := plaidpkg.SyncTransactions(context.Background(), item.PlaidAccessToken, item.TransactionsCursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to sync transactions: " + err.Error(),
		})
		return
	}

	// we need to add more but for now lets just return result for testing
	// temp reformat result into json for API
	c.JSON(http.StatusOK, gin.H{
		"added":      result.Added,
		"modified":   result.Modified,
		"removed":    result.Removed,
		"nextCursor": result.NextCursor,
	})
}
