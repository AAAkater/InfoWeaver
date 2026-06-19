import { localStg } from "@/utils/storage"
import { request } from "../request"

const BASE_URL = "/api/v1"

function getAuthHeaders(): Record<string, string> {
  const token = localStg.get("token")
  return token ? { Authorization: `Bearer ${token}` } : {}
}

/** Create a new chat session */
export function createChatSession(data: Api.Chat.CreateSessionReq) {
  return request<Api.Chat.Response<Api.Chat.CreateSessionResp>>({
    url: "/chat-session/create",
    method: "post",
    data,
  })
}

/** List all chat sessions for the current user */
export function listChatSessions() {
  return request<Api.Chat.Response<Api.Chat.ChatSessionListResp>>({
    url: "/chat-session",
    method: "get",
  })
}

/** Delete a chat session by ID */
export function deleteChatSession(sessionId: number) {
  return request({
    url: `/chat-session/delete/${sessionId}`,
    method: "post",
  })
}

/** List all messages in a session */
export function listSessionMessages(sessionId: number) {
  return request<Api.Chat.Response<Api.Chat.ChatMessageListResp>>({
    url: `/chat/session/${sessionId}`,
    method: "get",
  })
}

export interface StreamChunk {
  type: "thinking" | "text" | "source" | "error"
  content: string
}

function normalizeStreamPayload(payload: string): StreamChunk[] {
  const trimmed = payload.trim()
  if (!trimmed || trimmed === "[DONE]") return []

  try {
    const parsed = JSON.parse(trimmed)
    if (parsed.type && typeof parsed.content === "string") {
      return [{ type: parsed.type as StreamChunk["type"], content: parsed.content }]
    }
    if (typeof parsed.content === "string") {
      return [{ type: "text", content: parsed.content }]
    }
  } catch {
    // Fall through and render non-JSON payloads as text.
  }

  return [{ type: "text", content: trimmed }]
}

function extractJsonPayloads(input: string, flush = false) {
  const payloads: string[] = []
  let start = -1
  let depth = 0
  let inString = false
  let escaped = false

  for (let i = 0; i < input.length; i++) {
    const char = input[i]

    if (start === -1) {
      if (char === "{") {
        start = i
        depth = 1
      }
      continue
    }

    if (escaped) {
      escaped = false
      continue
    }

    if (char === "\\") {
      escaped = true
      continue
    }

    if (char === '"') {
      inString = !inString
      continue
    }

    if (inString) continue

    if (char === "{") {
      depth += 1
    } else if (char === "}") {
      depth -= 1
      if (depth === 0) {
        payloads.push(input.slice(start, i + 1))
        start = -1
      }
    }
  }

  const rest = start === -1 ? "" : input.slice(start)
  return {
    payloads,
    rest: flush ? "" : rest,
  }
}

function parseStreamBuffer(input: string, flush = false) {
  const chunks: StreamChunk[] = []
  const normalized = input.replace(/\r\n/g, "\n")

  if (!normalized.includes("\n\n")) {
    const { payloads, rest } = extractJsonPayloads(normalized, flush)
    for (const payload of payloads) {
      for (const chunk of normalizeStreamPayload(payload)) {
        chunks.push(chunk)
      }
    }

    if (chunks.length > 0) {
      return { chunks, rest }
    }

    if (flush) {
      for (const chunk of normalizeStreamPayload(normalized)) {
        chunks.push(chunk)
      }
      return { chunks, rest: "" }
    }

    return { chunks, rest: normalized }
  }

  const messages = normalized.split("\n\n")
  const tail = messages.pop() ?? ""

  for (const message of messages) {
    const dataLines = message
      .split("\n")
      .map((line) => line.trimEnd())
      .filter((line) => line.startsWith("data:"))
      .map((line) => line.replace(/^data:\s?/, ""))

    if (dataLines.length > 0) {
      for (const chunk of normalizeStreamPayload(dataLines.join("\n"))) {
        chunks.push(chunk)
      }
      continue
    }

    for (const line of message.split("\n")) {
      for (const chunk of normalizeStreamPayload(line)) {
        chunks.push(chunk)
      }
    }
  }

  if (!tail) {
    return { chunks, rest: "" }
  }

  const parsedTail = parseStreamBuffer(tail, flush)
  chunks.push(...parsedTail.chunks)

  return { chunks, rest: parsedTail.rest }
}

/**
 * Send chat message via SSE streaming.
 * Returns an async generator that yields typed content chunks.
 */
export async function* sendChatMessageStream(
  data: Api.Chat.SendChatStreamReq,
): AsyncGenerator<StreamChunk, void, unknown> {
  const response = await fetch(`${BASE_URL}/chat/send`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify(data),
  })

  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(errorText || `请求失败 (${response.status})`)
  }

  const reader = response.body?.getReader()
  if (!reader) throw new Error("无法读取响应流")

  const decoder = new TextDecoder()
  let buffer = ""

  try {
    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const parsed = parseStreamBuffer(buffer)
      buffer = parsed.rest

      for (const chunk of parsed.chunks) {
        yield chunk
      }
    }

    buffer += decoder.decode()
    const parsed = parseStreamBuffer(buffer, true)
    for (const chunk of parsed.chunks) {
      yield chunk
    }
  } finally {
    reader.releaseLock()
  }
}
