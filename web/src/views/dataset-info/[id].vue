<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue"
import { useRouter } from "vue-router"
import { NButton, NButtonGroup, NCard, NSpace, useMessage } from "naive-ui"
import type { DataTableRowKey } from "naive-ui"
import {
  getDatasetChunks,
  getDatasetFiles,
  splitDocument,
  uploadFiles,
  embedChunks,
} from "@/service/api/dataset"
import FileTable from "./modules/file-table.vue"
import ChunkTable from "./modules/chunk-table.vue"
import SplitConfigModal from "./modules/split-config-modal.vue"

interface Props {
  id: string
}

const props = defineProps<Props>()

const router = useRouter()
const message = useMessage()

type TableTab = "files" | "chunks"
const activeTab = ref<TableTab>("files")

// File table state
const fileLoading = ref(false)
const files = ref<Api.Dataset.SimpleFileInfo[]>([])
const fileTotal = ref(0)
const filePage = ref(1)
const filePageSize = ref(20)

// Chunk table state
const chunkLoading = ref(false)
const chunks = ref<Api.Dataset.ChunkInfo[]>([])
const chunkTotal = ref(0)
const chunkPage = ref(1)
const chunkPageSize = ref(20)

// File upload
const fileInputRef = ref<HTMLInputElement | null>(null)
const isUploading = ref(false)

// Split config
const isProcessing = ref(false)
const chunkSize = ref(512)
const chunkOverlap = ref(50)
const showSplitModal = ref(false)
const splitTargetFile = ref<Api.Dataset.SimpleFileInfo | null>(null)

const datasetId = computed(() => Number(props.id))

const anyLoading = computed(() => fileLoading.value || chunkLoading.value)

// ---- Data fetching ----

async function fetchFiles() {
  if (!Number.isFinite(datasetId.value)) {
    message.error("无效的数据集 ID")
    files.value = []
    fileTotal.value = 0
    return
  }
  fileLoading.value = true
  try {
    const { data: fileData } = await getDatasetFiles(
      datasetId.value,
      filePage.value,
      filePageSize.value,
    )
    if (fileData) {
      const data = fileData as unknown as Api.Dataset.FileListResp
      files.value = data.files ?? []
      fileTotal.value = data.total ?? 0
    } else {
      message.error("获取文件列表失败")
      files.value = []
      fileTotal.value = 0
    }
  } catch {
    message.error("获取文件列表失败")
    files.value = []
    fileTotal.value = 0
  } finally {
    fileLoading.value = false
  }
}

async function fetchChunks() {
  if (!Number.isFinite(datasetId.value)) {
    message.error("无效的数据集 ID")
    chunks.value = []
    chunkTotal.value = 0
    return
  }
  chunkLoading.value = true
  try {
    const { data: chunkData } = await getDatasetChunks(
      datasetId.value,
      chunkPage.value,
      chunkPageSize.value,
    )
    if (chunkData) {
      const data = chunkData as unknown as Api.Dataset.DatasetChunkListResp
      chunks.value = data.chunks ?? []
      chunkTotal.value = data.total ?? 0
    } else {
      chunks.value = []
      chunkTotal.value = 0
    }
  } catch {
    message.error("获取 Chunk 列表失败")
    chunks.value = []
    chunkTotal.value = 0
  } finally {
    chunkLoading.value = false
  }
}

function refresh() {
  if (activeTab.value === "files") fetchFiles()
  else fetchChunks()
}

function onTabChange(tab: TableTab) {
  activeTab.value = tab
  if (tab === "files" && files.value.length === 0) fetchFiles()
  else if (tab === "chunks" && chunks.value.length === 0) fetchChunks()
}

function onFilePageChange(page: number) {
  filePage.value = page
  fetchFiles()
}
function onFilePageSizeChange(size: number) {
  filePageSize.value = size
  filePage.value = 1
  fetchFiles()
}
function onChunkPageChange(page: number) {
  chunkPage.value = page
  fetchChunks()
}
function onChunkPageSizeChange(size: number) {
  chunkPageSize.value = size
  chunkPage.value = 1
  fetchChunks()
}

// ---- File upload ----

function triggerFileUpload() {
  fileInputRef.value?.click()
}

function getErrorMessage(error: unknown, fallback: string) {
  return (error as { response?: { data?: { msg?: string } } })?.response?.data?.msg || fallback
}

async function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const selectedFiles = input.files
  if (!selectedFiles || selectedFiles.length === 0) return
  if (!Number.isFinite(datasetId.value)) {
    message.error("无效的数据集 ID")
    input.value = ""
    return
  }

  isUploading.value = true
  try {
    const { data, error } = await uploadFiles(datasetId.value, Array.from(selectedFiles))
    if (!error && data) {
      message.success(`成功上传 ${selectedFiles.length} 个文件`)
      await fetchFiles()
    } else {
      message.error(getErrorMessage(error, "文件上传失败"))
    }
  } finally {
    input.value = ""
    isUploading.value = false
  }
}

// ---- Split ----

function onOpenSplit(file: Api.Dataset.SimpleFileInfo) {
  splitTargetFile.value = file
  showSplitModal.value = true
}

async function onConfirmSplit(chunkSizeVal: number, chunkOverlapVal: number) {
  const file = splitTargetFile.value
  if (!file || !Number.isFinite(datasetId.value)) {
    message.error("无效的数据")
    return
  }

  showSplitModal.value = false
  isProcessing.value = true
  try {
    const splitPayload: Api.Dataset.SplitDocReq = {
      file_id: file.id,
      dataset_id: datasetId.value,
      minio_path: `documents/${file.name}`,
      chunk_size: chunkSizeVal,
      chunk_overlap: chunkOverlapVal,
    }
    message.info("正在拆分文档...")
    const { data: splitData, error: splitError } = await splitDocument(splitPayload)
    if (splitError || !splitData) {
      message.error(getErrorMessage(splitError, "文档拆分失败"))
      return
    }
    message.success(`文档拆分完成，生成 ${splitData.chunks_count ?? 0} 个片段`)
    activeTab.value = "chunks"
    await fetchChunks()
  } catch {
    message.error("处理文件时发生错误")
  } finally {
    isProcessing.value = false
    splitTargetFile.value = null
  }
}

// ---- Embed ----

const embedLoading = ref(false)
const checkedChunkKeys = ref<DataTableRowKey[]>([])

async function onEmbedChunks() {
  if (checkedChunkKeys.value.length === 0) {
    message.warning("请先选择要向量化的分块")
    return
  }
  embedLoading.value = true
  try {
    const { data: embedData, error: embedError } = await embedChunks({
      chunk_ids: checkedChunkKeys.value.map(Number),
    })
    if (embedError || !embedData) {
      message.error(getErrorMessage(embedError, "向量化失败"))
      return
    }
    message.success(`向量化完成，处理 ${embedData.chunks_count ?? 0} 个分块`)
    checkedChunkKeys.value = []
    await fetchChunks()
  } catch {
    message.error("向量化处理时发生错误")
  } finally {
    embedLoading.value = false
  }
}

function onDeleteFile(_file: Api.Dataset.SimpleFileInfo) {
  message.success("删除成功")
}

// ---- Lifecycle ----

watch(
  () => props.id,
  () => {
    filePage.value = 1
    chunkPage.value = 1
    if (activeTab.value === "files") fetchFiles()
    else fetchChunks()
  },
)

onMounted(() => {
  fetchFiles()
})
</script>

<template>
  <NSpace vertical :size="16">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <div class="text-18px font-600">数据集详情</div>
        <div class="mt-4px text-12px text-gray-500">Dataset ID: {{ props.id }}</div>
      </div>
      <NSpace>
        <NButton
          type="primary"
          :loading="isUploading"
          :disabled="isProcessing"
          @click="triggerFileUpload"
        >
          {{ isUploading ? "上传中..." : "上传文件" }}
        </NButton>
        <NButton :loading="anyLoading" @click="refresh">刷新</NButton>
        <NButton @click="router.back()">返回</NButton>
      </NSpace>
    </div>

    <input ref="fileInputRef" type="file" multiple class="hidden" @change="handleFileChange" />

    <!-- Split Config Modal -->
    <SplitConfigModal
      v-model:chunk-size="chunkSize"
      v-model:chunk-overlap="chunkOverlap"
      :show="showSplitModal"
      :loading="isProcessing"
      @confirm="onConfirmSplit"
      @cancel="showSplitModal = false"
    />

    <!-- Table Card -->
    <NCard :bordered="false" size="small">
      <template #header>
        <NButtonGroup size="small">
          <NButton
            :type="activeTab === 'files' ? 'primary' : 'default'"
            @click="onTabChange('files')"
          >
            文件列表
          </NButton>
          <NButton
            :type="activeTab === 'chunks' ? 'primary' : 'default'"
            @click="onTabChange('chunks')"
          >
            分块列表
          </NButton>
        </NButtonGroup>
      </template>

      <FileTable
        v-if="activeTab === 'files'"
        :files="files"
        :loading="fileLoading"
        :total="fileTotal"
        :page="filePage"
        :page-size="filePageSize"
        @update-page="onFilePageChange"
        @update-page-size="onFilePageSizeChange"
        @split="onOpenSplit"
        @delete="onDeleteFile"
      />

      <ChunkTable
        v-else
        :chunks="chunks"
        :loading="chunkLoading"
        :total="chunkTotal"
        :page="chunkPage"
        :page-size="chunkPageSize"
        :embed-loading="embedLoading"
        :checked-row-keys="checkedChunkKeys"
        @update-page="onChunkPageChange"
        @update-page-size="onChunkPageSizeChange"
        @update-checked-row-keys="checkedChunkKeys = $event"
        @embed="onEmbedChunks"
      />
    </NCard>
  </NSpace>
</template>
