# Go Server Rebuild Roadmap

A step-by-step guide to rebuilding the Plaid Pattern server in Go from scratch.

## Overview

Build a Go server that:
- Integrates with Plaid API to link bank accounts
- Syncs transactions from Plaid
- Stores data in PostgreSQL
- Serves a REST API for the frontend

**Estimated Time:** 3-5 days (depending on Go experience)

---

## Project Structure

```
go-server/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
├── internal/
│   ├── config/                     # Environment configuration
│   ├── db/                         # Database queries
│   ├── handlers/                   # HTTP handlers
│   ├── middleware/                 # HTTP middleware
│   ├── plaid/                      # Plaid client wrapper
│   └── services/                   # Business logic
├── pkg/
│   └── models/                     # Data models
├── go.mod
└── .env
```

---

## Phase 1: Setup & Hello World (30 min)

### Goals
- Initialize Go module
- Set up basic HTTP server
- Understand Go project structure

### Tasks
- [ ] Create `go-server/` directory structure
- [ ] Run `go mod init github.com/bnowak/pattern-go`
- [ ] Create `cmd/server/main.go` with basic HTTP server
- [ ] Test with `/health` endpoint that returns "OK"

**Key Concepts:** Packages, imports, basic HTTP handling

**Resources:**
- https://go.dev/doc/tutorial/getting-started
- https://gobyexample.com/http-servers

---

## Phase 2: Web Framework & Routing (1-2 hours)

### Goals
- Add Gin web framework
- Set up routing
- Add environment configuration

### Dependencies
```bash
go get github.com/gin-gonic/gin
go get github.com/joho/godotenv
```

### Tasks
- [ ] Create `internal/config/config.go` to load env vars from `.env`
- [ ] Refactor `main.go` to use Gin router
- [ ] Create `internal/handlers/health.go` for health check
- [ ] Set up API route group `/api`
- [ ] Add CORS middleware (frontend runs on :3001)

**Config should load:**
- `PORT` (default: 8000)
- `PLAID_CLIENT_ID`
- `PLAID_SECRET`
- `PLAID_ENV` (sandbox/production)
- `DATABASE_URL`

**Key Concepts:** Packages, structs, methods, error handling

**Resources:**
- https://gin-gonic.com/docs/
- https://github.com/joho/godotenv

---

## Phase 3: PostgreSQL Integration (2-3 hours)

### Goals
- Connect to PostgreSQL
- Create database models
- Write CRUD queries
- Learn connection pooling

### Dependencies
```bash
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool
```

### Tasks
- [ ] Create `internal/db/db.go` with `Connect()` and connection pool
- [ ] Create models in `pkg/models/`:
  - `user.go` (id, username, created_at, updated_at)
  - `item.go` (id, user_id, plaid_access_token, plaid_item_id, plaid_institution_id, status, transactions_cursor, timestamps)
  - `account.go` (id, item_id, plaid_account_id, name, mask, balances, type, subtype, timestamps)
  - `transaction.go` (id, account_id, plaid_transaction_id, category, type, name, amount, date, pending, timestamps)
- [ ] Create `internal/db/users.go` with functions:
  - `CreateUser(ctx, username)` → User
  - `GetUserByID(ctx, id)` → User
  - `GetUserByUsername(ctx, username)` → User
- [ ] Create `internal/handlers/users.go` with endpoints:
  - `POST /api/users` (create user)
  - `GET /api/users/:id` (get user)
- [ ] Initialize database connection in `main.go`

**Database Schema:** Use existing PostgreSQL schema from `/database/init/create.sql`

**Key Concepts:** Context, pointers, database scanning, struct tags

**Resources:**
- https://github.com/jackc/pgx
- https://pkg.go.dev/github.com/jackc/pgx/v5

---

## Phase 4: Plaid Integration (3-4 hours)

### Goals
- Integrate Plaid Go SDK
- Create link tokens
- Exchange public tokens for access tokens
- Store items in database

### Dependencies
```bash
go get github.com/plaid/plaid-go/v31/plaid
```

### Tasks
- [ ] Create `internal/plaid/client.go`:
  - `Initialize(clientID, secret, env)` - set up Plaid client
  - `GetClient()` - return client instance
- [ ] Create `internal/handlers/link_token.go`:
  - `POST /api/link-token` - create Plaid link token for user
- [ ] Create `internal/db/items.go` with functions:
  - `CreateItem(ctx, userID, accessToken, plaidItemID, institutionID)` → Item
  - `GetItemByPlaidItemID(ctx, plaidItemID)` → Item
  - `UpdateTransactionsCursor(ctx, plaidItemID, cursor)` → error
- [ ] Create `internal/handlers/items.go`:
  - `POST /api/items` - exchange public token, save item to DB
- [ ] Initialize Plaid client in `main.go`

**Plaid API Calls Needed:**
- `LinkTokenCreate` - create link token
- `ItemPublicTokenExchange` - exchange public token for access token

**Key Concepts:** API clients, configuration, JSON binding

**Resources:**
- https://github.com/plaid/plaid-go
- https://plaid.com/docs/api/tokens/

---

## Phase 5: Transaction Sync (3-4 hours)

### Goals
- Implement transaction sync using cursor
- Save transactions to database
- Handle added, modified, and removed transactions

### Tasks
- [ ] Create `internal/db/accounts.go` with functions:
  - `CreateOrUpdateAccount(ctx, itemID, account)` → error (uses UPSERT)
  - `GetAccountByPlaidID(ctx, plaidAccountID)` → Account
- [ ] Create `internal/db/transactions.go` with functions:
  - `CreateOrUpdateTransaction(ctx, accountID, transaction)` → error (uses UPSERT)
  - `DeleteTransaction(ctx, plaidTransactionID)` → error
  - `GetTransactionsByUserID(ctx, userID)` → []Transaction
- [ ] Create `internal/services/transaction_sync.go`:
  - `SyncTransactions(ctx, plaidItemID)` → SyncResult
  - This should:
    1. Get item from DB to retrieve access token and cursor
    2. Call Plaid `TransactionsSync` API (paginated loop)
    3. Call Plaid `AccountsGet` to get updated accounts
    4. Save/update accounts to DB
    5. Save/update added and modified transactions to DB
    6. Delete removed transactions from DB
    7. Update cursor in DB
    8. Return counts of added/modified/removed
- [ ] Create `internal/handlers/transactions.go`:
  - `POST /api/transactions/sync` - trigger sync for an item
  - `GET /api/users/:user_id/transactions` - get all transactions for user
- [ ] Add helper functions to convert Plaid models to your models

**Plaid API Calls Needed:**
- `TransactionsSync` - get transaction updates (use cursor for pagination)
- `AccountsGet` - get account balances

**Transaction Sync Logic:**
Reference the Node.js version: `/server/update_transactions.js`

**Key Concepts:** Slices, loops, service layer pattern, data transformation

**Resources:**
- https://plaid.com/docs/api/products/transactions/#transactionssync

---

## Phase 6: Testing & Refinement (2-3 hours)

### Goals
- Test complete flow end-to-end
- Add proper error handling
- Add logging
- Add CORS for frontend

### Tasks
- [ ] Create `internal/handlers/errors.go`:
  - Helper function to handle Plaid errors
  - Helper function for generic errors
- [ ] Create `internal/middleware/logging.go`:
  - Request logger middleware
- [ ] Add CORS middleware (use `github.com/gin-contrib/cors`)
- [ ] Test complete flow:
  1. Create user
  2. Create link token
  3. Use frontend to link bank (get public_token)
  4. Exchange public token
  5. Sync transactions
  6. Retrieve transactions
- [ ] Verify data in PostgreSQL

**Testing Commands:**
```bash
# Create user
curl -X POST http://localhost:8000/api/users \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser"}'

# Create link token
curl -X POST http://localhost:8000/api/link-token \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1}'

# After using Link in frontend, exchange token
curl -X POST http://localhost:8000/api/items \
  -H "Content-Type: application/json" \
  -d '{"public_token": "...", "institution_id": "...", "user_id": 1}'

# Sync transactions
curl -X POST http://localhost:8000/api/transactions/sync \
  -H "Content-Type: application/json" \
  -d '{"item_id": "..."}'

# Get transactions
curl http://localhost:8000/api/users/1/transactions
```

**Key Concepts:** Error handling, middleware, logging, testing

---

## Phase 7: Additional Features (Optional, 2-4 hours)

### Tasks
- [ ] Add `GET /api/users/:user_id/accounts` endpoint
- [ ] Add `GET /api/institutions/:id` endpoint (use Plaid `InstitutionsGetById`)
- [ ] Add graceful shutdown on SIGINT/SIGTERM
- [ ] Add request timeout middleware
- [ ] Add API event logging (similar to Node.js version)

---

## API Endpoints Summary

**Users:**
- `POST /api/users` - Create user
- `GET /api/users/:id` - Get user

**Plaid Link:**
- `POST /api/link-token` - Create link token

**Items:**
- `POST /api/items` - Exchange public token, create item

**Transactions:**
- `POST /api/transactions/sync` - Sync transactions for item
- `GET /api/users/:user_id/transactions` - Get user's transactions

**Optional:**
- `GET /api/users/:user_id/accounts` - Get user's accounts
- `GET /api/institutions/:id` - Get institution info
- `GET /api/health` - Health check

---

## Key Go Concepts to Learn

### Basics
- Packages and imports
- Capitalization for exports
- Variables and types
- Pointers (`*` and `&`)

### Error Handling
- Multiple return values
- `if err != nil` pattern
- Error wrapping with `fmt.Errorf`

### Structs & Methods
- Struct definitions
- Struct tags (for JSON)
- Methods vs functions

### Web Development
- HTTP handlers
- Request/response
- JSON marshaling/unmarshaling
- Middleware pattern

### Database
- Connection pooling
- `QueryRow` vs `Query` vs `Exec`
- Scanning results
- Context usage

### Project Organization
- `internal/` - private packages
- `pkg/` - public/shared packages
- `cmd/` - executable entry points

---

## Useful Go Patterns

### Error Handling
```go
if err != nil {
    return nil, fmt.Errorf("failed to do something: %w", err)
}
```

### Defer for Cleanup
```go
rows, err := db.Query(...)
if err != nil {
    return err
}
defer rows.Close()
```

### Optional Fields with Pointers
```go
type User struct {
    Name  string  `json:"name"`
    Email *string `json:"email,omitempty"`  // Can be nil
}
```

### Context Usage
```go
func DoSomething(ctx context.Context, arg string) error {
    result, err := db.QueryRow(ctx, "SELECT ...", arg)
    // ...
}
```

---

## Common Pitfalls

1. **Forgetting to check for errors** - Always check `err != nil`
2. **Nil pointer dereference** - Check pointers before accessing
3. **Not closing database rows** - Always `defer rows.Close()`
4. **Ignoring context** - Pass context through function chains
5. **Capitalization** - Only capitalized functions/types are exported

---

## Resources

### Go Fundamentals
- [A Tour of Go](https://go.dev/tour/) - Interactive tutorial
- [Go by Example](https://gobyexample.com/) - Code examples
- [Effective Go](https://go.dev/doc/effective_go) - Best practices

### Libraries Documentation
- [Gin Framework](https://gin-gonic.com/docs/)
- [pgx Driver](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [Plaid Go SDK](https://github.com/plaid/plaid-go)

### Reference Code
- Existing Node.js server: `/Users/bnowak/dev/pattern/server/`
- Go quickstart: `/Users/bnowak/dev/pattern/quickstart/go/server.go`

---

## Progress Tracker

- [ ] Phase 1: Setup & Hello World (30 min)
- [ ] Phase 2: Web Framework & Routing (1-2 hours)
- [ ] Phase 3: PostgreSQL Integration (2-3 hours)
- [ ] Phase 4: Plaid Integration (3-4 hours)
- [ ] Phase 5: Transaction Sync (3-4 hours)
- [ ] Phase 6: Testing & Refinement (2-3 hours)
- [ ] Phase 7: Additional Features (Optional)

**Total Estimated Time:** 12-20 hours

---

## Getting Started

1. Start with Phase 1 and work through each phase sequentially
2. Check off tasks as you complete them
3. Reference the existing Node.js code when stuck
4. Don't hesitate to peek at the Go quickstart for hints
5. Test after each phase to ensure everything works

Good luck! Building it from scratch will teach you way more than copying code.