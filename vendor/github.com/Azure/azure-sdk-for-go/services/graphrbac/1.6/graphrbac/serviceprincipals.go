package graphrbac

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// ServicePrincipalsClient is the the Graph RBAC Management Client
type ServicePrincipalsClient struct {
	BaseClient
}

// NewServicePrincipalsClient creates an instance of the ServicePrincipalsClient client.
func NewServicePrincipalsClient(tenantID string) ServicePrincipalsClient {
	return NewServicePrincipalsClientWithBaseURI(DefaultBaseURI, tenantID)
}

// NewServicePrincipalsClientWithBaseURI creates an instance of the ServicePrincipalsClient client.
func NewServicePrincipalsClientWithBaseURI(baseURI string, tenantID string) ServicePrincipalsClient {
	return ServicePrincipalsClient{NewWithBaseURI(baseURI, tenantID)}
}

// Create creates a service principal in the directory.
// Parameters:
// parameters - parameters to create a service principal.
func (client ServicePrincipalsClient) Create(ctx context.Context, parameters ServicePrincipalCreateParameters) (result ServicePrincipal, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.Create")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.AppID", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("graphrbac.ServicePrincipalsClient", "Create", err.Error())
	}

	req, err := client.CreatePreparer(ctx, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Create", resp, "Failure responding to request")
	}

	return
}

// CreatePreparer prepares the Create request.
func (client ServicePrincipalsClient) CreatePreparer(ctx context.Context, parameters ServicePrincipalCreateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/servicePrincipals", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) CreateSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) CreateResponder(resp *http.Response) (result ServicePrincipal, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes a service principal from the directory.
// Parameters:
// objectID - the object ID of the service principal to delete.
func (client ServicePrincipalsClient) Delete(ctx context.Context, objectID string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.Delete")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, objectID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client ServicePrincipalsClient) DeletePreparer(ctx context.Context, objectID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"objectId": autorest.Encode("path", objectID),
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/servicePrincipals/{objectId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets service principal information from the directory. Query by objectId or pass a filter to query by appId
// Parameters:
// objectID - the object ID of the service principal to get.
func (client ServicePrincipalsClient) Get(ctx context.Context, objectID string) (result ServicePrincipal, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, objectID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ServicePrincipalsClient) GetPreparer(ctx context.Context, objectID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"objectId": autorest.Encode("path", objectID),
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/servicePrincipals/{objectId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) GetSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) GetResponder(resp *http.Response) (result ServicePrincipal, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets a list of service principals from the current tenant.
// Parameters:
// filter - the filter to apply to the operation.
func (client ServicePrincipalsClient) List(ctx context.Context, filter string) (result ServicePrincipalListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.List")
		defer func() {
			sc := -1
			if result.splr.Response.Response != nil {
				sc = result.splr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = func(ctx context.Context, lastResult ServicePrincipalListResult) (ServicePrincipalListResult, error) {
		if lastResult.OdataNextLink == nil || len(to.String(lastResult.OdataNextLink)) < 1 {
			return ServicePrincipalListResult{}, nil
		}
		return client.ListNext(ctx, *lastResult.OdataNextLink)
	}
	req, err := client.ListPreparer(ctx, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.splr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "List", resp, "Failure sending request")
		return
	}

	result.splr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client ServicePrincipalsClient) ListPreparer(ctx context.Context, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/servicePrincipals", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) ListSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) ListResponder(resp *http.Response) (result ServicePrincipalListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client ServicePrincipalsClient) ListComplete(ctx context.Context, filter string) (result ServicePrincipalListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, filter)
	return
}

// ListKeyCredentials get the keyCredentials associated with the specified service principal.
// Parameters:
// objectID - the object ID of the service principal for which to get keyCredentials.
func (client ServicePrincipalsClient) ListKeyCredentials(ctx context.Context, objectID string) (result KeyCredentialListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.ListKeyCredentials")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListKeyCredentialsPreparer(ctx, objectID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListKeyCredentials", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListKeyCredentialsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListKeyCredentials", resp, "Failure sending request")
		return
	}

	result, err = client.ListKeyCredentialsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListKeyCredentials", resp, "Failure responding to request")
	}

	return
}

// ListKeyCredentialsPreparer prepares the ListKeyCredentials request.
func (client ServicePrincipalsClient) ListKeyCredentialsPreparer(ctx context.Context, objectID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"objectId": autorest.Encode("path", objectID),
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/servicePrincipals/{objectId}/keyCredentials", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListKeyCredentialsSender sends the ListKeyCredentials request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) ListKeyCredentialsSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// ListKeyCredentialsResponder handles the response to the ListKeyCredentials request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) ListKeyCredentialsResponder(resp *http.Response) (result KeyCredentialListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListNext gets a list of service principals from the current tenant.
// Parameters:
// nextLink - next link for the list operation.
func (client ServicePrincipalsClient) ListNext(ctx context.Context, nextLink string) (result ServicePrincipalListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.ListNext")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListNextPreparer(ctx, nextLink)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListNext", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListNextSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListNext", resp, "Failure sending request")
		return
	}

	result, err = client.ListNextResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListNext", resp, "Failure responding to request")
	}

	return
}

// ListNextPreparer prepares the ListNext request.
func (client ServicePrincipalsClient) ListNextPreparer(ctx context.Context, nextLink string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"nextLink": nextLink,
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/{nextLink}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListNextSender sends the ListNext request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) ListNextSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// ListNextResponder handles the response to the ListNext request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) ListNextResponder(resp *http.Response) (result ServicePrincipalListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListOwners the owners are a set of non-admin users who are allowed to modify this object.
// Parameters:
// objectID - the object ID of the service principal for which to get owners.
func (client ServicePrincipalsClient) ListOwners(ctx context.Context, objectID string) (result DirectoryObjectListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.ListOwners")
		defer func() {
			sc := -1
			if result.dolr.Response.Response != nil {
				sc = result.dolr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listOwnersNextResults
	req, err := client.ListOwnersPreparer(ctx, objectID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListOwners", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListOwnersSender(req)
	if err != nil {
		result.dolr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListOwners", resp, "Failure sending request")
		return
	}

	result.dolr, err = client.ListOwnersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListOwners", resp, "Failure responding to request")
	}

	return
}

// ListOwnersPreparer prepares the ListOwners request.
func (client ServicePrincipalsClient) ListOwnersPreparer(ctx context.Context, objectID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"objectId": autorest.Encode("path", objectID),
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/servicePrincipals/{objectId}/owners", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListOwnersSender sends the ListOwners request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) ListOwnersSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// ListOwnersResponder handles the response to the ListOwners request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) ListOwnersResponder(resp *http.Response) (result DirectoryObjectListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listOwnersNextResults retrieves the next set of results, if any.
func (client ServicePrincipalsClient) listOwnersNextResults(ctx context.Context, lastResults DirectoryObjectListResult) (result DirectoryObjectListResult, err error) {
	req, err := lastResults.directoryObjectListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "listOwnersNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListOwnersSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "listOwnersNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListOwnersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "listOwnersNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListOwnersComplete enumerates all values, automatically crossing page boundaries as required.
func (client ServicePrincipalsClient) ListOwnersComplete(ctx context.Context, objectID string) (result DirectoryObjectListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.ListOwners")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListOwners(ctx, objectID)
	return
}

// ListPasswordCredentials gets the passwordCredentials associated with a service principal.
// Parameters:
// objectID - the object ID of the service principal.
func (client ServicePrincipalsClient) ListPasswordCredentials(ctx context.Context, objectID string) (result PasswordCredentialListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.ListPasswordCredentials")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListPasswordCredentialsPreparer(ctx, objectID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListPasswordCredentials", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListPasswordCredentialsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListPasswordCredentials", resp, "Failure sending request")
		return
	}

	result, err = client.ListPasswordCredentialsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "ListPasswordCredentials", resp, "Failure responding to request")
	}

	return
}

// ListPasswordCredentialsPreparer prepares the ListPasswordCredentials request.
func (client ServicePrincipalsClient) ListPasswordCredentialsPreparer(ctx context.Context, objectID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"objectId": autorest.Encode("path", objectID),
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/servicePrincipals/{objectId}/passwordCredentials", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListPasswordCredentialsSender sends the ListPasswordCredentials request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) ListPasswordCredentialsSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// ListPasswordCredentialsResponder handles the response to the ListPasswordCredentials request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) ListPasswordCredentialsResponder(resp *http.Response) (result PasswordCredentialListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Update updates a service principal in the directory.
// Parameters:
// objectID - the object ID of the service principal to delete.
// parameters - parameters to update a service principal.
func (client ServicePrincipalsClient) Update(ctx context.Context, objectID string, parameters ServicePrincipalUpdateParameters) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.Update")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdatePreparer(ctx, objectID, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Update", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Update", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "Update", resp, "Failure responding to request")
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client ServicePrincipalsClient) UpdatePreparer(ctx context.Context, objectID string, parameters ServicePrincipalUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"objectId": autorest.Encode("path", objectID),
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/servicePrincipals/{objectId}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) UpdateSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) UpdateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// UpdateKeyCredentials update the keyCredentials associated with a service principal.
// Parameters:
// objectID - the object ID for which to get service principal information.
// parameters - parameters to update the keyCredentials of an existing service principal.
func (client ServicePrincipalsClient) UpdateKeyCredentials(ctx context.Context, objectID string, parameters KeyCredentialsUpdateParameters) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.UpdateKeyCredentials")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdateKeyCredentialsPreparer(ctx, objectID, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "UpdateKeyCredentials", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateKeyCredentialsSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "UpdateKeyCredentials", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateKeyCredentialsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "UpdateKeyCredentials", resp, "Failure responding to request")
	}

	return
}

// UpdateKeyCredentialsPreparer prepares the UpdateKeyCredentials request.
func (client ServicePrincipalsClient) UpdateKeyCredentialsPreparer(ctx context.Context, objectID string, parameters KeyCredentialsUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"objectId": autorest.Encode("path", objectID),
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/servicePrincipals/{objectId}/keyCredentials", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateKeyCredentialsSender sends the UpdateKeyCredentials request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) UpdateKeyCredentialsSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// UpdateKeyCredentialsResponder handles the response to the UpdateKeyCredentials request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) UpdateKeyCredentialsResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// UpdatePasswordCredentials updates the passwordCredentials associated with a service principal.
// Parameters:
// objectID - the object ID of the service principal.
// parameters - parameters to update the passwordCredentials of an existing service principal.
func (client ServicePrincipalsClient) UpdatePasswordCredentials(ctx context.Context, objectID string, parameters PasswordCredentialsUpdateParameters) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServicePrincipalsClient.UpdatePasswordCredentials")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdatePasswordCredentialsPreparer(ctx, objectID, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "UpdatePasswordCredentials", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdatePasswordCredentialsSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "UpdatePasswordCredentials", resp, "Failure sending request")
		return
	}

	result, err = client.UpdatePasswordCredentialsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "graphrbac.ServicePrincipalsClient", "UpdatePasswordCredentials", resp, "Failure responding to request")
	}

	return
}

// UpdatePasswordCredentialsPreparer prepares the UpdatePasswordCredentials request.
func (client ServicePrincipalsClient) UpdatePasswordCredentialsPreparer(ctx context.Context, objectID string, parameters PasswordCredentialsUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"objectId": autorest.Encode("path", objectID),
		"tenantID": autorest.Encode("path", client.TenantID),
	}

	const APIVersion = "1.6"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{tenantID}/servicePrincipals/{objectId}/passwordCredentials", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdatePasswordCredentialsSender sends the UpdatePasswordCredentials request. The method will close the
// http.Response Body if it receives an error.
func (client ServicePrincipalsClient) UpdatePasswordCredentialsSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	return autorest.SendWithSender(client, req, sd...)
}

// UpdatePasswordCredentialsResponder handles the response to the UpdatePasswordCredentials request. The method always
// closes the http.Response Body.
func (client ServicePrincipalsClient) UpdatePasswordCredentialsResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}
