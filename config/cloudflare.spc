connection "cloudflare" {
  plugin = "cloudflare"

  # API Token for your Cloudflare account
  # See https://support.cloudflare.com/hc/en-us/articles/200167836-Managing-API-Tokens-and-Keys#12345680
  #token   = "YOUR_CLOUDFLARE_API_TOKEN"

  # Alternatively, use your (legacy) email and API Key.
  #email   = "YOUR_CLOUDFLARE_EMAIL"
  #api_key = "YOUR_CLOUDFLARE_API_KEY"

  # Access Key ID and Secret Access Key to access Cloudflare R2
  # See https://developers.cloudflare.com/r2/data-access/s3-api/tokens/
  # access_key = "YOUR_R2_ACCESS_KEY_ID"
  # secret_key = "YOUR_R2_SECRET_ACCESS_KEY"
}
