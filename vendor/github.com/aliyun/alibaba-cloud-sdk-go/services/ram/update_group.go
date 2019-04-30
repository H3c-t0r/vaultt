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

// UpdateGroup invokes the ram.UpdateGroup API synchronously
// api document: https://help.aliyun.com/api/ram/updategroup.html
func (client *Client) UpdateGroup(request *UpdateGroupRequest) (response *UpdateGroupResponse, err error) {
	response = CreateUpdateGroupResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateGroupWithChan invokes the ram.UpdateGroup API asynchronously
// api document: https://help.aliyun.com/api/ram/updategroup.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateGroupWithChan(request *UpdateGroupRequest) (<-chan *UpdateGroupResponse, <-chan error) {
	responseChan := make(chan *UpdateGroupResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateGroup(request)
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

// UpdateGroupWithCallback invokes the ram.UpdateGroup API asynchronously
// api document: https://help.aliyun.com/api/ram/updategroup.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateGroupWithCallback(request *UpdateGroupRequest, callback func(response *UpdateGroupResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateGroupResponse
		var err error
		defer close(result)
		response, err = client.UpdateGroup(request)
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

// UpdateGroupRequest is the request struct for api UpdateGroup
type UpdateGroupRequest struct {
	*requests.RpcRequest
	NewGroupName string `position:"Query" name:"NewGroupName"`
	NewComments  string `position:"Query" name:"NewComments"`
	GroupName    string `position:"Query" name:"GroupName"`
}

// UpdateGroupResponse is the response struct for api UpdateGroup
type UpdateGroupResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Group     Group  `json:"Group" xml:"Group"`
}

// CreateUpdateGroupRequest creates a request to invoke UpdateGroup API
func CreateUpdateGroupRequest() (request *UpdateGroupRequest) {
	request = &UpdateGroupRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ram", "2015-05-01", "UpdateGroup", "ram", "openAPI")
	return
}

// CreateUpdateGroupResponse creates a response to parse from UpdateGroup response
func CreateUpdateGroupResponse() (response *UpdateGroupResponse) {
	response = &UpdateGroupResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
