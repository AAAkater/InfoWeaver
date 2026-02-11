import { request } from '../request';
export function getDatasets() {
  return request<Api.Dataset.DatasetItem[]>({
    url: '/dataset',
    method: 'get'
  });
}
export function createDataset(formModel: Api.Dataset.FormModel) {
  return request<Api.Dataset.FormModel>({
    url: '/dataset/create',
    method: 'post',
    data: formModel
  });
}

export function deleteDataset(id: number) {
  return request<Api.Dataset.DeleteResponse>({
    url: `/dataset/delete/${id}`,
    method: 'post'
  });
}
