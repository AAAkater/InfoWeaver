<script setup lang="ts">
import { computed, h } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NIcon, NSwitch } from 'naive-ui';
import { Add12Filled } from '@vicons/fluent';

const props = defineProps<{
  loading?: boolean;
  modelFilter: 'all' | 'enabled' | 'disabled';
  page: number;
  pageSize: number;
  paginatedModels: Api.Provider.ModelInfo[];
  providerName?: string;
  show: boolean;
  totalModelsCount: number;
}>();

const emit = defineEmits<{
  (event: 'add-model'): void;
  (event: 'toggle-model', modelId: string, enabled: boolean): void;
  (event: 'update:modelFilter', value: 'all' | 'enabled' | 'disabled'): void;
  (event: 'update:page', value: number): void;
  (event: 'update:pageSize', value: number): void;
  (event: 'update:show', value: boolean): void;
}>();

const bodyStyle = {
  width: '600px'
};

const visible = computed({
  get: () => props.show,
  set: value => emit('update:show', value)
});

const columns: DataTableColumns<Api.Provider.ModelInfo> = [
  {
    title: '模型名',
    key: 'id',
    minWidth: 220,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '类型',
    key: 'object',
    width: 140
  },
  {
    title: '所属',
    key: 'owned_by',
    width: 140
  },
  {
    title: '状态',
    key: 'enabled',
    width: 120,
    render(row) {
      return h(NSwitch, {
        value: !!row.enabled,
        onUpdateValue: (value: boolean) => emit('toggle-model', row.id, value)
      });
    }
  }
];

const filterOptions = [
  { label: '全部', value: 'all' },
  { label: '已启用', value: 'enabled' },
  { label: '已禁用', value: 'disabled' }
];
</script>

<template>
  <NModal v-model:show="visible" :mask-closable="true" preset="dialog" :show-icon="false" :style="bodyStyle">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center; width: 100%">
        <div style="font-weight: bold">{{ props.providerName }} - 模型列表</div>
        <NSpace :size="8">
          <NSelect
            :value="props.modelFilter"
            :options="filterOptions"
            size="small"
            style="width: 100px"
            @update:value="value => emit('update:modelFilter', value)"
          />
          <NButton type="primary" size="small" @click="emit('add-model')">
            <template #icon>
              <NIcon :component="Add12Filled" />
            </template>
            添加模型
          </NButton>
        </NSpace>
      </div>
    </template>

    <NDataTable
      :columns="columns"
      :data="props.paginatedModels"
      :loading="props.loading"
      :bordered="false"
      :single-line="false"
      size="small"
    />

    <template #action>
      <NSpace justify="space-between" align="center">
        <NPagination
          :page="props.page"
          :page-count="Math.max(1, Math.ceil(props.totalModelsCount / props.pageSize))"
          :page-size="props.pageSize"
          show-size-picker
          :page-sizes="[10, 20, 50]"
          @update:page="value => emit('update:page', value)"
          @update:page-size="value => emit('update:pageSize', value)"
        />
        <NButton @click="emit('update:show', false)">关闭</NButton>
      </NSpace>
    </template>
  </NModal>
</template>
