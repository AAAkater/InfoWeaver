<script setup lang="ts">
import { computed, h, onMounted, ref, watch } from "vue"
import { useRouter } from "vue-router"
import type { DataTableColumns } from "naive-ui"
import { NButton, NButtonGroup, NInputNumber, NTag, useMessage } from "naive-ui"
import dayjs from "dayjs"
import {
  getDatasetChunks,
  getDatasetFiles,
  splitDocument,
  uploadFiles,
} from "@/service/api/dataset"

interface Props {
  id: string
}

const props = defineProps<Props>()

const router = useRouter()
const message = useMessage()

// Toggle between file table and chunk table
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
const uploadedFiles = ref<Api.Dataset.FileUploadInfo[]>([])
const chunkSize = ref(512)
const chunkOverlap = ref(50)
const isProcessing = ref(false)

const datasetId = computed(() => Number(props.id))

const scrollX = 80 + 360 + 100 + 100 + 180 + 180 + 170 + 170

const fileScrollX = 80 + 240 + 120 + 120 + 180 + 80

// ---- File table columns ----
const fileColumns: DataTableColumns<Api.Dataset.SimpleFileInfo> = [
  {
    title: "ID",
    key: "id",
    width: 80,
  },
  {
    title: "文件名",
    key: "name",
    minWidth: 240,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: "类型",
    key: "type",
    width: 120,
    render(row) {
      return row.type || "-"
    },
  },
  {
    title: "大小",
    key: "size",
    width: 120,
    render(row) {
      if (!row.size) return "-"
      if (row.size < 1024) return `${row.size} B`
      if (row.size < 1024 * 1024) return `${(row.size / 1024).toFixed(1)} KB`
      return `${(row.size / (1024 * 1024)).toFixed(2)} MB`
    },
  },
  {
    title: "创建时间",
    key: "createdAt",
    width: 180,
    render(row) {
      return row.createdAt ? dayjs(row.createdAt).format("YYYY-MM-DD HH:mm") : "-"
    },
  },
  {
    title: "操作",
    key: "actions",
    width: 80,
    render(row) {
      return h(
        NButton,
        { size: "tiny", type: "error", onClick: () => handleDeleteFile(row) },
        { default: () => "删除" },
      )
    },
  },
]

// ---- Chunk table columns ----
const chunkColumns: DataTableColumns<Api.Dataset.ChunkInfo> = [
  {
    title: "ID",
    key: "id",
    width: 80,
  },
  {
    title: "内容",
    key: "content",
    minWidth: 360,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: "状态",
    key: "status",
    width: 100,
    render(row) {
      const type = row.status === "completed" || row.status === "success" ? "success" : "default"

      return h(NTag, { size: "small", type }, { default: () => row.status || "-" })
    },
  },
  {
    title: "File ID",
    key: "file_id",
    width: 100,
  },
  {
    title: "Vector ID",
    key: "vector_id",
    minWidth: 180,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: "元数据",
    key: "chunk_metadata",
    minWidth: 180,
    ellipsis: {
      tooltip: true,
    },
    render(row) {
      return row.chunk_metadata ? JSON.stringify(row.chunk_metadata) : "-"
    },
  },
  {
    title: "创建时间",
    key: "created_at",
    width: 170,
    render(row) {
      return row.created_at ? dayjs(row.created_at).format("YYYY-MM-DD HH:mm") : "-"
    },
  },
  {
    title: "更新时间",
    key: "updated_at",
    width: 170,
    render(row) {
      return row.updated_at ? dayjs(row.updated_at).format("YYYY-MM-DD HH:mm") : "-"
    },
  },
]

function handleDeleteFile(_file: Api.Dataset.SimpleFileInfo) {
  message.success("删除成功")
}

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
    const { data: fileData, error: _fetchErr } = await getDatasetFiles(
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

function refreshCurrentTab() {
  if (activeTab.value === "files") {
    fetchFiles()
  } else {
    fetchChunks()
  }
}

function handleTabChange(tab: TableTab) {
  activeTab.value = tab
  if (tab === "files" && files.value.length === 0) {
    fetchFiles()
  } else if (tab === "chunks" && chunks.value.length === 0) {
    fetchChunks()
  }
}

function handleFilePageSizeChange(size: number) {
  filePageSize.value = size
  filePage.value = 1
  fetchFiles()
}

function handleChunkPageSizeChange(size: number) {
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

  if (!selectedFiles || selectedFiles.length === 0) {
    return
  }

  if (!Number.isFinite(datasetId.value)) {
    message.error("无效的数据集 ID")
    input.value = ""
    return
  }

  const fileList = Array.from(selectedFiles)
  isUploading.value = true

  try {
    const { data, error } = await uploadFiles(datasetId.value, fileList)

    if (!error && data?.code === 0) {
      uploadedFiles.value = data.data?.files ?? []
      message.success(`成功上传 ${fileList.length} 个文件`)
      // Refresh file list
      await fetchFiles()
    } else {
      message.error(getErrorMessage(error, "文件上传失败"))
    }
  } finally {
    input.value = ""
    isUploading.value = false
  }
}

async function handleSplitAndEmbed(file: Api.Dataset.FileUploadInfo) {
  if (!Number.isFinite(datasetId.value)) {
    message.error("无效的数据集 ID")
    return
  }

  isProcessing.value = true

  try {
    const splitPayload: Api.Dataset.SplitDocReq = {
      file_id: file.id,
      dataset_id: datasetId.value,
      minio_path: `documents/${file.name}`,
      chunk_size: chunkSize.value,
      chunk_overlap: chunkOverlap.value,
    }

    message.info("正在拆分文档...")
    const { data: splitData, error: splitError } = await splitDocument(splitPayload)

    if (splitError || splitData?.code !== 0) {
      message.error(getErrorMessage(splitError, "文档拆分失败"))
      return
    }

    message.success(`文档拆分完成，生成 ${splitData.data?.chunks_count ?? 0} 个片段`)

    // Switch to chunks tab and refresh
    activeTab.value = "chunks"
    await fetchChunks()
  } catch {
    message.error("处理文件时发生错误")
  } finally {
    isProcessing.value = false
  }
}

// ---- Lifecycle ----

watch(
  () => props.id,
  () => {
    filePage.value = 1
    chunkPage.value = 1
    if (activeTab.value === "files") {
      fetchFiles()
    } else {
      fetchChunks()
    }
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
        <NButton :loading="fileLoading || chunkLoading" @click="refreshCurrentTab">刷新</NButton>
        <NButton @click="router.back()">返回</NButton>
      </NSpace>
    </div>

    <!-- File Upload Card -->
    <NCard title="文件上传" :bordered="false" size="small">
      <NSpace vertical :size="12">
        <NSpace align="end">
          <div>
            <div class="mb-4px text-12px text-gray-500">Chunk Size</div>
            <NInputNumber
              v-model:value="chunkSize"
              :min="64"
              :max="4096"
              size="small"
              style="width: 140px"
            />
          </div>
          <div>
            <div class="mb-4px text-12px text-gray-500">Chunk Overlap</div>
            <NInputNumber
              v-model:value="chunkOverlap"
              :min="0"
              :max="2048"
              size="small"
              style="width: 140px"
            />
          </div>
          <NButton
            type="primary"
            :loading="isUploading"
            :disabled="isProcessing"
            @click="triggerFileUpload"
          >
            {{ isUploading ? "上传中..." : "选择文件" }}
          </NButton>
        </NSpace>

        <!-- Uploaded file list -->
        <div v-if="uploadedFiles.length > 0" class="flex flex-col gap-8px">
          <div
            v-for="file in uploadedFiles"
            :key="file.id"
            class="flex items-center justify-between rounded-6px bg-gray-50 px-12px py-8px"
          >
            <div class="flex items-center gap-8px">
              <span class="text-14px font-500">{{ file.name }}</span>
              <NTag size="small" :bordered="false">{{ (file.size / 1024).toFixed(1) }} KB</NTag>
            </div>
            <NButton
              size="small"
              type="primary"
              ghost
              :loading="isProcessing"
              @click="handleSplitAndEmbed(file)"
            >
              拆分处理
            </NButton>
          </div>
        </div>
      </NSpace>

      <input ref="fileInputRef" type="file" multiple class="hidden" @change="handleFileChange" />
    </NCard>

    <!-- Toggle + Table Card -->
    <NCard :bordered="false" size="small">
      <template #header>
        <div class="flex items-center justify-between">
          <NButtonGroup size="small">
            <NButton
              :type="activeTab === 'files' ? 'primary' : 'default'"
              @click="handleTabChange('files')"
            >
              文件列表
            </NButton>
            <NButton
              :type="activeTab === 'chunks' ? 'primary' : 'default'"
              @click="handleTabChange('chunks')"
            >
              分块列表
            </NButton>
          </NButtonGroup>
        </div>
      </template>

      <!-- File Table -->
      <div v-if="activeTab === 'files'" style="max-height: 60vh; overflow: auto">
        <NDataTable
          :columns="fileColumns"
          :data="files"
          :loading="fileLoading"
          :bordered="false"
          size="small"
        />
      </div>

      <!-- Chunk Table -->
      <div v-else style="max-height: 60vh; overflow: auto">
        <NDataTable
          :columns="chunkColumns"
          :data="chunks"
          :loading="chunkLoading"
          :bordered="false"
          size="small"
          :scroll-x="scrollX"
        />
      </div>

      <template #footer>
        <!-- File Pagination -->
        <NSpace v-if="activeTab === 'files'" justify="end">
          <NPagination
            v-model:page="filePage"
            :item-count="fileTotal"
            :page-size="filePageSize"
            show-size-picker
            :page-sizes="[10, 20, 50, 100]"
            @update:page="fetchFiles"
            @update:page-size="handleFilePageSizeChange"
          />
        </NSpace>
        <!-- Chunk Pagination -->
        <NSpace v-else justify="end">
          <NPagination
            v-model:page="chunkPage"
            :item-count="chunkTotal"
            :page-size="chunkPageSize"
            show-size-picker
            :page-sizes="[10, 20, 50, 100]"
            @update:page="fetchChunks"
            @update:page-size="handleChunkPageSizeChange"
          />
        </NSpace>
      </template>
    </NCard>
  </NSpace>
</template>
