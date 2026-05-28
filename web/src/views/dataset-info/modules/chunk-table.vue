<script setup lang="ts">
import { h } from "vue"
import type { DataTableColumns, DataTableRowKey } from "naive-ui"
import { NButton, NTag, NDataTable } from "naive-ui"
import dayjs from "dayjs"

defineOptions({ name: "ChunkTable" })

interface Props {
  chunks: Api.Dataset.ChunkInfo[]
  loading: boolean
  total: number
  page: number
  pageSize: number
  embedLoading: boolean
  checkedRowKeys: DataTableRowKey[]
}

defineProps<Props>()

const emit = defineEmits<{
  updatePage: [page: number]
  updatePageSize: [size: number]
  updateCheckedRowKeys: [keys: DataTableRowKey[]]
  embed: []
}>()

const scrollX = 40 + 80 + 360 + 100 + 180 + 170

const columns: DataTableColumns<Api.Dataset.ChunkInfo> = [
  {
    type: "selection",
    width: 40,
  },
  {
    title: "ID",
    key: "id",
    width: 80,
  },
  {
    title: "内容",
    key: "content",
    minWidth: 360,
    ellipsis: { tooltip: true },
  },
  {
    title: "状态",
    key: "status",
    width: 100,
    render(row) {
      const type = row.status === "completed" || row.status === "success" ? "success" : "default"
      return h(NTag, { size: "small", type }, { default: () => [row.status || "-"] })
    },
  },
  {
    title: "元数据",
    key: "chunk_metadata",
    minWidth: 180,
    ellipsis: { tooltip: true },
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
]
</script>

<template>
  <div>
    <div class="mb-12px flex justify-end">
      <NButton
        type="primary"
        size="small"
        :disabled="checkedRowKeys.length === 0"
        :loading="embedLoading"
        @click="emit('embed')"
      >
        向量化 {{ checkedRowKeys.length > 0 ? `(${checkedRowKeys.length})` : "" }}
      </NButton>
    </div>

    <div style="max-height: 60vh; overflow: auto">
      <NDataTable
        :columns="columns"
        :data="chunks"
        :loading="loading"
        :bordered="false"
        size="small"
        :scroll-x="scrollX"
        :row-key="(row: Api.Dataset.ChunkInfo) => row.id"
        :checked-row-keys="checkedRowKeys"
        @update:checked-row-keys="emit('updateCheckedRowKeys', $event)"
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
