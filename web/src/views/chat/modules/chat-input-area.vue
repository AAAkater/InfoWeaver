<script setup lang="ts">
import { NButton, NIcon, NInput, NPopover, NInputNumber, NSwitch } from "naive-ui"
import {
  Connector24Filled as ConnectorIcon,
  Send24Filled as SendIcon,
  Settings24Regular as SettingsIcon,
  Sparkle24Filled as SparkleIcon,
} from "@vicons/fluent"

interface ModelOption {
  label: string
  value: string
  provider: string
  providerId: number
  providerMode: string
}

interface RagSettings {
  topK: number
  similarityThreshold: number
  enableRerank: boolean
}

interface LlmSettings {
  model: string
  temperature: number
  topP: number
  maxTokens: number
  frequencyPenalty: number
  presencePenalty: number
}

const props = defineProps<{
  showSettings: boolean
  showModelSelect: boolean
  modelOptions: ModelOption[]
  ragSettings: RagSettings
  llmSettings: LlmSettings
  mcpEnabled: boolean
  isLoading: boolean
  inputMessage: string
}>()

const emit = defineEmits<{
  (e: "send"): void
  (e: "update:showSettings", v: boolean): void
  (e: "update:showModelSelect", v: boolean): void
  (e: "update:mcpEnabled", v: boolean): void
  (e: "update:inputMessage", v: string): void
  (e: "update:ragSettings", v: RagSettings): void
  (e: "update:llmSettings", v: LlmSettings): void
  (e: "selectModel", model: string): void
  (e: "keydown", ev: KeyboardEvent): void
}>()

function toggleSettings() {
  emit("update:showSettings", !props.showSettings)
}

function toggleMcp() {
  emit("update:mcpEnabled", !props.mcpEnabled)
}

function onSelectModel(model: string) {
  emit("selectModel", model)
  emit("update:showModelSelect", false)
}
</script>

<template>
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
              :value="ragSettings.topK"
              :min="1"
              :max="20"
              size="tiny"
              style="width: 100px"
              @update:value="(v) => emit('update:ragSettings', { ...ragSettings, topK: v ?? 3 })"
            />
          </div>
          <div class="flex items-center justify-between">
            <span class="text-xs text-gray-500">相似度阈值</span>
            <NInputNumber
              :value="ragSettings.similarityThreshold"
              :min="0"
              :max="1"
              :step="0.05"
              size="tiny"
              style="width: 100px"
              @update:value="
                (v) => emit('update:ragSettings', { ...ragSettings, similarityThreshold: v ?? 0.7 })
              "
            />
          </div>
          <div class="flex items-center justify-between">
            <span class="text-xs text-gray-500">重排序</span>
            <NSwitch
              :value="ragSettings.enableRerank"
              size="small"
              @update:value="(v) => emit('update:ragSettings', { ...ragSettings, enableRerank: v })"
            />
          </div>
        </div>
      </div>

      <!-- LLM 采样参数 -->
      <div>
        <div class="mb-2 text-xs text-gray-500 font-bold">LLM 采样参数</div>
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <span class="text-xs text-gray-500">Temperature</span>
            <NInputNumber
              :value="llmSettings.temperature"
              :min="0"
              :max="2"
              :step="0.05"
              size="tiny"
              style="width: 100px"
              @update:value="
                (v) => emit('update:llmSettings', { ...llmSettings, temperature: v ?? 0.7 })
              "
            />
          </div>
          <div class="flex items-center justify-between">
            <span class="text-xs text-gray-500">Top P</span>
            <NInputNumber
              :value="llmSettings.topP"
              :min="0"
              :max="1"
              :step="0.05"
              size="tiny"
              style="width: 100px"
              @update:value="(v) => emit('update:llmSettings', { ...llmSettings, topP: v ?? 0.9 })"
            />
          </div>
          <div class="flex items-center justify-between">
            <span class="text-xs text-gray-500">Max Tokens</span>
            <NInputNumber
              :value="llmSettings.maxTokens"
              :min="1"
              :max="32768"
              :step="256"
              size="tiny"
              style="width: 100px"
              @update:value="
                (v) => emit('update:llmSettings', { ...llmSettings, maxTokens: v ?? 2048 })
              "
            />
          </div>
          <div class="flex items-center justify-between">
            <span class="text-xs text-gray-500">频率惩罚</span>
            <NInputNumber
              :value="llmSettings.frequencyPenalty"
              :min="-2"
              :max="2"
              :step="0.1"
              size="tiny"
              style="width: 100px"
              @update:value="
                (v) => emit('update:llmSettings', { ...llmSettings, frequencyPenalty: v ?? 0 })
              "
            />
          </div>
          <div class="flex items-center justify-between">
            <span class="text-xs text-gray-500">存在惩罚</span>
            <NInputNumber
              :value="llmSettings.presencePenalty"
              :min="-2"
              :max="2"
              :step="0.1"
              size="tiny"
              style="width: 100px"
              @update:value="
                (v) => emit('update:llmSettings', { ...llmSettings, presencePenalty: v ?? 0 })
              "
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 输入行 -->
    <div class="flex items-end space-x-2">
      <!-- 设置按钮 -->
      <NButton
        size="small"
        :type="showSettings ? 'primary' : 'default'"
        quaternary
        @click="toggleSettings"
      >
        <template #icon>
          <NIcon :component="SettingsIcon" size="18" />
        </template>
      </NButton>

      <!-- 模型选择 -->
      <NPopover
        :show="showModelSelect"
        trigger="click"
        placement="top-start"
        @clickoutside="emit('update:showModelSelect', false)"
      >
        <template #trigger>
          <NButton
            size="small"
            quaternary
            @click="emit('update:showModelSelect', !showModelSelect)"
          >
            <template #icon>
              <NIcon :component="SparkleIcon" size="18" />
            </template>
          </NButton>
        </template>
        <div class="w-56 py-1 max-h-60 overflow-y-auto">
          <div v-if="modelOptions.length === 0" class="px-3 py-2 text-xs text-gray-400">
            暂无可用模型
          </div>
          <div
            v-for="m in modelOptions"
            :key="m.value"
            class="cursor-pointer rounded px-3 py-2 text-sm transition-colors hover:bg-blue-50"
            :class="
              llmSettings.model === m.value
                ? 'bg-blue-50 text-blue-600 font-medium'
                : 'text-gray-700'
            "
            @click="onSelectModel(m.value)"
          >
            <div>{{ m.label }}</div>
            <div class="text-xs text-gray-400">{{ m.provider }}</div>
          </div>
        </div>
      </NPopover>

      <!-- MCP 开关 -->
      <NButton
        size="small"
        :type="mcpEnabled ? 'primary' : 'default'"
        quaternary
        @click="toggleMcp"
      >
        <template #icon>
          <NIcon :component="ConnectorIcon" size="18" />
        </template>
      </NButton>

      <!-- 输入框 -->
      <NInput
        :value="inputMessage"
        type="textarea"
        placeholder="输入消息... (Enter发送, Shift+Enter换行)"
        :autosize="{ minRows: 1, maxRows: 4 }"
        :disabled="isLoading"
        @update:value="(v) => emit('update:inputMessage', v)"
        @keydown="(e: KeyboardEvent) => emit('keydown', e)"
      />

      <!-- 发送 -->
      <NButton
        type="primary"
        :disabled="!inputMessage.trim() || isLoading"
        :loading="isLoading"
        @click="emit('send')"
      >
        <template #icon>
          <NIcon :component="SendIcon" />
        </template>
        发送
      </NButton>
    </div>
  </div>
</template>
