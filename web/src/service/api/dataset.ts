import { request } from '../request';
export function getDatasets(id?: number) {
  const url = id !== undefined ? `/dataset/${id}` : '/dataset';
  return request<Api.Dataset.GetDatasetResponse>({
    url,
    method: 'get'
  });
}
export function createDataset(formModel: Api.Dataset.FormModel) {
  return request<Api.Dataset.DatasetResponse>({
    url: '/dataset/create',
    method: 'post',
    data: formModel
  });
}

export function deleteDataset(id: number) {
  return request<Api.Dataset.DatasetResponse>({
    url: `/dataset/delete/${id}`,
    method: 'post'
  });
}

export function updateDataset(formModel: Api.Dataset.FormModel) {
  return request<Api.Dataset.DatasetResponse>({
    url: '/dataset/update',
    method: 'post',
    data: formModel
  });
}
