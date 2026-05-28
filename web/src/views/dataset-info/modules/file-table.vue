<script setup lang="ts">
import { h } from "vue"
import type { DataTableColumns } from "naive-ui"
import { NButton, NButtonGroup } from "naive-ui"
import dayjs from "dayjs"

defineOptions({ name: "FileTable" })

interface Props {
  files: Api.Dataset.SimpleFileInfo[]
  loading: boolean
  total: number
  page: number
  pageSize: number
}

defineProps<Props>()

const emit = defineEmits<{
  updatePage: [page: number]
  updatePageSize: [size: number]
  split: [file: Api.Dataset.SimpleFileInfo]
  delete: [file: Api.Dataset.SimpleFileInfo]
}>()

/** Map MIME type to short display name */
function formatFileType(type: string): string {
  const mimeMap: Record<string, string> = {
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": "docx",
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": "xlsx",
    "application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
    "application/msword": "doc",
    "application/vnd.ms-excel": "xls",
    "application/vnd.ms-powerpoint": "ppt",
    "application/pdf": "pdf",
    "text/plain": "txt",
    "text/csv": "csv",
    "text/markdown": "md",
    "application/json": "json",
    "text/html": "html",
    "image/png": "png",
    "image/jpeg": "jpg",
    "image/gif": "gif",
    "image/webp": "webp",
  }
  return mimeMap[type] || type
}

const columns: DataTableColumns<Api.Dataset.SimpleFileInfo> = [
  {
    title: "ID",
    key: "id",
    width: 80,
  },
  {
    title: "文件名",
    key: "name",
    minWidth: 160,
    ellipsis: { tooltip: true },
  },
  {
    title: "类型",
    key: "type",
    width: 100,
    ellipsis: { tooltip: true },
    render(row) {
      return row.type ? formatFileType(row.type) : "-"
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
    key: "created_at",
    width: 180,
    render(row) {
      return row.created_at ? dayjs(row.created_at).format("YYYY-MM-DD HH:mm") : "-"
    },
  },
  {
    title: "操作",
    key: "actions",
    width: 160,
    render(row) {
      return h(
        NButtonGroup,
        { size: "tiny" },
        {
          default: () => [
            h(
              NButton,
              { type: "primary", onClick: () => emit("split", row) },
              { default: () => ["分块"] },
            ),
            h(
              NButton,
              { type: "error", onClick: () => emit("delete", row) },
              { default: () => ["删除"] },
            ),
          ],
        },
      )
    },
  },
]
</script>

<template>
  <div>
    <div style="max-height: 60vh; overflow: auto">
      <NDataTable
        :columns="columns"
        :data="files"
        :loading="loading"
        :bordered="false"
        size="small"
      />
    </div>

    <div class="mt-12px flex justify-end">
      <NPagination
        :page="page"
        :item-count="total"
        :page-size="pageSize"
        show-size-picker
        :page-sizes="[10, 20, 50, 100]"
        @update:page="emit('updatePage', $event)"
        @update:page-size="emit('updatePageSize', $event)"
      />
    </div>
  </div>
</template>
