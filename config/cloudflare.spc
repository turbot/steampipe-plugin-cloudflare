connection "cloudflare" {
  plugin   = "cloudflare"

  # The account ID to retrieve account scoped resources for
  # May alternatively be set via the CLOUDFLARE_ACCOUNT_ID env variable
  #account_id = "YOUR_ACCOUNT_ID"

  # API Token for your Cloudflare account
  # See https://support.cloudflare.com/hc/en-us/articles/200167836-Managing-API-Tokens-and-Keys#12345680
  #token   = "YOUR_CLOUDFLARE_API_TOKEN"

  # Alternatively, use your (legacy) email and API Key.
  #email   = "YOUR_CLOUDFLARE_EMAIL"
  #api_key = "YOUR_CLOUDFLARE_API_KEY"
}
