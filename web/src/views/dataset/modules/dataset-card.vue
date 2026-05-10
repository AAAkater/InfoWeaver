<script setup lang="ts">
import { h } from 'vue';
import type { Component } from 'vue';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';
import {
  Delete16Regular as DeleteIcon,
  Edit24Regular as EditIcon,
  Eye24Regular as EyeIcon,
  MoreHorizontal28Regular
} from '@vicons/fluent';

dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

const props = defineProps<{
  dataset: Api.Dataset.DatasetItem;
}>();

const emit = defineEmits<{
  (event: 'select', key: string, dataset: Api.Dataset.DatasetItem): void;
}>();

function renderIcon(icon: Component) {
  return () =>
    h(NIcon, null, {
      default: () => h(icon)
    });
}

const dropdownOptions = [
  {
    label: '查看详情',
    key: 'view',
    icon: renderIcon(EyeIcon)
  },
  {
    label: '编辑',
    key: 'edit',
    icon: renderIcon(EditIcon)
  },
  {
    label: '删除',
    key: 'delete',
    icon: renderIcon(DeleteIcon)
  }
];

function formatTime(isoString: string) {
  return dayjs(isoString).fromNow();
}
</script>

<template>
  <NCard hoverable size="huge" style="cursor: pointer" @click="emit('select', 'view', props.dataset)">
    <NSpace vertical>
      <div style="display: flex; align-items: flex-start; gap: 8px; width: 100%">
        <NAvatar
          size="large"
          :style="{
            color: 'black',
            backgroundColor: '#E0F2FE',
            cursor: 'pointer'
          }"
        >
          {{ props.dataset.icon || '🤖' }}
        </NAvatar>

        <div style="flex: 1">
          <div :style="{ fontWeight: 'bold' }">{{ props.dataset.name }}</div>
          <div style="color: #949494; font-size: 10px; line-height: 1.2">
            {{ dayjs(props.dataset.created_at).format('YYYY-MM-DD') }}
          </div>
        </div>

        <NDropdown
          :options="dropdownOptions"
          trigger="click"
          size="small"
          @select="key => emit('select', key, props.dataset)"
        >
          <NButton size="small" secondary @click.stop>
            <NIcon>
              <MoreHorizontal28Regular />
            </NIcon>
          </NButton>
        </NDropdown>
      </div>
      <div style="color: #666; font-size: 10px; line-height: 1.4">
        {{ props.dataset.description }}
      </div>
      <div style="color: #949494; font-size: 10px; transform: scale(0.85); transform-origin: left top">
        更新于 · {{ formatTime(props.dataset.updated_at) }}
      </div>
    </NSpace>
  </NCard>
</template>
