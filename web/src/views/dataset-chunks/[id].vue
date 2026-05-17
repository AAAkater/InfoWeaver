<script setup lang="ts">
import { computed, h, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NInputNumber, NTag, useMessage } from 'naive-ui';
import dayjs from 'dayjs';
import { embedChunks, getDatasetChunks, splitDocument, uploadFiles } from '@/service/api/dataset';

interface Props {
  id: string;
}

const props = defineProps<Props>();

const router = useRouter();
const message = useMessage();

const loading = ref(false);
const chunks = ref<Api.Dataset.ChunkInfo[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(20);

// File upload
const fileInputRef = ref<HTMLInputElement | null>(null);
const isUploading = ref(false);
const uploadedFiles = ref<Api.Dataset.FileUploadInfo[]>([]);
const chunkSize = ref(512);
const chunkOverlap = ref(50);
const isProcessing = ref(false);

const datasetId = computed(() => Number(props.id));

const columns: DataTableColumns<Api.Dataset.ChunkInfo> = [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: '内容',
    key: 'content',
    minWidth: 360,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row) {
      const type = row.status === 'completed' || row.status === 'success' ? 'success' : 'default';

      return h(NTag, { size: 'small', type }, { default: () => row.status || '-' });
    }
  },
  {
    title: 'File ID',
    key: 'file_id',
    width: 100
  },
  {
    title: 'Vector ID',
    key: 'vector_id',
    minWidth: 180,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '元数据',
    key: 'chunk_metadata',
    minWidth: 180,
    ellipsis: {
      tooltip: true
    },
    render(row) {
      return row.chunk_metadata ? JSON.stringify(row.chunk_metadata) : '-';
    }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 170,
    render(row) {
      return row.created_at ? dayjs(row.created_at).format('YYYY-MM-DD HH:mm') : '-';
    }
  },
  {
    title: '更新时间',
    key: 'updated_at',
    width: 170,
    render(row) {
      return row.updated_at ? dayjs(row.updated_at).format('YYYY-MM-DD HH:mm') : '-';
    }
  }
];

const scrollX = 80 + 360 + 100 + 100 + 180 + 180 + 170 + 170;

async function fetchChunks() {
  if (!Number.isFinite(datasetId.value)) {
    message.error('无效的数据集 ID');
    chunks.value = [];
    total.value = 0;
    return;
  }

  loading.value = true;

  try {
    const { response: res } = await getDatasetChunks(datasetId.value, page.value, pageSize.value);

    if (res?.data?.code === 0) {
      chunks.value = res.data.data?.chunks ?? [];
      total.value = res.data.data?.total ?? 0;
    } else {
      message.error(res?.data?.msg || '获取 Chunk 列表失败');
      chunks.value = [];
      total.value = 0;
    }
  } finally {
    loading.value = false;
  }
}

function handlePageSizeChange(size: number) {
  pageSize.value = size;
  page.value = 1;
  fetchChunks();
}

function triggerFileUpload() {
  fileInputRef.value?.click();
}

function getErrorMessage(error: unknown, fallback: string) {
  return (error as { response?: { data?: { msg?: string } } })?.response?.data?.msg || fallback;
}

async function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const files = input.files;

  if (!files || files.length === 0) {
    return;
  }

  if (!Number.isFinite(datasetId.value)) {
    message.error('无效的数据集 ID');
    input.value = '';
    return;
  }

  const fileList = Array.from(files);
  isUploading.value = true;

  try {
    const { data, error } = await uploadFiles(datasetId.value, fileList);

    if (!error && data?.code === 0) {
      uploadedFiles.value = data.data?.files ?? [];
      message.success(`成功上传 ${fileList.length} 个文件`);
    } else {
      message.error(getErrorMessage(error, '文件上传失败'));
    }
  } finally {
    input.value = '';
    isUploading.value = false;
  }
}

async function handleSplitAndEmbed(file: Api.Dataset.FileUploadInfo) {
  if (!Number.isFinite(datasetId.value)) {
    message.error('无效的数据集 ID');
    return;
  }

  isProcessing.value = true;

  try {
    // Step 1: Split the document
    const splitPayload: Api.Dataset.SplitDocReq = {
      file_id: file.id,
      dataset_id: datasetId.value,
      minio_path: `documents/${file.name}`,
      chunk_size: chunkSize.value,
      chunk_overlap: chunkOverlap.value
    };

    message.info('正在拆分文档...');
    const { data: splitData, error: splitError } = await splitDocument(splitPayload);

    if (splitError || splitData?.code !== 0) {
      message.error(getErrorMessage(splitError, '文档拆分失败'));
      return;
    }

    message.success(`文档拆分完成，生成 ${splitData.data?.chunks_count ?? 0} 个片段`);

    // Refresh chunk list
    await fetchChunks();
  } catch (err) {
    message.error('处理文件时发生错误');
  } finally {
    isProcessing.value = false;
  }
}

watch(
  () => props.id,
  () => {
    page.value = 1;
    fetchChunks();
  }
);

onMounted(() => {
  fetchChunks();
});
</script>

<template>
  <NSpace vertical :size="16">
    <div class="flex items-center justify-between">
      <div>
        <div class="text-18px font-600">Dataset Chunks</div>
        <div class="mt-4px text-12px text-gray-500">Dataset ID: {{ props.id }}</div>
      </div>
      <NSpace>
        <NButton :loading="loading" @click="fetchChunks">刷新</NButton>
        <NButton @click="router.back()">返回</NButton>
      </NSpace>
    </div>

    <!-- File Upload Card -->
    <NCard title="文件上传" :bordered="false" size="small">
      <NSpace vertical :size="12">
        <NSpace align="end">
          <div>
            <div class="mb-4px text-12px text-gray-500">Chunk Size</div>
            <NInputNumber v-model:value="chunkSize" :min="64" :max="4096" size="small" style="width: 140px" />
          </div>
          <div>
            <div class="mb-4px text-12px text-gray-500">Chunk Overlap</div>
            <NInputNumber v-model:value="chunkOverlap" :min="0" :max="2048" size="small" style="width: 140px" />
          </div>
          <NButton type="primary" :loading="isUploading" :disabled="isProcessing" @click="triggerFileUpload">
            {{ isUploading ? '上传中...' : '选择文件' }}
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
            <NButton size="small" type="primary" ghost :loading="isProcessing" @click="handleSplitAndEmbed(file)">
              拆分处理
            </NButton>
          </div>
        </div>
      </NSpace>

      <input ref="fileInputRef" type="file" multiple class="hidden" @change="handleFileChange" />
    </NCard>

    <NCard :bordered="false">
      <NDataTable
        :columns="columns"
        :data="chunks"
        :loading="loading"
        :bordered="false"
        size="small"
        remote
        :scroll-x="scrollX"
        :flex-height="true"
      />

      <template #footer>
        <NSpace justify="end">
          <NPagination
            v-model:page="page"
            :item-count="total"
            :page-size="pageSize"
            show-size-picker
            :page-sizes="[10, 20, 50, 100]"
            @update:page="fetchChunks"
            @update:page-size="handlePageSizeChange"
          />
        </NSpace>
      </template>
    </NCard>
  </NSpace>
</template>
