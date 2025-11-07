# CLAUDE.md - Compound Project Guide

**Project Name:** Compound (formerly Plaid Pattern)
**Status:** Active Development
**Last Updated:** 2025-11-03

---

## Project Overview

Compound is a full-stack Personal Finance Manager application that demonstrates integration with the Plaid API. The app allows users to:
- Link their bank accounts via Plaid Link
- View account information and balances
- Browse transaction history with categorization
- See net worth calculations
- Experience real-time data updates

**Current State:** The Node.js/Express backend is functional. A Go backend migration is in progress (early phases).

**Important:** This is primarily a development/reference application. While it can be deployed for another user to access, it is not production-hardened.

---

## Tech Stack

### Frontend
- **Framework:** React 16.14 (JavaScript/TypeScript)
- **Routing:** React Router
- **HTTP Client:** Axios
- **Visualization:** Recharts
- **Real-time:** Socket.io Client (for future use)
- **Styling:** Sass
- **Port:** 3001
- **Status:** Stable - likely to remain React unless maintenance becomes prohibitive

### Backend (Current - Node.js)
- **Runtime:** Node.js with Express.js
- **Plaid SDK:** v30.0.0 (Node.js)
- **Database Driver:** pg (node-postgres)
- **Port:** 5001
- **Real-time:** Socket.io (for future use)
- **Status:** Production-ready, but being replaced by Go

### Backend (In Development - Go)
- **Language:** Go 1.25.3
- **Web Framework:** Gin
- **Database Driver:** pgx/v5
- **Plaid SDK:** v40.1.0 (Go)
- **Port:** 8000 (planned)
- **Status:** Early phases (Phase 1-2 of 7-phase roadmap)

### Database
- **Type:** PostgreSQL 11.2
- **Deployment:** Docker container
- **Port:** 5432
- **Schema:** `users`, `items`, `accounts`, `transactions`, `link_events_table`, `plaid_api_events_table`

### Infrastructure
- **Orchestration:** Docker Compose
- **Dev Tooling:** Make (build commands)
- **~~Webhooks:~~ Deprecated** - Will use pull-based data fetching instead

---

## Project Structure

```
compound/
├── client/                    # React frontend (port 3001)
│   ├── src/
│   │   ├── components/       # React components
│   │   ├── services/api.js   # HTTP API calls to backend
│   │   ├── services/        # Plaid Link, socket.io, etc.
│   │   └── hooks/           # Custom React hooks
│   ├── public/              # Static assets
│   └── package.json
│
├── server/                    # Node.js/Express backend (port 5001) - CURRENT
│   ├── routes/              # API endpoints
│   ├── db/                  # Database functions
│   ├── webhookHandlers/     # Webhook logic (DEPRECATED - will be removed)
│   ├── index.js             # Express app entry point
│   ├── plaid.js             # Plaid client configuration
│   ├── update_transactions.js # Transaction sync logic
│   └── package.json
│
├── go-server/                # Go backend (IN DEVELOPMENT)
│   ├── cmd/server/
│   │   └── main.go          # Server entry point, routing
│   ├── internal/
│   │   └── db/              # Database operations
│   │       ├── db.go        # Connection, User queries
│   │       └── (expanding...)
│   ├── go.mod               # Dependencies
│   └── go.sum
│
├── database/                 # Database setup
│   └── init/
│       └── create.sql       # Schema definition
│
├── docs/                     # Documentation
│   └── troubleshooting.md
│
├── quickstart-unused/        # Legacy examples (NOT USED)
│
├── docker-compose.yml        # Service orchestration
├── Makefile                  # Build/run commands
├── GO_REBUILD_ROADMAP.md     # 7-phase Go migration plan
├── COMPOUND_README.md        # Project notes
├── README.md                 # Full documentation (from Plaid Pattern)
└── CLAUDE.md                 # This file
```

---

## Getting Started (Local Development)

### Prerequisites
1. **Docker** (v2.0.0.3 or higher) - running and signed in
2. **Plaid API Keys** - from [dashboard.plaid.com](https://dashboard.plaid.com) (free Sandbox account available)
3. **Make** - available on Unix/Mac (Windows users: use WSL)

### Quick Start
```bash
# 1. Clone and navigate
git clone https://github.com/plaid/pattern.git
cd compound

# 2. Set up environment
# Create .env in project root with your Plaid credentials
# (See .env.template as reference if available)

# 3. Start all services
make start

# 4. Open app
# Frontend: http://localhost:3001
# Logs: make logs

# 5. Stop when done
make stop
```

### Database Management
```bash
make sql          # Start interactive psql session
make clear-db     # Clear all data and restart
make logs         # View service logs
```

---

## Architecture

### Service Dependencies (Docker Compose)
```
Database (PostgreSQL)
        ↓
    Node.js Server (port 5001)
        ↓
React Client (port 3001)
```

**Note:** ngrok service is defined in docker-compose.yml but will be deprecated (no webhook support needed).

### Data Flow (Current)
1. **User opens app** → React client (port 3001)
2. **User links account** → Plaid Link UI
3. **Frontend exchanges token** → POST to Node.js server (port 5001)
4. **Server stores credentials** → PostgreSQL database
5. **Server fetches data** → Plaid API (via SDK)
6. **Server stores data** → PostgreSQL
7. **Frontend displays** ← REST API responses

### Data Flow (Future with Go + Pull-based)
Same as above, but:
- Node.js server (5001) → Go server (8000)
- No webhook support (will use periodic polling)
- No ngrok needed

---

## Key Concepts & Implementation Details

### Plaid Integration

#### Token Exchange Flow
- **Plaid Link** generates a `public_token` client-side
- Backend exchanges it for `access_token` and `item_id` (never exposed client-side)
- These are stored securely in PostgreSQL
- User is associated with their items via the `items` table

**File:** [server/routes/items.js](server/routes/items.js) (Node.js) → will be ported to Go

#### Transaction Syncing
- Uses **Plaid TransactionsSync API** (cursor-based pagination)
- Currently triggered:
  - On new item creation
  - Via webhooks (DEPRECATED - will be removed)
  - Manual polling (future)
- Handles added, modified, and removed transactions
- Stores in `transactions` table

**File:** [server/update_transactions.js](server/update_transactions.js) (Node.js)

#### Preventing Duplicate Item Linkage
- Check `institution_id` when creating new items
- Prevent users from linking the same bank twice (optional, currently disabled)

### Database Schema

**Key Tables:**
- `users` - App users
- `items` - Plaid items (one per linked bank account)
- `accounts` - Accounts within items (checking, savings, etc.)
- `transactions` - Individual transactions (synced from Plaid)
- `link_events_table` - Link UI events (onExit, onSuccess, errors)
- `plaid_api_events_table` - API requests/responses (for debugging/auditing)

**Security:** Never store `access_tokens` client-side. Keep them server-only.

**Plaid Identifiers for Debugging:**
- `request_id` - All API responses
- `link_session_id` - Link callbacks
- Store these for troubleshooting with Plaid Support

See [database/init/create.sql](database/init/create.sql) for full schema.

### Real-time Updates (Socket.io)
- Currently implemented but webhook-dependent (will be refactored for pull-based)
- Future: Replace webhooks with periodic polling
- Frontend listens in [client/src/components/Sockets.jsx](client/src/components/Sockets.jsx)

### OAuth Testing (Sandbox)
- Redirect URI: `http://localhost:3001/oauth-link`
- Configure in Plaid Dashboard
- Test institution: "Playtypus OAuth Bank"

### Plaid Sandbox Testing Credentials
**General Testing:**
- Username: `user_good`
- Password: `pass_good`

**Transaction Testing:**
- Username: `user_transactions_dynamic`
- Password: any value

**Persona-based Testing (detailed transaction data):**
- `user_ewa_user` - Standard individual
- `user_yuppie` - Higher balance account
- `user_small_business` - Business account transactions

---

## Development Notes

### Important Distinctions

#### ✅ Currently Used
- **Frontend:** React (port 3001)
- **Backend:** Node.js/Express (port 5001)
- **Database:** PostgreSQL
- **Data Fetching:** Plaid API calls (pull-based)
- **Real-time:** Socket.io (for websocket infrastructure)

#### ⏳ In Progress
- **Go Backend:** Early migration (Phase 1-2 of 7)
- Basic structure in place: Gin routes, pgx DB connections
- Will gradually replace Node.js server

#### 🚫 Deprecated / Not Used
- **ngrok:** Will be removed (no webhook support needed)
- **quickstart-unused/:** Legacy examples (ignore)
- **Webhook handlers:** Logic will be removed in Go version

#### 🔄 Future Plans
- Go backend to replace Node.js (phases 3-7)
- React frontend rewrite (if maintenance becomes hard)
- Pull-based data fetching only (no webhooks)
- Deployment-ready configuration

### Go Migration Status

**Completed (Phase 1-4):**
- ✅ Basic Go project structure (cmd/server, internal/db, internal/handlers, internal/plaid, pkg/models)
- ✅ Gin web framework setup with CORS middleware
- ✅ PostgreSQL connection (pgx/v5) with user CRUD operations
- ✅ Data models (User, Item, Account, Transaction)
- ✅ Plaid client wrapper (`/internal/plaid/client.go`)
- ✅ Link token generation endpoint (`POST /api/link-token`)
- ✅ Token exchange and item creation endpoint (`POST /api/items`)
- ✅ React test app for full Plaid Link flow testing (`/client-go`)

**Pending (Phase 5-7):**
- Transaction sync logic
- Account and transaction endpoints
- Error handling and logging enhancements
- Full API endpoint coverage
- Tests and documentation

**Reference:** [GO_REBUILD_ROADMAP.md](GO_REBUILD_ROADMAP.md)

### Go Development Workflow

**Quick Setup:**
```bash
# Make sure database is running
make start

# In another terminal, hot-reload with Air
make go-air
```

**Available Make Targets:**
- `make go-build` - Compile Go server to binary (`go-server/server`)
- `make go-run` - Run server once (`go run ./cmd/server`)
- `make go-test` - Run Go tests
- `make go-air` - **RECOMMENDED for development** - Watch files and auto-rebuild/restart (hot-reload)

**What is Air?**
[Air](https://github.com/air-verse/air) is a live-reload tool that automatically rebuilds and restarts your Go server whenever you save changes. Config is in `go-server/.air.toml`.

**Testing API Endpoints:**
```bash
# Create user (POST with JSON body)
curl -X POST http://localhost:8000/api/users \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser"}'

# Get user by ID
curl http://localhost:8000/api/users/1

# Get user by username
curl http://localhost:8000/api/users/username/testuser
```

### Go Server Testing (Plaid Link Flow)

**Prerequisites:**
- `.env` file configured with Plaid credentials
- Database running (`make start`)
- Node packages installed in `/client-go` (`npm install`)

**Quick Start:**
```bash
# Launch Go server + React test app (single command, handles Ctrl+C cleanup)
make go-test-frontend
```

**Manual Launch (if preferred):**
```bash
# Terminal 1: Database
make start

# Terminal 2: Go server with hot-reload
make go-air

# Terminal 3: React test app
cd client-go
npm run dev
```

**Testing Flow:**
1. Open http://localhost:5173 in browser
2. Create test user (e.g., "testuser")
3. Click "Get Link Token" - should generate token and open Plaid Link modal
4. Use Plaid Sandbox credentials (Username: `user_good`, Password: `pass_good`)
5. Select institution and authorize
6. System exchanges token, fetches accounts, stores item in database
7. View linked items and accounts in test app

**Expected Response from POST /api/items:**
```json
{
  "item_id": "123",
  "plaid_item_id": "item-sandbox-...",
  "institution_name": "Plaid Sandbox Bank",
  "access_token": "access-sandbox-...",
  "accounts": [
    {
      "id": "account-id",
      "name": "Checking Account",
      "mask": "1234",
      "type": "depository",
      "subtype": "checking"
    }
  ]
}
```

**Verify Data in Database:**
```bash
make sql
# In psql:
SELECT id, username FROM users;
SELECT id, user_id, plaid_item_id FROM items;
SELECT id, item_id, name FROM accounts;
```

**Troubleshooting:**
- **"Link modal won't open"** - Check that Plaid client ID is set in `/client-go/.env`
- **"Token exchange fails"** - Verify `.env` has correct PLAID_CLIENT_ID and PLAID_SECRET
- **"Port already in use"** - Check for existing processes: `lsof -i :8000` or `lsof -i :5173`
- **"Database connection error"** - Ensure `make start` completed successfully, check `make logs`

---

## Common Tasks

### Add a New API Endpoint (Current: Node.js)
1. Create route in [server/routes/](server/routes/)
2. Add database query if needed in [server/db/](server/db/)
3. Call Plaid API via [server/plaid.js](server/plaid.js)
4. Test with curl or frontend integration

### Add a React Component
1. Create `.jsx` file in [client/src/components/](client/src/components/)
2. Use existing components as reference (e.g., Dashboard, AccountsList)
3. Call backend APIs via [client/src/services/api.js](client/src/services/api.js)
4. Style with Sass in accompanying `.scss` file

### Migrate Node.js Endpoint to Go
1. Reference the Node.js implementation
2. Implement in Go using Gin + pgx
3. Update docker-compose.yml to route to port 8000 (when ready)
4. Test thoroughly before removing Node.js version

### Debug the App
- **Frontend:** React DevTools browser extension
- **Backend (Node.js):** VS Code Docker debugger (port 9229)
- **Backend (Go):** Go debugger (delve) - TBD
- **Database:** `make sql` for psql access
- **Network:** Check API calls in browser DevTools

### View Data
```bash
make sql
# In psql:
SELECT * FROM users;
SELECT * FROM items;
SELECT * FROM transactions LIMIT 10;
SELECT * FROM plaid_api_events_table ORDER BY created_at DESC LIMIT 5;
```

---

## Environment Configuration (.env)

**Required Variables:**
```
PLAID_CLIENT_ID=<your-client-id>
PLAID_SECRET_SANDBOX=<sandbox-secret>
PLAID_SECRET_PRODUCTION=<production-secret>
PLAID_ENV=sandbox               # or production
PLAID_SANDBOX_REDIRECT_URI=http://localhost:3001/oauth-link
PLAID_PRODUCTION_REDIRECT_URI=https://localhost:3001/oauth-link (for https testing)

# Database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
DB_HOST_NAME=db
DB_PORT=5432

# Services
PORT=5001                       # Node.js server port
```

**Note:** Use `.env.template` as a starting point.

---

## Deployment (Multi-user)

When deploying for another user to access:

1. **Set Production Credentials**
   - Use Plaid Production keys (not Sandbox)
   - Set `PLAID_ENV=production` in `.env`

2. **Configure Domain**
   - Update `PLAID_SANDBOX_REDIRECT_URI` to your domain
   - Update `PLAID_PRODUCTION_REDIRECT_URI` to your domain
   - Register URIs in Plaid Dashboard

3. **HTTPS**
   - Required for production
   - Use Let's Encrypt or similar (not self-signed)

4. **Database**
   - Use managed PostgreSQL (AWS RDS, DigitalOcean, etc.)
   - Update DB connection string in `.env`
   - Run migrations before deploying

5. **Docker Compose**
   - Modify docker-compose.yml for production
   - Remove ngrok service
   - Use environment-specific volumes

6. **Security**
   - Never commit `.env` file
   - Use secrets management (AWS Secrets Manager, etc.)
   - Rotate Plaid credentials regularly

---

## Troubleshooting

### Services Won't Start
```bash
make logs              # Check what failed
make clear-db          # Reset everything
make start             # Try again
```

### Port Already in Use
- 3001 (client), 5001 (server), 5432 (db), 4040 (ngrok)
- Check: `lsof -i :<PORT>` or modify docker-compose.yml

### Database Schema Issues
```bash
make clear-db          # Rebuilds schema from init/create.sql
```

### Plaid Connection Errors
- Verify credentials in `.env`
- Check `plaid_api_events_table` for error details
- Enable debug logging in server code

### Socket.io Connection Issues
- Currently only works with webhooks (deprecated feature)
- Will be refactored for pull-based updates

See [docs/troubleshooting.md](docs/troubleshooting.md) for more.

---

## Useful Files to Know

| File | Purpose | Language |
|------|---------|----------|
| [Makefile](Makefile) | Build/run commands | Make |
| [database/init/create.sql](database/init/create.sql) | Database schema | SQL |
| [server/index.js](server/index.js) | Express server entry | Node.js |
| [server/plaid.js](server/plaid.js) | Plaid client setup | Node.js |
| [server/update_transactions.js](server/update_transactions.js) | Transaction sync | Node.js |
| [client/src/services/api.js](client/src/services/api.js) | Frontend API calls | React/JS |
| [go-server/cmd/server/main.go](go-server/cmd/server/main.go) | Go server entry point, routing | Go |
| [go-server/internal/db/db.go](go-server/internal/db/db.go) | Database connection & lifecycle | Go |
| [go-server/internal/db/user.go](go-server/internal/db/user.go) | User CRUD operations | Go |
| [go-server/.air.toml](go-server/.air.toml) | Air hot-reload config | TOML |
| [GO_REBUILD_ROADMAP.md](GO_REBUILD_ROADMAP.md) | Migration plan | Markdown |
| [README.md](README.md) | Full documentation | Markdown |

---

## Questions? Tips for Claude

When asking Claude to help with Compound:

- **Specify which component:** "Fix bug in Node.js server" vs. "Add feature to React frontend"
- **For Go work:** Reference the Node.js implementation as a template
- **For database issues:** Include relevant SQL schema snippets
- **For deployment:** Mention if it's for local dev vs. multi-user access
- **For debugging:** Share error messages and relevant log sections

Example prompt:
> "The transaction sync endpoint in the Node.js server isn't returning recent transactions. Look at server/update_transactions.js and server/plaid.js to debug."

---

## Quick Links

- **Plaid Docs:** https://plaid.com/docs/
- **Plaid Dashboard:** https://dashboard.plaid.com/
- **PostgreSQL Docs:** https://www.postgresql.org/docs/
- **React Docs:** https://react.dev/
- **Go Docs:** https://golang.org/doc/
- **Gin Framework:** https://gin-gonic.com/
- **pgx Driver:** https://github.com/jackc/pgx

---

**Last Updated:** 2025-11-05
**Node.js Server Status:** ✅ Functional
**Go Server Status:** ⏳ In Development (Phase 4/7 Complete)
**Frontend Status:** ✅ Stable
**Test App Status:** ✅ Created and Ready
