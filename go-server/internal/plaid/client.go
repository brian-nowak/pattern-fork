package plaid

import (
	"context"
	"fmt"

	plaid "github.com/plaid/plaid-go/v40/plaid"
)

var (
	apiClient *plaid.APIClient
)

// Initialize sets up the Plaid API client with credentials
func Initialize(clientID, secret, env string) error {
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	configuration.AddDefaultHeader("PLAID-SECRET", secret)

	environments := map[string]plaid.Environment{
		"sandbox":    plaid.Sandbox,
		"production": plaid.Production,
	}

	if envType, ok := environments[env]; ok {
		configuration.UseEnvironment(envType)
	} else {
		return fmt.Errorf("invalid environment: %s", env)
	}

	apiClient = plaid.NewAPIClient(configuration)
	return nil
}

// GetClient returns the initialized Plaid API client
func GetClient() *plaid.APIClient {
	return apiClient
}

// CreateLinkToken creates a new Plaid Link token for account linking
func CreateLinkToken(
	ctx context.Context,
	userID int,
	products []string,
	countryCodes []string,
	redirectURI string,
) (string, error) {
	if apiClient == nil {
		return "", fmt.Errorf("plaid client not initialized")
	}

	// Convert string country codes to plaid.CountryCode enum
	countryCodeEnums := convertCountryCodes(countryCodes)

	// Convert string products to plaid.Products enum
	productEnums := convertProducts(products)

	// Create user identifier
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: fmt.Sprintf("user_%d", userID),
	}

	// Create request with required parameters
	request := plaid.NewLinkTokenCreateRequest(
		"Compound",           // client_name
		"en",                 // language
		countryCodeEnums,     // country_codes
	)

	// Set optional user and products
	request.SetUser(user)
	request.SetProducts(productEnums)

	// Set redirect URI if provided
	if redirectURI != "" {
		request.SetRedirectUri(redirectURI)
	}

	resp, _, err := apiClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to create link token: %w", err)
	}

	return resp.GetLinkToken(), nil
}

// ExchangePublicToken exchanges a public token from Link for an access token and item ID
func ExchangePublicToken(ctx context.Context, publicToken string) (string, string, error) {
	if apiClient == nil {
		return "", "", fmt.Errorf("plaid client not initialized")
	}

	request := plaid.NewItemPublicTokenExchangeRequest(publicToken)

	resp, _, err := apiClient.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(*request).Execute()
	if err != nil {
		return "", "", fmt.Errorf("failed to exchange public token: %w", err)
	}

	return resp.GetAccessToken(), resp.GetItemId(), nil
}

// GetAccounts retrieves accounts for an item using the access token
func GetAccounts(ctx context.Context, accessToken string) ([]plaid.AccountBase, error) {
	if apiClient == nil {
		return nil, fmt.Errorf("plaid client not initialized")
	}

	request := plaid.NewAccountsGetRequest(accessToken)

	resp, _, err := apiClient.PlaidApi.AccountsGet(ctx).AccountsGetRequest(*request).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return resp.GetAccounts(), nil
}

// SyncTransactionsResult contains the results from a transaction sync operation
type SyncTransactionsResult struct {
	Added      []plaid.Transaction
	Modified   []plaid.Transaction
	Removed    []plaid.RemovedTransaction
	NextCursor string
	HasMore    bool
}

// SyncTransactions uses the Transactions Sync API to fetch transactions with cursor-based pagination
func SyncTransactions(
	ctx context.Context,
	accessToken string,
	cursor *string,
) (SyncTransactionsResult, error) {
	if apiClient == nil {
		return SyncTransactionsResult{}, fmt.Errorf("plaid client not initialized")
	}

	result := SyncTransactionsResult{}

	request := plaid.NewTransactionsSyncRequest(accessToken)
	if cursor != nil && *cursor != "" {
		request.SetCursor(*cursor)
	}

	resp, _, err := apiClient.PlaidApi.TransactionsSync(ctx).TransactionsSyncRequest(*request).Execute()
	if err != nil {
		return result, fmt.Errorf("failed to sync transactions: %w", err)
	}

	result.Added = resp.GetAdded()
	result.Modified = resp.GetModified()
	result.Removed = resp.GetRemoved()
	result.NextCursor = resp.GetNextCursor()
	result.HasMore = resp.GetHasMore()

	return result, nil
}

// GetItem retrieves item details
func GetItem(ctx context.Context, accessToken string) (plaid.ItemWithConsentFields, error) {
	if apiClient == nil {
		return plaid.ItemWithConsentFields{}, fmt.Errorf("plaid client not initialized")
	}

	request := plaid.NewItemGetRequest(accessToken)

	resp, _, err := apiClient.PlaidApi.ItemGet(ctx).ItemGetRequest(*request).Execute()
	if err != nil {
		return plaid.ItemWithConsentFields{}, fmt.Errorf("failed to get item: %w", err)
	}

	return resp.GetItem(), nil
}

// InstitutionsGetByID retrieves institution details by ID
func InstitutionsGetByID(ctx context.Context, institutionID string) (plaid.Institution, error) {
	if apiClient == nil {
		return plaid.Institution{}, fmt.Errorf("plaid client not initialized")
	}

	// Convert country codes string to slice
	countryCodes := []plaid.CountryCode{plaid.COUNTRYCODE_US}

	request := plaid.NewInstitutionsGetByIdRequest(institutionID, countryCodes)

	resp, _, err := apiClient.PlaidApi.InstitutionsGetById(ctx).InstitutionsGetByIdRequest(*request).Execute()
	if err != nil {
		return plaid.Institution{}, fmt.Errorf("failed to get institution: %w", err)
	}

	return resp.GetInstitution(), nil
}

// Helper functions to convert string slices to enum types

func convertCountryCodes(countryCodeStrs []string) []plaid.CountryCode {
	countryCodes := []plaid.CountryCode{}
	for _, cc := range countryCodeStrs {
		countryCodes = append(countryCodes, plaid.CountryCode(cc))
	}
	return countryCodes
}

func convertProducts(productStrs []string) []plaid.Products {
	products := []plaid.Products{}
	for _, p := range productStrs {
		products = append(products, plaid.Products(p))
	}
	return products
}
