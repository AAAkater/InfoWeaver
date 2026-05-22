<script setup lang="ts">
import { h } from "vue"
import type { Component } from "vue"
import type { DataTableColumns } from "naive-ui"
import { NButton, NIcon, NTag } from "naive-ui"
import dayjs from "dayjs"
import {
  Delete16Regular as DeleteIcon,
  Edit24Regular as EditIcon,
  Server24Regular,
} from "@vicons/fluent"

const props = defineProps<{
  providers: Api.Provider.ProviderInfo[]
}>()

const emit = defineEmits<{
  (event: "select", key: "edit" | "models" | "delete", provider: Api.Provider.ProviderInfo): void
}>()

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const columns: DataTableColumns<Api.Provider.ProviderInfo> = [
  {
    title: "名称",
    key: "name",
    width: 180,
  },
  {
    title: "Base URL",
    key: "base_url",
    minWidth: 280,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: "模式",
    key: "mode",
    width: 150,
    render(row) {
      return h(
        NTag,
        {
          size: "small",
          type: "info",
          bordered: false,
        },
        { default: () => row.mode },
      )
    },
  },
  {
    title: "创建时间",
    key: "created_at",
    width: 180,
    render(row) {
      return dayjs(row.created_at).format("YYYY-MM-DD HH:mm")
    },
  },
  {
    title: "操作",
    key: "actions",
    width: 220,
    render(row) {
      return h(
        "div",
        {
          style: "display: flex; gap: 8px; flex-wrap: wrap;",
        },
        [
          h(
            NButton,
            {
              size: "small",
              secondary: true,
              onClick: () => emit("select", "edit", row),
            },
            {
              default: () => "编辑",
              icon: renderIcon(EditIcon),
            },
          ),
          h(
            NButton,
            {
              size: "small",
              secondary: true,
              onClick: () => emit("select", "models", row),
            },
            {
              default: () => "模型",
              icon: renderIcon(Server24Regular),
            },
          ),
          h(
            NButton,
            {
              size: "small",
              secondary: true,
              type: "error",
              onClick: () => emit("select", "delete", row),
            },
            {
              default: () => "删除",
              icon: renderIcon(DeleteIcon),
            },
          ),
        ],
      )
    },
  },
]
</script>

<template>
  <NDataTable
    :columns="columns"
    :data="props.providers"
    :bordered="false"
    :single-line="false"
    size="small"
  />
</template>
