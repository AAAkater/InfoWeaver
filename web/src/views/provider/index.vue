<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue"
import { useDialog, useMessage } from "naive-ui"
import {
  createProvider,
  deleteProvider,
  getProviderList,
  getProviderModels,
  setProviderModelEnable,
  updateProvider,
} from "@/service/api/provider"
import ProviderAddModelModal from "./modules/provider-add-model-modal.vue"
import ProviderFormModal from "./modules/provider-form-modal.vue"
import ProviderModelsModal from "./modules/provider-models-modal.vue"
import ProviderSearch from "./modules/provider-search.vue"
import ProviderTable from "./modules/provider-table.vue"

const dialog = useDialog()
const message = useMessage()

const providers = ref<Api.Provider.ProviderInfo[]>([])
const searchKey = ref("")
const loading = ref(false)

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showModelsModal = ref(false)
const showAddModelModal = ref(false)

const currentProvider = ref<Api.Provider.ProviderInfo | null>(null)
const providerModels = ref<Api.Provider.ModelInfo[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const modelFilter = ref<"all" | "enabled" | "disabled">("all")

const createModel = reactive<Api.Provider.ProviderCreateReq>({
  api_key: "",
  base_url: "",
  mode: "openai",
  name: "",
})

const editModel = reactive<Api.Provider.ProviderUpdateReq>({
  id: 0,
  api_key: "",
  base_url: "",
  mode: "openai",
  name: "",
})

function resetCreateModel() {
  createModel.api_key = ""
  createModel.base_url = ""
  createModel.mode = "openai"
  createModel.name = ""
}

function openCreateModal() {
  resetCreateModel()
  showCreateModal.value = true
}

function getFilteredProviders() {
  if (!searchKey.value.trim()) {
    return providers.value
  }

  const keyword = searchKey.value.trim().toLowerCase()
  return providers.value.filter(
    (provider) =>
      provider.name.toLowerCase().includes(keyword) ||
      provider.base_url.toLowerCase().includes(keyword),
  )
}

const filteredProviders = computed(() => getFilteredProviders())

async function fetchProviders() {
  loading.value = true

  try {
    const { response: res } = await getProviderList()
    if (res?.data?.code === 0) {
      providers.value = res.data.data?.providers ?? []
    } else {
      providers.value = []
      message.error("获取 Provider 列表失败")
    }
  } finally {
    loading.value = false
  }
}

async function handleCreate(payload: Api.Provider.ProviderCreateReq) {
  const { response: res } = await createProvider(payload)
  if (res.data.code === 0) {
    message.success("创建成功")
    showCreateModal.value = false
    await fetchProviders()
  } else {
    message.error(res.data.msg || "创建失败")
  }
}

async function handleUpdate(payload: Api.Provider.ProviderUpdateReq) {
  const { response: res } = await updateProvider(payload)
  if (res.data.code === 0) {
    message.success("更新成功")
    showEditModal.value = false
    await fetchProviders()
  } else {
    message.error(res.data.msg || "更新失败")
  }
}

async function handleDelete(id: number) {
  const { response: res } = await deleteProvider(id)
  if (res.data.code === 0) {
    message.success("删除成功")
    await fetchProviders()
  } else {
    message.error(res.data.msg || "删除失败")
  }
}

async function handleGetModels(providerId: number) {
  const { response: res } = await getProviderModels(providerId)
  if (res.data.code === 0) {
    providerModels.value = res.data.data?.models ?? []
    currentPage.value = 1
  } else {
    providerModels.value = []
    message.error(res.data.msg || "获取模型列表失败")
  }
}

const paginatedModels = computed(() => {
  let filteredModels = providerModels.value
  if (modelFilter.value === "enabled") {
    filteredModels = filteredModels.filter((m) => m.enabled)
  } else if (modelFilter.value === "disabled") {
    filteredModels = filteredModels.filter((m) => !m.enabled)
  }

  const sortedModels = [...filteredModels].sort((a, b) => {
    if (a.enabled === b.enabled) return 0
    return a.enabled ? -1 : 1
  })

  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return sortedModels.slice(start, end)
})

const totalModelsCount = computed(() => {
  if (modelFilter.value === "enabled") {
    return providerModels.value.filter((m) => m.enabled).length
  }
  if (modelFilter.value === "disabled") {
    return providerModels.value.filter((m) => !m.enabled).length
  }
  return providerModels.value.length
})

async function handleSetModelEnable(modelId: string, enabled: boolean) {
  if (!currentProvider.value) return

  const { response: res } = await setProviderModelEnable({
    id: currentProvider.value.id,
    model_id: modelId,
    enabled,
  })
  if (res.data.code === 0) {
    message.success(enabled ? "模型已启用" : "模型已禁用")
    await handleGetModels(currentProvider.value.id)
  } else {
    message.error(res.data.msg || "设置失败")
  }
}

async function handleAddModel(modelName: string) {
  if (!currentProvider.value || !modelName) {
    message.warning("请输入模型名称")
    return
  }

  const { response: res } = await setProviderModelEnable({
    id: currentProvider.value.id,
    model_id: modelName,
    enabled: true,
  })
  if (res.data.code === 0) {
    message.success("添加模型成功")
    showAddModelModal.value = false
    await handleGetModels(currentProvider.value.id)
  } else {
    message.error(res.data.msg || "添加模型失败")
  }
}

function handleSelect(key: "edit" | "models" | "delete", provider: Api.Provider.ProviderInfo) {
  currentProvider.value = provider

  switch (key) {
    case "edit":
      editModel.id = provider.id
      editModel.name = provider.name
      editModel.base_url = provider.base_url
      editModel.mode = provider.mode as Api.Provider.ProviderMode
      editModel.api_key = ""
      showEditModal.value = true
      break
    case "models":
      modelFilter.value = "all"
      handleGetModels(provider.id)
      showModelsModal.value = true
      break
    case "delete":
      dialog.warning({
        title: "确认删除",
        content: `确定要删除 Provider "${provider.name}" 吗？此操作不可恢复。`,
        positiveText: "删除",
        negativeText: "取消",
        onPositiveClick: () => handleDelete(provider.id),
      })
      break
  }
}

function handlePageChange(value: number) {
  currentPage.value = value
}

function handlePageSizeChange(value: number) {
  pageSize.value = value
  currentPage.value = 1
}

watch(modelFilter, () => {
  currentPage.value = 1
})

onMounted(() => {
  fetchProviders()
})
</script>

<template>
  <NSpace vertical :size="16">
    <ProviderSearch v-model="searchKey" :loading="loading" @create="openCreateModal" />

    <NCard :bordered="false">
      <ProviderTable :providers="filteredProviders" @select="handleSelect" />
      <template #footer>
        <div style="color: #949494; font-size: 12px">
          共 {{ filteredProviders.length }} 个 Provider
        </div>
      </template>
    </NCard>

    <ProviderFormModal
      v-model:show="showCreateModal"
      :model="createModel"
      :is-edit="false"
      @submit="handleCreate"
    />

    <ProviderFormModal
      v-model:show="showEditModal"
      :model="editModel"
      :is-edit="true"
      @submit="handleUpdate"
    />

    <ProviderModelsModal
      v-model:show="showModelsModal"
      v-model:modelFilter="modelFilter"
      :loading="loading"
      :page="currentPage"
      :page-size="pageSize"
      :paginated-models="paginatedModels"
      :provider-name="currentProvider?.name"
      :total-models-count="totalModelsCount"
      @add-model="showAddModelModal = true"
      @toggle-model="handleSetModelEnable"
      @update:page="handlePageChange"
      @update:pageSize="handlePageSizeChange"
    />

    <ProviderAddModelModal v-model:show="showAddModelModal" @submit="handleAddModel" />
  </NSpace>
</template>
