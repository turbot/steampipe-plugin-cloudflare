package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	cloudflare4 "github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type Organization struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"`
}

type UserDetails struct {
	ID                   string         `json:"id"`                                // ID of the user.
	Email                string         `json:"email"`                             // Email of the user.
	Username             string         `json:"username"`                          // Username (often a hashed ID) of the user.
	FirstName            *string        `json:"first_name"`                        // First name of the user (nullable).
	LastName             *string        `json:"last_name"`                         // Last name of the user (nullable).
	Telephone            *string        `json:"telephone"`                         // Telephone number (nullable).
	Country              *string        `json:"country"`                           // Country (nullable).
	Zipcode              *string        `json:"zipcode"`                           // Zipcode (nullable).
	TwoFactorAuthEnabled bool           `json:"two_factor_authentication_enabled"` // True if 2FA is enabled.
	TwoFactorAuthLocked  bool           `json:"two_factor_authentication_locked"`  // True if 2FA is locked.
	CreatedOn            string         `json:"created_on"`                        // ISO8601 timestamp.
	ModifiedOn           string         `json:"modified_on"`                       // ISO8601 timestamp.
	Organizations        []Organization `json:"organizations"`                     // List of organizations the user is part of.
	HasProZones          bool           `json:"has_pro_zones"`
	HasBusinessZones     bool           `json:"has_business_zones"`
	HasEnterpriseZones   bool           `json:"has_enterprise_zones"`
	Suspended            bool           `json:"suspended"`
	Betas                []string       `json:"betas"` // List of beta features enabled.
}

// Create Cloudflare R2 API client
func getR2Client(ctx context.Context, d *plugin.QueryData, accountID string) (*s3.Client, error) {
	sessionCacheKey := fmt.Sprintf("session-v2-%s", accountID)
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*s3.Client), nil
	}

	cloudflareConfig := GetConfig(d.Connection)
	var accessKey, secret string

	if cloudflareConfig.AccessKey != nil {
		accessKey = *cloudflareConfig.AccessKey
	}

	if cloudflareConfig.SecretKey != nil {
		secret = *cloudflareConfig.SecretKey
	}

	if accessKey == "" || secret == "" {
		return nil, errors.New("cloudflare R2 API credentials not found. Edit your connection to configure AccessKey and Secret, and then restart Steampipe")
	}

	r2EndpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithEndpointResolverWithOptions(r2EndpointResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secret, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}
	client := s3.NewFromConfig(cfg)

	d.ConnectionManager.Cache.Set(sessionCacheKey, client)

	return client, nil
}

func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		errMsg := err.Error()
		for _, msg := range notFoundErrors {
			if strings.Contains(errMsg, msg) {
				return true
			}
		}
		return false
	}
}

func shouldRetryError(err error) bool {
	if err.Error() == "rate limit exceeded" {
		return true
	}
	return false
}

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize
// since getUser is a call, caching should be per connection
var getUserMemoized = plugin.HydrateFunc(getUserUncached).Memoize(memoize.WithCacheKeyFunction(getUserCacheKey))

// Build a cache key for the call to getUser.
func getUserCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getUserInfo"
	return key, nil
}

func getUserInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	userDetails, err := getUserMemoized(ctx, d, h)
	if err != nil {
		return nil, err
	}

	return userDetails, nil
}

func getUserUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get user details using the v4 API
	resp, err := conn.User.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Log response for debugging
	plugin.Logger(ctx).Debug("cloudflare_utils.getUserUncached", "user_response", resp)

	// Create a default UserDetails object
	userDetails := &UserDetails{}

	// Try to extract ID from the response
	if resp != nil {
		// For debugging purposes
		plugin.Logger(ctx).Debug("cloudflare_utils.getUserUncached", "resp_type", fmt.Sprintf("%T", resp))

		// Try to extract values using type assertion
		if idValue, ok := getValueFromInterface(*resp, "id"); ok {
			userDetails.ID = idValue.(string)
		}
		if emailValue, ok := getValueFromInterface(*resp, "email"); ok {
			userDetails.Email = emailValue.(string)
		}
		if usernameValue, ok := getValueFromInterface(*resp, "username"); ok {
			userDetails.Username = usernameValue.(string)
		}
	}

	plugin.Logger(ctx).Debug("cloudflare_utils.getUserUncached", "user_details", userDetails)
	return userDetails, nil
}

// Helper function to safely extract values from interfaces
func getValueFromInterface(obj interface{}, key string) (interface{}, bool) {
	// Try to convert to map[string]interface{}
	if m, ok := obj.(map[string]interface{}); ok {
		value, exists := m[key]
		return value, exists
	}

	// If it's a struct, try reflection to get the field
	// For now we'll just return false, but we could add reflection here if needed
	return nil, false
}

func getUserId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	res, err := getUserInfo(ctx, d, h)
	if err != nil {
		return nil, err
	}

	userDetails, ok := res.(*UserDetails)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T in getUserId, expected UserDetails", res)
	}

	if userDetails.ID == "" {
		return nil, fmt.Errorf("user ID is empty or could not be retrieved")
	}

	return userDetails.ID, nil
}

// getRequestTimeout returns the request timeout from config or environment variables
func getRequestTimeout(config cloudflareConfig) time.Duration {
	// Default timeout
	defaultTimeout := 30 * time.Second

	// Get max MaxRequestTimeout
	if config.MaxRequestTimeout != nil {
		return time.Duration(*config.MaxRequestTimeout) * time.Second
	}

	if timeoutStr := os.Getenv("CLOUDFLARE_MAX_REQUEST_TIMEOUT"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil && timeout > 0 {
			return time.Duration(timeout) * time.Second
		}
	}

	return defaultTimeout
}

// getMaxRetries returns the max retries from config or environment variables
func getMaxRetries(config cloudflareConfig) int {
	// Default retries
	defaultRetries := 3

	if config.MaxRetries != nil {
		return *config.MaxRetries
	}

	// Check environment variable
	if retriesStr := os.Getenv("CLOUDFLARE_MAX_RETRIES"); retriesStr != "" {
		if retries, err := strconv.Atoi(retriesStr); err == nil && retries >= 0 {
			return retries
		}
	}

	return defaultRetries
}

func connectV4(ctx context.Context, d *plugin.QueryData) (*cloudflare4.Client, error) {
	// Get the config
	cloudflareConfig := GetConfig(d.Connection)

	// Get timeout and retry configuration
	requestTimeout := getRequestTimeout(cloudflareConfig)
	maxRetries := getMaxRetries(cloudflareConfig)

	// Prepare client options with timeout and retries
	var clientOptions []option.RequestOption
	clientOptions = append(clientOptions, option.WithRequestTimeout(requestTimeout))
	clientOptions = append(clientOptions, option.WithMaxRetries(maxRetries))

	// First: check for the token in config
	if cloudflareConfig.Token != nil {
		clientOptions = append(clientOptions, option.WithAPIToken(*cloudflareConfig.Token))
		client := cloudflare4.NewClient(clientOptions...)
		return client, nil
	}

	// Second: Email + API Key from config
	if cloudflareConfig.Email != nil && cloudflareConfig.APIKey != nil {
		clientOptions = append(clientOptions, option.WithAPIKey(*cloudflareConfig.APIKey))
		clientOptions = append(clientOptions, option.WithAPIEmail(*cloudflareConfig.Email))
		client := cloudflare4.NewClient(clientOptions...)
		return client, nil
	}

	// Third: CLOUDFLARE_API_TOKEN (like Terraform)
	token, ok := os.LookupEnv("CLOUDFLARE_API_TOKEN")
	if ok && token != "" {
		clientOptions = append(clientOptions, option.WithAPIToken(token))
		client := cloudflare4.NewClient(clientOptions...)
		return client, nil
	}

	// Fourth: CLOUDFLARE_EMAIL / CLOUDFLARE_API_KEY (like Terraform)
	email, ok := os.LookupEnv("CLOUDFLARE_EMAIL")
	if ok && email != "" {
		key, ok := os.LookupEnv("CLOUDFLARE_API_KEY")
		if ok && key != "" {
			clientOptions = append(clientOptions, option.WithAPIKey(key))
			clientOptions = append(clientOptions, option.WithAPIEmail(email))
			client := cloudflare4.NewClient(clientOptions...)
			return client, nil
		}
	}

	// Fifth: CF_API_TOKEN (like flarectl and Go SDK)
	token, ok = os.LookupEnv("CF_API_TOKEN")
	if ok && token != "" {
		clientOptions = append(clientOptions, option.WithAPIToken(token))
		client := cloudflare4.NewClient(clientOptions...)
		return client, nil
	}

	// Sixth: CF_EMAIL / CF_API_KEY (like flarectl / Go SDK)
	email, ok = os.LookupEnv("CF_API_EMAIL")
	if ok && email != "" {
		key, ok := os.LookupEnv("CF_API_KEY")
		if ok && key != "" {
			clientOptions = append(clientOptions, option.WithAPIKey(key))
			clientOptions = append(clientOptions, option.WithAPIEmail(email))
			client := cloudflare4.NewClient(clientOptions...)
			return client, nil
		}
	}

	return nil, errors.New("Cloudflare API credentials not found. Edit your connection configuration file and then restart Steampipe.")
}

// This function is used to extract the extra field that are not in the part of API response structure.
// To keep the table column backward compatible, we need to extract the extra field from the API response.
// This function is used by the tables:
//    - cloudflare_account
//    - cloudflare_zone
func toMap(in interface{}) (map[string]interface{}, error) {
	// Nil guard – return nil map, nil error so callers can test for nil safely.
	if in == nil {
		return nil, nil
	}

	// Handle the most common cases first.
	switch v := in.(type) {
	case map[string]interface{}: // already what we want
		return v, nil
	case string: // JSON in string
		return unmarshalString(v)
	case []byte: // JSON in raw bytes
		return unmarshalBytes(v)
	case json.RawMessage: // alias of []byte
		return unmarshalBytes(v)
	}

	// If we’re given a pointer, de-reference once and try again
	// (covers *accounts.accountJSON, *SomeStruct, etc.).
	rv := reflect.ValueOf(in)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		return toMap(rv.Elem().Interface())
	}

	// Catch-all: try json.Marshal.  This works for:
	//   • structs with exported fields
	//   • types that implement json.Marshaler (e.g. apijson.Field)
	//   • maps/arrays/slices of the above
	b, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("toMap: unsupported type %T: %w", in, err)
	}
	return unmarshalBytes(b)
}

func unmarshalString(s string) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	return m, err
}

func unmarshalBytes(b []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	return m, err
}
