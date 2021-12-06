package cloudflare

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
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
