# \DefaultAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Connect**](DefaultAPI.md#Connect) | **Put** /connect | 
[**Deprovision**](DefaultAPI.md#Deprovision) | **Put** /deprovision | 
[**Disconnect**](DefaultAPI.md#Disconnect) | **Put** /disconnect | 
[**GetConfiguration**](DefaultAPI.md#GetConfiguration) | **Get** /get_configuration | 
[**GetConnectionStatus**](DefaultAPI.md#GetConnectionStatus) | **Get** /get_state | 
[**ReportActualConfiguration**](DefaultAPI.md#ReportActualConfiguration) | **Put** /report_actual_configuration | 
[**ReportStatus**](DefaultAPI.md#ReportStatus) | **Put** /report_status | 



## Connect

> ConnectResponseContent Connect(ctx).ConnectRequestContent(connectRequestContent).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

func main() {
	connectRequestContent := *openapiclient.NewConnectRequestContent(*openapiclient.NewDeviceRegistration([]openapiclient.Channel{*openapiclient.NewChannel("Name_example", "Id_example")}), "HostId_example") // ConnectRequestContent | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.Connect(context.Background()).ConnectRequestContent(connectRequestContent).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.Connect``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Connect`: ConnectResponseContent
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.Connect`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiConnectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **connectRequestContent** | [**ConnectRequestContent**](ConnectRequestContent.md) |  | 

### Return type

[**ConnectResponseContent**](ConnectResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Deprovision

> DeprovisionResponseContent Deprovision(ctx).HostId(hostId).Force(force).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

func main() {
	hostId := "hostId_example" // string | 
	force := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.Deprovision(context.Background()).HostId(hostId).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.Deprovision``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Deprovision`: DeprovisionResponseContent
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.Deprovision`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeprovisionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **hostId** | **string** |  | 
 **force** | **bool** |  | 

### Return type

[**DeprovisionResponseContent**](DeprovisionResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Disconnect

> DisconnectResponseContent Disconnect(ctx).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.Disconnect(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.Disconnect``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Disconnect`: DisconnectResponseContent
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.Disconnect`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDisconnectRequest struct via the builder pattern


### Return type

[**DisconnectResponseContent**](DisconnectResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetConfiguration

> GetConfigurationResponseContent GetConfiguration(ctx).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetConfiguration(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetConfiguration``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetConfiguration`: GetConfigurationResponseContent
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetConfiguration`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetConfigurationRequest struct via the builder pattern


### Return type

[**GetConfigurationResponseContent**](GetConfigurationResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetConnectionStatus

> GetConnectionStatusResponseContent GetConnectionStatus(ctx).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetConnectionStatus(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetConnectionStatus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetConnectionStatus`: GetConnectionStatusResponseContent
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetConnectionStatus`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetConnectionStatusRequest struct via the builder pattern


### Return type

[**GetConnectionStatusResponseContent**](GetConnectionStatusResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReportActualConfiguration

> ReportActualConfigurationResponseContent ReportActualConfiguration(ctx).ReportActualConfigurationRequestContent(reportActualConfigurationRequestContent).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

func main() {
	reportActualConfigurationRequestContent := *openapiclient.NewReportActualConfigurationRequestContent(*openapiclient.NewDeviceConfiguration("ConfigurationId_example", []openapiclient.ChannelConfiguration{*openapiclient.NewChannelConfiguration("Id_example", "ConfigurationId_example", openapiclient.ChannelState("ACTIVE"))})) // ReportActualConfigurationRequestContent | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ReportActualConfiguration(context.Background()).ReportActualConfigurationRequestContent(reportActualConfigurationRequestContent).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ReportActualConfiguration``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReportActualConfiguration`: ReportActualConfigurationResponseContent
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ReportActualConfiguration`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiReportActualConfigurationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **reportActualConfigurationRequestContent** | [**ReportActualConfigurationRequestContent**](ReportActualConfigurationRequestContent.md) |  | 

### Return type

[**ReportActualConfigurationResponseContent**](ReportActualConfigurationResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReportStatus

> ReportStatusResponseContent ReportStatus(ctx).ReportStatusRequestContent(reportStatusRequestContent).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

func main() {
	reportStatusRequestContent := *openapiclient.NewReportStatusRequestContent(*openapiclient.NewDeviceStatus([]openapiclient.StatusValue{*openapiclient.NewStatusValue("Name_example", "Info_example", "Value_example")})) // ReportStatusRequestContent | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ReportStatus(context.Background()).ReportStatusRequestContent(reportStatusRequestContent).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ReportStatus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReportStatus`: ReportStatusResponseContent
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ReportStatus`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiReportStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **reportStatusRequestContent** | [**ReportStatusRequestContent**](ReportStatusRequestContent.md) |  | 

### Return type

[**ReportStatusResponseContent**](ReportStatusResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

