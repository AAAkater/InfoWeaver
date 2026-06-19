declare namespace Api {
  namespace Chat {
    /** Create chat session request */
    interface CreateSessionReq {
      title: string
    }

    /** Create chat session response (map[string]uint) */
    interface CreateSessionResp {
      id: number
    }

    /** Chat session info */
    interface ChatSessionInfo {
      id: number
      owner_id: number
      title: string
      created_at: string
      updated_at: string
    }

    /** Chat session list response */
    interface ChatSessionListResp {
      sessions: ChatSessionInfo[]
      total: number
    }

    /** Chat message info */
    interface ChatMessageInfo {
      id: number
      session_id: number
      role: string
      content: string
      created_at: string
      updated_at: string
    }

    /** Chat message list response */
    interface ChatMessageListResp {
      messages: ChatMessageInfo[]
      total: number
    }

    /** LLM sampling params */
    interface SamplingParams {
      temperature?: number
      top_p?: number
      max_tokens?: number
      frequency_penalty?: number
      presence_penalty?: number
    }

    /** LLM config */
    interface LLMConfig {
      model_name: string
      provider_id: number
      sampling_params?: SamplingParams
    }

    /** Retrieval config */
    interface RetrievalConfig {
      top_k?: number
    }

    /** Embedding config */
    interface EmbeddingConfig {
      api_key: string
      base_url: string
      embed_type: string
      model_name: string
      provider_type: string
    }

    /** Send chat stream request */
    interface SendChatStreamReq {
      session_id: number
      query: string
      llm_config: LLMConfig
      retrieval_config: RetrievalConfig
      dataset_id: number
      system_prompt?: string
    }

    /** Common response wrapper */
    interface Response<T = unknown> {
      code: number
      data: T
      msg: string
    }
  }
}
