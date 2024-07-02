package cloudflare

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareR2Bucket(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "cloudflare_r2_bucket",
		Description:      "Cloudflare R2 Buckets",
		DefaultTransform: transform.FromCamel().NullIfZero(),
		List: &plugin.ListConfig{
			Hydrate:    listR2Buckets,
			KeyColumns: plugin.SingleColumn("account_id"),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "account_id"}),
			Hydrate:    getR2Bucket,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The user friendly name of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The date and time when bucket was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "server_side_encryption_configuration",
				Description: "The default encryption configuration for the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketEncryption,
				Transform:   transform.FromField("ServerSideEncryptionConfiguration"),
			},
			{
				Name:        "cors",
				Description: "The CORS configuration for the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBucketCORS,
				Transform:   transform.FromValue(),
			},

			// steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "account_id",
				Description: "ID of the account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("account_id"),
			},
			{
				Name:        "region",
				Description: "The date and time when bucket was created.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBucketLocation,
				Transform:   transform.FromField("LocationConstraint"),
			},
		}),
	}
}

type BucketData = struct {
	types.Bucket
	AccountId string
}

//// LIST FUNCTION

func listR2Buckets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Get cloudflare account data
	accountID := d.EqualsQualString("account_id")
	if accountID == "" {
		return nil, nil
	}

	// What happens if we try to access a bucket in an account where the R2 plan is not enabled?
	// Can we do something with this error?
	// Error: operation error S3: ListBuckets, exceeded maximum number of attempts, 3, https response error StatusCode: 0, RequestID: , HostID: , request send failed, Get "https://<account>.r2.cloudflarestorage.com/": remote error: tls: handshake failure (SQLSTATE HV000)

	// Get R2 client
	conn, err := getR2Client(ctx, d, accountID)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.ListBucketsInput{}
	bucketsResult, err := conn.ListBuckets(ctx, input)
	if err != nil {
		// Get "https://<account>.r2.cloudflarestorage.com/": remote error: tls: handshake failure (SQLSTATE HV000)
		if strings.Contains(err.Error(), "tls: handshake failure") {
			return nil, nil
		}

		plugin.Logger(ctx).Error("cloudflare_r2_bucket.listR2Buckets", "api_error", err)
		return nil, err
	}

	for _, bucket := range bucketsResult.Buckets {
		d.StreamListItem(ctx, BucketData{bucket, accountID})

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

// do not have a get call for R2 bucket.
// using list api call to create get function
func getR2Bucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountID := d.EqualsQualString("account_id")
	bucketName := d.EqualsQualString("name")

	// Return nil if either of the required params is not passed in the qual
	if accountID == "" || bucketName == "" {
		return nil, nil
	}

	// Get R2 client
	conn, err := getR2Client(ctx, d, accountID)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.ListBucketsInput{}
	listBucketsOutput, err := conn.ListBuckets(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_r2_bucket.getR2Bucket", "api_error", err)
		return nil, err
	}

	for _, bucket := range listBucketsOutput.Buckets {
		if bucket.Name == aws.String(bucketName) {
			return BucketData{bucket, accountID}, nil
		}
	}

	return nil, nil
}

func getBucketLocation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get cloudflare account data
	bucketData := h.Item.(BucketData)

	// Get R2 client
	conn, err := getR2Client(ctx, d, bucketData.AccountId)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.GetBucketLocationInput{
		Bucket: bucketData.Name,
	}
	location, err := conn.GetBucketLocation(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_r2_bucket.getBucketLocation", "api_error", err)
		return nil, err
	}

	if location != nil && location.LocationConstraint != "" {
		return location, nil
	}

	// Buckets in us-east-1 have a LocationConstraint of null
	return &s3.GetBucketLocationOutput{
		LocationConstraint: "us-east-1",
	}, nil
}

func getBucketEncryption(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get cloudflare account data
	bucketData := h.Item.(BucketData)

	// Get R2 client
	conn, err := getR2Client(ctx, d, bucketData.AccountId)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.GetBucketEncryptionInput{
		Bucket: bucketData.Name,
	}
	encryption, err := conn.GetBucketEncryption(ctx, input)
	if err != nil {
		var a smithy.APIError
		if errors.As(err, &a) {
			if a.ErrorCode() == "ServerSideEncryptionConfigurationNotFoundError" {
				return nil, nil
			}
		}
		return nil, err
	}

	return encryption, nil
}

func getBucketCORS(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get cloudflare account data
	bucketData := h.Item.(BucketData)

	// Get R2 client
	conn, err := getR2Client(ctx, d, bucketData.AccountId)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.GetBucketCorsInput{
		Bucket: bucketData.Name,
	}
	cors, err := conn.GetBucketCors(ctx, input)
	if err != nil {
		var a smithy.APIError
		if errors.As(err, &a) {
			if a.ErrorCode() == "NoSuchCORSConfiguration" {
				return nil, nil
			}
		}
		return nil, err
	}

	return cors, nil
}
