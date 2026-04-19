<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue';
import type { Component } from 'vue';
import { NIcon, NInputNumber, NSelect, useDialog, useMessage } from 'naive-ui';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';
import {
  Add12Filled,
  Delete16Regular as DeleteIcon,
  Edit24Regular as EditIcon,
  Eye24Regular as EyeIcon,
  MoreHorizontal28Regular,
  Search48Filled
} from '@vicons/fluent';
import { createDataset, deleteDataset, getDatasetById, getDatasets, updateDataset } from '@/service/api/dataset';

const dialog = useDialog();
const message = useMessage();
dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

// Loading state
const loading = ref(false);

// Detail modal
const showDetailModal = ref(false);
const detailDataset = ref<Api.Dataset.DatasetItem | null>(null);

function formatTime(isoString: string) {
  return dayjs(isoString).fromNow();
}
const showModal = ref(false);
const isEdit = ref(false);
const bodyStyle = {
  width: '500px'
};
const searchKey = ref();
const datasets = ref<Api.Dataset.DatasetItem[]>();
const model: Api.Dataset.FormModel = reactive({
  id: undefined,
  icon: '🤖',
  description: '',
  name: '',
  search_type: 'hybrid',
  embedding_model: '',
  provider_id: 0
});

const searchTypeOptions = [
  { label: '关键词检索', value: 'sparse' },
  { label: '语义检索', value: 'dense' },
  { label: '混合检索', value: 'hybrid' }
];
function renderIcon(icon: Component) {
  return () => {
    return h(NIcon, null, {
      default: () => h(icon)
    });
  };
}
function openCreateModal() {
  model.id = undefined;
  model.icon = '🤖';
  model.name = '';
  model.description = '';
  model.search_type = 'hybrid';
  model.embedding_model = '';
  model.provider_id = 0;
  isEdit.value = false;
  showModal.value = true;
}
async function fetchDatasets() {
  loading.value = true;
  const keyword = searchKey.value && searchKey.value.trim() ? searchKey.value.trim() : undefined;
  const { response: res } = await getDatasets(keyword);
  loading.value = false;
  if (res?.data?.code === 0) {
    const datasetList = res.data.data?.datasets ?? [];
    datasets.value = datasetList;
  } else {
    message.error('获取数据失败');
    datasets.value = [];
  }
}
async function handleViewDetail(datasetId: number) {
  const { response: res } = await getDatasetById(datasetId);
  if (res?.data?.code === 0) {
    detailDataset.value = res.data.data;
    showDetailModal.value = true;
  } else {
    message.error(res.data.msg || '获取详情失败');
  }
}
async function handleCreateDataset(md: Api.Dataset.FormModel) {
  const { response: res } = await createDataset(md);
  if (res.data.code === 0) {
    message.success('创建成功');
  } else {
    const errorMsg = res.data.msg;
    window.$message?.error(errorMsg);
  }
  fetchDatasets();
}
async function handleDeleteDataset(id: number) {
  const { response: res } = await deleteDataset(id);
  if (res.data.code === 0) {
    message.success('删除成功');
  } else {
    const errorMsg = res.data.msg;
    window.$message?.error(errorMsg);
  }
  fetchDatasets();
}
async function handleEditDataset(md: Api.Dataset.FormModel) {
  const { response: res } = await updateDataset(md);
  if (res.data.code === 0) {
    message.success('编辑成功');
  } else {
    const errorMsg = res.data.msg;
    window.$message?.error(errorMsg);
  }
  fetchDatasets();
}
const emojiList = ['😀', '😂', '😎', '🤖', '🐱', '🦊', '🐶', '🦄', '🐼', '🦉'];

const hovered = ref<string | null>(null);
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
function handleSelect(key: string, dataset: Api.Dataset.DatasetItem) {
  switch (key) {
    case 'view':
      handleViewDetail(dataset.id);
      break;
    case 'edit':
      showModal.value = true;
      isEdit.value = true;
      Object.assign(model, {
        id: dataset.id,
        icon: dataset.icon || '🤖',
        name: dataset.name,
        description: dataset.description,
        search_type: dataset.search_type || 'hybrid',
        embedding_model: dataset.embedding_model || '',
        provider_id: dataset.provider_id || 0
      });
      break;
    case 'delete':
      dialog.warning({
        title: '确认删除',
        content: `确定要删除知识库 “${dataset.name}” 吗？此操作不可恢复。`,
        positiveText: '删除',
        negativeText: '取消',
        onPositiveClick: () => handleDeleteDataset(dataset.id)
      });
      break;
    default:
  }
}
function handleSubmit() {
  if (isEdit.value && model.id !== null) {
    handleEditDataset(model);
  } else {
    handleCreateDataset(model);
  }
}
function selectEmoji(emoji: string) {
  model.icon = emoji;
}

onMounted(() => {
  fetchDatasets();
});
</script>

<template>
  <NSpace vertical :size="16">
    <NSpace>
      <NInput v-model:value="searchKey" round placeholder="请输入关键字" clearable :loading="loading" @blur="fetchDatasets()"
        @keyup.enter="fetchDatasets()">
        <template #prefix>
          <NIcon :component="Search48Filled" />
        </template>
      </NInput>
    </NSpace>

    <NSpin :show="loading">
      <NGrid cols="3" x-gap="16" y-gap="17" responsive="screen" item-responsive>
        <NGridItem class="grid-item-equal-height">
          <NCard size="huge" class="bg-[#E8ECEF]" content-style="padding-left: 5px;padding-right: 5px; height: 100%;">
            <NButton quaternary class="w-full justify-start" @click="openCreateModal">
              <NIcon :component="Add12Filled" class="mr-2" />
              创建知识库
            </NButton>
            <NModal v-model:show="showModal" :mask-closable="false" preset="dialog" :show-icon="false"
              :style="bodyStyle" positive-text="确认" negative-text="取消" title="知识库设置" @positive-click="handleSubmit">
              <NSpace :size="10" vertical>
                <NCard title="知识库名称" :bordered="false" size="small" content-style="display:flex;gap: 8px;">
                  <NPopover trigger="click" placement="bottom-start">
                    <template #trigger>
                      <NAvatar :style="{
                        color: 'black',
                        backgroundColor: '#FFEAD5',
                        cursor: 'pointer'
                      }">
                        {{ model.icon }}
                      </NAvatar>
                    </template>
                    <!-- Popover 内容：Emoji 选择器 -->
                    <div style="display: flex; gap: 8px; padding: 5px; flex-wrap: wrap">
                      <span v-for="emoji in emojiList" :key="emoji"
                        style="font-size: 20px; cursor: pointer; padding: 4px; border-radius: 4px" :style="{
                          backgroundColor: hovered === emoji ? '#e0e0e0' : 'transparent'
                        }" @click="selectEmoji(emoji)" @mouseenter="hovered = emoji" @mouseleave="hovered = null">
                        {{ emoji }}
                      </span>
                    </div>
                  </NPopover>

                  <NInput v-model:value="model.name" type="text" style="background-color: #f1f3f6" size="tiny"
                    placeholder="请输入知识库名称" />
                </NCard>
                <NCard title="描述" :bordered="false" size="small">
                  <NInput v-model:value="model.description" type="textarea" size="tiny"
                    style="background-color: #f1f3f6" placeholder="描述该数据集的内容。详细描述可以让AI更快地访问数据集的内容。如果为空，将使用默认的命中策略。" />
                </NCard>
                <NCard title="检索方式" :bordered="false" size="small">
                  <NSelect v-model:value="model.search_type" :options="searchTypeOptions" size="tiny"
                    style="background-color: #f1f3f6" />
                </NCard>
                <NCard title="Embedding模型" :bordered="false" size="small">
                  <NInput v-model:value="model.embedding_model" type="text" size="tiny"
                    style="background-color: #f1f3f6" placeholder="请输入Embedding模型名称" />
                </NCard>
                <NCard title="Provider ID" :bordered="false" size="small">
                  <NInputNumber v-model:value="model.provider_id" size="tiny" style="background-color: #f1f3f6"
                    placeholder="请输入Provider ID" />
                </NCard>
              </NSpace>
              <NSpace></NSpace>
            </NModal>
          </NCard>
        </NGridItem>
        <NGridItem v-for="dataset in datasets" :key="dataset.id" class="grid-item-equal-height">
          <NCard hoverable size="huge" style="cursor: pointer" @click="handleViewDetail(dataset.id)">
            <NSpace vertical>
              <div style="display: flex; align-items: flex-start; gap: 8px; width: 100%">
                <NAvatar size="large" :style="{
                  color: 'black',
                  backgroundColor: '#E0F2FE',
                  cursor: 'pointer'
                }">
                  {{ dataset.icon || '🤖' }}
                </NAvatar>

                <div style="flex: 1">
                  <div :style="{ fontWeight: 'bold' }">{{ dataset.name }}</div>
                  <div style="color: #949494; font-size: 10px; line-height: 1.2">
                    {{ dayjs(dataset.created_at).format('YYYY-MM-DD') }}
                  </div>
                </div>

                <NDropdown :options="dropdownOptions" trigger="click" size="small"
                  @select="key => handleSelect(key, dataset)">
                  <NButton size="small" secondary @click.stop>
                    <NIcon>
                      <MoreHorizontal28Regular />
                    </NIcon>
                  </NButton>
                </NDropdown>
              </div>
              <div style="color: #666; font-size: 10px; line-height: 1.4">
                {{ dataset.description }}
              </div>
              <div style="color: #949494; font-size: 10px; transform: scale(0.85); transform-origin: left top">
                更新于 · {{ formatTime(dataset.updated_at) }}
              </div>
            </NSpace>
          </NCard>
        </NGridItem>
      </NGrid>
    </NSpin>

    <!-- Detail Modal -->
    <NModal v-model:show="showDetailModal" :mask-closable="true" :show-icon="false"
      style="width: 700px; max-width: 90vw">
      <NCard :bordered="false" size="large" style="border-radius: 12px">
        <NSpace vertical :size="20">
          <!-- Header Section -->
          <div style="
              display: flex;
              align-items: center;
              gap: 16px;
              padding-bottom: 16px;
              border-bottom: 1px solid #e8e8e8;
            ">
            <NAvatar :size="64" :style="{
              color: '#333',
              backgroundColor: '#f5f5f5',
              fontSize: '32px'
            }">
              {{ detailDataset?.icon || '🤖' }}
            </NAvatar>
            <div style="flex: 1">
              <h2 style="margin: 0; color: #333; font-size: 24px; font-weight: 600">{{ detailDataset?.name }}</h2>
              <div style="color: #999; font-size: 12px; margin-top: 4px">
                创建于 {{ dayjs(detailDataset?.created_at).format('YYYY-MM-DD HH:mm') }}
              </div>
            </div>
            <NButton text style="color: #666" @click="showDetailModal = false">
              <template #icon>
                <NIcon size="24">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                    <path
                      d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z" />
                  </svg>
                </NIcon>
              </template>
            </NButton>
          </div>

          <!-- Description Section -->
          <div v-if="detailDataset?.description" style="background: #f8f9fa; padding: 16px; border-radius: 8px">
            <div style="color: #666; font-size: 12px; margin-bottom: 8px; font-weight: 500">描述</div>
            <div style="color: #333; font-size: 14px; line-height: 1.6">{{ detailDataset.description }}</div>
          </div>

          <!-- Info Cards -->
          <NGrid cols="2" x-gap="12" y-gap="12" responsive="screen">
            <NGridItem>
              <div class="info-card">
                <div class="info-label">检索方式</div>
                <div class="info-value">
                  <NTag :type="detailDataset?.search_type === 'hybrid'
                      ? 'success'
                      : detailDataset?.search_type === 'dense'
                        ? 'info'
                        : 'warning'
                    " size="small">
                    {{
                      searchTypeOptions.find(o => o.value === detailDataset?.search_type)?.label ||
                      detailDataset?.search_type
                    }}
                  </NTag>
                </div>
              </div>
            </NGridItem>
            <NGridItem>
              <div class="info-card">
                <div class="info-label">Embedding 模型</div>
                <div class="info-value">{{ detailDataset?.embedding_model || '未设置' }}</div>
              </div>
            </NGridItem>
            <NGridItem>
              <div class="info-card">
                <div class="info-label">Provider ID</div>
                <div class="info-value">{{ detailDataset?.provider_id || '-' }}</div>
              </div>
            </NGridItem>
            <NGridItem>
              <div class="info-card">
                <div class="info-label">Owner ID</div>
                <div class="info-value">{{ detailDataset?.owner_id || '-' }}</div>
              </div>
            </NGridItem>
          </NGrid>

          <!-- Footer -->
          <div style="
              display: flex;
              justify-content: space-between;
              align-items: center;
              padding-top: 16px;
              border-top: 1px solid #e8e8e8;
            ">
            <div style="color: #999; font-size: 12px">ID: {{ detailDataset?.id }}</div>
            <div style="color: #999; font-size: 12px">更新于 {{ formatTime(detailDataset?.updated_at || '') }}</div>
          </div>
        </NSpace>
      </NCard>
    </NModal>
  </NSpace>
</template>

<style scoped>
.grid-item-equal-height {
  display: flex;
  min-height: 170px;
}

.info-card {
  background: #f8f9fa;
  padding: 12px 16px;
  border-radius: 8px;
  border: 1px solid #e8e8e8;
}

.info-label {
  color: #666;
  font-size: 12px;
  margin-bottom: 6px;
  font-weight: 500;
}

.info-value {
  color: #333;
  font-size: 14px;
  font-weight: 500;
}
</style>
