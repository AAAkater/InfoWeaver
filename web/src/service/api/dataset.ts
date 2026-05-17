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

/** List chunks for a dataset */
export function getDatasetChunks(datasetId: number, page = 1, pageSize = 20) {
  return request<Api.Dataset.Response<Api.Dataset.DatasetChunkListResp>>({
    url: `/dataset/${datasetId}/chunk`,
    method: 'get',
    params: {
      page,
      page_size: pageSize
    }
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

/** Upload files to a dataset */
export function uploadFiles(datasetId: number, files: File[]) {
  const formData = new FormData();
  formData.append('id', String(datasetId));
  files.forEach(file => {
    formData.append('files', file);
  });
  return request<Api.Dataset.Response<Api.Dataset.FileUploadResp>>({
    url: '/file/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
}

/** Split a file into chunks for RAG processing */
export function splitDocument(data: Api.Dataset.SplitDocReq) {
  return request<Api.Dataset.Response<Api.Dataset.SplitDocResp>>({
    url: '/file/split',
    method: 'post',
    data
  });
}

/** Compute embeddings for document chunks */
export function embedChunks(data: Api.Dataset.EmbeddingReq) {
  return request<Api.Dataset.Response<Api.Dataset.EmbeddingResp>>({
    url: '/file/embedding',
    method: 'post',
    data
  });
}

/** Get file list for a dataset */
export function getDatasetFiles(datasetId: number, page = 1, pageSize = 20) {
  return request<Api.Dataset.Response<Api.Dataset.FileListResp>>({
    url: '/file/list',
    method: 'get',
    params: {
      dataset_id: datasetId,
      page,
      page_size: pageSize
    }
  });
}
