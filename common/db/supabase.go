package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// SupabaseClient represents a client for interacting with Supabase
type SupabaseClient struct {
	URL         string
	AnonKey     string
	ServiceRole string
	TenantID    string
	HTTPClient  *http.Client
}

// NewSupabaseClient creates a new Supabase client with the provided configuration
func NewSupabaseClient(tenantID string) (*SupabaseClient, error) {
	url := os.Getenv("SUPABASE_URL")
	anonKey := os.Getenv("SUPABASE_ANON_KEY")
	serviceRole := os.Getenv("SUPABASE_SERVICE_ROLE")

	if url == "" || anonKey == "" || serviceRole == "" {
		return nil, errors.New("missing required Supabase environment variables")
	}

	return &SupabaseClient{
		URL:         url,
		AnonKey:     anonKey,
		ServiceRole: serviceRole,
		TenantID:    tenantID,
		HTTPClient:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// QueryBuilder represents a builder for Supabase queries
type QueryBuilder struct {
	client     *SupabaseClient
	table      string
	selects    []string
	filters    []filter
	limitCount int
	orderBy    string
	orderDesc  bool
}

type filter struct {
	Column   string
	Operator string
	Value    interface{}
}

// Query creates a new query builder for the specified table
func (s *SupabaseClient) Query(table string) *QueryBuilder {
	return &QueryBuilder{
		client:  s,
		table:   table,
		selects: []string{"*"},
	}
}

// Select adds a select clause to the query
func (q *QueryBuilder) Select(columns ...string) *QueryBuilder {
	q.selects = columns
	return q
}

// Where adds a where clause to the query
func (q *QueryBuilder) Where(column, operator string, value interface{}) *QueryBuilder {
	q.filters = append(q.filters, filter{
		Column:   column,
		Operator: operator,
		Value:    value,
	})
	return q
}

// Limit adds a limit clause to the query
func (q *QueryBuilder) Limit(count int) *QueryBuilder {
	q.limitCount = count
	return q
}

// Order adds an order clause to the query
func (q *QueryBuilder) Order(column string, desc bool) *QueryBuilder {
	q.orderBy = column
	q.orderDesc = desc
	return q
}

// Execute executes the query and returns the results
func (q *QueryBuilder) Execute(result interface{}) error {
	// Build the URL for the query
	url := fmt.Sprintf("%s/rest/v1/%s", q.client.URL, q.table)

	// Create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Add("apikey", q.client.AnonKey)
	req.Header.Add("Authorization", "Bearer "+q.client.ServiceRole)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Prefer", "return=representation")

	// Add query parameters
	query := req.URL.Query()

	// Add select
	if len(q.selects) > 0 {
		query.Add("select", strings.Join(q.selects, ","))
	}

	// Add tenant filter
	if q.client.TenantID != "" {
		query.Add("tenant_id", "eq."+q.client.TenantID)
	}

	// Add filters
	for _, f := range q.filters {
		query.Add(f.Column, fmt.Sprintf("%s.%v", f.Operator, f.Value))
	}

	// Add limit
	if q.limitCount > 0 {
		query.Add("limit", fmt.Sprintf("%d", q.limitCount))
	}

	// Add order
	if q.orderBy != "" {
		order := q.orderBy
		if q.orderDesc {
			order = order + ".desc"
		} else {
			order = order + ".asc"
		}
		query.Add("order", order)
	}

	req.URL.RawQuery = query.Encode()

	// Execute the request
	resp, err := q.client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase request failed with status: %d", resp.StatusCode)
	}

	// Decode the response
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(result)
}

// Insert inserts a new record into the table
func (s *SupabaseClient) Insert(ctx context.Context, table string, data interface{}) error {
	url := fmt.Sprintf("%s/rest/v1/%s", s.URL, table)

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create the request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	// Add context
	req = req.WithContext(ctx)

	// Add headers
	req.Header.Add("apikey", s.AnonKey)
	req.Header.Add("Authorization", "Bearer "+s.ServiceRole)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Prefer", "return=representation")

	// Execute the request
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase insert failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Update updates an existing record in the table
func (s *SupabaseClient) Update(ctx context.Context, table, idColumn, id string, data interface{}) error {
	url := fmt.Sprintf("%s/rest/v1/%s?%s=eq.%s", s.URL, table, idColumn, id)

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create the request
	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	// Add context
	req = req.WithContext(ctx)

	// Add headers
	req.Header.Add("apikey", s.AnonKey)
	req.Header.Add("Authorization", "Bearer "+s.ServiceRole)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Prefer", "return=representation")

	// Execute the request
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase update failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Delete deletes a record from the table
func (s *SupabaseClient) Delete(ctx context.Context, table, idColumn, id string) error {
	url := fmt.Sprintf("%s/rest/v1/%s?%s=eq.%s", s.URL, table, idColumn, id)

	// Create the request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Add context
	req = req.WithContext(ctx)

	// Add headers
	req.Header.Add("apikey", s.AnonKey)
	req.Header.Add("Authorization", "Bearer "+s.ServiceRole)
	req.Header.Add("Content-Type", "application/json")

	// Execute the request
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase delete failed with status: %d", resp.StatusCode)
	}

	return nil
}
