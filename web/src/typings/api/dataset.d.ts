declare namespace Api {
  namespace Dataset {
    interface FormModel {
      icon: string;
      description: string;
      name: string;
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
