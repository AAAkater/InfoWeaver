<script setup lang="ts">
import { computed, h, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NTag, useMessage } from 'naive-ui';
import dayjs from 'dayjs';
import { getDatasetChunks } from '@/service/api/dataset';

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
