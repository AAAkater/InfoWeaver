declare namespace Api {
  namespace Mcp {
    /** Transport type */
    type TransportType = 'stdio' | 'sse' | 'streamable_http';

    /** MCP create request */
    interface McpCreateReq {
      name: string;
      transport: TransportType;
      command?: string;
      args?: string;
      url?: string;
      enabled?: boolean;
      env_vars?: Record<string, any>;
      headers?: Record<string, any>;
    }

    /** MCP update request */
    interface McpUpdateReq {
      id: number;
      name: string;
      transport: TransportType;
      command?: string;
      args?: string;
      url?: string;
      enabled?: boolean;
      env_vars?: Record<string, any>;
      headers?: Record<string, any>;
    }

    /** MCP info (single server) */
    interface McpInfo {
      id: number;
      name: string;
      transport: TransportType;
      command?: string;
      args?: string;
      url?: string;
      enabled: boolean;
      env_vars?: Record<string, any>;
      headers?: Record<string, any>;
      created_at: string;
      updated_at: string;
    }

    /** MCP list response */
    interface McpListResp {
      mcps: McpInfo[];
      page: number;
      total: number;
    }

    /** Common response wrapper */
    interface Response<T = any> {
      code: number;
      data: T;
      msg: string;
    }
  }
}
