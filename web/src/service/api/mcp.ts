import { request } from '../request';

/** Get all MCP servers (paginated) */
export function getMcpList(page = 1, pageSize = 20) {
  return request<Api.Mcp.Response<Api.Mcp.McpListResp>>({
    url: '/mcp/list',
    method: 'get',
    params: {
      page,
      page_size: pageSize
    }
  });
}

/** Get a specific MCP server by ID */
export function getMcpById(mcpId: number) {
  return request<Api.Mcp.Response<Api.Mcp.McpInfo>>({
    url: `/mcp/info/${mcpId}`,
    method: 'get'
  });
}

/** Create a new MCP server */
export function createMcp(data: Api.Mcp.McpCreateReq) {
  return request<Api.Mcp.Response>({
    url: '/mcp',
    method: 'post',
    data
  });
}

/** Update an existing MCP server */
export function updateMcp(data: Api.Mcp.McpUpdateReq) {
  return request<Api.Mcp.Response>({
    url: '/mcp/update',
    method: 'post',
    data
  });
}

/** Delete an MCP server by ID */
export function deleteMcp(mcpId: number) {
  return request<Api.Mcp.Response>({
    url: `/mcp/delete/${mcpId}`,
    method: 'post'
  });
}
