package cloudflare

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	// "github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*cloudflare.API, error) {

	cloudflareConfig := GetConfig(d.Connection)

	// First: check for the token
	if cloudflareConfig.Token != nil {
		return cloudflare.NewWithAPIToken(*cloudflareConfig.Token)
	}

	// Second: Email + API Key
	if cloudflareConfig.Email != nil && cloudflareConfig.APIKey != nil {
		return cloudflare.New(*cloudflareConfig.APIKey, *cloudflareConfig.Email)
	}

	// Third: CLOUDFLARE_API_TOKEN (like Terraform)
	token, ok := os.LookupEnv("CLOUDFLARE_API_TOKEN")
	if ok && token != "" {
		return cloudflare.NewWithAPIToken(token)
	}

	// Fourth: CLOUDFLARE_EMAIL / CLOUDFLARE_API_KEY (like Terraform)
	email, ok := os.LookupEnv("CLOUDFLARE_EMAIL")
	if ok && email != "" {
		key, ok := os.LookupEnv("CLOUDFLARE_API_KEY")
		if ok && key != "" {
			return cloudflare.New(key, email)
		}
	}

	// Fifth: CF_API_TOKEN (like flarectl and Go SDK)
	token, ok = os.LookupEnv("CF_API_TOKEN")
	if ok && token != "" {
		return cloudflare.NewWithAPIToken(token)
	}

	// Sixth: CF_EMAIL / CF_API_KEY (like flarectl / Go SDK)
	email, ok = os.LookupEnv("CF_API_EMAIL")
	if ok && email != "" {
		key, ok := os.LookupEnv("CF_API_KEY")
		if ok && key != "" {
			return cloudflare.New(key, email)
		}
	}

	return nil, errors.New("Cloudflare API credentials not found. Edit your connection configuration file and then restart Steampipe.")
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
	if cloudflareErr, ok := err.(*cloudflare.APIRequestError); ok {
		return cloudflareErr.ClientRateLimited()
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

func getUserId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	res, err := getUserInfo(ctx, d, h)
	if err != nil {
		return nil, err
	}

	data := res.(cloudflare.User)
	return data.ID, nil
}

func getUserInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	UserId, err := getUserMemoized(ctx, d, h)
	if err != nil {
		return nil, err
	}

	return UserId, nil
}

func getUserUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	item, err := conn.UserDetails(ctx)
	if err != nil {
		return nil, err
	}

	return item, nil
}
