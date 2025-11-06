package main

import (
	"compound/go-server/internal/db"
	"compound/go-server/internal/handlers"
	plaidpkg "compound/go-server/internal/plaid"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	PLAID_CLIENT_ID     = ""
	PLAID_SECRET_ID     = ""
	PLAID_ENV           = ""
	PLAID_PRODUCTS      = ""
	PLAID_COUNTRY_CODES = ""
	PLAID_REDIRECT_URI  = ""
	DATABASE_URL        = "" // remove later
	APP_PORT            = ""
)

func init() {
	// load env vars from .env file
	err := godotenv.Load("../.env") // env file located at proj root
	if err != nil {
		fmt.Println("Error when loading environment variables from .env file %w", err)
	}

	// set constants from env
	PLAID_ENV = os.Getenv("PLAID_ENV")
	PLAID_CLIENT_ID = os.Getenv("PLAID_CLIENT_ID")
	// set secret based on env
	if PLAID_ENV == "production" {
		PLAID_SECRET_ID = os.Getenv("PLAID_SECRET_PRODUCTION")
	} else {
		PLAID_SECRET_ID = os.Getenv("PLAID_SECRET_SANDBOX")
	}

	PLAID_PRODUCTS = os.Getenv("PLAID_PRODUCTS")
	PLAID_COUNTRY_CODES = os.Getenv("PLAID_COUNTRY_CODES")
	PLAID_REDIRECT_URI = os.Getenv("PLAID_REDIRECT_URI")
	DATABASE_URL = os.Getenv("DATABASE_URL")
	APP_PORT = os.Getenv("APP_PORT")

	// set defaults if env not present
	if PLAID_PRODUCTS == "" {
		PLAID_PRODUCTS = "transactions"
	}
	if PLAID_COUNTRY_CODES == "" {
		PLAID_COUNTRY_CODES = "US"
	}
	if PLAID_ENV == "" {
		PLAID_ENV = "sandbox"
	}
	if APP_PORT == "" {
		APP_PORT = "8000"
	}
	if PLAID_CLIENT_ID == "" {
		log.Fatal("PLAID_CLIENT_ID is not set. Make sure to fill out the .env file")
	}
	if PLAID_SECRET_ID == "" {
		log.Fatal("PLAID_SECRET_ID is not set. Make sure to fill out the .env file")
	}

	// initialize Plaid client
	err = plaidpkg.Initialize(PLAID_CLIENT_ID, PLAID_SECRET_ID, PLAID_ENV)
	if err != nil {
		log.Fatal("Failed to initialize Plaid client:", err)
	}
}

func main() {
	fmt.Printf("Plaid client initialized successfully, using environment: %s \n", PLAID_ENV)

	// initialize the DB
	err := db.Connect(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.CloseConnection(context.Background())

	// // test db connection with username query
	// user, err := db.GetUserByUsername(context.Background(), "browak")
	// if err != nil {
	// 	log.Fatal("Query failed:", err)
	// }
	// fmt.Printf("User found! ID: %d, Username: %s, Created at: %s\n", user.ID, user.Username, user.CreatedAt)

	// create a Gin router with default middleware (logger and recovery)
	// test
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3001", "http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
	}))

	// -------------------------------------------------
	// start API endpoints
	// -------------------------------------------------

	// define a simple test endpoint - pointing to a func
	router.GET("/api/ping", returnPong)

	// User endpoints
	router.POST("/api/users", handlers.CreateUser)
	router.GET("/api/users/:id", handlers.GetUser)
	router.GET("/api/users/username/:username", handlers.GetUserByUsername)

	// Link Token endpoint
	router.POST("/api/link-token", handlers.MakeLinkTokenHandler(PLAID_PRODUCTS, PLAID_COUNTRY_CODES, PLAID_REDIRECT_URI))

	// Item endpoints
	router.POST("/api/items", handlers.ExchangeToken)
	router.GET("/api/items/:id/accounts", handlers.GetItemAccounts)

	// Transaction endpoints
	router.POST("/api/items/:itemID/sync-transactions", handlers.SyncTransactionsForItem)

	// -------------------------------------------------
	// end API endpoints
	// -------------------------------------------------

	// // call fx outside API for testing
	// // this works - y API no work?
	// ctx := context.Background()
	// user, err := db.CreateUser(ctx, "julia")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(user)

	// start server on port 8000 w/ error handling
	err = router.Run(":" + APP_PORT)
	if err != nil {
		panic("unable to start server")
	}

}

// Handler functions for testing

func returnPong(c *gin.Context) {
	// return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "ping pong ding dong",
	})
}
