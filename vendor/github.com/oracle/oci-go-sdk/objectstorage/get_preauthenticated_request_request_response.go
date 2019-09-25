// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetPreauthenticatedRequestRequest wrapper for the GetPreauthenticatedRequest operation
type GetPreauthenticatedRequestRequest struct {

	// The Object Storage namespace used for the request.
	NamespaceName *string `mandatory:"true" contributesTo:"path" name:"namespaceName"`

	// The name of the bucket. Avoid entering confidential information.
	// Example: `my-new-bucket1`
	BucketName *string `mandatory:"true" contributesTo:"path" name:"bucketName"`

	// The unique identifier for the pre-authenticated request. This can be used to manage operations against
	// the pre-authenticated request, such as GET or DELETE.
	ParId *string `mandatory:"true" contributesTo:"path" name:"parId"`

	// The client request ID for tracing.
	OpcClientRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-client-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetPreauthenticatedRequestRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetPreauthenticatedRequestRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetPreauthenticatedRequestRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetPreauthenticatedRequestResponse wrapper for the GetPreauthenticatedRequest operation
type GetPreauthenticatedRequestResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The PreauthenticatedRequestSummary instance
	PreauthenticatedRequestSummary `presentIn:"body"`

	// Echoes back the value passed in the opc-client-request-id header, for use by clients when debugging.
	OpcClientRequestId *string `presentIn:"header" name:"opc-client-request-id"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular
	// request, provide this request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response GetPreauthenticatedRequestResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetPreauthenticatedRequestResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}
