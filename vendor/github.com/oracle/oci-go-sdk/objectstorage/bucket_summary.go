// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object Storage and Archive Storage APIs for managing buckets, objects, and related resources.
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// BucketSummary To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type BucketSummary struct {

	// The Object Storage namespace in which the bucket lives.
	Namespace *string `mandatory:"true" json:"namespace"`

	// The name of the bucket. Avoid entering confidential information.
	// Example: my-new-bucket1
	Name *string `mandatory:"true" json:"name"`

	// The compartment ID in which the bucket is authorized.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the user who created the bucket.
	CreatedBy *string `mandatory:"true" json:"createdBy"`

	// The date and time the bucket was created, as described in RFC 2616 (https://tools.ietf.org/rfc/rfc2616), section 14.29.
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The entity tag (ETag) for the bucket.
	Etag *string `mandatory:"true" json:"etag"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m BucketSummary) String() string {
	return common.PointerString(m)
}
