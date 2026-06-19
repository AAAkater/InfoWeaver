<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref } from "vue"
import { NIcon, NSpin } from "naive-ui"
import {
  ChatMultiple24Filled as ChatIcon,
  ArrowDown24Filled as ArrowDownIcon,
} from "@vicons/fluent"
import { renderMarkdown } from "@/hooks/markdown"
import "@/styles/css/markdown.css"

interface Message {
  id: string
  role: "user" | "assistant"
  content: string
  thinking: string
  timestamp: Date
  sources?: Array<{
    id: number
    title: string
    content: string
    score: number
    datasetName: string
  }>
}

const props = defineProps<{
  messages: Message[]
  isLoading: boolean
}>()

const containerRef = ref<HTMLElement>()
const expandedThinking = ref<Record<string, boolean>>({})
const showScrollButton = ref(false)

const hasMessages = computed(() => props.messages.length > 0)

const lastAssistantIndex = computed(() => {
  for (let i = props.messages.length - 1; i >= 0; i--) {
    if (props.messages[i].role === "assistant") return i
  }
  return -1
})

function isThinkingExpanded(msgId: string, index: number) {
  // Auto-expand thinking while streaming the last assistant message
  if (props.isLoading && index === lastAssistantIndex.value) return true
  return !!expandedThinking.value[msgId]
}

function scrollToBottom() {
  nextTick(() => {
    if (containerRef.value) {
      containerRef.value.scrollTop = containerRef.value.scrollHeight
    }
  })
}

function handleScroll() {
  if (!containerRef.value) return
  const { scrollTop, scrollHeight, clientHeight } = containerRef.value
  // Show button when scrolled up more than 100px from bottom
  showScrollButton.value = scrollHeight - scrollTop - clientHeight > 100
}

function toggleThinking(msgId: string) {
  expandedThinking.value[msgId] = !expandedThinking.value[msgId]
}

function formatTime(date: Date) {
  return date.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" })
}

onMounted(() => {
  containerRef.value?.addEventListener("scroll", handleScroll)
})

onUnmounted(() => {
  containerRef.value?.removeEventListener("scroll", handleScroll)
})

defineExpose({ scrollToBottom })
</script>

<template>
  <div ref="containerRef" class="flex-1 overflow-y-auto p-4 space-y-4 relative">
    <!-- 空状态 -->
    <div v-if="!hasMessages" class="h-full flex items-center justify-center">
      <div class="text-center text-gray-400">
        <NIcon :component="ChatIcon" size="48" class="mb-4" />
        <p>开始一个新的对话吧</p>
      </div>
    </div>

    <!-- 消息列表 -->
    <div
      v-for="(msg, idx) in messages"
      :key="msg.id"
      class="flex"
      :class="msg.role === 'user' ? 'justify-end' : 'justify-start'"
    >
      <div class="max-w-[70%]">
        <!-- RAG 检索来源 -->
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
          <!-- 思考过程（仅 assistant） -->
          <div v-if="msg.role === 'assistant' && msg.thinking" class="mb-2">
            <button
              class="flex items-center gap-1 text-xs text-purple-600 hover:text-purple-800 transition-colors"
              @click="toggleThinking(msg.id)"
            >
              <span>{{ isThinkingExpanded(msg.id, idx) ? "▾" : "▸" }}</span>
              <span>🧠 思考过程</span>
            </button>
            <div
              v-if="isThinkingExpanded(msg.id, idx)"
              class="mt-2 p-2 rounded bg-purple-50 border border-purple-200 text-xs text-purple-800 whitespace-pre-wrap max-h-48 overflow-y-auto"
            >
              {{ msg.thinking }}
            </div>
          </div>

          <div
            class="text-sm markdown-body"
            :class="msg.role === 'user' ? 'text-white' : 'text-gray-800'"
            v-html="renderMarkdown(msg.content)"
          />
          <div
            class="mt-1 text-xs"
            :class="msg.role === 'user' ? 'text-blue-200' : 'text-gray-500'"
          >
            {{ formatTime(msg.timestamp) }}
          </div>
        </div>
      </div>
    </div>

    <!-- 加载中 -->
    <div v-if="isLoading" class="flex justify-start">
      <div class="rounded-lg bg-gray-100 p-3">
        <div class="flex items-center space-x-2">
          <NSpin size="small" />
          <span class="text-sm text-gray-500">正在思考...</span>
        </div>
      </div>
    </div>

    <!-- 一键返回底部 -->
    <Transition name="scroll-fade">
      <button
        v-if="showScrollButton"
        class="sticky bottom-2 left-1/2 -translate-x-1/2 z-10 flex items-center justify-center w-9 h-9 rounded-full bg-white border border-gray-200 shadow-md hover:shadow-lg hover:bg-gray-50 transition-all"
        @click="scrollToBottom"
      >
        <NIcon :component="ArrowDownIcon" size="18" class="text-gray-500" />
      </button>
    </Transition>
  </div>
</template>
