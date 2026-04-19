import { request } from '../request';

/** Create a new provider */
export function createProvider(data: Api.Provider.ProviderCreateReq) {
  return request<Api.Provider.Response>({
    url: '/provider',
    method: 'post',
    data
  });
}

/** Delete a provider by ID */
export function deleteProvider(providerId: number) {
  return request<Api.Provider.Response>({
    url: `/provider/delete/${providerId}`,
    method: 'post'
  });
}

/** Get a provider by ID */
export function getProviderById(providerId: number) {
  return request<Api.Provider.Response<Api.Provider.ProviderInfo>>({
    url: `/provider/info/${providerId}`,
    method: 'get'
  });
}

/** Get all providers */
export function getProviderList() {
  return request<Api.Provider.Response<Api.Provider.ProviderListResp>>({
    url: '/provider/list',
    method: 'get'
  });
}

/** Get available models from a provider */
export function getProviderModels(providerId: number) {
  return request<Api.Provider.Response<Api.Provider.ProviderModelsResp>>({
    url: `/provider/models/${providerId}`,
    method: 'get'
  });
}

/** Update an existing provider */
export function updateProvider(data: Api.Provider.ProviderUpdateReq) {
  return request<Api.Provider.Response>({
    url: '/provider/update',
    method: 'post',
    data
  });
}

/** Set model enable status */
export function setProviderModelEnable(data: Api.Provider.ProviderSetModelEnableReq) {
  return request<Api.Provider.Response>({
    url: '/provider/models/enable',
    method: 'post',
    data
  });
}
