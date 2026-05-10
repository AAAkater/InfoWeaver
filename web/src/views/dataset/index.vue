<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { useDialog, useMessage } from 'naive-ui';
import { createDataset, deleteDataset, getDatasetById, getDatasets, updateDataset } from '@/service/api/dataset';
import DatasetCard from './modules/dataset-card.vue';
import DatasetCreateCard from './modules/dataset-create-card.vue';
import DatasetDetailModal from './modules/dataset-detail-modal.vue';
import DatasetFormModal from './modules/dataset-form-modal.vue';
import DatasetSearch from './modules/dataset-search.vue';

const dialog = useDialog();
const message = useMessage();

const loading = ref(false);
const showDetailModal = ref(false);
const detailDataset = ref<Api.Dataset.DatasetItem | null>(null);
const showModal = ref(false);
const isEdit = ref(false);
const searchKey = ref('');
const datasets = ref<Api.Dataset.DatasetItem[]>([]);

const model: Api.Dataset.FormModel = reactive({
  id: undefined,
  icon: '🤖',
  description: '',
  name: '',
  search_type: 'hybrid',
  embedding_model: '',
  provider_id: 0
});

function resetModel() {
  Object.assign(model, {
    id: undefined,
    icon: '🤖',
    description: '',
    name: '',
    search_type: 'hybrid',
    embedding_model: '',
    provider_id: 0
  });
}

function fillModel(dataset: Api.Dataset.DatasetItem) {
  Object.assign(model, {
    id: dataset.id,
    icon: dataset.icon || '🤖',
    name: dataset.name,
    description: dataset.description,
    search_type: dataset.search_type || 'hybrid',
    embedding_model: dataset.embedding_model || '',
    provider_id: dataset.provider_id || 0
  });
}

function openCreateModal() {
  resetModel();
  isEdit.value = false;
  showModal.value = true;
}

async function fetchDatasets() {
  loading.value = true;

  try {
    const keyword = searchKey.value.trim() || undefined;
    const { response: res } = await getDatasets(keyword);

    if (res?.data?.code === 0) {
      datasets.value = res.data.data?.datasets ?? [];
    } else {
      message.error('获取数据失败');
      datasets.value = [];
    }
  } finally {
    loading.value = false;
  }
}

async function handleViewDetail(datasetId: number) {
  const { response: res } = await getDatasetById(datasetId);

  if (res?.data?.code === 0) {
    detailDataset.value = res.data.data;
    showDetailModal.value = true;
  } else {
    message.error(res?.data?.msg || '获取详情失败');
  }
}

async function handleCreateDataset(md: Api.Dataset.FormModel) {
  const { response: res } = await createDataset(md);

  if (res.data.code === 0) {
    message.success('创建成功');
    showModal.value = false;
    await fetchDatasets();
  } else {
    message.error(res.data.msg || '创建失败');
  }
}

async function handleDeleteDataset(id: number) {
  const { response: res } = await deleteDataset(id);

  if (res.data.code === 0) {
    message.success('删除成功');
    await fetchDatasets();
  } else {
    message.error(res.data.msg || '删除失败');
  }
}

async function handleEditDataset(md: Api.Dataset.FormModel) {
  const { response: res } = await updateDataset(md);

  if (res.data.code === 0) {
    message.success('编辑成功');
    showModal.value = false;
    await fetchDatasets();
  } else {
    message.error(res.data.msg || '编辑失败');
  }
}

function handleSelect(key: string, dataset: Api.Dataset.DatasetItem) {
  switch (key) {
    case 'view':
      handleViewDetail(dataset.id);
      break;
    case 'edit':
      isEdit.value = true;
      fillModel(dataset);
      showModal.value = true;
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

function handleSubmit(form: Api.Dataset.FormModel) {
  if (isEdit.value && form.id !== undefined) {
    handleEditDataset(form);
  } else {
    handleCreateDataset(form);
  }
}

onMounted(() => {
  fetchDatasets();
});
</script>

<template>
  <NSpace vertical :size="16">
    <DatasetSearch v-model="searchKey" :loading="loading" @search="fetchDatasets" />

    <NSpin :show="loading">
      <NGrid cols="3" x-gap="16" y-gap="17" responsive="screen" item-responsive>
        <NGridItem class="grid-item-equal-height">
          <DatasetCreateCard @create="openCreateModal" />
        </NGridItem>

        <NGridItem v-for="dataset in datasets" :key="dataset.id" class="grid-item-equal-height">
          <DatasetCard :dataset="dataset" @select="handleSelect" />
        </NGridItem>
      </NGrid>
    </NSpin>

    <DatasetFormModal v-model:show="showModal" :model="model" :is-edit="isEdit" @submit="handleSubmit" />

    <DatasetDetailModal v-model:show="showDetailModal" :dataset="detailDataset" />
  </NSpace>
</template>

<style scoped>
.grid-item-equal-height {
  display: flex;
  min-height: 170px;
}
</style>
