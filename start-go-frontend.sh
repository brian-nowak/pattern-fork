#!/bin/bash
# Start both Go server and React frontend for testing Plaid Link flow

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🚀 Starting Go Server + React Frontend Test Environment"
echo "=========================================================="
echo ""
echo "Backend: http://localhost:8000"
echo "Frontend: http://localhost:5173"
echo ""
echo "Press Ctrl+C to stop both servers"
echo ""

# Trap Ctrl+C to kill both processes
trap 'echo ""; echo "Stopping servers..."; kill $GO_PID $REACT_PID 2>/dev/null || true; exit' INT

# Start Go server in background
echo "Starting Go server (port 8000)..."
cd "$PROJECT_ROOT/go-server"
~/go/bin/air -c .air.toml &
GO_PID=$!

# Give Go server a moment to start
sleep 2

# Start React frontend
echo "Starting React frontend (port 5173)..."
cd "$PROJECT_ROOT/client-go"
npm run dev &
REACT_PID=$!

# Wait for both processes
wait
