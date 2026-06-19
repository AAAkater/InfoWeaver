<script setup lang="ts">
import { NButton, NIcon, NInput } from "naive-ui"
import { ChatMultiple24Filled as ChatIcon, Delete24Regular as DeleteIcon } from "@vicons/fluent"

defineProps<{
  conversations: Array<{
    id: number
    title: string
    messages: Array<{ id: string }>
  }>
  activeConversationId?: number
  searchQuery: string
}>()

const emit = defineEmits<{
  (e: "create"): void
  (e: "select", id: number): void
  (e: "delete", id: number): void
  (e: "update:searchQuery", value: string): void
}>()
</script>

<template>
  <div class="w-64 flex flex-col border-r border-gray-200 bg-gray-50">
    <div class="border-b border-gray-200 p-4">
      <NButton type="primary" block @click="emit('create')">
        <template #icon>
          <NIcon :component="ChatIcon" />
        </template>
        新对话
      </NButton>
      <div class="mt-3">
        <NInput
          :value="searchQuery"
          size="small"
          placeholder="搜索会话..."
          clearable
          @update:value="(v) => emit('update:searchQuery', v)"
        />
      </div>
    </div>
    <div class="flex-1 overflow-y-auto p-2">
      <div
        v-for="conv in conversations"
        :key="conv.id"
        class="mb-2 cursor-pointer rounded-lg p-3 transition-colors"
        :class="
          activeConversationId === conv.id
            ? 'bg-blue-100 text-blue-700'
            : 'bg-white hover:bg-gray-100'
        "
        @click="emit('select', conv.id)"
      >
        <div class="flex items-center justify-between">
          <div class="flex-1 truncate text-sm font-medium">{{ conv.title }}</div>
          <NButton size="tiny" quaternary circle @click.stop="emit('delete', conv.id)">
            <template #icon>
              <NIcon :component="DeleteIcon" size="14" />
            </template>
          </NButton>
        </div>
        <div class="mt-1 text-xs text-gray-500">{{ conv.messages.length }} 条消息</div>
      </div>
      <div v-if="conversations.length === 0" class="py-8 text-center text-xs text-gray-400">
        未找到匹配的会话
      </div>
    </div>
  </div>
</template>
