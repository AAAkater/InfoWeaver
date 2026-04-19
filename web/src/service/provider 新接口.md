---
title: 默认模块
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# 默认模块

Multi-modal RAG system API

Base URLs:

# Authentication

# v1/Provider

## POST Create Provider

POST /provider

Create a new provider

> Body 请求参数

```json
{
  "api_key": "string",
  "base_url": "string",
  "mode": "openai",
  "name": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|
|body|body|[models.ProviderCreateReq](#schemamodels.providercreatereq)| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "data": null,
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Provider created successfully|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid request parameters|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Invalid or expired token|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Provider name already exists|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[response.ResponseBase-any](#schemaresponse.responsebase-any)|

## POST Delete Provider

POST /provider/delete/{provider_id}

Delete a provider by ID

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|provider_id|path|integer| 是 |Provider ID|
|Authorization|header|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "data": null,
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Provider deleted successfully|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid request parameters|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Invalid or expired token|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Provider not found|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[response.ResponseBase-any](#schemaresponse.responsebase-any)|

## GET Get Provider by ID

GET /provider/info/{provider_id}

Get a provider by its ID

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|provider_id|path|integer| 是 |Provider ID|
|Authorization|header|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "data": {
    "base_url": "string",
    "created_at": "string",
    "id": 0,
    "mode": "string",
    "name": "string"
  },
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Provider retrieved successfully|[response.ResponseBase-models_ProviderInfo](#schemaresponse.responsebase-models_providerinfo)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid request parameters|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Invalid or expired token|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Provider not found|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[response.ResponseBase-any](#schemaresponse.responsebase-any)|

## GET Get All Providers

GET /provider/list

Get a list of all providers

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "data": {
    "providers": [
      {
        "base_url": "string",
        "created_at": "string",
        "id": 0,
        "mode": "string",
        "name": "string"
      }
    ],
    "total": 0
  },
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Providers retrieved successfully|[response.ResponseBase-models_ProviderListResp](#schemaresponse.responsebase-models_providerlistresp)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Invalid or expired token|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[response.ResponseBase-any](#schemaresponse.responsebase-any)|

## GET List Models

GET /provider/models/{provider_id}

Get available models from a provider

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|provider_id|path|integer| 是 |Provider ID|
|Authorization|header|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "data": {
    "models": [
      {
        "enabled": true,
        "id": "string",
        "object": "string",
        "owned_by": "string"
      }
    ]
  },
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Models retrieved successfully|[response.ResponseBase-models_ProviderModelsResp](#schemaresponse.responsebase-models_providermodelsresp)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid request parameters|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Invalid or expired token|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Provider not found|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[response.ResponseBase-any](#schemaresponse.responsebase-any)|

## POST Update Provider

POST /provider/update

Update an existing provider

> Body 请求参数

```json
{
  "api_key": "string",
  "base_url": "string",
  "id": 0,
  "mode": "openai",
  "name": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|
|body|body|[models.ProviderUpdateReq](#schemamodels.providerupdatereq)| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "data": null,
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Provider updated successfully|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid request parameters|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Invalid or expired token|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Provider name already exists|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Provider not found|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[response.ResponseBase-any](#schemaresponse.responsebase-any)|

## POST Set Model Enable Status

POST /provider/models/enable

Enable or disable a single model for a provider

> Body 请求参数

```json
{
  "enabled": true,
  "id": 0,
  "model_id": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|
|body|body|[models.ProviderSetModelEnableReq](#schemamodels.providersetmodelenablereq)| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "data": null,
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Model enable status set successfully|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid request parameters|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Invalid or expired token|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Provider not found|[response.ResponseBase-any](#schemaresponse.responsebase-any)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[response.ResponseBase-any](#schemaresponse.responsebase-any)|

# 数据模型

<h2 id="tocS_response.ResponseBase-any">response.ResponseBase-any</h2>

<a id="schemaresponse.responsebase-any"></a>
<a id="schema_response.ResponseBase-any"></a>
<a id="tocSresponse.responsebase-any"></a>
<a id="tocsresponse.responsebase-any"></a>

```json
{
  "code": 0,
  "data": null,
  "msg": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|code|integer|false|none||none|
|data|any|false|none||none|
|msg|string|false|none||none|

<h2 id="tocS_models.ProviderCreateReq">models.ProviderCreateReq</h2>

<a id="schemamodels.providercreatereq"></a>
<a id="schema_models.ProviderCreateReq"></a>
<a id="tocSmodels.providercreatereq"></a>
<a id="tocsmodels.providercreatereq"></a>

```json
{
  "api_key": "string",
  "base_url": "string",
  "mode": "openai",
  "name": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|api_key|string|true|none||none|
|base_url|string|true|none||none|
|mode|string|true|none||none|
|name|string|true|none||none|

#### 枚举值

|属性|值|
|---|---|
|mode|openai|
|mode|openai_response|
|mode|gemini|
|mode|anthropic|
|mode|ollama|

<h2 id="tocS_models.ModelInfo">models.ModelInfo</h2>

<a id="schemamodels.modelinfo"></a>
<a id="schema_models.ModelInfo"></a>
<a id="tocSmodels.modelinfo"></a>
<a id="tocsmodels.modelinfo"></a>

```json
{
  "enabled": true,
  "id": "string",
  "object": "string",
  "owned_by": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|enabled|boolean|false|none||Enable status for this model|
|id|string|false|none||Model identifier (e.g., "gpt-4", "text-embedding-3-small")|
|object|string|false|none||Object type (usually "model")|
|owned_by|string|false|none||Owner of the model (e.g., "openai")|

<h2 id="tocS_models.ProviderInfo">models.ProviderInfo</h2>

<a id="schemamodels.providerinfo"></a>
<a id="schema_models.ProviderInfo"></a>
<a id="tocSmodels.providerinfo"></a>
<a id="tocsmodels.providerinfo"></a>

```json
{
  "base_url": "string",
  "created_at": "string",
  "id": 0,
  "mode": "string",
  "name": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|base_url|string|false|none||none|
|created_at|string|false|none||none|
|id|integer|false|none||none|
|mode|string|false|none||none|
|name|string|false|none||none|

<h2 id="tocS_models.ProviderListResp">models.ProviderListResp</h2>

<a id="schemamodels.providerlistresp"></a>
<a id="schema_models.ProviderListResp"></a>
<a id="tocSmodels.providerlistresp"></a>
<a id="tocsmodels.providerlistresp"></a>

```json
{
  "providers": [
    {
      "base_url": "string",
      "created_at": "string",
      "id": 0,
      "mode": "string",
      "name": "string"
    }
  ],
  "total": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|providers|[[models.ProviderInfo](#schemamodels.providerinfo)]|false|none||none|
|total|integer|false|none||none|

<h2 id="tocS_models.ProviderUpdateReq">models.ProviderUpdateReq</h2>

<a id="schemamodels.providerupdatereq"></a>
<a id="schema_models.ProviderUpdateReq"></a>
<a id="tocSmodels.providerupdatereq"></a>
<a id="tocsmodels.providerupdatereq"></a>

```json
{
  "api_key": "string",
  "base_url": "string",
  "id": 0,
  "mode": "openai",
  "name": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|api_key|string|true|none||none|
|base_url|string|true|none||none|
|id|integer|true|none||none|
|mode|string|true|none||none|
|name|string|true|none||none|

#### 枚举值

|属性|值|
|---|---|
|mode|openai|
|mode|openai_response|
|mode|gemini|
|mode|anthropic|
|mode|ollama|

<h2 id="tocS_models.ProviderModelsResp">models.ProviderModelsResp</h2>

<a id="schemamodels.providermodelsresp"></a>
<a id="schema_models.ProviderModelsResp"></a>
<a id="tocSmodels.providermodelsresp"></a>
<a id="tocsmodels.providermodelsresp"></a>

```json
{
  "models": [
    {
      "enabled": true,
      "id": "string",
      "object": "string",
      "owned_by": "string"
    }
  ]
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|models|[[models.ModelInfo](#schemamodels.modelinfo)]|false|none||none|

<h2 id="tocS_models.ProviderSetModelEnableReq">models.ProviderSetModelEnableReq</h2>

<a id="schemamodels.providersetmodelenablereq"></a>
<a id="schema_models.ProviderSetModelEnableReq"></a>
<a id="tocSmodels.providersetmodelenablereq"></a>
<a id="tocsmodels.providersetmodelenablereq"></a>

```json
{
  "enabled": true,
  "id": 0,
  "model_id": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|enabled|boolean|true|none||Enable status|
|id|integer|true|none||Provider ID|
|model_id|string|true|none||Model ID to enable/disable|

<h2 id="tocS_response.ResponseBase-models_ProviderInfo">response.ResponseBase-models_ProviderInfo</h2>

<a id="schemaresponse.responsebase-models_providerinfo"></a>
<a id="schema_response.ResponseBase-models_ProviderInfo"></a>
<a id="tocSresponse.responsebase-models_providerinfo"></a>
<a id="tocsresponse.responsebase-models_providerinfo"></a>

```json
{
  "code": 0,
  "data": {
    "base_url": "string",
    "created_at": "string",
    "id": 0,
    "mode": "string",
    "name": "string"
  },
  "msg": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|code|integer|false|none||none|
|data|[models.ProviderInfo](#schemamodels.providerinfo)|false|none||none|
|msg|string|false|none||none|

<h2 id="tocS_response.ResponseBase-models_ProviderListResp">response.ResponseBase-models_ProviderListResp</h2>

<a id="schemaresponse.responsebase-models_providerlistresp"></a>
<a id="schema_response.ResponseBase-models_ProviderListResp"></a>
<a id="tocSresponse.responsebase-models_providerlistresp"></a>
<a id="tocsresponse.responsebase-models_providerlistresp"></a>

```json
{
  "code": 0,
  "data": {
    "providers": [
      {
        "base_url": "string",
        "created_at": "string",
        "id": 0,
        "mode": "string",
        "name": "string"
      }
    ],
    "total": 0
  },
  "msg": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|code|integer|false|none||none|
|data|[models.ProviderListResp](#schemamodels.providerlistresp)|false|none||none|
|msg|string|false|none||none|

<h2 id="tocS_response.ResponseBase-models_ProviderModelsResp">response.ResponseBase-models_ProviderModelsResp</h2>

<a id="schemaresponse.responsebase-models_providermodelsresp"></a>
<a id="schema_response.ResponseBase-models_ProviderModelsResp"></a>
<a id="tocSresponse.responsebase-models_providermodelsresp"></a>
<a id="tocsresponse.responsebase-models_providermodelsresp"></a>

```json
{
  "code": 0,
  "data": {
    "models": [
      {
        "enabled": true,
        "id": "string",
        "object": "string",
        "owned_by": "string"
      }
    ]
  },
  "msg": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|code|integer|false|none||none|
|data|[models.ProviderModelsResp](#schemamodels.providermodelsresp)|false|none||none|
|msg|string|false|none||none|

