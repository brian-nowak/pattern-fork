# Client-Go: Plaid Link Test App

A minimal React frontend for testing the Go backend implementation of the Compound (Plaid Pattern) application.

## Purpose

This app serves as a **testing and validation tool** for the Go server endpoints. It provides a simple UI to:
- Create test users
- Generate Plaid Link tokens
- Open the Plaid Link modal
- Exchange public tokens for access tokens
- View linked bank accounts

**Note:** This is NOT the full Compound application. It's specifically designed for testing the Plaid integration flow with the Go backend.

## Quick Start

### Prerequisites
- Node.js 16+ installed
- Go server running on `http://localhost:8000`
- Plaid API credentials configured in Go server's `.env`

### Setup

```bash
# Install dependencies (already done)
npm install

# Set up environment variables
# Copy the .env file and add your Plaid public key
# VITE_API_URL=http://localhost:8000
# VITE_PLAID_CLIENT_ID=<your-plaid-public-key>

# Start development server
npm run dev

# App will be available at http://localhost:5173
```

## Architecture

### Technology Stack
- **React 18** with Vite
- **Axios** for HTTP requests
- **react-plaid-link** for Plaid Link modal integration
- **Plain CSS** for styling

### File Structure
```
client-go/
├── src/
│   ├── App.jsx        # Main component with state management
│   ├── App.css        # Styling
│   ├── api.js         # API helper functions for Go endpoints
│   ├── main.jsx       # React entry point
│   └── index.css      # Global styles
├── .env              # Environment configuration
├── package.json
└── vite.config.js
```

### State Management
Uses React hooks (`useState`) for simplicity:
- `currentUser` - Currently selected/created user
- `linkToken` - Generated Plaid link token
- `linkedItems` - Array of linked bank items
- `mode` - "normal" (new account) or "update" (re-link) mode
- `loading` & `error` - Request and error states

## API Endpoints (Called)

All requests target the Go server on `http://localhost:8000`:

### User Management
- `POST /api/users` - Create a new user
- `GET /api/users/:id` - Fetch user details

### Link Token
- `POST /api/link-token` - Generate a link token for account linking

### Item Management
- `POST /api/items` - Exchange public token for access token and fetch accounts
- `GET /api/items/:id/accounts` - Get accounts for an item (future)

## User Flow

1. **Create User**
   - Enter username → Click "Create User"
   - User ID is displayed and automatically selected

2. **Generate Link Token**
   - Select mode: "Normal" (new account) or "Update" (re-link existing)
   - If update mode, select an already-linked item
   - Click "Get Link Token"
   - System calls `POST /api/link-token` on Go server

3. **Open Plaid Link Modal**
   - Click "Open Plaid Link" button
   - Modal opens with Plaid's bank connection UI
   - User selects their bank and logs in

4. **Link Account Successfully**
   - Plaid returns a `publicToken` to the frontend
   - Frontend automatically calls `POST /api/items` to exchange it
   - Go server returns `accessToken`, `itemId`, and linked accounts
   - Results displayed in "Linked Items" section

5. **View Linked Accounts**
   - Displays institution name and all connected accounts
   - Can re-link same account by selecting "Update Mode"

## Environment Variables

```
VITE_API_URL=http://localhost:8000
VITE_PLAID_CLIENT_ID=<public-key-from-plaid-dashboard>
```

**Note:** The Plaid Client ID here is for the Link modal initialization. The secret key remains on the Go server side (never exposed to frontend).

## Debugging

### Browser Console
The app logs API calls and state changes. Check the browser DevTools console for:
- API request/response details
- Error messages from the Go server
- Link modal events (onSuccess, onExit)

### Network Tab
- Check actual HTTP requests to `http://localhost:8000`
- Verify CORS headers are correct
- See response bodies from Go server

### Debug Info Section
At the bottom of the page, a "Debug Info" card shows:
- Current API URL
- Current user
- Count of users created
- Count of items linked

## Common Issues

### CORS Errors
**Problem:** "No 'Access-Control-Allow-Origin' header"

**Solution:** Ensure the Go server has CORS middleware enabled for `http://localhost:5173` (or the port your dev server uses).

### Link Token Not Generating
**Problem:** 400 or 500 errors when clicking "Get Link Token"

**Steps:**
1. Check browser console for error details
2. Verify Go server is running on port 8000
3. Confirm Plaid credentials are set in Go server's `.env`
4. Check Go server logs for more details

### Public Token Exchange Fails
**Problem:** Link modal opens but token exchange fails

**Steps:**
1. Check that `POST /api/items` endpoint exists on Go server
2. Verify the Go server can successfully call Plaid's `ItemPublicTokenExchange` API
3. Check Go server logs for Plaid API errors

## Next Steps

Once all endpoints are working and tested here:
1. Copy working patterns into the main `/client` app
2. Update `/client/src/services/api.tsx` to support both backends
3. Integrate user and item management into full app
4. Add dashboard for viewing transactions and net worth

## Useful Commands

```bash
# Start dev server with hot reload
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Check code for errors
npm run lint

# Format code
npm run format
```

## Related Files

- Go Server: `/go-server/cmd/server/main.go`
- Main App (existing): `/client/src/services/api.tsx`
- Go Plaid Client: `/go-server/internal/plaid/client.go`
- Go Models: `/go-server/pkg/models/`
