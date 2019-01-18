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

// ListUsersForGroup invokes the ram.ListUsersForGroup API synchronously
// api document: https://help.aliyun.com/api/ram/listusersforgroup.html
func (client *Client) ListUsersForGroup(request *ListUsersForGroupRequest) (response *ListUsersForGroupResponse, err error) {
	response = CreateListUsersForGroupResponse()
	err = client.DoAction(request, response)
	return
}

// ListUsersForGroupWithChan invokes the ram.ListUsersForGroup API asynchronously
// api document: https://help.aliyun.com/api/ram/listusersforgroup.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListUsersForGroupWithChan(request *ListUsersForGroupRequest) (<-chan *ListUsersForGroupResponse, <-chan error) {
	responseChan := make(chan *ListUsersForGroupResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListUsersForGroup(request)
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

// ListUsersForGroupWithCallback invokes the ram.ListUsersForGroup API asynchronously
// api document: https://help.aliyun.com/api/ram/listusersforgroup.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListUsersForGroupWithCallback(request *ListUsersForGroupRequest, callback func(response *ListUsersForGroupResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListUsersForGroupResponse
		var err error
		defer close(result)
		response, err = client.ListUsersForGroup(request)
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

// ListUsersForGroupRequest is the request struct for api ListUsersForGroup
type ListUsersForGroupRequest struct {
	*requests.RpcRequest
	GroupName string `position:"Query" name:"GroupName"`
}

// ListUsersForGroupResponse is the response struct for api ListUsersForGroup
type ListUsersForGroupResponse struct {
	*responses.BaseResponse
	RequestId string                   `json:"RequestId" xml:"RequestId"`
	Users     UsersInListUsersForGroup `json:"Users" xml:"Users"`
}

// CreateListUsersForGroupRequest creates a request to invoke ListUsersForGroup API
func CreateListUsersForGroupRequest() (request *ListUsersForGroupRequest) {
	request = &ListUsersForGroupRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ram", "2015-05-01", "ListUsersForGroup", "", "")
	return
}

// CreateListUsersForGroupResponse creates a response to parse from ListUsersForGroup response
func CreateListUsersForGroupResponse() (response *ListUsersForGroupResponse) {
	response = &ListUsersForGroupResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
