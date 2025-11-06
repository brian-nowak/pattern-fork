package handlers

import (
	"compound/go-server/internal/db"
	plaidpkg "compound/go-server/internal/plaid"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LinkTokenRequest represents the request body for creating a link token
type LinkTokenRequest struct {
	UserID int  `json:"userId" binding:"required"`
	ItemID *int `json:"itemId"`
}

// MakeLinkTokenHandler creates a handler for POST /api/link-token
// Generates a Plaid Link token for account linking (normal mode) or updating (update mode)
//
// Request body:
// {
//   "userId": 123,
//   "itemId": null  // omit or null for normal mode, set for update mode
// }
//
// Response:
// {
//   "link_token": "link-sandbox-..."
// }
//
// The handler uses closures to capture plaid configuration from main.go
func MakeLinkTokenHandler(plaidProducts, plaidCountryCodes, plaidRedirectURI string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LinkTokenRequest

		// Parse request body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "userId is required",
			})
			return
		}

		// Determine mode based on presence of itemId
		var products []string

		if req.ItemID != nil && *req.ItemID != 0 {
			// Update mode: re-linking existing item
			// Fetch the item to validate it exists
			item, err := db.GetItemByID(context.Background(), *req.ItemID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "item not found",
				})
				return
			}

			// Verify the item belongs to this user
			if item.UserID != req.UserID {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "item does not belong to this user",
				})
				return
			}

			// Empty products array for update mode - only updating connection, not adding new products
			products = []string{}
		} else {
			// Normal mode: new account linking
			// Must include transactions product for webhooks
			products = []string{"transactions"}
		}

		// Parse country codes
		countryCodes := []string{plaidCountryCodes}

		// Create the link token via plaid package
		linkToken, err := plaidpkg.CreateLinkToken(
			context.Background(),
			req.UserID,
			products,
			countryCodes,
			plaidRedirectURI,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return the link token
		c.JSON(http.StatusOK, gin.H{
			"link_token": linkToken,
		})
	}
}