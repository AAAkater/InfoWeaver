<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns } from 'naive-ui';
import { NButton, NTag, useDialog, useMessage } from 'naive-ui';
import { createMcp, deleteMcp, getMcpList, updateMcp } from '@/service/api/mcp';
import McpFormModal from './modules/mcp-form-modal.vue';

const dialog = useDialog();
const message = useMessage();

const loading = ref(false);
const showModal = ref(false);
const isEdit = ref(false);
const mcpList = ref<Api.Mcp.McpInfo[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(20);

const model: Api.Mcp.McpUpdateReq = reactive({
  id: 0,
  name: '',
  transport: 'stdio',
  command: '',
  args: '',
  url: '',
  enabled: true,
  env_vars: {},
  headers: {}
});

function resetModel() {
  Object.assign(model, {
    id: 0,
    name: '',
    transport: 'stdio' as const,
    command: '',
    args: '',
    url: '',
    enabled: true,
    env_vars: {},
    headers: {}
  });
}

function fillModel(mcp: Api.Mcp.McpInfo) {
  Object.assign(model, {
    id: mcp.id,
    name: mcp.name,
    transport: (mcp.transport || 'stdio') as Api.Mcp.TransportType,
    command: mcp.command || '',
    args: mcp.args || '',
    url: mcp.url || '',
    enabled: mcp.enabled,
    env_vars: mcp.env_vars ?? {},
    headers: mcp.headers ?? {}
  });
}

function openCreateModal() {
  resetModel();
  isEdit.value = false;
  showModal.value = true;
}

async function fetchList() {
  loading.value = true;
  try {
    const { response: res } = await getMcpList(page.value, pageSize.value);
    if (res?.data?.code === 0) {
      mcpList.value = res.data.data?.mcps ?? [];
      total.value = res.data.data?.total ?? 0;
    } else {
      message.error(res?.data?.msg || '获取 MCP 列表失败');
      mcpList.value = [];
      total.value = 0;
    }
  } finally {
    loading.value = false;
  }
}

const transportTagMap: Record<string, 'info' | 'success' | 'warning'> = {
  stdio: 'info',
  sse: 'success',
  streamable_http: 'warning'
};

const columns: DataTableColumns<Api.Mcp.McpInfo> = [
  { title: 'ID', key: 'id', width: 80 },
  {
    title: '名称',
    key: 'name',
    minWidth: 160,
    ellipsis: { tooltip: true }
  },
  {
    title: '传输方式',
    key: 'transport',
    width: 140,
    render(row) {
      return h(
        NTag,
        { size: 'small', type: transportTagMap[row.transport] || 'default', bordered: false },
        { default: () => row.transport }
      );
    }
  },
  {
    title: '命令/URL',
    key: 'command',
    minWidth: 180,
    ellipsis: { tooltip: true },
    render(row) {
      return row.transport === 'stdio' ? row.command || '-' : row.url || '-';
    }
  },
  {
    title: '状态',
    key: 'enabled',
    width: 80,
    render(row) {
      return h(
        NTag,
        { size: 'small', type: row.enabled ? 'success' : 'default', bordered: false },
        { default: () => (row.enabled ? '启用' : '禁用') }
      );
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    render(row) {
      return h('div', { style: { display: 'flex', gap: '8px' } }, [
        h(NButton, { size: 'tiny', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'tiny', type: 'error', onClick: () => handleDelete(row) }, { default: () => '删除' })
      ]);
    }
  }
];

function handleEdit(mcp: Api.Mcp.McpInfo) {
  fillModel(mcp);
  isEdit.value = true;
  showModal.value = true;
}

function handleDelete(mcp: Api.Mcp.McpInfo) {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除 MCP 服务器 "${mcp.name}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      const { response: res } = await deleteMcp(mcp.id);
      if (res.data.code === 0) {
        message.success('删除成功');
        await fetchList();
      } else {
        message.error(res.data.msg || '删除失败');
      }
    }
  });
}

async function handleCreate(form: Api.Mcp.McpUpdateReq) {
  const { response: res } = await createMcp({
    name: form.name,
    transport: form.transport,
    command: form.command || undefined,
    args: form.args || undefined,
    url: form.url || undefined,
    enabled: form.enabled,
    env_vars: form.env_vars || undefined,
    headers: form.headers || undefined
  });
  if (res.data.code === 0) {
    message.success('创建成功');
    showModal.value = false;
    await fetchList();
  } else {
    message.error(res.data.msg || '创建失败');
  }
}

async function handleUpdate(form: Api.Mcp.McpUpdateReq) {
  const { response: res } = await updateMcp(form);
  if (res.data.code === 0) {
    message.success('更新成功');
    showModal.value = false;
    await fetchList();
  } else {
    message.error(res.data.msg || '更新失败');
  }
}

function handleSubmit(form: Api.Mcp.McpUpdateReq) {
  if (isEdit.value && form.id !== 0) {
    handleUpdate(form);
  } else {
    handleCreate(form);
  }
}

function handlePageSizeChange(size: number) {
  pageSize.value = size;
  page.value = 1;
  fetchList();
}

onMounted(() => {
  fetchList();
});
</script>

<template>
  <NSpace vertical :size="16">
    <div class="flex items-center justify-between">
      <div class="text-18px font-600">MCP 服务器</div>
      <NSpace>
        <NButton type="primary" @click="openCreateModal">创建 MCP 服务器</NButton>
      </NSpace>
    </div>

    <NCard :bordered="false">
      <div style="max-height: 60vh; overflow: auto">
        <NDataTable :columns="columns" :data="mcpList" :loading="loading" :bordered="false" size="small" />
      </div>
      <template #footer>
        <NSpace justify="end">
          <NPagination
            v-model:page="page"
            :item-count="total"
            :page-size="pageSize"
            show-size-picker
            :page-sizes="[10, 20, 50, 100]"
            @update:page="fetchList"
            @update:page-size="handlePageSizeChange"
          />
        </NSpace>
      </template>
    </NCard>

    <McpFormModal v-model:show="showModal" :model="model" :is-edit="isEdit" @submit="handleSubmit" />
  </NSpace>
</template>
