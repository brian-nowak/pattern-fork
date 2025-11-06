package handlers

import (
	"compound/go-server/internal/db"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUser handles POST /api/users
func CreateUser(c *gin.Context) {
	// Get username from request body
	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username is required",
		})
		return
	}

	user, err := db.CreateUser(context.Background(), req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user not created successfully - check server logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"created_at": user.CreatedAt,
	})
}

// GetUser handles GET /api/users/:id
func GetUser(c *gin.Context) {
	// Get user ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	user, err := db.GetUserByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"created_at": user.CreatedAt,
	})
}

// GetUserByUsername handles GET /api/users/username/:username
func GetUserByUsername(c *gin.Context) {
	// Get username from URL parameter
	username := c.Param("username")

	user, err := db.GetUserByUsername(context.Background(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"created_at": user.CreatedAt,
	})
}
