package cloudflare

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var errRequestCaptured = errors.New("request captured")

type captureClient struct {
	req *http.Request
}

func (c *captureClient) Do(r *http.Request) (*http.Response, error) {
	c.req = r
	return nil, errRequestCaptured
}

// AWS region, FIPS and dual-stack settings a machine may carry for the AWS
// plugin must not change the R2 endpoint or signing region.
func TestNewR2ClientIgnoresAWSEnvironment(t *testing.T) {
	cases := map[string][]func(*config.LoadOptions) error{
		"no region":  nil,
		"aws region": {config.WithRegion("eu-west-1")},
		"fips":       {config.WithUseFIPSEndpoint(aws.FIPSEndpointStateEnabled)},
		"dual-stack": {config.WithUseDualStackEndpoint(aws.DualStackEndpointStateEnabled)},
	}
	for name, opts := range cases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("AWS_PROFILE", "")
			t.Setenv("AWS_DEFAULT_PROFILE", "")
			t.Setenv("AWS_REGION", "")
			t.Setenv("AWS_DEFAULT_REGION", "")
			t.Setenv("AWS_USE_FIPS_ENDPOINT", "")
			t.Setenv("AWS_USE_DUALSTACK_ENDPOINT", "")
			t.Setenv("AWS_CONFIG_FILE", "/dev/null")
			t.Setenv("AWS_SHARED_CREDENTIALS_FILE", "/dev/null")

			capture := &captureClient{}
			opts = append(opts,
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AK", "SK", "")),
				config.WithHTTPClient(capture),
			)
			cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
			if err != nil {
				t.Fatal(err)
			}
			client := newR2Client(cfg, "acct")

			_, err = client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{Bucket: aws.String("mybucket")})
			if capture.req == nil {
				t.Fatalf("request was never sent: %v", err)
			}
			if got, want := capture.req.URL.Host, "mybucket.acct.r2.cloudflarestorage.com"; got != want {
				t.Errorf("host = %q, want %q", got, want)
			}
			if auth := capture.req.Header.Get("Authorization"); !strings.Contains(auth, "/auto/s3/aws4_request") {
				t.Errorf("request not signed for region auto: %q", auth)
			}
		})
	}
}
