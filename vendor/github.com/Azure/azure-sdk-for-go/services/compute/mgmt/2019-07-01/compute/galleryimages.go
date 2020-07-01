package compute

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
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// GalleryImagesClient is the compute Client
type GalleryImagesClient struct {
	BaseClient
}

// NewGalleryImagesClient creates an instance of the GalleryImagesClient client.
func NewGalleryImagesClient(subscriptionID string) GalleryImagesClient {
	return NewGalleryImagesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewGalleryImagesClientWithBaseURI creates an instance of the GalleryImagesClient client.
func NewGalleryImagesClientWithBaseURI(baseURI string, subscriptionID string) GalleryImagesClient {
	return GalleryImagesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate create or update a gallery Image Definition.
// Parameters:
// resourceGroupName - the name of the resource group.
// galleryName - the name of the Shared Image Gallery in which the Image Definition is to be created.
// galleryImageName - the name of the gallery Image Definition to be created or updated. The allowed characters
// are alphabets and numbers with dots, dashes, and periods allowed in the middle. The maximum length is 80
// characters.
// galleryImage - parameters supplied to the create or update gallery image operation.
func (client GalleryImagesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, galleryImage GalleryImage) (result GalleryImagesCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryImagesClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: galleryImage,
			Constraints: []validation.Constraint{{Target: "galleryImage.GalleryImageProperties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "galleryImage.GalleryImageProperties.Identifier", Name: validation.Null, Rule: true,
					Chain: []validation.Constraint{{Target: "galleryImage.GalleryImageProperties.Identifier.Publisher", Name: validation.Null, Rule: true, Chain: nil},
						{Target: "galleryImage.GalleryImageProperties.Identifier.Offer", Name: validation.Null, Rule: true, Chain: nil},
						{Target: "galleryImage.GalleryImageProperties.Identifier.Sku", Name: validation.Null, Rule: true, Chain: nil},
					}},
				}}}}}); err != nil {
		return result, validation.NewError("compute.GalleryImagesClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, galleryName, galleryImageName, galleryImage)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client GalleryImagesClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string, galleryImage GalleryImage) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"galleryImageName":  autorest.Encode("path", galleryImageName),
		"galleryName":       autorest.Encode("path", galleryName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images/{galleryImageName}", pathParameters),
		autorest.WithJSON(galleryImage),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client GalleryImagesClient) CreateOrUpdateSender(req *http.Request) (future GalleryImagesCreateOrUpdateFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client GalleryImagesClient) CreateOrUpdateResponder(resp *http.Response) (result GalleryImage, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete delete a gallery image.
// Parameters:
// resourceGroupName - the name of the resource group.
// galleryName - the name of the Shared Image Gallery in which the Image Definition is to be deleted.
// galleryImageName - the name of the gallery Image Definition to be deleted.
func (client GalleryImagesClient) Delete(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string) (result GalleryImagesDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryImagesClient.Delete")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, galleryName, galleryImageName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "Delete", result.Response(), "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client GalleryImagesClient) DeletePreparer(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"galleryImageName":  autorest.Encode("path", galleryImageName),
		"galleryName":       autorest.Encode("path", galleryName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images/{galleryImageName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client GalleryImagesClient) DeleteSender(req *http.Request) (future GalleryImagesDeleteFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client GalleryImagesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get retrieves information about a gallery Image Definition.
// Parameters:
// resourceGroupName - the name of the resource group.
// galleryName - the name of the Shared Image Gallery from which the Image Definitions are to be retrieved.
// galleryImageName - the name of the gallery Image Definition to be retrieved.
func (client GalleryImagesClient) Get(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string) (result GalleryImage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryImagesClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, galleryName, galleryImageName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client GalleryImagesClient) GetPreparer(ctx context.Context, resourceGroupName string, galleryName string, galleryImageName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"galleryImageName":  autorest.Encode("path", galleryImageName),
		"galleryName":       autorest.Encode("path", galleryName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images/{galleryImageName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client GalleryImagesClient) GetSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client GalleryImagesClient) GetResponder(resp *http.Response) (result GalleryImage, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByGallery list gallery Image Definitions in a gallery.
// Parameters:
// resourceGroupName - the name of the resource group.
// galleryName - the name of the Shared Image Gallery from which Image Definitions are to be listed.
func (client GalleryImagesClient) ListByGallery(ctx context.Context, resourceGroupName string, galleryName string) (result GalleryImageListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryImagesClient.ListByGallery")
		defer func() {
			sc := -1
			if result.gil.Response.Response != nil {
				sc = result.gil.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listByGalleryNextResults
	req, err := client.ListByGalleryPreparer(ctx, resourceGroupName, galleryName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "ListByGallery", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByGallerySender(req)
	if err != nil {
		result.gil.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "ListByGallery", resp, "Failure sending request")
		return
	}

	result.gil, err = client.ListByGalleryResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "ListByGallery", resp, "Failure responding to request")
	}

	return
}

// ListByGalleryPreparer prepares the ListByGallery request.
func (client GalleryImagesClient) ListByGalleryPreparer(ctx context.Context, resourceGroupName string, galleryName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"galleryName":       autorest.Encode("path", galleryName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/images", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByGallerySender sends the ListByGallery request. The method will close the
// http.Response Body if it receives an error.
func (client GalleryImagesClient) ListByGallerySender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// ListByGalleryResponder handles the response to the ListByGallery request. The method always
// closes the http.Response Body.
func (client GalleryImagesClient) ListByGalleryResponder(resp *http.Response) (result GalleryImageList, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByGalleryNextResults retrieves the next set of results, if any.
func (client GalleryImagesClient) listByGalleryNextResults(ctx context.Context, lastResults GalleryImageList) (result GalleryImageList, err error) {
	req, err := lastResults.galleryImageListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "listByGalleryNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByGallerySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "listByGalleryNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByGalleryResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryImagesClient", "listByGalleryNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByGalleryComplete enumerates all values, automatically crossing page boundaries as required.
func (client GalleryImagesClient) ListByGalleryComplete(ctx context.Context, resourceGroupName string, galleryName string) (result GalleryImageListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryImagesClient.ListByGallery")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByGallery(ctx, resourceGroupName, galleryName)
	return
}
