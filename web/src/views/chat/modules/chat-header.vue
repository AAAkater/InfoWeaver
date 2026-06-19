<script setup lang="ts">
import { computed, ref } from "vue"
import { NIcon, NTag } from "naive-ui"
import { Database20Regular as DatabaseIcon, Sparkle24Filled as SparkleIcon } from "@vicons/fluent"

const props = defineProps<{
  modelName: string
  datasetName: string
  datasetOptions: { label: string; value: number }[]
  mcpEnabled: boolean
}>()

const emit = defineEmits<{
  (e: "selectDataset", dataset: { label: string; value: number }): void
}>()

const showDatasetDropdown = ref(false)

function toggleDatasetDropdown() {
  showDatasetDropdown.value = !showDatasetDropdown.value
}

function onSelectDataset(dataset: { label: string; value: number }) {
  emit("selectDataset", dataset)
  showDatasetDropdown.value = false
}
</script>

<template>
  <div class="flex items-center justify-between border-b border-gray-100 px-4 py-2">
    <div class="flex items-center gap-3 text-xs">
      <div class="flex items-center gap-1 text-gray-500">
        <NIcon :component="SparkleIcon" size="14" />
        <span>{{ modelName || "未选择模型" }}</span>
      </div>
      <span class="text-gray-300">|</span>
      <div class="relative">
        <button
          class="flex items-center gap-1 text-gray-400 hover:text-gray-600 transition-colors"
          @click="toggleDatasetDropdown"
        >
          <NIcon :component="DatabaseIcon" size="13" />
          <span>{{ datasetName || "未选择知识库" }}</span>
          <span class="text-xs ml-0.5">▾</span>
        </button>
        <!-- Dropdown -->
        <div
          v-if="showDatasetDropdown"
          class="absolute left-0 top-full mt-1 z-50 min-w-[160px] rounded-lg border border-gray-200 bg-white shadow-lg py-1"
        >
          <button
            v-for="ds in datasetOptions"
            :key="ds.value"
            class="w-full text-left px-3 py-1.5 text-xs text-gray-700 hover:bg-gray-100 transition-colors"
            :class="{ 'bg-blue-50 text-blue-600 font-medium': ds.label === datasetName }"
            @click="onSelectDataset(ds)"
          >
            {{ ds.label }}
          </button>
          <div v-if="datasetOptions.length === 0" class="px-3 py-1.5 text-xs text-gray-400">
            暂无知识库
          </div>
        </div>
        <!-- Backdrop to close dropdown -->
        <div
          v-if="showDatasetDropdown"
          class="fixed inset-0 z-40"
          @click="showDatasetDropdown = false"
        />
      </div>
    </div>
    <NTag v-if="mcpEnabled" type="success" size="small" :bordered="false">MCP</NTag>
  </div>
</template>
