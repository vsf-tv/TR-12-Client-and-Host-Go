# .DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**connect**](DefaultApi.md#connect) | **PUT** /connect | 
[**deprovision**](DefaultApi.md#deprovision) | **PUT** /deprovision | 
[**disconnect**](DefaultApi.md#disconnect) | **PUT** /disconnect | 
[**getConfiguration**](DefaultApi.md#getConfiguration) | **GET** /get_configuration | 
[**getConnectionStatus**](DefaultApi.md#getConnectionStatus) | **GET** /get_state | 
[**reportActualConfiguration**](DefaultApi.md#reportActualConfiguration) | **PUT** /report_actual_configuration | 
[**reportStatus**](DefaultApi.md#reportStatus) | **PUT** /report_status | 


# **connect**
> ConnectResponseContent connect(connectRequestContent)


### Example


```typescript
import { createConfiguration, DefaultApi } from 'CddServiceSDK';
import type { DefaultApiConnectRequest } from 'CddServiceSDK';

const configuration = createConfiguration();
const apiInstance = new DefaultApi(configuration);

const request: DefaultApiConnectRequest = {
  
  connectRequestContent: {
    registration: {
      channels: [
        {
          name: "name_example",
          id: "id_example",
          channelType: "SOURCE",
          standardSettings: [
            {
              id: "id_example",
              name: "name_example",
              description: "description_example",
              enums: {
                values: [
                  "values_example",
                ],
                defaultValue: "defaultValue_example",
              },
              ranges: {
                minimum: 3.14,
                maximum: 3.14,
                defaultValue: 3.14,
              },
            },
          ],
          profiles: [
            {
              name: "name_example",
              id: "id_example",
              description: "description_example",
            },
          ],
          connectionProtocols: [
            "SRT_LISTENER",
          ],
        },
      ],
      standardSettings: [
        {
          id: "id_example",
          name: "name_example",
          description: "description_example",
          enums: {
            values: [
              "values_example",
            ],
            defaultValue: "defaultValue_example",
          },
          ranges: {
            minimum: 3.14,
            maximum: 3.14,
            defaultValue: 3.14,
          },
        },
      ],
      thumbnails: [
        {
          name: "name_example",
          id: "id_example",
          localPath: "localPath_example",
        },
      ],
    },
    hostId: "hostId_example",
  },
};

const data = await apiInstance.connect(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **connectRequestContent** | **ConnectRequestContent**|  |


### Return type

**ConnectResponseContent**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Connect 200 response |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **deprovision**
> DeprovisionResponseContent deprovision()


### Example


```typescript
import { createConfiguration, DefaultApi } from 'CddServiceSDK';
import type { DefaultApiDeprovisionRequest } from 'CddServiceSDK';

const configuration = createConfiguration();
const apiInstance = new DefaultApi(configuration);

const request: DefaultApiDeprovisionRequest = {
  
  hostId: "host_id_example",
  
  force: true,
};

const data = await apiInstance.deprovision(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **hostId** | [**string**] |  | defaults to undefined
 **force** | [**boolean**] |  | (optional) defaults to undefined


### Return type

**DeprovisionResponseContent**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Deprovision 200 response |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **disconnect**
> DisconnectResponseContent disconnect()


### Example


```typescript
import { createConfiguration, DefaultApi } from 'CddServiceSDK';

const configuration = createConfiguration();
const apiInstance = new DefaultApi(configuration);

const request = {};

const data = await apiInstance.disconnect(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters
This endpoint does not need any parameter.


### Return type

**DisconnectResponseContent**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Disconnect 200 response |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **getConfiguration**
> GetConfigurationResponseContent getConfiguration()


### Example


```typescript
import { createConfiguration, DefaultApi } from 'CddServiceSDK';

const configuration = createConfiguration();
const apiInstance = new DefaultApi(configuration);

const request = {};

const data = await apiInstance.getConfiguration(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters
This endpoint does not need any parameter.


### Return type

**GetConfigurationResponseContent**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | GetConfiguration 200 response |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **getConnectionStatus**
> GetConnectionStatusResponseContent getConnectionStatus()


### Example


```typescript
import { createConfiguration, DefaultApi } from 'CddServiceSDK';

const configuration = createConfiguration();
const apiInstance = new DefaultApi(configuration);

const request = {};

const data = await apiInstance.getConnectionStatus(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters
This endpoint does not need any parameter.


### Return type

**GetConnectionStatusResponseContent**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | GetConnectionStatus 200 response |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **reportActualConfiguration**
> ReportActualConfigurationResponseContent reportActualConfiguration(reportActualConfigurationRequestContent)


### Example


```typescript
import { createConfiguration, DefaultApi } from 'CddServiceSDK';
import type { DefaultApiReportActualConfigurationRequest } from 'CddServiceSDK';

const configuration = createConfiguration();
const apiInstance = new DefaultApi(configuration);

const request: DefaultApiReportActualConfigurationRequest = {
  
  reportActualConfigurationRequestContent: {
    configuration: {
      configurationId: "configurationId_example",
      channels: [
        {
          id: "id_example",
          configurationId: "configurationId_example",
          state: "ACTIVE",
          settings: null,
          connection: {
            transportProtocol: null,
          },
          health: null,
        },
      ],
      standardSettings: [
        {
          key: "key_example",
          value: "value_example",
        },
      ],
      health: null,
    },
  },
};

const data = await apiInstance.reportActualConfiguration(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **reportActualConfigurationRequestContent** | **ReportActualConfigurationRequestContent**|  |


### Return type

**ReportActualConfigurationResponseContent**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | ReportActualConfiguration 200 response |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **reportStatus**
> ReportStatusResponseContent reportStatus(reportStatusRequestContent)


### Example


```typescript
import { createConfiguration, DefaultApi } from 'CddServiceSDK';
import type { DefaultApiReportStatusRequest } from 'CddServiceSDK';

const configuration = createConfiguration();
const apiInstance = new DefaultApi(configuration);

const request: DefaultApiReportStatusRequest = {
  
  reportStatusRequestContent: {
    status: {
      status: [
        {
          name: "name_example",
          info: "info_example",
          value: "value_example",
        },
      ],
      channels: [
        {
          id: "id_example",
          state: "ACTIVE",
          status: [
            {
              name: "name_example",
              info: "info_example",
              value: "value_example",
            },
          ],
        },
      ],
    },
  },
};

const data = await apiInstance.reportStatus(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **reportStatusRequestContent** | **ReportStatusRequestContent**|  |


### Return type

**ReportStatusResponseContent**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | ReportStatus 200 response |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


