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

// SetAccountAlias invokes the ram.SetAccountAlias API synchronously
// api document: https://help.aliyun.com/api/ram/setaccountalias.html
func (client *Client) SetAccountAlias(request *SetAccountAliasRequest) (response *SetAccountAliasResponse, err error) {
	response = CreateSetAccountAliasResponse()
	err = client.DoAction(request, response)
	return
}

// SetAccountAliasWithChan invokes the ram.SetAccountAlias API asynchronously
// api document: https://help.aliyun.com/api/ram/setaccountalias.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetAccountAliasWithChan(request *SetAccountAliasRequest) (<-chan *SetAccountAliasResponse, <-chan error) {
	responseChan := make(chan *SetAccountAliasResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetAccountAlias(request)
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

// SetAccountAliasWithCallback invokes the ram.SetAccountAlias API asynchronously
// api document: https://help.aliyun.com/api/ram/setaccountalias.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetAccountAliasWithCallback(request *SetAccountAliasRequest, callback func(response *SetAccountAliasResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetAccountAliasResponse
		var err error
		defer close(result)
		response, err = client.SetAccountAlias(request)
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

// SetAccountAliasRequest is the request struct for api SetAccountAlias
type SetAccountAliasRequest struct {
	*requests.RpcRequest
	AccountAlias string `position:"Query" name:"AccountAlias"`
}

// SetAccountAliasResponse is the response struct for api SetAccountAlias
type SetAccountAliasResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetAccountAliasRequest creates a request to invoke SetAccountAlias API
func CreateSetAccountAliasRequest() (request *SetAccountAliasRequest) {
	request = &SetAccountAliasRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ram", "2015-05-01", "SetAccountAlias", "ram", "openAPI")
	return
}

// CreateSetAccountAliasResponse creates a response to parse from SetAccountAlias response
func CreateSetAccountAliasResponse() (response *SetAccountAliasResponse) {
	response = &SetAccountAliasResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
