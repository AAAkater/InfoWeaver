<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, watch } from 'vue';
import type { Component } from 'vue';
import { NIcon, NSwitch, useDialog, useMessage } from 'naive-ui';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';
import {
  Add12Filled,
  Delete16Regular as DeleteIcon,
  Edit24Regular as EditIcon,
  MoreHorizontal28Regular,
  Search48Filled,
  Server24Regular
} from '@vicons/fluent';
import {
  createProvider,
  deleteProvider,
  getProviderList,
  getProviderModels,
  setProviderModelEnable,
  updateProvider
} from '@/service/api/provider';

const dialog = useDialog();
const message = useMessage();
dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

// Provider list
const providers = ref<Api.Provider.ProviderInfo[]>([]);
const searchKey = ref('');

// Modal states
const showCreateModal = ref(false);
const showEditModal = ref(false);
const showModelsModal = ref(false);
const showAddModelModal = ref(false);

// Current provider for operations
const currentProvider = ref<Api.Provider.ProviderInfo | null>(null);
const providerModels = ref<Api.Provider.ModelInfo[]>([]);
const newModelName = ref('');

// Pagination
const currentPage = ref(1);
const pageSize = ref(10);

// Filter state
const modelFilter = ref<'all' | 'enabled' | 'disabled'>('all');

// Form model
const formModel = reactive<Api.Provider.ProviderCreateReq>({
  api_key: '',
  base_url: '',
  mode: 'openai',
  name: ''
});

const editModel = reactive<Api.Provider.ProviderUpdateReq>({
  id: 0,
  api_key: '',
  base_url: '',
  mode: 'openai',
  name: ''
});

// Mode options
const modeOptions = [
  { label: 'OpenAI', value: 'openai' },
  { label: 'OpenAI Response', value: 'openai_response' },
  { label: 'Gemini', value: 'gemini' },
  { label: 'Anthropic', value: 'anthropic' },
  { label: 'Ollama', value: 'ollama' }
];

// Render icon helper
function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

// Dropdown options
const dropdownOptions = [
  {
    label: '编辑',
    key: 'edit',
    icon: renderIcon(EditIcon)
  },
  {
    label: '模型列表',
    key: 'models',
    icon: renderIcon(Server24Regular)
  },
  {
    label: '删除',
    key: 'delete',
    icon: renderIcon(DeleteIcon)
  }
];

// Fetch providers
async function fetchProviders() {
  const { response: res } = await getProviderList();
  if (res?.data?.code === 0) {
    providers.value = res.data.data?.providers ?? [];
  } else {
    message.error('获取 Provider 列表失败');
    providers.value = [];
  }
}

// Create provider
async function handleCreate() {
  const { response: res } = await createProvider(formModel);
  if (res.data.code === 0) {
    message.success('创建成功');
    showCreateModal.value = false;
    resetCreateForm();
    fetchProviders();
  } else {
    message.error(res.data.msg || '创建失败');
  }
}

// Update provider
async function handleUpdate() {
  const { response: res } = await updateProvider(editModel);
  if (res.data.code === 0) {
    message.success('更新成功');
    showEditModal.value = false;
    fetchProviders();
  } else {
    message.error(res.data.msg || '更新失败');
  }
}

// Delete provider
async function handleDelete(id: number) {
  const { response: res } = await deleteProvider(id);
  if (res.data.code === 0) {
    message.success('删除成功');
    fetchProviders();
  } else {
    message.error(res.data.msg || '删除失败');
  }
}

// Get provider models
async function handleGetModels(providerId: number) {
  const { response: res } = await getProviderModels(providerId);
  if (res.data.code === 0) {
    providerModels.value = res.data.data?.models ?? [];
    currentPage.value = 1; // Reset to first page
  } else {
    message.error(res.data.msg || '获取模型列表失败');
    providerModels.value = [];
  }
}

// Computed paginated models with sorting (enabled models first) and filtering
const paginatedModels = computed(() => {
  // Filter models based on filter state
  let filteredModels = providerModels.value;
  if (modelFilter.value === 'enabled') {
    filteredModels = providerModels.value.filter(m => m.enabled);
  } else if (modelFilter.value === 'disabled') {
    filteredModels = providerModels.value.filter(m => !m.enabled);
  }

  // Sort models: enabled models first, then disabled models
  const sortedModels = [...filteredModels].sort((a, b) => {
    // If both are enabled or both are disabled, keep original order
    if (a.enabled === b.enabled) return 0;
    // Enabled models come first
    return a.enabled ? -1 : 1;
  });

  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return sortedModels.slice(start, end);
});

// Computed total count for pagination
const totalModelsCount = computed(() => {
  if (modelFilter.value === 'enabled') {
    return providerModels.value.filter(m => m.enabled).length;
  } else if (modelFilter.value === 'disabled') {
    return providerModels.value.filter(m => !m.enabled).length;
  }
  return providerModels.value.length;
});

// Set model enable status
async function handleSetModelEnable(modelId: string, enabled: boolean) {
  if (!currentProvider.value) return;

  const { response: res } = await setProviderModelEnable({
    id: currentProvider.value.id,
    model_id: modelId,
    enabled
  });
  if (res.data.code === 0) {
    message.success(enabled ? '模型已启用' : '模型已禁用');
    // Refresh models list
    handleGetModels(currentProvider.value.id);
  } else {
    message.error(res.data.msg || '设置失败');
  }
}

// Add single model to provider
async function handleAddModel() {
  if (!currentProvider.value || !newModelName.value.trim()) {
    message.warning('请输入模型名称');
    return;
  }
  const { response: res } = await setProviderModelEnable({
    id: currentProvider.value.id,
    model_id: newModelName.value.trim(),
    enabled: true // 默认启用
  });
  if (res.data.code === 0) {
    message.success('添加模型成功');
    newModelName.value = '';
    showAddModelModal.value = false;
    // Refresh models list
    handleGetModels(currentProvider.value.id);
  } else {
    message.error(res.data.msg || '添加模型失败');
  }
}

// Open add model modal
function openAddModelModal() {
  newModelName.value = '';
  showAddModelModal.value = true;
}

// Reset create form
function resetCreateForm() {
  formModel.api_key = '';
  formModel.base_url = '';
  formModel.mode = 'openai';
  formModel.name = '';
}

// Open create modal
function openCreateModal() {
  resetCreateForm();
  showCreateModal.value = true;
}

// Handle dropdown select
function handleSelect(key: string, provider: Api.Provider.ProviderInfo) {
  currentProvider.value = provider;
  switch (key) {
    case 'edit':
      editModel.id = provider.id;
      editModel.name = provider.name;
      editModel.base_url = provider.base_url;
      editModel.mode = provider.mode as Api.Provider.ProviderMode;
      editModel.api_key = '';
      showEditModal.value = true;
      break;
    case 'models':
      modelFilter.value = 'all'; // Reset filter
      handleGetModels(provider.id);
      showModelsModal.value = true;
      break;
    case 'delete':
      dialog.warning({
        title: '确认删除',
        content: `确定要删除 Provider "${provider.name}" 吗？此操作不可恢复。`,
        positiveText: '删除',
        negativeText: '取消',
        onPositiveClick: () => handleDelete(provider.id)
      });
      break;
    default:
  }
}

// Filter providers by search key
const filteredProviders = ref<Api.Provider.ProviderInfo[]>([]);

// Watch providers and searchKey to update filteredProviders
watch(
  [providers, searchKey],
  () => {
    if (!searchKey.value || searchKey.value.trim() === '') {
      filteredProviders.value = providers.value;
    } else {
      const keyword = searchKey.value.trim().toLowerCase();
      filteredProviders.value = providers.value.filter(
        p => p.name.toLowerCase().includes(keyword) || p.base_url.toLowerCase().includes(keyword)
      );
    }
  },
  { immediate: true }
);

// Watch modelFilter to reset pagination
watch(modelFilter, () => {
  currentPage.value = 1;
});

onMounted(() => {
  fetchProviders();
});
</script>

<template>
  <NSpace vertical :size="16">
    <!-- Header with search and create button -->
    <NSpace justify="space-between">
      <NInput v-model:value="searchKey" round placeholder="搜索 Provider 名称或 URL" clearable style="width: 300px">
        <template #prefix>
          <NIcon :component="Search48Filled" />
        </template>
      </NInput>
      <NButton type="primary" @click="openCreateModal">
        <NIcon :component="Add12Filled" class="mr-2" />
        创建 Provider
      </NButton>
    </NSpace>

    <!-- Provider list -->
    <NGrid cols="3" x-gap="16" y-gap="16" responsive="screen" item-responsive>
      <NGridItem v-for="provider in filteredProviders" :key="provider.id" class="grid-item-equal-height">
        <NCard hoverable size="huge" style="cursor: pointer">
          <NSpace vertical>
            <div style="display: flex; align-items: flex-start; gap: 8px; width: 100%">
              <NAvatar
                size="large"
                :style="{
                  color: 'black',
                  backgroundColor: '#E0F2FE'
                }"
              >
                <NIcon :component="Server24Regular" />
              </NAvatar>

              <div style="flex: 1">
                <div style="font-weight: bold">{{ provider.name }}</div>
                <div style="color: #949494; font-size: 12px">{{ provider.mode }}</div>
              </div>

              <NDropdown
                :options="dropdownOptions"
                trigger="click"
                size="small"
                @select="key => handleSelect(key, provider)"
              >
                <NButton size="small" secondary>
                  <NIcon>
                    <MoreHorizontal28Regular />
                  </NIcon>
                </NButton>
              </NDropdown>
            </div>
            <div style="color: #666; font-size: 12px; line-height: 1.4">
              {{ provider.base_url }}
            </div>
            <div style="color: #949494; font-size: 10px">
              创建于 · {{ dayjs(provider.created_at).format('YYYY-MM-DD HH:mm') }}
            </div>
          </NSpace>
        </NCard>
      </NGridItem>
    </NGrid>

    <!-- Create Modal -->
    <NModal
      v-model:show="showCreateModal"
      :mask-closable="false"
      preset="dialog"
      :show-icon="false"
      style="width: 500px"
    >
      <template #header>
        <div style="font-weight: bold">创建 Provider</div>
      </template>
      <NSpace :size="16" vertical>
        <NFormItem label="名称" required>
          <NInput v-model:value="formModel.name" placeholder="请输入 Provider 名称" />
        </NFormItem>
        <NFormItem label="Base URL" required>
          <NInput v-model:value="formModel.base_url" placeholder="请输入 API Base URL" />
        </NFormItem>
        <NFormItem label="API Key" required>
          <NInput
            v-model:value="formModel.api_key"
            type="password"
            placeholder="请输入 API Key"
            show-password-on="click"
          />
        </NFormItem>
        <NFormItem label="模式" required>
          <NSelect v-model:value="formModel.mode" :options="modeOptions" />
        </NFormItem>
      </NSpace>
      <template #action>
        <NSpace justify="end">
          <NButton @click="showCreateModal = false">取消</NButton>
          <NButton type="primary" @click="handleCreate">创建</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- Edit Modal -->
    <NModal v-model:show="showEditModal" :mask-closable="false" preset="dialog" :show-icon="false" style="width: 500px">
      <template #header>
        <div style="font-weight: bold">编辑 Provider</div>
      </template>
      <NSpace :size="16" vertical>
        <NFormItem label="名称" required>
          <NInput v-model:value="editModel.name" placeholder="请输入 Provider 名称" />
        </NFormItem>
        <NFormItem label="Base URL" required>
          <NInput v-model:value="editModel.base_url" placeholder="请输入 API Base URL" />
        </NFormItem>
        <NFormItem label="API Key" required>
          <NInput
            v-model:value="editModel.api_key"
            type="password"
            placeholder="请输入新的 API Key（如不修改可留空）"
            show-password-on="click"
          />
        </NFormItem>
        <NFormItem label="模式" required>
          <NSelect v-model:value="editModel.mode" :options="modeOptions" />
        </NFormItem>
      </NSpace>
      <template #action>
        <NSpace justify="end">
          <NButton @click="showEditModal = false">取消</NButton>
          <NButton type="primary" @click="handleUpdate">保存</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- Models Modal -->
    <NModal
      v-model:show="showModelsModal"
      :mask-closable="true"
      preset="dialog"
      :show-icon="false"
      style="width: 600px"
    >
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center; width: 100%">
          <div style="font-weight: bold">{{ currentProvider?.name }} - 模型列表</div>
          <NSpace :size="8">
            <NSelect
              v-model:value="modelFilter"
              :options="[
                { label: '全部', value: 'all' },
                { label: '已启用', value: 'enabled' },
                { label: '已禁用', value: 'disabled' }
              ]"
              size="small"
              style="width: 100px"
            />
            <NButton type="primary" size="small" @click="openAddModelModal">
              <template #icon>
                <NIcon :component="Add12Filled" />
              </template>
              添加模型
            </NButton>
          </NSpace>
        </div>
      </template>
      <NDataTable
        :columns="[
          { title: '模型名', key: 'id' },
          { title: '类型', key: 'object' },
          { title: '所属', key: 'owned_by' },
          {
            title: '状态',
            key: 'enabled',
            render(row) {
              return h(NSwitch, {
                value: row.enabled,
                onUpdateValue: (value: boolean) => handleSetModelEnable(row.id, value)
              });
            }
          }
        ]"
        :data="paginatedModels"
        :bordered="false"
        size="small"
      />
      <template #action>
        <NSpace justify="space-between" align="center">
          <NPagination
            v-model:page="currentPage"
            :page-count="Math.ceil(totalModelsCount / pageSize)"
            :page-size="pageSize"
            show-size-picker
            :page-sizes="[10, 20, 50]"
            @update:page-size="
              (size: number) => {
                pageSize = size;
                currentPage = 1;
              }
            "
          />
          <NButton @click="showModelsModal = false">关闭</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- Add Model Modal -->
    <NModal
      v-model:show="showAddModelModal"
      :mask-closable="false"
      preset="dialog"
      :show-icon="false"
      style="width: 400px"
    >
      <template #header>
        <div style="font-weight: bold">添加模型</div>
      </template>
      <NSpace :size="16" vertical>
        <NFormItem label="模型名称" required>
          <NInput v-model:value="newModelName" placeholder="请输入模型名称" />
        </NFormItem>
      </NSpace>
      <template #action>
        <NSpace justify="end">
          <NButton @click="showAddModelModal = false">取消</NButton>
          <NButton type="primary" @click="handleAddModel">添加</NButton>
        </NSpace>
      </template>
    </NModal>
  </NSpace>
</template>

<style scoped>
.grid-item-equal-height {
  display: flex;
  min-height: 120px;
}
</style>
