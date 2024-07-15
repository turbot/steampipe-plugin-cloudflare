package cloudflare

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareR2Object(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "cloudflare_r2_object",
		Description:      "List Cloudflare R2 Objects by bucket name",
		DefaultTransform: transform.FromCamel().NullIfZero(),
		List: &plugin.ListConfig{
			Hydrate: listR2Objects,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "account_id", Require: plugin.Required},
				{Name: "bucket", Require: plugin.Required},
				{Name: "key", Require: plugin.Optional},
				{Name: "prefix", Require: plugin.Optional},
			},
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "key",
				Description: "The name assigned to an object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Object.Key"),
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
				Name:        "etag",
				Description: "A a hash of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Object.ETag"),
			},
			{
				Name:        "version",
				Description: "Version of the object.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectData,
			},
			{
				Name:        "bucket_key_enabled",
				Description: "Indicates whether the object uses an R2 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS).",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getR2ObjectData,
			},
			{
				Name:        "content_encoding",
				Description: "Specifies what content encodings have been applied to the object.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectData,
			},
			{
				Name:        "content_type",
				Description: "A standard MIME type describing the format of the object data.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectData,
			},
			{
				Name:        "content_length",
				Description: "Size of the body in bytes.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getR2ObjectData,
			},
			{
				Name:        "delete_marker",
				Description: "Specifies whether the object retrieved was (true) or was not (false) a delete marker.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getR2ObjectData,
			},
			{
				Name:        "restore",
				Description: "Provides information about object restoration action and expiration time of the restored object copy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectData,
			},
			{
				Name:        "server_side_encryption",
				Description: "The server-side encryption algorithm used when storing this object in Cloudflare R2.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectData,
			},
			{
				Name:        "website_redirection_location",
				Description: "If the bucket is configured as a website, redirects requests for this object  to another object in the same bucket or to an external URL.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectData,
			},
			{
				Name:        "metadata",
				Description: "A map of metadata to store with the object in R2.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getR2ObjectData,
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
				Transform:   transform.FromField("Object.Key"),
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
			{
				Name:        "prefix",
				Description: "The prefix of the key of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("prefix"),
			},
		}),
	}
}

type s3ObjectMetadata = struct {
	types.Object
	AccountID *string
	Bucket    *string
}

//// LIST FUNCTION

func listR2Objects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountID := d.EqualsQualString("account_id")
	bucketName := d.EqualsQualString("bucket")

	// return nil, if either of the required columns is missing
	if accountID == "" || bucketName == "" {
		return nil, nil
	}

	// get R2 client
	conn, err := getR2Client(ctx, d, accountID)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}

	// set prefixes
	prefix := d.EqualsQualString("prefix")
	key := d.EqualsQualString("key")

	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	// if key is provided in qual, use key as prefix
	// also, set the max limit to 1 for exact match
	if key != "" {
		input.Prefix = aws.String(key)
		input.MaxKeys = 1
	}

	// the owner field is not present in listV2 by default
	// set it true if the column is passed in qual
	if helpers.StringSliceContains(d.QueryContext.Columns, "owner") {
		input.FetchOwner = true
	}

	object, err := conn.ListObjectsV2(ctx, input)
	if err != nil {
		var a smithy.APIError
		if errors.As(err, &a) {
			if a.ErrorCode() == "InvalidBucketName" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("cloudflare_r2_object.getR2Object", "api_error", err)
		return nil, err
	}

	for _, i := range object.Contents {
		d.StreamListItem(ctx, &s3ObjectMetadata{
			Object:    i,
			AccountID: aws.String(accountID),
			Bucket:    aws.String(bucketName),
		})

		// context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
