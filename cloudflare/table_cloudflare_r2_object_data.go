package cloudflare

import (
	"context"
	"encoding/base64"
	"io"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareR2ObjectData(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "cloudflare_r2_object_data",
		Description:      "List Cloudflare R2 Objects in buckets by bucket name",
		DefaultTransform: transform.FromCamel().NullIfZero(),
		Get: &plugin.GetConfig{
			Hydrate: getR2ObjectData,
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
				Name:        "sse_customer_algorithm",
				Description: "If server-side encryption with a customer-provided encryption key was requested, the response will include this header confirming the encryption algorithm used.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectDataBody,
				Transform:   transform.FromField("SSECustomerAlgorithm"),
			},
			{
				Name:        "sse_kms_key_id",
				Description: "If present, specifies the ID of the Amazon Web Services Key Management Service (Amazon Web Services KMS) symmetric customer managed key that was used for the object.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectDataBody,
				Transform:   transform.FromField("SSEKMSKeyId"),
			},
			{
				Name:        "data",
				Description: "The raw bytes of the object as a string. An UTF8 encoded string is sent, if the bytes entirely consists of valid UTF8 runes, an UTF8 is sent otherwise the bas64 encoding of the bytes is sent.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getR2ObjectDataBody,
				Transform:   transform.FromField("Body").Transform(parseBody),
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

type s3ObjectContent struct {
	types.Object
	AccountID *string
	Bucket    *string
}

//// LIST FUNCTION

func getR2ObjectData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	object, err := conn.ListObjectsV2(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_r2_object_data.getR2ObjectData", "api_error", err)
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

//// HYDRATE FUNCTIONS

func getR2ObjectDataBody(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	data := h.Item.(*s3ObjectContent)

	// Get R2 client
	conn, err := getR2Client(ctx, d, *data.AccountID)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.GetObjectInput{
		Bucket: data.Bucket,
		Key:    data.Object.Key,
	}

	object, err := conn.GetObject(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_r2_object_data.getR2ObjectDataBody", "api_error", err)
		return nil, err
	}

	return object, nil
}

//// TRANSFORM FUNCTIONS

func parseBody(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*s3.GetObjectOutput)

	body, err := io.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}

	if utf8.Valid(body) {
		return string(body), nil
	}

	return base64.StdEncoding.EncodeToString(body), nil
}
