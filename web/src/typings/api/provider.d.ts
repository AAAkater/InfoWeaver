declare namespace Api {
  /**
   * namespace Provider
   *
   * backend api module: "provider"
   */
  namespace Provider {
    /** Provider mode enum */
    type ProviderMode = 'openai' | 'openai_response' | 'gemini' | 'anthropic' | 'ollama';

    /** Provider create request */
    interface ProviderCreateReq {
      api_key: string;
      base_url: string;
      mode: ProviderMode;
      name: string;
    }

    /** Provider update request */
    interface ProviderUpdateReq {
      id: number;
      api_key: string;
      base_url: string;
      mode: ProviderMode;
      name: string;
    }

    /** Provider info */
    interface ProviderInfo {
      id: number;
      name: string;
      base_url: string;
      mode: string;
      created_at: string;
    }

    /** Provider list response */
    interface ProviderListResp {
      providers: ProviderInfo[];
      total: number;
    }

    /** Model info */
    interface ModelInfo {
      id: string;
      object: string;
      owned_by: string;
      enabled?: boolean;
    }

    /** Provider models response */
    interface ProviderModelsResp {
      models: ModelInfo[];
    }

    /** Provider add models request */
    interface ProviderAddModelsReq {
      id: number;
      available_models: string[];
    }

    /** Provider set model enable request */
    interface ProviderSetModelEnableReq {
      id: number;
      model_id: string;
      enabled: boolean;
    }

    /** Common response */
    interface Response<T = any> {
      code: number;
      data: T;
      msg: string;
    }
  }
}
