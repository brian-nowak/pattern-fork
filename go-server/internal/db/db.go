package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Global connection variable
var conn *pgx.Conn

// TODO - add some context driven timeouts (for connections and queries)

// Connect establishes a connection to the PostgreSQL database
func Connect(databaseURL string) error {
	var err error
	// Connect to the database
	conn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	// Test the connection
	err = conn.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to database successfully")
	return nil
}

// CloseConnection closes the database connection
func CloseConnection(ctx context.Context) error {
	if conn != nil {
		return conn.Close(ctx)
	}
	return nil
}
