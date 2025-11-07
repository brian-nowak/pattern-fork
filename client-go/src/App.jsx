import { useState } from 'react';
import { PlaidLink } from 'react-plaid-link';
import './App.css';
import * as api from './api';

function App() {
  const [users, setUsers] = useState([]);
  const [currentUser, setCurrentUser] = useState(null);
  const [username, setUsername] = useState('');
  const [linkToken, setLinkToken] = useState(null);
  const [linkedItems, setLinkedItems] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [mode, setMode] = useState('normal'); // 'normal' or 'update'
  const [selectedItemId, setSelectedItemId] = useState(null);
  const [transactions, setTransactions] = useState([]);
  const [showTransactions, setShowTransactions] = useState(false);

  // Create a new user or load existing user by ID
  const handleCreateUser = async (e) => {
    e.preventDefault();
    if (!username.trim()) return;

    setLoading(true);
    setError(null);
    try {
      const isNumeric = /^\d+$/.test(username.trim());
      let user;

      if (isNumeric) {
        // Try to load user by ID
        try {
          user = await api.getUser(parseInt(username));
          // Check if user already in list, if not add them
          const userExists = users.some(u => u.id === user.id);
          if (!userExists) {
            setUsers([...users, user]);
          }
          setCurrentUser(user);
          setUsername('');
        } catch (err) {
          // If not found by ID, create as new user with this username
          user = await api.createUser(username);
          setUsers([...users, user]);
          setCurrentUser(user);
          setUsername('');
        }
      } else {
        // Create user with given username
        user = await api.createUser(username);
        setUsers([...users, user]);
        setCurrentUser(user);
        setUsername('');
      }
    } catch (err) {
      setError(`Failed to create user: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Get link token
  const handleGetLinkToken = async () => {
    if (!currentUser) {
      setError('Please create or select a user first');
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const itemId = mode === 'update' ? selectedItemId : null;
      const token = await api.getLinkToken(currentUser.id, itemId);
      setLinkToken(token);
    } catch (err) {
      setError(`Failed to get link token: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Handle successful Plaid Link completion
  const handlePlaidSuccess = async (publicToken, metadata) => {
    setLoading(true);
    setError(null);
    try {
      const result = await api.exchangePublicToken(publicToken, currentUser.id);

      // Add to linked items
      const newItem = {
        id: result.item_id,
        institution: metadata.institution?.name || 'Unknown Institution',
        accounts: result.accounts || [],
      };
      setLinkedItems([...linkedItems, newItem]);

      // Reset
      setLinkToken(null);
      setMode('normal');
      setSelectedItemId(null);

      setError(null); // Clear any previous errors
    } catch (err) {
      setError(`Failed to exchange token: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handlePlaidExit = (err) => {
    if (err) {
      setError(`Plaid Link closed with error: ${err.message}`);
    } else {
      setLinkToken(null);
    }
  };

  // Get transactions for current user
  const handleGetTransactions = async () => {
    if (!currentUser) {
      setError('Please create or select a user first');
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const data = await api.getUserTransactions(currentUser.id);
      setTransactions(data.transactions || []);
      setShowTransactions(true);
    } catch (err) {
      setError(`Failed to get transactions: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app">
      <header className="header">
        <h1>Go Server - Plaid Link Test</h1>
        <p className="subtitle">Testing Plaid integration with Go backend</p>
      </header>

      <main className="main">
        {error && <div className="error-banner">{error}</div>}

        {/* User Creation Section */}
        <section className="card">
          <h2>1. Create New User (or select from existing)</h2>
          <form onSubmit={handleCreateUser} className="form">
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter username"
              disabled={loading}
            />
            <button type="submit" disabled={loading}>
              {loading ? 'Creating...' : 'Create User'}
            </button>
          </form>

          {users.length > 0 && (
            <div className="user-list">
              <h3>Users Created:</h3>
              <ul>
                {users.map((user) => (
                  <li
                    key={user.id}
                    className={currentUser?.id === user.id ? 'active' : ''}
                    onClick={() => setCurrentUser(user)}
                  >
                    {user.username} (ID: {user.id})
                    {currentUser?.id === user.id && ' ✓'}
                  </li>
                ))}
              </ul>
            </div>
          )}

          {currentUser && (
            <div className="current-user">
              <strong>Current User:</strong> {currentUser.username} (ID: {currentUser.id})
            </div>
          )}
        </section>

        {/* Link Token Section */}
        {currentUser && (
          <section className="card">
            <h2>2. Generate Link Token</h2>

            <div className="mode-selector">
              <label>
                <input
                  type="radio"
                  name="mode"
                  value="normal"
                  checked={mode === 'normal'}
                  onChange={(e) => {
                    setMode(e.target.value);
                    setSelectedItemId(null);
                  }}
                />
                Normal Mode (New Account)
              </label>
              <label>
                <input
                  type="radio"
                  name="mode"
                  value="update"
                  checked={mode === 'update'}
                  onChange={(e) => setMode(e.target.value)}
                />
                Update Mode (Re-link Account)
              </label>
            </div>

            {mode === 'update' && linkedItems.length > 0 && (
              <div className="item-selector">
                <label>
                  Select Item to Update:
                  <select
                    value={selectedItemId || ''}
                    onChange={(e) => setSelectedItemId(Number(e.target.value))}
                  >
                    <option value="">Choose an item...</option>
                    {linkedItems.map((item) => (
                      <option key={item.id} value={item.id}>
                        {item.institution}
                      </option>
                    ))}
                  </select>
                </label>
              </div>
            )}

            <button onClick={handleGetLinkToken} disabled={loading || (mode === 'update' && !selectedItemId)}>
              {loading ? 'Getting Token...' : 'Get Link Token'}
            </button>

            {linkToken && (
              <div className="token-display">
                <p>
                  <strong>Link Token Generated!</strong>
                </p>
                <code className="token">{linkToken.substring(0, 20)}...</code>
              </div>
            )}
          </section>
        )}

        {/* Plaid Link Modal Section */}
        {linkToken && (
          <section className="card">
            <h2>3. Open Plaid Link</h2>
            <p>Click the button below to open Plaid Link and connect a bank account.</p>

            <PlaidLink
              token={linkToken}
              onSuccess={handlePlaidSuccess}
              onExit={handlePlaidExit}
            >
              {({ open, ready }) => (
                <button
                  onClick={() => open()}
                  disabled={!ready || loading}
                  className="plaid-button"
                >
                  {loading ? 'Processing...' : 'Open Plaid Link'}
                </button>
              )}
            </PlaidLink>
          </section>
        )}

        {/* Linked Items Section */}
        {linkedItems.length > 0 && (
          <section className="card">
            <h2>4. Linked Items</h2>
            <div className="items-list">
              {linkedItems.map((item) => (
                <div key={item.id} className="item-card">
                  <h3>{item.institution}</h3>
                  <p>Item ID: <code>{item.id}</code></p>

                  {item.accounts && item.accounts.length > 0 && (
                    <div className="accounts">
                      <h4>Accounts ({item.accounts.length}):</h4>
                      <ul>
                        {item.accounts.map((account) => (
                          <li key={account.id}>
                            <strong>{account.name}</strong>
                            {account.mask && ` (...${account.mask})`}
                            <br />
                            <small>Type: {account.type}</small>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </section>
        )}

        {/* Get Transactions Section */}
        {currentUser && (
          <section className="card">
            <h2>5. View Transactions</h2>
            <button onClick={handleGetTransactions} disabled={loading}>
              {loading ? 'Loading...' : 'Get All Transactions'}
            </button>

            {showTransactions && transactions.length > 0 && (
              <div className="transactions-list">
                <h3>Transactions ({transactions.length})</h3>
                <div className="transactions-table">
                  <table>
                    <thead>
                      <tr>
                        <th>Date</th>
                        <th>Name</th>
                        <th>Amount</th>
                        <th>Type</th>
                        <th>Category</th>
                        <th>Pending</th>
                      </tr>
                    </thead>
                    <tbody>
                      {transactions.map((tx) => (
                        <tr key={tx.id}>
                          <td>{new Date(tx.date).toLocaleDateString()}</td>
                          <td>{tx.name}</td>
                          <td>${Math.abs(tx.amount).toFixed(2)}</td>
                          <td>{tx.type}</td>
                          <td>{tx.category || 'N/A'}</td>
                          <td>{tx.pending ? 'Yes' : 'No'}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}

            {showTransactions && transactions.length === 0 && (
              <div className="no-transactions">
                <p>No transactions found for this user.</p>
              </div>
            )}
          </section>
        )}

        {/* Debug Info */}
        <section className="card debug">
          <h3>Debug Info</h3>
          <ul>
            <li>API URL: {import.meta.env.VITE_API_URL}</li>
            <li>Current User: {currentUser ? currentUser.username : 'None'}</li>
            <li>Users Created: {users.length}</li>
            <li>Items Linked: {linkedItems.length}</li>
          </ul>
        </section>
      </main>
    </div>
  );
}

export default App;
