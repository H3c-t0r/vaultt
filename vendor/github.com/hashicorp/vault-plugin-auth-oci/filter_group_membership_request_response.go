// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// Do not edit this file. This is based on standard OCI GO SDK format

// Stores the request body and meta-data required for filtering the group membership
type FilterGroupMembershipRequest struct {

	// Request object for FilterGroupMembershipRequest
	FilterGroupMembershipDetails `contributesTo:"body"`

	// A token that uniquely identifies a request so it can be retried in case of a timeout or
	// server error without risk of executing that same action again. Retry tokens expire after 24
	// hours, but can be invalidated before then due to conflicting operations (e.g., if a resource
	// has been deleted and purged from the system, then a retry of the original creation request
	// may be rejected).
	OpcRetryToken *string `mandatory:"false" contributesTo:"header" name:"opc-retry-token"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request FilterGroupMembershipRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request FilterGroupMembershipRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request FilterGroupMembershipRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// Stores the response of the FilterGroupMembership request, including meta-data.
type FilterGroupMembershipResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The FilterGroupMembershipResult instance
	FilterGroupMembershipResult `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `presentIn:"header" name:"etag"`
}

func (response FilterGroupMembershipResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response FilterGroupMembershipResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}
