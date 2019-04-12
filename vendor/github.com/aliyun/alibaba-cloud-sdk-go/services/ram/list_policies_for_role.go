package ram

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ListPoliciesForRole invokes the ram.ListPoliciesForRole API synchronously
// api document: https://help.aliyun.com/api/ram/listpoliciesforrole.html
func (client *Client) ListPoliciesForRole(request *ListPoliciesForRoleRequest) (response *ListPoliciesForRoleResponse, err error) {
	response = CreateListPoliciesForRoleResponse()
	err = client.DoAction(request, response)
	return
}

// ListPoliciesForRoleWithChan invokes the ram.ListPoliciesForRole API asynchronously
// api document: https://help.aliyun.com/api/ram/listpoliciesforrole.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListPoliciesForRoleWithChan(request *ListPoliciesForRoleRequest) (<-chan *ListPoliciesForRoleResponse, <-chan error) {
	responseChan := make(chan *ListPoliciesForRoleResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListPoliciesForRole(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// ListPoliciesForRoleWithCallback invokes the ram.ListPoliciesForRole API asynchronously
// api document: https://help.aliyun.com/api/ram/listpoliciesforrole.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListPoliciesForRoleWithCallback(request *ListPoliciesForRoleRequest, callback func(response *ListPoliciesForRoleResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListPoliciesForRoleResponse
		var err error
		defer close(result)
		response, err = client.ListPoliciesForRole(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// ListPoliciesForRoleRequest is the request struct for api ListPoliciesForRole
type ListPoliciesForRoleRequest struct {
	*requests.RpcRequest
	RoleName string `position:"Query" name:"RoleName"`
}

// ListPoliciesForRoleResponse is the response struct for api ListPoliciesForRole
type ListPoliciesForRoleResponse struct {
	*responses.BaseResponse
	RequestId string                        `json:"RequestId" xml:"RequestId"`
	Policies  PoliciesInListPoliciesForRole `json:"Policies" xml:"Policies"`
}

// CreateListPoliciesForRoleRequest creates a request to invoke ListPoliciesForRole API
func CreateListPoliciesForRoleRequest() (request *ListPoliciesForRoleRequest) {
	request = &ListPoliciesForRoleRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ram", "2015-05-01", "ListPoliciesForRole", "ram", "openAPI")
	return
}

// CreateListPoliciesForRoleResponse creates a response to parse from ListPoliciesForRole response
func CreateListPoliciesForRoleResponse() (response *ListPoliciesForRoleResponse) {
	response = &ListPoliciesForRoleResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
