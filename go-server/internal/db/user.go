package db

import (
	"compound/go-server/pkg/models"
	"context"
	"fmt"
)

// GetUserByID retrieves a user from the database by ID
func GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := "SELECT id, username, created_at, updated_at FROM users WHERE id=$1"

	user := &models.User{}
	err := conn.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt, // time.Time handles TIMESTAMPTZ
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return user, nil
}

// GetUserByUsername retrieves a user from the database by username
func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := "SELECT id, username, created_at, updated_at FROM users WHERE username=$1"

	user := &models.User{}
	err := conn.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt, // time.Time handles TIMESTAMPTZ
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return user, nil
}

// CreateUser creates a new user in the database
func CreateUser(ctx context.Context, username string) (*models.User, error) {
	query := `INSERT INTO users (username, created_at, updated_at)
				VALUES ($1, NOW(), NOW())
				RETURNING id, username, created_at, updated_at`

	user := &models.User{}
	err := conn.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt, // time.Time handles TIMESTAMPTZ
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user from the database by username
func DeleteUser(ctx context.Context, username string) error {
	query := "DELETE FROM users WHERE username=$1"

	result, err := conn.Exec(ctx, query, username)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	// Check if any rows were affected
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
