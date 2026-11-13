package cloudflare

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
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

func TestNewR2ClientRegion(t *testing.T) {
	for _, envRegion := range []string{"", "eu-west-1"} {
		t.Run("env region "+envRegion, func(t *testing.T) {
			capture := &captureClient{}
			cfg := aws.Config{
				Region:      envRegion,
				Credentials: credentials.NewStaticCredentialsProvider("AK", "SK", ""),
				HTTPClient:  capture,
			}
			client := newR2Client(cfg, "acct")

			_, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{Bucket: aws.String("mybucket")})
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
