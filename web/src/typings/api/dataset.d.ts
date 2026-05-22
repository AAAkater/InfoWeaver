declare namespace Api {
  namespace Dataset {
    /** Search type enum */
    type SearchType = "sparse" | "dense" | "hybrid"

    /** Dataset create request */
    interface DatasetCreateReq {
      name: string
      icon: string
      description?: string
      search_type: SearchType
      embedding_model: string
      provider_id: number
    }

    /** Dataset update request */
    interface DatasetUpdateReq {
      id: number
      name: string
      icon: string
      description?: string
      search_type: SearchType
      embedding_model: string
      provider_id?: number
    }

    /** Dataset info (single dataset) */
    interface DatasetInfo {
      id: number
      name: string
      icon: string
      description: string
      search_type: string
      embedding_model: string
      provider_id: number
      owner_id: number
      created_at: string
      updated_at: string
    }

    /** Dataset list response */
    interface DatasetListResp {
      datasets: DatasetInfo[]
      total: number
    }

    /** Dataset chunk info */
    interface ChunkInfo {
      chunk_metadata: Record<string, string>
      content: string
      created_at: string
      file_id: number
      id: number
      status: string
      updated_at: string
      vector_id: string
    }

    /** Dataset chunk list response */
    interface DatasetChunkListResp {
      chunks: ChunkInfo[]
      total: number
    }

    /** Common response wrapper */
    interface Response<T = any> {
      code: number
      data: T
      msg: string
    }

    /** Form model for create/edit (legacy compatibility) */
    interface FormModel {
      icon: string
      name: string
      description: string
      search_type: SearchType
      embedding_model: string
      provider_id: number
      id?: number
    }

    /** Dataset item (legacy compatibility) */
    interface DatasetItem {
      created_at: string
      description: string
      icon: string
      id: number
      name: string
      owner_id: number
      updated_at: string
      search_type: SearchType
      embedding_model: string
      provider_id: number
    }

    /** Get dataset response (legacy compatibility) */
    interface GetDatasetResponse {
      datasets: DatasetItem[]
    }

    /** File upload info */
    interface FileUploadInfo {
      id: number
      name: string
      type: string
      size: number
      dataset_id: number
      owner_id: number
    }

    /** File upload response */
    interface FileUploadResp {
      files: FileUploadInfo[]
    }

    /** File split request */
    interface SplitDocReq {
      file_id: number
      dataset_id: number
      minio_path: string
      chunk_size: number
      chunk_overlap: number
    }

    /** File split response */
    interface SplitDocResp {
      file_id: number
      dataset_id: number
      file_name: string
      chunks_count: number
    }

    /** Embedding config */
    interface EmbeddingConfig {
      model_name: string
      provider_type: string
      api_key: string
      base_url: string
    }

    /** Embedding request */
    interface EmbeddingReq {
      chunk_ids: number[]
      embedding_config: EmbeddingConfig
    }

    /** Embedding response */
    interface EmbeddingResp {
      chunk_ids: number[]
      chunks_count: number
    }

    /** Simple file info (for list) */
    interface SimpleFileInfo {
      id: number
      name: string
      type?: string
      size?: number
      createdAt?: string
    }

    /** File list response */
    interface FileListResp {
      files: SimpleFileInfo[]
      total: number
    }

    /** Dataset response (legacy compatibility) */
    interface DatasetResponse {
      code: number
      data: any
      msg: string
    }
  }
}
