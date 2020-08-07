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
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// VirtualMachineScaleSetRollingUpgradesClient is the compute Client
type VirtualMachineScaleSetRollingUpgradesClient struct {
	BaseClient
}

// NewVirtualMachineScaleSetRollingUpgradesClient creates an instance of the
// VirtualMachineScaleSetRollingUpgradesClient client.
func NewVirtualMachineScaleSetRollingUpgradesClient(subscriptionID string) VirtualMachineScaleSetRollingUpgradesClient {
	return NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI creates an instance of the
// VirtualMachineScaleSetRollingUpgradesClient client.
func NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI(baseURI string, subscriptionID string) VirtualMachineScaleSetRollingUpgradesClient {
	return VirtualMachineScaleSetRollingUpgradesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Cancel cancels the current virtual machine scale set rolling upgrade.
// Parameters:
// resourceGroupName - the name of the resource group.
// VMScaleSetName - the name of the VM scale set.
func (client VirtualMachineScaleSetRollingUpgradesClient) Cancel(ctx context.Context, resourceGroupName string, VMScaleSetName string) (result VirtualMachineScaleSetRollingUpgradesCancelFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualMachineScaleSetRollingUpgradesClient.Cancel")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CancelPreparer(ctx, resourceGroupName, VMScaleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "Cancel", nil, "Failure preparing request")
		return
	}

	result, err = client.CancelSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "Cancel", result.Response(), "Failure sending request")
		return
	}

	return
}

// CancelPreparer prepares the Cancel request.
func (client VirtualMachineScaleSetRollingUpgradesClient) CancelPreparer(ctx context.Context, resourceGroupName string, VMScaleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vmScaleSetName":    autorest.Encode("path", VMScaleSetName),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{vmScaleSetName}/rollingUpgrades/cancel", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CancelSender sends the Cancel request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineScaleSetRollingUpgradesClient) CancelSender(req *http.Request) (future VirtualMachineScaleSetRollingUpgradesCancelFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// CancelResponder handles the response to the Cancel request. The method always
// closes the http.Response Body.
func (client VirtualMachineScaleSetRollingUpgradesClient) CancelResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// GetLatest gets the status of the latest virtual machine scale set rolling upgrade.
// Parameters:
// resourceGroupName - the name of the resource group.
// VMScaleSetName - the name of the VM scale set.
func (client VirtualMachineScaleSetRollingUpgradesClient) GetLatest(ctx context.Context, resourceGroupName string, VMScaleSetName string) (result RollingUpgradeStatusInfo, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualMachineScaleSetRollingUpgradesClient.GetLatest")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetLatestPreparer(ctx, resourceGroupName, VMScaleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "GetLatest", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetLatestSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "GetLatest", resp, "Failure sending request")
		return
	}

	result, err = client.GetLatestResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "GetLatest", resp, "Failure responding to request")
	}

	return
}

// GetLatestPreparer prepares the GetLatest request.
func (client VirtualMachineScaleSetRollingUpgradesClient) GetLatestPreparer(ctx context.Context, resourceGroupName string, VMScaleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vmScaleSetName":    autorest.Encode("path", VMScaleSetName),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{vmScaleSetName}/rollingUpgrades/latest", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetLatestSender sends the GetLatest request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineScaleSetRollingUpgradesClient) GetLatestSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// GetLatestResponder handles the response to the GetLatest request. The method always
// closes the http.Response Body.
func (client VirtualMachineScaleSetRollingUpgradesClient) GetLatestResponder(resp *http.Response) (result RollingUpgradeStatusInfo, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// StartExtensionUpgrade starts a rolling upgrade to move all extensions for all virtual machine scale set instances to
// the latest available extension version. Instances which are already running the latest extension versions are not
// affected.
// Parameters:
// resourceGroupName - the name of the resource group.
// VMScaleSetName - the name of the VM scale set.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartExtensionUpgrade(ctx context.Context, resourceGroupName string, VMScaleSetName string) (result VirtualMachineScaleSetRollingUpgradesStartExtensionUpgradeFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualMachineScaleSetRollingUpgradesClient.StartExtensionUpgrade")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.StartExtensionUpgradePreparer(ctx, resourceGroupName, VMScaleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "StartExtensionUpgrade", nil, "Failure preparing request")
		return
	}

	result, err = client.StartExtensionUpgradeSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "StartExtensionUpgrade", result.Response(), "Failure sending request")
		return
	}

	return
}

// StartExtensionUpgradePreparer prepares the StartExtensionUpgrade request.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartExtensionUpgradePreparer(ctx context.Context, resourceGroupName string, VMScaleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vmScaleSetName":    autorest.Encode("path", VMScaleSetName),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{vmScaleSetName}/extensionRollingUpgrade", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// StartExtensionUpgradeSender sends the StartExtensionUpgrade request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartExtensionUpgradeSender(req *http.Request) (future VirtualMachineScaleSetRollingUpgradesStartExtensionUpgradeFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// StartExtensionUpgradeResponder handles the response to the StartExtensionUpgrade request. The method always
// closes the http.Response Body.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartExtensionUpgradeResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// StartOSUpgrade starts a rolling upgrade to move all virtual machine scale set instances to the latest available
// Platform Image OS version. Instances which are already running the latest available OS version are not affected.
// Parameters:
// resourceGroupName - the name of the resource group.
// VMScaleSetName - the name of the VM scale set.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartOSUpgrade(ctx context.Context, resourceGroupName string, VMScaleSetName string) (result VirtualMachineScaleSetRollingUpgradesStartOSUpgradeFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/VirtualMachineScaleSetRollingUpgradesClient.StartOSUpgrade")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.StartOSUpgradePreparer(ctx, resourceGroupName, VMScaleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "StartOSUpgrade", nil, "Failure preparing request")
		return
	}

	result, err = client.StartOSUpgradeSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachineScaleSetRollingUpgradesClient", "StartOSUpgrade", result.Response(), "Failure sending request")
		return
	}

	return
}

// StartOSUpgradePreparer prepares the StartOSUpgrade request.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartOSUpgradePreparer(ctx context.Context, resourceGroupName string, VMScaleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vmScaleSetName":    autorest.Encode("path", VMScaleSetName),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{vmScaleSetName}/osRollingUpgrade", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// StartOSUpgradeSender sends the StartOSUpgrade request. The method will close the
// http.Response Body if it receives an error.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartOSUpgradeSender(req *http.Request) (future VirtualMachineScaleSetRollingUpgradesStartOSUpgradeFuture, err error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, sd...)
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// StartOSUpgradeResponder handles the response to the StartOSUpgrade request. The method always
// closes the http.Response Body.
func (client VirtualMachineScaleSetRollingUpgradesClient) StartOSUpgradeResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}
