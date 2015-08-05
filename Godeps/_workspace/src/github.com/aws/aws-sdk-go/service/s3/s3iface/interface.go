// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

// Package s3iface provides an interface for the Amazon Simple Storage Service.
package s3iface

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3API is the interface type for s3.S3.
type S3API interface {
	AbortMultipartUploadRequest(*s3.AbortMultipartUploadInput) (*aws.Request, *s3.AbortMultipartUploadOutput)

	AbortMultipartUpload(*s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error)

	CompleteMultipartUploadRequest(*s3.CompleteMultipartUploadInput) (*aws.Request, *s3.CompleteMultipartUploadOutput)

	CompleteMultipartUpload(*s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error)

	CopyObjectRequest(*s3.CopyObjectInput) (*aws.Request, *s3.CopyObjectOutput)

	CopyObject(*s3.CopyObjectInput) (*s3.CopyObjectOutput, error)

	CreateBucketRequest(*s3.CreateBucketInput) (*aws.Request, *s3.CreateBucketOutput)

	CreateBucket(*s3.CreateBucketInput) (*s3.CreateBucketOutput, error)

	CreateMultipartUploadRequest(*s3.CreateMultipartUploadInput) (*aws.Request, *s3.CreateMultipartUploadOutput)

	CreateMultipartUpload(*s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error)

	DeleteBucketRequest(*s3.DeleteBucketInput) (*aws.Request, *s3.DeleteBucketOutput)

	DeleteBucket(*s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error)

	DeleteBucketCORSRequest(*s3.DeleteBucketCORSInput) (*aws.Request, *s3.DeleteBucketCORSOutput)

	DeleteBucketCORS(*s3.DeleteBucketCORSInput) (*s3.DeleteBucketCORSOutput, error)

	DeleteBucketLifecycleRequest(*s3.DeleteBucketLifecycleInput) (*aws.Request, *s3.DeleteBucketLifecycleOutput)

	DeleteBucketLifecycle(*s3.DeleteBucketLifecycleInput) (*s3.DeleteBucketLifecycleOutput, error)

	DeleteBucketPolicyRequest(*s3.DeleteBucketPolicyInput) (*aws.Request, *s3.DeleteBucketPolicyOutput)

	DeleteBucketPolicy(*s3.DeleteBucketPolicyInput) (*s3.DeleteBucketPolicyOutput, error)

	DeleteBucketReplicationRequest(*s3.DeleteBucketReplicationInput) (*aws.Request, *s3.DeleteBucketReplicationOutput)

	DeleteBucketReplication(*s3.DeleteBucketReplicationInput) (*s3.DeleteBucketReplicationOutput, error)

	DeleteBucketTaggingRequest(*s3.DeleteBucketTaggingInput) (*aws.Request, *s3.DeleteBucketTaggingOutput)

	DeleteBucketTagging(*s3.DeleteBucketTaggingInput) (*s3.DeleteBucketTaggingOutput, error)

	DeleteBucketWebsiteRequest(*s3.DeleteBucketWebsiteInput) (*aws.Request, *s3.DeleteBucketWebsiteOutput)

	DeleteBucketWebsite(*s3.DeleteBucketWebsiteInput) (*s3.DeleteBucketWebsiteOutput, error)

	DeleteObjectRequest(*s3.DeleteObjectInput) (*aws.Request, *s3.DeleteObjectOutput)

	DeleteObject(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)

	DeleteObjectsRequest(*s3.DeleteObjectsInput) (*aws.Request, *s3.DeleteObjectsOutput)

	DeleteObjects(*s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error)

	GetBucketACLRequest(*s3.GetBucketACLInput) (*aws.Request, *s3.GetBucketACLOutput)

	GetBucketACL(*s3.GetBucketACLInput) (*s3.GetBucketACLOutput, error)

	GetBucketCORSRequest(*s3.GetBucketCORSInput) (*aws.Request, *s3.GetBucketCORSOutput)

	GetBucketCORS(*s3.GetBucketCORSInput) (*s3.GetBucketCORSOutput, error)

	GetBucketLifecycleRequest(*s3.GetBucketLifecycleInput) (*aws.Request, *s3.GetBucketLifecycleOutput)

	GetBucketLifecycle(*s3.GetBucketLifecycleInput) (*s3.GetBucketLifecycleOutput, error)

	GetBucketLocationRequest(*s3.GetBucketLocationInput) (*aws.Request, *s3.GetBucketLocationOutput)

	GetBucketLocation(*s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error)

	GetBucketLoggingRequest(*s3.GetBucketLoggingInput) (*aws.Request, *s3.GetBucketLoggingOutput)

	GetBucketLogging(*s3.GetBucketLoggingInput) (*s3.GetBucketLoggingOutput, error)

	GetBucketNotificationRequest(*s3.GetBucketNotificationConfigurationRequest) (*aws.Request, *s3.NotificationConfigurationDeprecated)

	GetBucketNotification(*s3.GetBucketNotificationConfigurationRequest) (*s3.NotificationConfigurationDeprecated, error)

	GetBucketNotificationConfigurationRequest(*s3.GetBucketNotificationConfigurationRequest) (*aws.Request, *s3.NotificationConfiguration)

	GetBucketNotificationConfiguration(*s3.GetBucketNotificationConfigurationRequest) (*s3.NotificationConfiguration, error)

	GetBucketPolicyRequest(*s3.GetBucketPolicyInput) (*aws.Request, *s3.GetBucketPolicyOutput)

	GetBucketPolicy(*s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error)

	GetBucketReplicationRequest(*s3.GetBucketReplicationInput) (*aws.Request, *s3.GetBucketReplicationOutput)

	GetBucketReplication(*s3.GetBucketReplicationInput) (*s3.GetBucketReplicationOutput, error)

	GetBucketRequestPaymentRequest(*s3.GetBucketRequestPaymentInput) (*aws.Request, *s3.GetBucketRequestPaymentOutput)

	GetBucketRequestPayment(*s3.GetBucketRequestPaymentInput) (*s3.GetBucketRequestPaymentOutput, error)

	GetBucketTaggingRequest(*s3.GetBucketTaggingInput) (*aws.Request, *s3.GetBucketTaggingOutput)

	GetBucketTagging(*s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error)

	GetBucketVersioningRequest(*s3.GetBucketVersioningInput) (*aws.Request, *s3.GetBucketVersioningOutput)

	GetBucketVersioning(*s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error)

	GetBucketWebsiteRequest(*s3.GetBucketWebsiteInput) (*aws.Request, *s3.GetBucketWebsiteOutput)

	GetBucketWebsite(*s3.GetBucketWebsiteInput) (*s3.GetBucketWebsiteOutput, error)

	GetObjectRequest(*s3.GetObjectInput) (*aws.Request, *s3.GetObjectOutput)

	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)

	GetObjectACLRequest(*s3.GetObjectACLInput) (*aws.Request, *s3.GetObjectACLOutput)

	GetObjectACL(*s3.GetObjectACLInput) (*s3.GetObjectACLOutput, error)

	GetObjectTorrentRequest(*s3.GetObjectTorrentInput) (*aws.Request, *s3.GetObjectTorrentOutput)

	GetObjectTorrent(*s3.GetObjectTorrentInput) (*s3.GetObjectTorrentOutput, error)

	HeadBucketRequest(*s3.HeadBucketInput) (*aws.Request, *s3.HeadBucketOutput)

	HeadBucket(*s3.HeadBucketInput) (*s3.HeadBucketOutput, error)

	HeadObjectRequest(*s3.HeadObjectInput) (*aws.Request, *s3.HeadObjectOutput)

	HeadObject(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error)

	ListBucketsRequest(*s3.ListBucketsInput) (*aws.Request, *s3.ListBucketsOutput)

	ListBuckets(*s3.ListBucketsInput) (*s3.ListBucketsOutput, error)

	ListMultipartUploadsRequest(*s3.ListMultipartUploadsInput) (*aws.Request, *s3.ListMultipartUploadsOutput)

	ListMultipartUploads(*s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error)

	ListMultipartUploadsPages(*s3.ListMultipartUploadsInput, func(*s3.ListMultipartUploadsOutput, bool) bool) error

	ListObjectVersionsRequest(*s3.ListObjectVersionsInput) (*aws.Request, *s3.ListObjectVersionsOutput)

	ListObjectVersions(*s3.ListObjectVersionsInput) (*s3.ListObjectVersionsOutput, error)

	ListObjectVersionsPages(*s3.ListObjectVersionsInput, func(*s3.ListObjectVersionsOutput, bool) bool) error

	ListObjectsRequest(*s3.ListObjectsInput) (*aws.Request, *s3.ListObjectsOutput)

	ListObjects(*s3.ListObjectsInput) (*s3.ListObjectsOutput, error)

	ListObjectsPages(*s3.ListObjectsInput, func(*s3.ListObjectsOutput, bool) bool) error

	ListPartsRequest(*s3.ListPartsInput) (*aws.Request, *s3.ListPartsOutput)

	ListParts(*s3.ListPartsInput) (*s3.ListPartsOutput, error)

	ListPartsPages(*s3.ListPartsInput, func(*s3.ListPartsOutput, bool) bool) error

	PutBucketACLRequest(*s3.PutBucketACLInput) (*aws.Request, *s3.PutBucketACLOutput)

	PutBucketACL(*s3.PutBucketACLInput) (*s3.PutBucketACLOutput, error)

	PutBucketCORSRequest(*s3.PutBucketCORSInput) (*aws.Request, *s3.PutBucketCORSOutput)

	PutBucketCORS(*s3.PutBucketCORSInput) (*s3.PutBucketCORSOutput, error)

	PutBucketLifecycleRequest(*s3.PutBucketLifecycleInput) (*aws.Request, *s3.PutBucketLifecycleOutput)

	PutBucketLifecycle(*s3.PutBucketLifecycleInput) (*s3.PutBucketLifecycleOutput, error)

	PutBucketLoggingRequest(*s3.PutBucketLoggingInput) (*aws.Request, *s3.PutBucketLoggingOutput)

	PutBucketLogging(*s3.PutBucketLoggingInput) (*s3.PutBucketLoggingOutput, error)

	PutBucketNotificationRequest(*s3.PutBucketNotificationInput) (*aws.Request, *s3.PutBucketNotificationOutput)

	PutBucketNotification(*s3.PutBucketNotificationInput) (*s3.PutBucketNotificationOutput, error)

	PutBucketNotificationConfigurationRequest(*s3.PutBucketNotificationConfigurationInput) (*aws.Request, *s3.PutBucketNotificationConfigurationOutput)

	PutBucketNotificationConfiguration(*s3.PutBucketNotificationConfigurationInput) (*s3.PutBucketNotificationConfigurationOutput, error)

	PutBucketPolicyRequest(*s3.PutBucketPolicyInput) (*aws.Request, *s3.PutBucketPolicyOutput)

	PutBucketPolicy(*s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error)

	PutBucketReplicationRequest(*s3.PutBucketReplicationInput) (*aws.Request, *s3.PutBucketReplicationOutput)

	PutBucketReplication(*s3.PutBucketReplicationInput) (*s3.PutBucketReplicationOutput, error)

	PutBucketRequestPaymentRequest(*s3.PutBucketRequestPaymentInput) (*aws.Request, *s3.PutBucketRequestPaymentOutput)

	PutBucketRequestPayment(*s3.PutBucketRequestPaymentInput) (*s3.PutBucketRequestPaymentOutput, error)

	PutBucketTaggingRequest(*s3.PutBucketTaggingInput) (*aws.Request, *s3.PutBucketTaggingOutput)

	PutBucketTagging(*s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error)

	PutBucketVersioningRequest(*s3.PutBucketVersioningInput) (*aws.Request, *s3.PutBucketVersioningOutput)

	PutBucketVersioning(*s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error)

	PutBucketWebsiteRequest(*s3.PutBucketWebsiteInput) (*aws.Request, *s3.PutBucketWebsiteOutput)

	PutBucketWebsite(*s3.PutBucketWebsiteInput) (*s3.PutBucketWebsiteOutput, error)

	PutObjectRequest(*s3.PutObjectInput) (*aws.Request, *s3.PutObjectOutput)

	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)

	PutObjectACLRequest(*s3.PutObjectACLInput) (*aws.Request, *s3.PutObjectACLOutput)

	PutObjectACL(*s3.PutObjectACLInput) (*s3.PutObjectACLOutput, error)

	RestoreObjectRequest(*s3.RestoreObjectInput) (*aws.Request, *s3.RestoreObjectOutput)

	RestoreObject(*s3.RestoreObjectInput) (*s3.RestoreObjectOutput, error)

	UploadPartRequest(*s3.UploadPartInput) (*aws.Request, *s3.UploadPartOutput)

	UploadPart(*s3.UploadPartInput) (*s3.UploadPartOutput, error)

	UploadPartCopyRequest(*s3.UploadPartCopyInput) (*aws.Request, *s3.UploadPartCopyOutput)

	UploadPartCopy(*s3.UploadPartCopyInput) (*s3.UploadPartCopyOutput, error)
}
