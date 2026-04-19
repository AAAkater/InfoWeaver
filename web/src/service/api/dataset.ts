import { request } from '../request';

/** Get all datasets (with optional name filter) */
export function getDatasets(name?: string) {
  return request<Api.Dataset.Response<Api.Dataset.DatasetListResp>>({
    url: '/dataset',
    method: 'get',
    params: {
      name
    }
  });
}

/** Get a specific dataset by ID */
export function getDatasetById(datasetId: number) {
  return request<Api.Dataset.Response<Api.Dataset.DatasetInfo>>({
    url: `/dataset/${datasetId}`,
    method: 'get'
  });
}

/** Create a new dataset */
export function createDataset(formModel: Api.Dataset.FormModel) {
  return request<Api.Dataset.Response>({
    url: '/dataset/create',
    method: 'post',
    data: formModel
  });
}

/** Delete a dataset by ID */
export function deleteDataset(id: number) {
  return request<Api.Dataset.Response>({
    url: `/dataset/delete/${id}`,
    method: 'post'
  });
}

/** Update an existing dataset */
export function updateDataset(formModel: Api.Dataset.FormModel) {
  return request<Api.Dataset.Response>({
    url: '/dataset/update',
    method: 'post',
    data: formModel
  });
}
