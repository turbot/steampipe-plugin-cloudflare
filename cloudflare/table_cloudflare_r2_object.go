package cloudflare

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareR2Object(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "cloudflare_r2_object",
		Description:      "List Cloudflare R2 Objects in R2 buckets by bucket name",
		DefaultTransform: transform.FromCamel().NullIfZero(),
		Get: &plugin.GetConfig{
			Hydrate: getR2Object,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "account_id", Require: plugin.Required},
				{Name: "bucket", Require: plugin.Required},
				{Name: "key", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "key",
				Description: "The name assigned to an object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("key"),
			},
			{
				Name:        "etag",
				Description: "A a hash of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Object.ETag"),
			},
			{
				Name:        "last_modified",
				Description: "Specifies the time when the object is last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Object.LastModified"),
			},
			{
				Name:        "size",
				Description: "Size in bytes of the object",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Object.Size"),
			},
			{
				Name:        "storage_class",
				Description: "Provides storage class information of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Object.StorageClass"),
			},
			{
				Name:        "version",
				Description: "Version of the object.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "bucket_key_enabled",
				Description: "Indicates whether the object uses an S3 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS).",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "content_encoding",
				Description: "Specifies what content encodings have been applied to the object.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "content_type",
				Description: "A standard MIME type describing the format of the object data.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "content_length",
				Description: "Size of the body in bytes.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "delete_marker",
				Description: "Specifies whether the object retrieved was (true) or was not (false) a delete marker.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "replication_status",
				Description: "Amazon S3 can return this if your request involves a bucket that is either a source or destination in a replication rule.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "restore",
				Description: "Provides information about object restoration action and expiration time of the restored object copy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "server_side_encryption",
				Description: "The server-side encryption algorithm used when storing this object in Amazon S3.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "website_redirection_location",
				Description: "If the bucket is configured as a website, redirects requests for this object  to another object in the same bucket or to an external URL.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "metadata",
				Description: "A map of metadata to store with the object in S3.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getR2ObjectDataBody,
			},
			{
				Name:        "owner",
				Description: "The owner of the object",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Object.Owner"),
			},

			// steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("key"),
			},
			{
				Name:        "account_id",
				Description: "ID of the account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("account_id"),
			},
			{
				Name:        "bucket",
				Description: "The name of the container bucket of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("bucket"),
			},
		},
	}
}

//// LIST FUNCTION

func getR2Object(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountID := d.KeyColumnQualString("account_id")
	bucketName := d.KeyColumnQualString("bucket")
	key := d.KeyColumnQualString("key")

	if accountID == "" || bucketName == "" || key == "" {
		return nil, nil
	}

	// Get R2 client
	conn, err := getR2Client(ctx, d, accountID)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(key),
	}

	// The owner field is not present in listV2 by default
	// Set it true if the column is passed in qual
	if helpers.StringSliceContains(d.QueryContext.Columns, "owner") {
		input.FetchOwner = true
	}

	object, err := conn.ListObjectsV2(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_r2_object.getR2Object", "api_error", err)
		return nil, err
	}

	for _, i := range object.Contents {
		if *i.Key == key {
			return &s3ObjectContent{
				Object:    object.Contents[0],
				AccountID: aws.String(accountID),
				Bucket:    aws.String(bucketName),
			}, nil
		}
	}

	return nil, nil
}
