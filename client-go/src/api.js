import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL;

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// User endpoints
export const createUser = async (username) => {
  const response = await api.post('/api/users', { username });
  return response.data;
};

export const getUser = async (userId) => {
  const response = await api.get(`/api/users/${userId}`);
  return response.data;
};

// Link token endpoint
export const getLinkToken = async (userId, itemId = null) => {
  const response = await api.post('/api/link-token', {
    userId,
    itemId,
  });
  return response.data.link_token;
};

// Item endpoints (token exchange)
export const exchangePublicToken = async (publicToken, userId) => {
  const response = await api.post('/api/items', {
    publicToken,
    userId,
  });
  return response.data;
};

// Get user's items
export const getUserItems = async (userId) => {
  const response = await api.get(`/api/users/${userId}/items`);
  return response.data;
};

// Get item accounts
export const getItemAccounts = async (itemId) => {
  const response = await api.get(`/api/items/${itemId}/accounts`);
  return response.data;
};

// Get user's transactions
export const getUserTransactions = async (userId) => {
  const response = await api.get(`/api/transactions/${userId}`);
  return response.data;
};

export default api;
