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

// SetSecurityPreference invokes the ram.SetSecurityPreference API synchronously
// api document: https://help.aliyun.com/api/ram/setsecuritypreference.html
func (client *Client) SetSecurityPreference(request *SetSecurityPreferenceRequest) (response *SetSecurityPreferenceResponse, err error) {
	response = CreateSetSecurityPreferenceResponse()
	err = client.DoAction(request, response)
	return
}

// SetSecurityPreferenceWithChan invokes the ram.SetSecurityPreference API asynchronously
// api document: https://help.aliyun.com/api/ram/setsecuritypreference.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetSecurityPreferenceWithChan(request *SetSecurityPreferenceRequest) (<-chan *SetSecurityPreferenceResponse, <-chan error) {
	responseChan := make(chan *SetSecurityPreferenceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetSecurityPreference(request)
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

// SetSecurityPreferenceWithCallback invokes the ram.SetSecurityPreference API asynchronously
// api document: https://help.aliyun.com/api/ram/setsecuritypreference.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetSecurityPreferenceWithCallback(request *SetSecurityPreferenceRequest, callback func(response *SetSecurityPreferenceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetSecurityPreferenceResponse
		var err error
		defer close(result)
		response, err = client.SetSecurityPreference(request)
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

// SetSecurityPreferenceRequest is the request struct for api SetSecurityPreference
type SetSecurityPreferenceRequest struct {
	*requests.RpcRequest
	AllowUserToManageAccessKeys requests.Boolean `position:"Query" name:"AllowUserToManageAccessKeys"`
	AllowUserToManageMFADevices requests.Boolean `position:"Query" name:"AllowUserToManageMFADevices"`
	AllowUserToManagePublicKeys requests.Boolean `position:"Query" name:"AllowUserToManagePublicKeys"`
	EnableSaveMFATicket         requests.Boolean `position:"Query" name:"EnableSaveMFATicket"`
	LoginNetworkMasks           string           `position:"Query" name:"LoginNetworkMasks"`
	AllowUserToChangePassword   requests.Boolean `position:"Query" name:"AllowUserToChangePassword"`
	LoginSessionDuration        requests.Integer `position:"Query" name:"LoginSessionDuration"`
}

// SetSecurityPreferenceResponse is the response struct for api SetSecurityPreference
type SetSecurityPreferenceResponse struct {
	*responses.BaseResponse
	RequestId          string             `json:"RequestId" xml:"RequestId"`
	SecurityPreference SecurityPreference `json:"SecurityPreference" xml:"SecurityPreference"`
}

// CreateSetSecurityPreferenceRequest creates a request to invoke SetSecurityPreference API
func CreateSetSecurityPreferenceRequest() (request *SetSecurityPreferenceRequest) {
	request = &SetSecurityPreferenceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ram", "2015-05-01", "SetSecurityPreference", "", "")
	return
}

// CreateSetSecurityPreferenceResponse creates a response to parse from SetSecurityPreference response
func CreateSetSecurityPreferenceResponse() (response *SetSecurityPreferenceResponse) {
	response = &SetSecurityPreferenceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
