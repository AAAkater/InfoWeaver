declare namespace Api {
  namespace Dataset {
    interface FormModel {
      icon: string;
      name: string;
      description: string;
      search_type: 'sparse' | 'dense' | 'hybrid';
      embedding_model: string;
      provider_id: number;
      id?: number;
    }
    interface DatasetItem {
      created_at: string;
      description: string;
      icon: string;
      id: number;
      name: string;
      owner_id: number;
      updated_at: string;
      search_type: 'sparse' | 'dense' | 'hybrid';
      embedding_model: string;
      provider_id: number;
    }
    interface GetDatasetResponse {
      datasets: DatasetItem[];
    }
    interface DatasetResponse {
      code: number;
      data: any;
      msg: string;
    }
  }
}
