declare namespace Api {
  namespace Dataset {
    /** Search type enum */
    type SearchType = 'sparse' | 'dense' | 'hybrid';

    /** Dataset create request */
    interface DatasetCreateReq {
      name: string;
      icon: string;
      description?: string;
      search_type: SearchType;
      embedding_model: string;
      provider_id: number;
    }

    /** Dataset update request */
    interface DatasetUpdateReq {
      id: number;
      name: string;
      icon: string;
      description?: string;
      search_type: SearchType;
      embedding_model: string;
      provider_id?: number;
    }

    /** Dataset info (single dataset) */
    interface DatasetInfo {
      id: number;
      name: string;
      icon: string;
      description: string;
      search_type: string;
      embedding_model: string;
      provider_id: number;
      owner_id: number;
      created_at: string;
      updated_at: string;
    }

    /** Dataset list response */
    interface DatasetListResp {
      datasets: DatasetInfo[];
      total: number;
    }

    /** Common response wrapper */
    interface Response<T = any> {
      code: number;
      data: T;
      msg: string;
    }

    /** Form model for create/edit (legacy compatibility) */
    interface FormModel {
      icon: string;
      name: string;
      description: string;
      search_type: SearchType;
      embedding_model: string;
      provider_id: number;
      id?: number;
    }

    /** Dataset item (legacy compatibility) */
    interface DatasetItem {
      created_at: string;
      description: string;
      icon: string;
      id: number;
      name: string;
      owner_id: number;
      updated_at: string;
      search_type: SearchType;
      embedding_model: string;
      provider_id: number;
    }

    /** Get dataset response (legacy compatibility) */
    interface GetDatasetResponse {
      datasets: DatasetItem[];
    }

    /** Dataset response (legacy compatibility) */
    interface DatasetResponse {
      code: number;
      data: any;
      msg: string;
    }
  }
}
