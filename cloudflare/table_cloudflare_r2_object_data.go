package cloudflare

import (
	"context"
	"encoding/base64"
	"io"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
		Columns: commonColumns([]*plugin.Column{
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
				Name:        "accept_ranges",
				Description: "Indicates that a range of bytes was specified.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bucket_key_enabled",
				Description: "Indicates whether the object uses an S3 Bucket Key for server-side encryption with Amazon Web Services KMS (SSE-KMS).",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "cache_control",
				Description: "Specifies caching behavior along the request/reply chain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "checksum_crc32",
				Description: "The base64-encoded, 32-bit CRC32 checksum of the object. This will only be present if it was uploaded with the object. With multipart uploads, this may not be a checksum value of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ChecksumCRC32"),
			},
			{
				Name:        "checksum_crc32c",
				Description: "The base64-encoded, 32-bit CRC32C checksum of the object. This will only be present if it was uploaded with the object. With multipart uploads, this may not be a checksum value of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ChecksumCRC32C"),
			},
			{
				Name:        "checksum_sha1",
				Description: "The base64-encoded, 160-bit SHA-1 digest of the object. This will only be present if it was uploaded with the object. With multipart uploads, this may not be a checksum value of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ChecksumSHA1"),
			},
			{
				Name:        "checksum_sha256",
				Description: "The base64-encoded, 256-bit SHA-256 digest of the object. This will only be present if it was uploaded with the object. With multipart uploads, this may not be a checksum value of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ChecksumSHA256"),
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
				Name:        "content_disposition",
				Description: "Specifies presentational information for the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_language",
				Description: "The language the content is in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_length",
				Description: "Size of the body in bytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "content_range",
				Description: "The portion of the object returned in the response.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_type",
				Description: "A standard MIME type describing the format of the object data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delete_marker",
				Description: "Specifies whether the object retrieved was (true) or was not (false) a Delete Marker. If false, this response header does not appear in the response.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "expiration",
				Description: "If the object expiration is configured (see PUT Bucket lifecycle), the response includes this header. It includes the expiry-date and rule-id key-value pairs providing object expiration information. The value of the rule-id is URL-encoded.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "expires",
				Description: "The date and time at which the object is no longer cacheable.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "metadata",
				Description: "A map of metadata to store with the object in S3.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "missing_meta",
				Description: "A map of metadata to store with the object in S3.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "object_lock_legal_hold_status",
				Description: "Indicates whether this object has an active legal hold. This field is only returned if you have permission to view an object's legal hold status.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "object_lock_mode",
				Description: "The Object Lock mode currently in place for this object.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "object_lock_retain_until_date",
				Description: "The date and time when this object's Object Lock will expire.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "parts_count",
				Description: "The count of parts this object has. This value is only returned if you specify partNumber in your request and the object was uploaded as a multipart upload.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "replication_status",
				Description: "Amazon S3 can return this if your request involves a bucket that is either a source or destination in a replication rule.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "request_charged",
				Description: "If present, indicates that the requester was successfully charged for the request.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "restore",
				Description: "Provides information about object restoration action and expiration time of the restored object copy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sse_customer_algorithm",
				Description: "If server-side encryption with a customer-provided encryption key was requested, the response will include this header confirming the encryption algorithm used.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSECustomerAlgorithm"),
			},
			{
				Name:        "sse_customer_key_md5",
				Description: "If server-side encryption with a customer-provided encryption key was requested, the response will include this header to provide round-trip message integrity verification of the customer-provided encryption key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSECustomerKeyMD5"),
			},
			{
				Name:        "sse_kms_key_id",
				Description: "If present, specifies the ID of the Amazon Web Services Key Management Service (Amazon Web Services KMS) symmetric customer managed key that was used for the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSEKMSKeyId"),
			},
			{
				Name:        "server_side_encryption",
				Description: "The server-side encryption algorithm used when storing this object in Amazon S3 (for example, AES256, aws:kms).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "storage_class",
				Description: "Provides storage class information of the object. Amazon S3 returns this header for all objects except for S3 Standard storage class objects.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tag_count",
				Description: "The number of tags, if any, on the object.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "version_id",
				Description: "Version of the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "website_redirection_loaction",
				Description: "If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. Amazon S3 stores the value of this header in the object metadata.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "result_metadata",
				Description: "Metadata pertaining to the operation's result.",
				Type:        proto.ColumnType_JSON,
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
		}),
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
		accountID = aws.String(d.EqualsQualString("account_id"))
		bucket = aws.String(d.EqualsQualString("bucket"))
		key = aws.String(d.EqualsQualString("key"))
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
