# Compound

## The best personal finance app not yet created

## Quick Start

### Running Services
```bash
# Start all services (client, Node.js server, database)
make start

# View logs
make logs

# Stop services
make stop
```

### Go Server Development

The Go backend is being developed to replace the Node.js server. For development:

```bash
# Terminal 1: Start database and other services
make start

# Terminal 2: Run Go server with hot-reload
make go-air
```

**Available Make Targets for Go:**
- `make go-build` - Compile to binary
- `make go-run` - Run once
- `make go-test` - Run tests
- `make go-air` - **Recommended** - Watch files and auto-rebuild (uses Air)
- `make go-test-frontend` - **NEW** - Launch Go server + React test app (one command)

**Test API Endpoints (Basic):**
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

### Testing Plaid Link Flow

A minimal React test app (`/client-go`) has been created to test the complete Plaid Link flow with the Go backend:

```bash
# One-command launch (includes database, Go server, and React frontend)
make go-test-frontend

# Then open http://localhost:5173 in your browser
```

**Full Testing Workflow:**
1. Create a test user in the app
2. Click "Get Link Token" to generate a Plaid Link token
3. Use Plaid Sandbox credentials (`user_good` / `pass_good`)
4. Authorize the test institution
5. Verify the linked bank account appears in the app
6. Check database: `make sql` then `SELECT * FROM items;`

See [CLAUDE.md](CLAUDE.md) "Go Server Testing" section for detailed troubleshooting and endpoints documentation.

## Useful Commands

### Database
With PG DB running in docker, connect to run queries in CLI:
```shell
make sql
# Or manually:
docker exec -it compound-db-1 psql -U postgres -d postgres
```

### Rebuild & Reset
```bash
make clear-db      # Clears all data, rebuilds schema
make go-build      # Compile Go server
```

## Development Notes

- **Node.js Server:** Running on port 5001 (currently active)
- **Go Server:** Port 8000 (in development with Air hot-reload)
- **Client:** Port 3001
- **Database:** Port 5432

See [CLAUDE.md](CLAUDE.md) for comprehensive project documentation.