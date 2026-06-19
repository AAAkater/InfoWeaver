<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from "vue"
import {
  sendChatMessageStream,
  createChatSession,
  listChatSessions,
  deleteChatSession,
  listSessionMessages,
} from "@/service/api/chat"
import { getDatasets } from "@/service/api/dataset"
import { getProviderList, getProviderModels } from "@/service/api/provider"
import { localStg } from "@/utils/storage"
import ChatHeader from "./modules/chat-header.vue"
import ChatSidebar from "./modules/chat-sidebar.vue"
import ChatMessageArea from "./modules/chat-message-area.vue"
import ChatInputArea from "./modules/chat-input-area.vue"

// ── Types ──

interface Message {
  id: string
  role: "user" | "assistant"
  content: string
  thinking: string
  timestamp: Date
  sources?: RetrievedSource[]
}

interface RetrievedSource {
  id: number
  title: string
  content: string
  score: number
  datasetName: string
}

interface Conversation {
  id: number
  title: string
  messages: Message[]
  createdAt: Date
}

interface ModelOption {
  label: string
  value: string
  provider: string
  providerId: number
  providerMode: string
}

// ── State ──

const conversations = ref<Conversation[]>([])
const activeConversationId = ref<number>()
const inputMessage = ref("")
const isLoading = ref(false)
const searchQuery = ref("")
const showSettings = ref(false)
const showModelSelect = ref(false)
const mcpEnabled = ref(false)

const datasetId = ref<number>()
const datasetName = ref("")
const datasetOptions = ref<{ label: string; value: number }[]>([])

const modelOptions = ref<ModelOption[]>([])

const ragSettings = ref({
  topK: 3,
  similarityThreshold: 0.7,
  enableRerank: true,
})

const llmSettings = ref({
  model: "",
  temperature: 0.7,
  topP: 0.9,
  maxTokens: 2048,
  frequencyPenalty: 0,
  presencePenalty: 0,
})

const messageAreaRef = ref<InstanceType<typeof ChatMessageArea>>()

// ── Computed ──

const activeConversation = computed(() =>
  conversations.value.find((c) => c.id === activeConversationId.value),
)

const messages = computed(() => activeConversation?.value?.messages || [])

const filteredConversations = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return conversations.value
  return conversations.value.filter((c) => c.title.toLowerCase().includes(q))
})

const selectedModelOption = computed(() =>
  modelOptions.value.find((m) => m.value === llmSettings.value.model),
)

// ── Helpers ──

function handleSelectModel(model: string) {
  llmSettings.value.model = model
  localStg.set("chatSelectedModel", model)
}

function handleSelectDataset(dataset: { label: string; value: number }) {
  datasetId.value = dataset.value
  datasetName.value = dataset.label
  localStg.set("chatSelectedDatasetId", dataset.value)
}

async function selectConversation(id: number) {
  activeConversationId.value = id
  const conv = activeConversation.value
  if (!conv || conv.messages.length > 0) return

  try {
    const { data, error } = await listSessionMessages(id)
    if (error || !data) return
    conv.messages = data.messages.map((m) => ({
      id: String(m.id),
      role: m.role as "user" | "assistant",
      content: m.content,
      thinking: "",
      timestamp: new Date(m.created_at),
    }))
  } catch {
    // failed to load messages
  }
}

function scrollToBottom() {
  nextTick(() => {
    messageAreaRef.value?.scrollToBottom()
  })
}

// ── Conversation CRUD ──

async function createNewConversation() {
  const { data, error } = await createChatSession({ title: "新对话" })
  if (error || !data?.id) {
    window.$message?.error("创建会话失败")
    return false
  }
  const newConv: Conversation = {
    id: data.id,
    title: "新对话",
    messages: [],
    createdAt: new Date(),
  }
  conversations.value.unshift(newConv)
  activeConversationId.value = newConv.id
  return true
}

async function deleteConversation(id: number) {
  const index = conversations.value.findIndex((c) => c.id === id)
  if (index === -1) return

  const { error } = await deleteChatSession(id)
  if (error) {
    window.$message?.error("删除会话失败")
    return
  }

  conversations.value.splice(index, 1)
  if (activeConversationId.value === id) {
    activeConversationId.value = conversations.value[0]?.id
  }
}

// ── Data fetching ──

async function fetchModelOptions() {
  try {
    const { data, error } = await getProviderList()
    if (error || !data) return
    const allModels: ModelOption[] = []
    for (const p of data.providers) {
      try {
        const modelRes = await getProviderModels(p.id)
        if (modelRes.error || !modelRes.data) continue
        for (const m of modelRes.data.models) {
          if (!m.enabled) continue
          allModels.push({
            label: m.id,
            value: m.id,
            provider: p.name,
            providerId: p.id,
            providerMode: p.mode,
          })
        }
      } catch {
        // skip providers that fail to load models
      }
    }
    modelOptions.value = allModels

    // Restore persisted model if still valid
    const savedModel = localStg.get("chatSelectedModel") as string | null
    if (savedModel && allModels.some((m) => m.value === savedModel)) {
      llmSettings.value.model = savedModel
    } else if (allModels.length > 0 && !llmSettings.value.model) {
      llmSettings.value.model = allModels[0].value
    }
  } catch {
    // no providers available
  }
}

async function fetchDatasets() {
  try {
    const { data, error } = await getDatasets()
    if (error || !data) return
    if (data.datasets.length > 0) {
      datasetOptions.value = data.datasets.map((d) => ({
        label: d.name,
        value: d.id,
      }))

      // Restore persisted dataset if still valid
      const savedDatasetId = localStg.get("chatSelectedDatasetId") as number | null
      if (savedDatasetId && data.datasets.some((d) => d.id === savedDatasetId)) {
        datasetId.value = savedDatasetId
        datasetName.value = data.datasets.find((d) => d.id === savedDatasetId)!.name
      } else {
        datasetId.value = data.datasets[0].id
        datasetName.value = data.datasets[0].name
      }
    }
  } catch {
    // no datasets available
  }
}

async function fetchConversations() {
  try {
    const { data, error } = await listChatSessions()
    if (error || !data) return
    conversations.value = data.sessions.map((s) => ({
      id: s.id,
      title: s.title,
      messages: [],
      createdAt: new Date(s.created_at),
    }))
  } catch {
    // no sessions available
  }
}

// ── Send message ──

function handleKeydown(e: KeyboardEvent) {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault()
    sendMessage()
  }
}

async function sendMessage() {
  if (!inputMessage.value.trim()) {
    if (!isLoading.value) {
      window.$message?.warning("请输入内容后再发送")
    }
    return
  }
  if (isLoading.value) return

  if (!activeConversationId.value) {
    const ok = await createNewConversation()
    if (!ok) return
  }

  const userMessage: Message = {
    id: crypto.randomUUID(),
    role: "user",
    content: inputMessage.value.trim(),
    thinking: "",
    timestamp: new Date(),
  }

  activeConversation!.value!.messages.push(userMessage)

  if (activeConversation!.value!.messages.length === 1) {
    const title = userMessage.content.slice(0, 20)
    activeConversation!.value!.title = userMessage.content.length > 20 ? `${title}...` : title
  }

  inputMessage.value = ""
  isLoading.value = true

  try {
    const conv = activeConversation!.value!

    const assistantMessage: Message = {
      id: crypto.randomUUID(),
      role: "assistant",
      content: "",
      thinking: "",
      timestamp: new Date(),
    }
    conv.messages.push(assistantMessage)
    const assistantIndex = conv.messages.length - 1

    scrollToBottom()

    const modelOpt = selectedModelOption.value
    const sendReq: Api.Chat.SendChatStreamReq = {
      session_id: conv.id,
      query: userMessage.content,
      dataset_id: datasetId.value ?? 0,
      llm_config: {
        model_name: llmSettings.value.model,
        provider_id: modelOpt?.providerId ?? 0,
        sampling_params: {
          temperature: llmSettings.value.temperature,
          top_p: llmSettings.value.topP,
          max_tokens: llmSettings.value.maxTokens,
          frequency_penalty: llmSettings.value.frequencyPenalty,
          presence_penalty: llmSettings.value.presencePenalty,
        },
      },
      retrieval_config: { top_k: ragSettings.value.topK },
    }

    for await (const chunk of sendChatMessageStream(sendReq)) {
      const streamingMessage = conv.messages[assistantIndex]
      if (!streamingMessage) continue

      if (chunk.type === "thinking") {
        streamingMessage.thinking += chunk.content
      } else if (chunk.type === "source") {
        try {
          streamingMessage.sources = JSON.parse(chunk.content)
        } catch {
          // ignore malformed source JSON
        }
      } else {
        streamingMessage.content += chunk.content
      }
      // Yield to browser so it can paint before processing the next chunk
      await new Promise((r) => requestAnimationFrame(r))
      scrollToBottom()
    }
  } catch {
    window.$message?.error("发送消息失败")
  } finally {
    isLoading.value = false
    scrollToBottom()
  }
}

// ── Lifecycle ──

onMounted(() => {
  fetchConversations()
  fetchModelOptions()
  fetchDatasets()
})
</script>

<template>
  <div class="h-full flex">
    <ChatSidebar
      :conversations="filteredConversations"
      :active-conversation-id="activeConversationId"
      :search-query="searchQuery"
      @create="createNewConversation"
      @select="selectConversation"
      @delete="deleteConversation"
      @update:search-query="searchQuery = $event"
    />

    <div class="flex flex-col flex-1">
      <ChatHeader
        :model-name="llmSettings.model"
        :dataset-name="datasetName"
        :dataset-options="datasetOptions"
        :mcp-enabled="mcpEnabled"
        @select-dataset="handleSelectDataset"
      />

      <ChatMessageArea ref="messageAreaRef" :messages="messages" :is-loading="isLoading" />

      <ChatInputArea
        :show-settings="showSettings"
        :show-model-select="showModelSelect"
        :model-options="modelOptions"
        :rag-settings="ragSettings"
        :llm-settings="llmSettings"
        :mcp-enabled="mcpEnabled"
        :is-loading="isLoading"
        :input-message="inputMessage"
        @send="sendMessage"
        @keydown="handleKeydown"
        @update:show-settings="showSettings = $event"
        @update:show-model-select="showModelSelect = $event"
        @update:mcp-enabled="mcpEnabled = $event"
        @update:input-message="inputMessage = $event"
        @update:rag-settings="ragSettings = $event"
        @update:llm-settings="llmSettings = $event"
        @select-model="handleSelectModel"
      />
    </div>
  </div>
</template>
