# DefaultApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**connect**](DefaultApi.md#connect) | **PUT** /connect |  |
| [**deprovision**](DefaultApi.md#deprovision) | **PUT** /deprovision |  |
| [**disconnect**](DefaultApi.md#disconnect) | **PUT** /disconnect |  |
| [**getConfiguration**](DefaultApi.md#getconfiguration) | **GET** /get_configuration |  |
| [**getConnectionStatus**](DefaultApi.md#getconnectionstatus) | **GET** /get_state |  |
| [**reportActualConfiguration**](DefaultApi.md#reportactualconfiguration) | **PUT** /report_actual_configuration |  |
| [**reportStatus**](DefaultApi.md#reportstatus) | **PUT** /report_status |  |



## connect

> ConnectResponseContent connect(connectRequestContent)



### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '';
import type { ConnectRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DefaultApi();

  const body = {
    // ConnectRequestContent
    connectRequestContent: ...,
  } satisfies ConnectRequest;

  try {
    const data = await api.connect(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **connectRequestContent** | [ConnectRequestContent](ConnectRequestContent.md) |  | |

### Return type

[**ConnectResponseContent**](ConnectResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Connect 200 response |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deprovision

> DeprovisionResponseContent deprovision(hostId, force)



### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '';
import type { DeprovisionRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DefaultApi();

  const body = {
    // string
    hostId: hostId_example,
    // boolean (optional)
    force: true,
  } satisfies DeprovisionRequest;

  try {
    const data = await api.deprovision(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **hostId** | `string` |  | [Defaults to `undefined`] |
| **force** | `boolean` |  | [Optional] [Defaults to `undefined`] |

### Return type

[**DeprovisionResponseContent**](DeprovisionResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Deprovision 200 response |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## disconnect

> DisconnectResponseContent disconnect()



### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '';
import type { DisconnectRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DefaultApi();

  try {
    const data = await api.disconnect();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**DisconnectResponseContent**](DisconnectResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Disconnect 200 response |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getConfiguration

> GetConfigurationResponseContent getConfiguration()



### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '';
import type { GetConfigurationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DefaultApi();

  try {
    const data = await api.getConfiguration();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**GetConfigurationResponseContent**](GetConfigurationResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | GetConfiguration 200 response |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getConnectionStatus

> GetConnectionStatusResponseContent getConnectionStatus()



### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '';
import type { GetConnectionStatusRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DefaultApi();

  try {
    const data = await api.getConnectionStatus();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**GetConnectionStatusResponseContent**](GetConnectionStatusResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | GetConnectionStatus 200 response |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## reportActualConfiguration

> ReportActualConfigurationResponseContent reportActualConfiguration(reportActualConfigurationRequestContent)



### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '';
import type { ReportActualConfigurationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DefaultApi();

  const body = {
    // ReportActualConfigurationRequestContent
    reportActualConfigurationRequestContent: ...,
  } satisfies ReportActualConfigurationRequest;

  try {
    const data = await api.reportActualConfiguration(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **reportActualConfigurationRequestContent** | [ReportActualConfigurationRequestContent](ReportActualConfigurationRequestContent.md) |  | |

### Return type

[**ReportActualConfigurationResponseContent**](ReportActualConfigurationResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | ReportActualConfiguration 200 response |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## reportStatus

> ReportStatusResponseContent reportStatus(reportStatusRequestContent)



### Example

```ts
import {
  Configuration,
  DefaultApi,
} from '';
import type { ReportStatusRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DefaultApi();

  const body = {
    // ReportStatusRequestContent
    reportStatusRequestContent: ...,
  } satisfies ReportStatusRequest;

  try {
    const data = await api.reportStatus(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **reportStatusRequestContent** | [ReportStatusRequestContent](ReportStatusRequestContent.md) |  | |

### Return type

[**ReportStatusResponseContent**](ReportStatusResponseContent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | ReportStatus 200 response |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

