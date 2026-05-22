<script setup lang="ts">
import { computed, h, nextTick, onMounted, ref } from "vue"
import type { Component } from "vue"
import { NIcon } from "naive-ui"
import {
  ChatMultiple24Filled as ChatIcon,
  Delete24Regular as DeleteIcon,
  Send24Filled as SendIcon,
} from "@vicons/fluent"

interface Message {
  id: string
  role: "user" | "assistant"
  content: string
  timestamp: Date
  // RAG相关字段
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
  id: string
  title: string
  messages: Message[]
  createdAt: Date
}

const conversations = ref<Conversation[]>([])
const activeConversationId = ref<string>()
const inputMessage = ref("")
const isLoading = ref(false)
const messageContainerRef = ref<HTMLElement>()
const searchQuery = ref("")
const showSettings = ref(false)

// Mock RAG 召回设置
const ragSettings = ref({
  topK: 3,
  similarityThreshold: 0.7,
  selectedDatasets: ["前端技术文档", "AI技术文档"],
  enableRerank: true,
})

// Mock LLM 采样参数
const llmSettings = ref({
  model: "gpt-4o-mini",
  temperature: 0.7,
  topP: 0.9,
  maxTokens: 2048,
  frequencyPenalty: 0,
  presencePenalty: 0,
})

// Mock 可用模型列表
const availableModels = [
  "gpt-4o",
  "gpt-4o-mini",
  "gpt-3.5-turbo",
  "claude-3-sonnet",
  "claude-3-haiku",
]

// Mock 可用知识库列表
const availableDatasets = ["前端技术文档", "AI技术文档", "项目开发规范", "API接口文档", "运维手册"]

function toggleDataset(ds: string) {
  const idx = ragSettings.value.selectedDatasets.indexOf(ds)
  if (idx === -1) {
    ragSettings.value.selectedDatasets.push(ds)
  } else {
    ragSettings.value.selectedDatasets.splice(idx, 1)
  }
}

const activeConversation = computed(() => {
  return conversations.value.find((c) => c.id === activeConversationId.value)
})

const messages = computed(() => {
  return activeConversation?.value?.messages || []
})

const filteredConversations = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return conversations.value
  return conversations.value.filter((c) => c.title.toLowerCase().includes(q))
})

function generateId() {
  return Math.random().toString(36).substring(2, 15)
}

function createNewConversation() {
  const newConv: Conversation = {
    id: generateId(),
    title: "新对话",
    messages: [],
    createdAt: new Date(),
  }
  conversations.value.unshift(newConv)
  activeConversationId.value = newConv.id
}

function deleteConversation(id: string) {
  const index = conversations.value.findIndex((c) => c.id === id)
  if (index !== -1) {
    conversations.value.splice(index, 1)
    if (activeConversationId.value === id) {
      activeConversationId.value = conversations.value[0]?.id
    }
  }
}

function selectConversation(id: string) {
  activeConversationId.value = id
}

// Mock知识库数据
const mockKnowledgeBase: RetrievedSource[] = [
  {
    id: 1,
    title: "Vue 3 组合式API指南",
    content:
      "Vue 3 的组合式API是一种新的编写组件逻辑的方式。它通过setup函数来组织代码，使用ref和reactive来创建响应式数据。组合式API提供了更好的代码组织能力和逻辑复用能力。",
    score: 0.95,
    datasetName: "前端技术文档",
  },
  {
    id: 2,
    title: "TypeScript 类型系统",
    content:
      "TypeScript是JavaScript的超集，添加了静态类型检查。主要类型包括：基础类型（string, number, boolean）、对象类型、数组类型、元组类型、枚举类型等。类型推断和类型断言是两个重要的概念。",
    score: 0.88,
    datasetName: "前端技术文档",
  },
  {
    id: 3,
    title: "Naive UI 组件库介绍",
    content:
      "Naive UI是一个Vue 3组件库，提供了丰富的UI组件。它支持主题定制、国际化、暗黑模式等功能。常用组件包括：NButton、NInput、NCard、NTable、NModal等。",
    score: 0.82,
    datasetName: "前端技术文档",
  },
  {
    id: 4,
    title: "RAG技术原理",
    content:
      "RAG（检索增强生成）是一种结合信息检索和文本生成的技术。它首先从知识库中检索相关文档，然后将检索结果作为上下文，辅助大语言模型生成更准确、更有依据的回答。RAG可以有效减少模型的幻觉问题。",
    score: 0.91,
    datasetName: "AI技术文档",
  },
  {
    id: 5,
    title: "向量数据库原理",
    content:
      "向量数据库用于存储和检索高维向量。在RAG系统中，文档会被转换为向量表示，然后通过向量相似度搜索来找到最相关的文档片段。常用的向量数据库包括：Pinecone、Milvus、Weaviate等。",
    score: 0.85,
    datasetName: "AI技术文档",
  },
]

// Mock API - 模拟RAG检索
function mockRetrieveDocuments(query: string): RetrievedSource[] {
  // 根据关键词匹配模拟检索
  const keywords = query.toLowerCase()
  const relevantDocs = mockKnowledgeBase.filter((doc) => {
    const docText = doc.title.toLowerCase() + doc.content.toLowerCase()
    return (
      docText.includes(keywords.split(" ")[0]) ||
      keywords.includes("vue") ||
      keywords.includes("typescript") ||
      keywords.includes("rag") ||
      keywords.includes("向量") ||
      keywords.includes("组件")
    )
  })

  // 返回最相关的2-3个文档
  return relevantDocs.slice(0, 3).map((doc) => ({
    ...doc,
    score: doc.score - Math.random() * 0.1, // 添加一些随机性
  }))
}

// Mock API - 模拟RAG回复
async function mockSendMessage(
  content: string,
): Promise<{ response: string; sources: RetrievedSource[] }> {
  // 模拟网络延迟
  await new Promise<void>((resolve) => {
    setTimeout(resolve, 1500 + Math.random() * 1000)
  })

  // 模拟RAG检索
  const sources = mockRetrieveDocuments(content)

  // 基于检索结果生成回答
  let response = ""

  if (sources.length > 0) {
    response = `根据知识库检索结果，我为您找到了以下相关信息：\n\n`

    // 模拟基于检索内容的回答
    if (content.toLowerCase().includes("vue") || content.toLowerCase().includes("组件")) {
      response += `关于Vue相关的问题，根据文档显示：\nVue 3 的组合式API提供了更好的代码组织能力。Naive UI是一个优秀的Vue 3组件库，提供了丰富的UI组件如NButton、NInput等。\n\n您可以参考上述检索到的文档获取更详细的信息。`
    } else if (
      content.toLowerCase().includes("typescript") ||
      content.toLowerCase().includes("类型")
    ) {
      response += `关于TypeScript的问题，根据文档显示：\nTypeScript是JavaScript的超集，添加了静态类型检查。主要类型包括基础类型、对象类型、数组类型等。类型推断和类型断言是两个重要的概念。\n\n您可以参考上述检索到的文档获取更详细的信息。`
    } else if (
      content.toLowerCase().includes("rag") ||
      content.toLowerCase().includes("检索") ||
      content.toLowerCase().includes("向量")
    ) {
      response += `关于RAG技术的问题，根据文档显示：\nRAG（检索增强生成）是一种结合信息检索和文本生成的技术。它首先从知识库中检索相关文档，然后将检索结果作为上下文辅助大语言模型生成回答。向量数据库用于存储和检索高维向量。\n\n您可以参考上述检索到的文档获取更详细的信息。`
    } else {
      response += `根据检索到的文档内容，我为您整理了相关信息。\n\n您的问题涉及到了${sources.map((s) => s.title).join("、")}等相关内容。\n\n您可以参考上述检索到的文档获取更详细的信息。`
    }
  } else {
    response = `抱歉，我在知识库中没有找到与您问题直接相关的文档。\n\n您可以尝试：\n1. 使用更具体的关键词\n2. 检查知识库中是否有相关内容\n3. 或者直接描述您想了解的具体内容\n\n我会尽力为您解答！`
  }

  return { response, sources }
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
    createNewConversation()
  }

  const userMessage: Message = {
    id: generateId(),
    role: "user",
    content: inputMessage.value.trim(),
    timestamp: new Date(),
  }

  activeConversation!.value!.messages.push(userMessage)

  // 更新对话标题（使用第一条消息）
  if (activeConversation!.value!.messages.length === 1) {
    const title = userMessage.content.slice(0, 20)
    activeConversation!.value!.title = userMessage.content.length > 20 ? `${title}...` : title
  }

  const currentInput = inputMessage.value
  inputMessage.value = ""
  isLoading.value = true

  try {
    const { response, sources } = await mockSendMessage(currentInput)

    const assistantMessage: Message = {
      id: generateId(),
      role: "assistant",
      content: response,
      timestamp: new Date(),
      sources,
    }

    activeConversation!.value!.messages.push(assistantMessage)
  } catch {
    window.$message?.error("发送消息失败")
  } finally {
    isLoading.value = false
    scrollToBottom()
  }
}

function scrollToBottom() {
  nextTick(() => {
    if (messageContainerRef.value) {
      messageContainerRef.value.scrollTop = messageContainerRef.value.scrollHeight
    }
  })
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault()
    sendMessage()
  }
}

function formatTime(date: Date) {
  return date.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" })
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

// Mock 历史会话数据
function createMockConversations() {
  const mockData: Conversation[] = [
    {
      id: "mock-1",
      title: "Vue3 组合式API 使用问题",
      messages: [
        {
          id: "m1-1",
          role: "user",
          content: "Vue3 的 ref 和 reactive 有什么区别？",
          timestamp: new Date(Date.now() - 86400000 * 2),
        },
        {
          id: "m1-2",
          role: "assistant",
          content:
            "ref 用于基本类型，reactive 用于对象类型。ref 需要通过 .value 访问，reactive 直接访问。",
          timestamp: new Date(Date.now() - 86400000 * 2 + 60000),
          sources: [mockKnowledgeBase[0]],
        },
      ],
      createdAt: new Date(Date.now() - 86400000 * 2),
    },
    {
      id: "mock-2",
      title: "RAG 检索增强生成原理",
      messages: [
        {
          id: "m2-1",
          role: "user",
          content: "请详细解释一下 RAG 技术的工作原理",
          timestamp: new Date(Date.now() - 86400000),
        },
        {
          id: "m2-2",
          role: "assistant",
          content:
            "RAG（检索增强生成）结合了信息检索和文本生成。首先从知识库中检索相关文档，然后将检索结果作为上下文，辅助 LLM 生成更准确的回答。",
          timestamp: new Date(Date.now() - 86400000 + 30000),
          sources: [mockKnowledgeBase[3]],
        },
        {
          id: "m2-3",
          role: "user",
          content: "RAG 相比传统微调有什么优势？",
          timestamp: new Date(Date.now() - 86400000 + 120000),
        },
        {
          id: "m2-4",
          role: "assistant",
          content:
            "RAG 可以实时更新知识库而无需重新训练模型，减少幻觉问题，并且答案可追溯至具体文档来源。",
          timestamp: new Date(Date.now() - 86400000 + 150000),
          sources: [mockKnowledgeBase[3], mockKnowledgeBase[4]],
        },
      ],
      createdAt: new Date(Date.now() - 86400000),
    },
    {
      id: "mock-3",
      title: "VectorDB 选型建议咨询",
      messages: [
        {
          id: "m3-1",
          role: "user",
          content: "我们项目想选一个向量数据库，有什么推荐吗？",
          timestamp: new Date(Date.now() - 43200000),
        },
        {
          id: "m3-2",
          role: "assistant",
          content:
            "常用的向量数据库包括：Pinecone（SaaS托管）、Milvus（开源高性能）、Weaviate（自带向量化）、Qdrant（Rust编写，性能优秀）。选择时需考虑数据量、延迟要求、是否自托管等因素。",
          timestamp: new Date(Date.now() - 43200000 + 45000),
          sources: [mockKnowledgeBase[4]],
        },
      ],
      createdAt: new Date(Date.now() - 43200000),
    },
    {
      id: "mock-4",
      title: "NaiveUI 表单校验方案",
      messages: [
        {
          id: "m4-1",
          role: "user",
          content: "NaiveUI 的表单校验怎么做？需要动态校验规则",
          timestamp: new Date(Date.now() - 36000000),
        },
        {
          id: "m4-2",
          role: "assistant",
          content:
            "NaiveUI 的 NForm 组件支持通过 rules 属性设置校验规则，可以使用 computed 动态生成规则。支持 required、min、max、pattern 等内置规则，也支持自定义 validator 函数。",
          timestamp: new Date(Date.now() - 36000000 + 30000),
          sources: [mockKnowledgeBase[2]],
        },
        {
          id: "m4-3",
          role: "user",
          content: "跨字段校验比如密码确认怎么实现？",
          timestamp: new Date(Date.now() - 36000000 + 120000),
        },
        {
          id: "m4-4",
          role: "assistant",
          content:
            "可以在 confirmPassword 字段的 validator 中通过 getFieldValue 获取 password 字段的值进行比对。或者使用 NForm 的 validate 回调中手动处理跨字段逻辑。",
          timestamp: new Date(Date.now() - 36000000 + 150000),
          sources: [mockKnowledgeBase[2]],
        },
      ],
      createdAt: new Date(Date.now() - 36000000),
    },
    {
      id: "mock-5",
      title: "TypeScript 高级类型体操",
      messages: [
        {
          id: "m5-1",
          role: "user",
          content: "TypeScript 中怎么实现一个 DeepReadonly 类型？",
          timestamp: new Date(Date.now() - 7200000),
        },
        {
          id: "m5-2",
          role: "assistant",
          content:
            "可以用递归条件类型实现：type DeepReadonly<T> = { readonly [K in keyof T]: T[K] extends object ? DeepReadonly<T[K]> : T[K] }。注意需要处理数组和函数等特殊情况。",
          timestamp: new Date(Date.now() - 7200000 + 20000),
          sources: [mockKnowledgeBase[1]],
        },
      ],
      createdAt: new Date(Date.now() - 7200000),
    },
    {
      id: "mock-6",
      title: "知识库文档分块策略探讨",
      messages: [
        {
          id: "m6-1",
          role: "user",
          content: "处理长文档时，chunk_size 和 chunk_overlap 怎么设置比较合理？",
          timestamp: new Date(Date.now() - 1800000),
        },
        {
          id: "m6-2",
          role: "assistant",
          content:
            "一般建议 chunk_size 512-1024 tokens，overlap 50-100 tokens。中文文档可以适当调大 chunk_size 到 800-1500。关键是要保证每个 chunk 的语义完整性，避免在句子中间截断。",
          timestamp: new Date(Date.now() - 1800000 + 25000),
          sources: [mockKnowledgeBase[3], mockKnowledgeBase[4]],
        },
      ],
      createdAt: new Date(Date.now() - 1800000),
    },
    {
      id: "mock-7",
      title: "前端性能优化实战经验",
      messages: [
        {
          id: "m7-1",
          role: "user",
          content: "页面加载很慢，有哪些前端性能优化手段？",
          timestamp: new Date(Date.now() - 600000),
        },
        {
          id: "m7-2",
          role: "assistant",
          content:
            "可以从几个方面优化：1. 代码分割和懒加载 2. 图片压缩和 WebP 格式 3. CDN 加速静态资源 4. 虚拟列表处理长列表 5. 防抖节流优化高频事件 6. 使用 Web Worker 处理密集计算。具体需要根据你的项目情况选择合适方案。",
          timestamp: new Date(Date.now() - 600000 + 30000),
        },
        {
          id: "m7-3",
          role: "user",
          content: "虚拟列表在 Vue3 中怎么实现？",
          timestamp: new Date(Date.now() - 600000 + 90000),
        },
        {
          id: "m7-4",
          role: "assistant",
          content:
            "可以使用 vue-virtual-scroller 或者基于 NaiveUI 的虚拟滚动组件。核心原理是只渲染可视区域内的 DOM 节点，通过计算 scrollTop 动态更新渲染范围。",
          timestamp: new Date(Date.now() - 600000 + 120000),
          sources: [mockKnowledgeBase[0], mockKnowledgeBase[2]],
        },
      ],
      createdAt: new Date(Date.now() - 600000),
    },
  ]

  conversations.value = mockData
  activeConversationId.value = mockData[0].id
}

onMounted(() => {
  createMockConversations()
})
</script>

<template>
  <div class="h-full flex">
    <!-- 左侧对话列表 -->
    <div class="w-64 flex flex-col border-r border-gray-200 bg-gray-50">
      <div class="border-b border-gray-200 p-4">
        <NButton type="primary" block @click="createNewConversation">
          <template #icon>
            <NIcon :component="ChatIcon" />
          </template>
          新对话
        </NButton>
        <div class="mt-3">
          <NInput v-model:value="searchQuery" size="small" placeholder="搜索会话..." clearable />
        </div>
      </div>
      <div class="flex-1 overflow-y-auto p-2">
        <div
          v-for="conv in filteredConversations"
          :key="conv.id"
          class="mb-2 cursor-pointer rounded-lg p-3 transition-colors"
          :class="
            activeConversationId === conv.id
              ? 'bg-blue-100 text-blue-700'
              : 'bg-white hover:bg-gray-100'
          "
          @click="selectConversation(conv.id)"
        >
          <div class="flex items-center justify-between">
            <div class="flex-1 truncate text-sm font-medium">{{ conv.title }}</div>
            <NButton size="tiny" quaternary circle @click.stop="deleteConversation(conv.id)">
              <template #icon>
                <NIcon :component="DeleteIcon" size="14" />
              </template>
            </NButton>
          </div>
          <div class="mt-1 text-xs text-gray-500">{{ conv.messages.length }} 条消息</div>
        </div>
        <div
          v-if="filteredConversations.length === 0"
          class="py-8 text-center text-xs text-gray-400"
        >
          未找到匹配的会话
        </div>
      </div>
    </div>

    <!-- 右侧对话内容 -->
    <div class="flex flex-col flex-1">
      <!-- 消息区域 -->
      <div ref="messageContainerRef" class="flex-1 overflow-y-auto p-4 space-y-4">
        <div v-if="messages.length === 0" class="h-full flex items-center justify-center">
          <div class="text-center text-gray-400">
            <NIcon :component="ChatIcon" size="48" class="mb-4" />
            <p>开始一个新的对话吧</p>
          </div>
        </div>
        <div
          v-for="msg in messages"
          :key="msg.id"
          class="flex"
          :class="msg.role === 'user' ? 'justify-end' : 'justify-start'"
        >
          <div class="max-w-[70%]">
            <!-- RAG检索来源 -->
            <div v-if="msg.sources && msg.sources.length > 0" class="mb-2 space-y-1">
              <div class="mb-1 text-xs text-gray-500">📚 检索来源：</div>
              <div
                v-for="source in msg.sources"
                :key="source.id"
                class="border border-amber-200 rounded bg-amber-50 p-2 text-xs"
              >
                <div class="flex items-center justify-between">
                  <span class="text-amber-700 font-medium">{{ source.title }}</span>
                  <span class="text-amber-500">相关度: {{ (source.score * 100).toFixed(0) }}%</span>
                </div>
                <div class="line-clamp-2 mt-1 text-gray-600">
                  {{ source.content.slice(0, 100) }}...
                </div>
                <div class="mt-1 text-gray-400">来源: {{ source.datasetName }}</div>
              </div>
            </div>
            <!-- 消息内容 -->
            <div
              class="rounded-lg p-3"
              :class="msg.role === 'user' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-800'"
            >
              <div class="whitespace-pre-wrap text-sm">{{ msg.content }}</div>
              <div
                class="mt-1 text-xs"
                :class="msg.role === 'user' ? 'text-blue-200' : 'text-gray-500'"
              >
                {{ formatTime(msg.timestamp) }}
              </div>
            </div>
          </div>
        </div>
        <div v-if="isLoading" class="flex justify-start">
          <div class="rounded-lg bg-gray-100 p-3">
            <div class="flex items-center space-x-2">
              <NSpin size="small" />
              <span class="text-sm text-gray-500">正在思考...</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="border-t border-gray-200 bg-white p-4">
        <!-- 设置面板 -->
        <div
          v-if="showSettings"
          class="mb-3 rounded-lg bg-gray-50 p-3"
          style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px"
        >
          <!-- RAG 召回设置 -->
          <div>
            <div class="mb-2 text-xs text-gray-500 font-bold">RAG 召回设置</div>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500">Top-K</span>
                <NInputNumber
                  v-model:value="ragSettings.topK"
                  :min="1"
                  :max="20"
                  size="tiny"
                  style="width: 100px"
                />
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500">相似度阈值</span>
                <NInputNumber
                  v-model:value="ragSettings.similarityThreshold"
                  :min="0"
                  :max="1"
                  :step="0.05"
                  size="tiny"
                  style="width: 100px"
                />
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500">重排序</span>
                <NSwitch v-model:value="ragSettings.enableRerank" size="small" />
              </div>
              <div>
                <div class="mb-1 text-xs text-gray-500">检索知识库</div>
                <div class="flex flex-wrap gap-1">
                  <NTag
                    v-for="ds in availableDatasets"
                    :key="ds"
                    size="small"
                    :type="ragSettings.selectedDatasets.includes(ds) ? 'primary' : 'default'"
                    :bordered="false"
                    :checkable="false"
                    class="cursor-pointer"
                    @click="toggleDataset(ds)"
                  >
                    {{ ds }}
                  </NTag>
                </div>
              </div>
            </div>
          </div>
          <!-- LLM 采样参数 -->
          <div>
            <div class="mb-2 text-xs text-gray-500 font-bold">LLM 采样参数</div>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500">模型</span>
                <NSelect
                  v-model:value="llmSettings.model"
                  :options="availableModels.map((m) => ({ label: m, value: m }))"
                  size="tiny"
                  style="width: 140px"
                />
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500">Temperature</span>
                <NInputNumber
                  v-model:value="llmSettings.temperature"
                  :min="0"
                  :max="2"
                  :step="0.05"
                  size="tiny"
                  style="width: 100px"
                />
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500">Top P</span>
                <NInputNumber
                  v-model:value="llmSettings.topP"
                  :min="0"
                  :max="1"
                  :step="0.05"
                  size="tiny"
                  style="width: 100px"
                />
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500">Max Tokens</span>
                <NInputNumber
                  v-model:value="llmSettings.maxTokens"
                  :min="1"
                  :max="32768"
                  :step="256"
                  size="tiny"
                  style="width: 100px"
                />
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500">频率惩罚</span>
                <NInputNumber
                  v-model:value="llmSettings.frequencyPenalty"
                  :min="-2"
                  :max="2"
                  :step="0.1"
                  size="tiny"
                  style="width: 100px"
                />
              </div>
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-500">存在惩罚</span>
                <NInputNumber
                  v-model:value="llmSettings.presencePenalty"
                  :min="-2"
                  :max="2"
                  :step="0.1"
                  size="tiny"
                  style="width: 100px"
                />
              </div>
            </div>
          </div>
        </div>

        <div class="flex items-end space-x-2">
          <NButton
            size="small"
            :type="showSettings ? 'primary' : 'default'"
            quaternary
            @click="showSettings = !showSettings"
          >
            <span class="text-14px" :class="showSettings ? 'i-mdi:cog' : 'i-mdi:cog-outline'" />
          </NButton>
          <NInput
            v-model:value="inputMessage"
            type="textarea"
            placeholder="输入消息... (Enter发送, Shift+Enter换行)"
            :autosize="{ minRows: 1, maxRows: 4 }"
            :disabled="isLoading"
            @keydown="handleKeydown"
          />
          <NButton
            type="primary"
            :disabled="!inputMessage.trim() || isLoading"
            :loading="isLoading"
            @click="sendMessage"
          >
            <template #icon>
              <NIcon :component="SendIcon" />
            </template>
            发送
          </NButton>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
