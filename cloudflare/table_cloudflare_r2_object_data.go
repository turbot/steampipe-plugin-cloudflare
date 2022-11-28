package cloudflare

import (
	"context"
	"encoding/base64"
	"io"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareR2ObjectData(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "cloudflare_r2_object_data",
		Description:      "List content of specific Cloudflare R2 objects by bucket name",
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
				Name:        "last_modified",
				Description: "Specifies the time when the object is last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "content_length",
				Description: "Size of the body in bytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "etag",
				Description: "A a hash of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ETag"),
			},
			{
				Name:        "content_encoding",
				Description: "Specifies what content encodings have been applied to the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_type",
				Description: "A standard MIME type describing the format of the object data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sse_customer_algorithm",
				Description: "If server-side encryption with a customer-provided encryption key was requested, the response will include this header confirming the encryption algorithm used.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSECustomerAlgorithm"),
			},
			{
				Name:        "sse_kms_key_id",
				Description: "If present, specifies the ID of the Amazon Web Services Key Management Service (Amazon Web Services KMS) symmetric customer managed key that was used for the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSEKMSKeyId"),
			},
			{
				Name:        "data",
				Description: "The raw bytes of the object as a string. An UTF8 encoded string is sent, if the bytes entirely consists of valid UTF8 runes, an UTF8 is sent otherwise the bas64 encoding of the bytes is sent.",
				Type:        proto.ColumnType_STRING,
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
	s3.GetObjectOutput
	AccountID *string
	Bucket    *string
}

//// HYDRATE FUNCTIONS

func getR2ObjectData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var accountID, bucket, key *string
	if h.Item != nil {
		data := h.Item.(*s3ObjectMetadata)
		accountID = data.AccountID
		bucket = data.Bucket
		key = data.Key
	} else {
		accountID = aws.String(d.KeyColumnQualString("account_id"))
		bucket = aws.String(d.KeyColumnQualString("bucket"))
		key = aws.String(d.KeyColumnQualString("key"))
	}

	// get R2 client
	conn, err := getR2Client(ctx, d, *accountID)
	if err != nil {
		return nil, err
	}

	// execute list call
	input := &s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	}

	object, err := conn.GetObject(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_r2_object_data.getR2ObjectData", "api_error", err)
		return nil, err
	}

	return &s3ObjectContent{*object, accountID, bucket}, nil
}

//// TRANSFORM FUNCTIONS

func parseBody(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*s3ObjectContent)

	body, err := io.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}

	if utf8.Valid(body) {
		return string(body), nil
	}

	return base64.StdEncoding.EncodeToString(body), nil
}
