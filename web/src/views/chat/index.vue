<script setup lang="ts">
import { computed, h, nextTick, onMounted, ref } from 'vue';
import type { Component } from 'vue';
import { NIcon } from 'naive-ui';
import {
  ChatMultiple24Filled as ChatIcon,
  Delete24Regular as DeleteIcon,
  Send24Filled as SendIcon
} from '@vicons/fluent';

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
  // RAG相关字段
  sources?: RetrievedSource[];
}

interface RetrievedSource {
  id: number;
  title: string;
  content: string;
  score: number;
  datasetName: string;
}

interface Conversation {
  id: string;
  title: string;
  messages: Message[];
  createdAt: Date;
}

const conversations = ref<Conversation[]>([]);
const activeConversationId = ref<string>();
const inputMessage = ref('');
const isLoading = ref(false);
const messageContainerRef = ref<HTMLElement>();

const activeConversation = computed(() => {
  return conversations.value.find(c => c.id === activeConversationId.value);
});

const messages = computed(() => {
  return activeConversation?.value?.messages || [];
});

function generateId() {
  return Math.random().toString(36).substring(2, 15);
}

function createNewConversation() {
  const newConv: Conversation = {
    id: generateId(),
    title: '新对话',
    messages: [],
    createdAt: new Date()
  };
  conversations.value.unshift(newConv);
  activeConversationId.value = newConv.id;
}

function deleteConversation(id: string) {
  const index = conversations.value.findIndex(c => c.id === id);
  if (index !== -1) {
    conversations.value.splice(index, 1);
    if (activeConversationId.value === id) {
      activeConversationId.value = conversations.value[0]?.id;
    }
  }
}

function selectConversation(id: string) {
  activeConversationId.value = id;
}

// Mock知识库数据
const mockKnowledgeBase: RetrievedSource[] = [
  {
    id: 1,
    title: 'Vue 3 组合式API指南',
    content:
      'Vue 3 的组合式API是一种新的编写组件逻辑的方式。它通过setup函数来组织代码，使用ref和reactive来创建响应式数据。组合式API提供了更好的代码组织能力和逻辑复用能力。',
    score: 0.95,
    datasetName: '前端技术文档'
  },
  {
    id: 2,
    title: 'TypeScript 类型系统',
    content:
      'TypeScript是JavaScript的超集，添加了静态类型检查。主要类型包括：基础类型（string, number, boolean）、对象类型、数组类型、元组类型、枚举类型等。类型推断和类型断言是两个重要的概念。',
    score: 0.88,
    datasetName: '前端技术文档'
  },
  {
    id: 3,
    title: 'Naive UI 组件库介绍',
    content:
      'Naive UI是一个Vue 3组件库，提供了丰富的UI组件。它支持主题定制、国际化、暗黑模式等功能。常用组件包括：NButton、NInput、NCard、NTable、NModal等。',
    score: 0.82,
    datasetName: '前端技术文档'
  },
  {
    id: 4,
    title: 'RAG技术原理',
    content:
      'RAG（检索增强生成）是一种结合信息检索和文本生成的技术。它首先从知识库中检索相关文档，然后将检索结果作为上下文，辅助大语言模型生成更准确、更有依据的回答。RAG可以有效减少模型的幻觉问题。',
    score: 0.91,
    datasetName: 'AI技术文档'
  },
  {
    id: 5,
    title: '向量数据库原理',
    content:
      '向量数据库用于存储和检索高维向量。在RAG系统中，文档会被转换为向量表示，然后通过向量相似度搜索来找到最相关的文档片段。常用的向量数据库包括：Pinecone、Milvus、Weaviate等。',
    score: 0.85,
    datasetName: 'AI技术文档'
  }
];

// Mock API - 模拟RAG检索
function mockRetrieveDocuments(query: string): RetrievedSource[] {
  // 根据关键词匹配模拟检索
  const keywords = query.toLowerCase();
  const relevantDocs = mockKnowledgeBase.filter(doc => {
    const docText = doc.title.toLowerCase() + doc.content.toLowerCase();
    return (
      docText.includes(keywords.split(' ')[0]) ||
      keywords.includes('vue') ||
      keywords.includes('typescript') ||
      keywords.includes('rag') ||
      keywords.includes('向量') ||
      keywords.includes('组件')
    );
  });

  // 返回最相关的2-3个文档
  return relevantDocs.slice(0, 3).map(doc => ({
    ...doc,
    score: doc.score - Math.random() * 0.1 // 添加一些随机性
  }));
}

// Mock API - 模拟RAG回复
async function mockSendMessage(content: string): Promise<{ response: string; sources: RetrievedSource[] }> {
  // 模拟网络延迟
  await new Promise<void>(resolve => {
    setTimeout(resolve, 1500 + Math.random() * 1000);
  });

  // 模拟RAG检索
  const sources = mockRetrieveDocuments(content);

  // 基于检索结果生成回答
  let response = '';

  if (sources.length > 0) {
    response = `根据知识库检索结果，我为您找到了以下相关信息：\n\n`;

    // 模拟基于检索内容的回答
    if (content.toLowerCase().includes('vue') || content.toLowerCase().includes('组件')) {
      response += `关于Vue相关的问题，根据文档显示：\nVue 3 的组合式API提供了更好的代码组织能力。Naive UI是一个优秀的Vue 3组件库，提供了丰富的UI组件如NButton、NInput等。\n\n您可以参考上述检索到的文档获取更详细的信息。`;
    } else if (content.toLowerCase().includes('typescript') || content.toLowerCase().includes('类型')) {
      response += `关于TypeScript的问题，根据文档显示：\nTypeScript是JavaScript的超集，添加了静态类型检查。主要类型包括基础类型、对象类型、数组类型等。类型推断和类型断言是两个重要的概念。\n\n您可以参考上述检索到的文档获取更详细的信息。`;
    } else if (
      content.toLowerCase().includes('rag') ||
      content.toLowerCase().includes('检索') ||
      content.toLowerCase().includes('向量')
    ) {
      response += `关于RAG技术的问题，根据文档显示：\nRAG（检索增强生成）是一种结合信息检索和文本生成的技术。它首先从知识库中检索相关文档，然后将检索结果作为上下文辅助大语言模型生成回答。向量数据库用于存储和检索高维向量。\n\n您可以参考上述检索到的文档获取更详细的信息。`;
    } else {
      response += `根据检索到的文档内容，我为您整理了相关信息。\n\n您的问题涉及到了${sources.map(s => s.title).join('、')}等相关内容。\n\n您可以参考上述检索到的文档获取更详细的信息。`;
    }
  } else {
    response = `抱歉，我在知识库中没有找到与您问题直接相关的文档。\n\n您可以尝试：\n1. 使用更具体的关键词\n2. 检查知识库中是否有相关内容\n3. 或者直接描述您想了解的具体内容\n\n我会尽力为您解答！`;
  }

  return { response, sources };
}

async function sendMessage() {
  if (!inputMessage.value.trim() || isLoading.value) return;

  if (!activeConversationId.value) {
    createNewConversation();
  }

  const userMessage: Message = {
    id: generateId(),
    role: 'user',
    content: inputMessage.value.trim(),
    timestamp: new Date()
  };

  activeConversation!.value!.messages.push(userMessage);

  // 更新对话标题（使用第一条消息）
  if (activeConversation!.value!.messages.length === 1) {
    const title = userMessage.content.slice(0, 20);
    activeConversation!.value!.title = userMessage.content.length > 20 ? `${title}...` : title;
  }

  const currentInput = inputMessage.value;
  inputMessage.value = '';
  isLoading.value = true;

  try {
    const { response, sources } = await mockSendMessage(currentInput);

    const assistantMessage: Message = {
      id: generateId(),
      role: 'assistant',
      content: response,
      timestamp: new Date(),
      sources
    };

    activeConversation!.value!.messages.push(assistantMessage);
  } catch {
    window.$message?.error('发送消息失败');
  } finally {
    isLoading.value = false;
    scrollToBottom();
  }
}

function scrollToBottom() {
  nextTick(() => {
    if (messageContainerRef.value) {
      messageContainerRef.value.scrollTop = messageContainerRef.value.scrollHeight;
    }
  });
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault();
    sendMessage();
  }
}

function formatTime(date: Date) {
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

onMounted(() => {
  // 初始化时创建一个对话
  createNewConversation();
});
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
      </div>
      <div class="flex-1 overflow-y-auto p-2">
        <div
          v-for="conv in conversations"
          :key="conv.id"
          class="mb-2 cursor-pointer rounded-lg p-3 transition-colors"
          :class="activeConversationId === conv.id ? 'bg-blue-100 text-blue-700' : 'bg-white hover:bg-gray-100'"
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
                <div class="line-clamp-2 mt-1 text-gray-600">{{ source.content.slice(0, 100) }}...</div>
                <div class="mt-1 text-gray-400">来源: {{ source.datasetName }}</div>
              </div>
            </div>
            <!-- 消息内容 -->
            <div
              class="rounded-lg p-3"
              :class="msg.role === 'user' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-800'"
            >
              <div class="whitespace-pre-wrap text-sm">{{ msg.content }}</div>
              <div class="mt-1 text-xs" :class="msg.role === 'user' ? 'text-blue-200' : 'text-gray-500'">
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
        <div class="flex items-end space-x-2">
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
