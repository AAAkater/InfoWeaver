<script setup lang="ts">
import { NAvatar, NButton, NCard, NIcon, NSpace } from "naive-ui"
import dayjs from "dayjs"
import relativeTime from "dayjs/plugin/relativeTime"
import "dayjs/locale/zh-cn"
import {
  Delete16Regular as DeleteIcon,
  Edit24Regular as EditIcon,
  Eye24Regular as EyeIcon,
} from "@vicons/fluent"
dayjs.extend(relativeTime)
dayjs.locale("zh-cn")

const props = defineProps<{
  dataset: Api.Dataset.DatasetItem
}>()

const emit = defineEmits<{
  (event: "select", key: string, dataset: Api.Dataset.DatasetItem): void
}>()

function formatTime(isoString: string) {
  return dayjs(isoString).fromNow()
}
</script>

<template>
  <NCard hoverable size="huge">
    <NSpace vertical>
      <div style="display: flex; align-items: flex-start; gap: 8px; width: 100%">
        <div style="flex: 1; cursor: pointer" @click="emit('select', 'view', props.dataset)">
          <NAvatar
            size="large"
            :style="{
              color: 'black',
              backgroundColor: '#E0F2FE',
              cursor: 'pointer',
            }"
          >
            {{ props.dataset.icon || "🤖" }}
          </NAvatar>

          <div style="margin-top: 8px">
            <div :style="{ fontWeight: 'bold' }">{{ props.dataset.name }}</div>
            <div style="color: #949494; font-size: 10px; line-height: 1.2">
              {{ dayjs(props.dataset.created_at).format("YYYY-MM-DD") }}
            </div>
          </div>
        </div>

        <div style="flex-shrink: 0; display: flex; gap: 8px">
          <NButton
            size="small"
            secondary
            circle
            @click.stop="emit('select', 'view', props.dataset)"
          >
            <NIcon>
              <EyeIcon />
            </NIcon>
          </NButton>
          <NButton
            size="small"
            secondary
            circle
            @click.stop="emit('select', 'edit', props.dataset)"
          >
            <NIcon>
              <EditIcon />
            </NIcon>
          </NButton>
          <NButton
            size="small"
            tertiary
            circle
            type="error"
            @click.stop="emit('select', 'delete', props.dataset)"
          >
            <NIcon>
              <DeleteIcon />
            </NIcon>
          </NButton>
        </div>
      </div>
      <div style="color: #666; font-size: 10px; line-height: 1.4">
        {{ props.dataset.description }}
      </div>
      <div
        style="color: #949494; font-size: 10px; transform: scale(0.85); transform-origin: left top"
      >
        更新于 · {{ formatTime(props.dataset.updated_at) }}
      </div>
    </NSpace>
  </NCard>
</template>
