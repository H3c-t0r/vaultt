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

// AttachPolicyToRole invokes the ram.AttachPolicyToRole API synchronously
// api document: https://help.aliyun.com/api/ram/attachpolicytorole.html
func (client *Client) AttachPolicyToRole(request *AttachPolicyToRoleRequest) (response *AttachPolicyToRoleResponse, err error) {
	response = CreateAttachPolicyToRoleResponse()
	err = client.DoAction(request, response)
	return
}

// AttachPolicyToRoleWithChan invokes the ram.AttachPolicyToRole API asynchronously
// api document: https://help.aliyun.com/api/ram/attachpolicytorole.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AttachPolicyToRoleWithChan(request *AttachPolicyToRoleRequest) (<-chan *AttachPolicyToRoleResponse, <-chan error) {
	responseChan := make(chan *AttachPolicyToRoleResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AttachPolicyToRole(request)
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

// AttachPolicyToRoleWithCallback invokes the ram.AttachPolicyToRole API asynchronously
// api document: https://help.aliyun.com/api/ram/attachpolicytorole.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AttachPolicyToRoleWithCallback(request *AttachPolicyToRoleRequest, callback func(response *AttachPolicyToRoleResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AttachPolicyToRoleResponse
		var err error
		defer close(result)
		response, err = client.AttachPolicyToRole(request)
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

// AttachPolicyToRoleRequest is the request struct for api AttachPolicyToRole
type AttachPolicyToRoleRequest struct {
	*requests.RpcRequest
	PolicyType string `position:"Query" name:"PolicyType"`
	RoleName   string `position:"Query" name:"RoleName"`
	PolicyName string `position:"Query" name:"PolicyName"`
}

// AttachPolicyToRoleResponse is the response struct for api AttachPolicyToRole
type AttachPolicyToRoleResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateAttachPolicyToRoleRequest creates a request to invoke AttachPolicyToRole API
func CreateAttachPolicyToRoleRequest() (request *AttachPolicyToRoleRequest) {
	request = &AttachPolicyToRoleRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ram", "2015-05-01", "AttachPolicyToRole", "", "")
	return
}

// CreateAttachPolicyToRoleResponse creates a response to parse from AttachPolicyToRole response
func CreateAttachPolicyToRoleResponse() (response *AttachPolicyToRoleResponse) {
	response = &AttachPolicyToRoleResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
