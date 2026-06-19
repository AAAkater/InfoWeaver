<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue"
import { getProviderList, getProviderModels } from "@/service/api/provider"

const props = defineProps<{
  isEdit?: boolean
  model: Api.Dataset.FormModel
  show: boolean
}>()

const emit = defineEmits<{
  (event: "submit", model: Api.Dataset.FormModel): void
  (event: "update:show", value: boolean): void
}>()

const bodyStyle = {
  width: "600px",
}

const searchTypeOptions = [
  { label: "关键词检索", value: "sparse" },
  { label: "语义检索", value: "dense" },
  { label: "混合检索", value: "hybrid" },
]

const emojiList = ["😀", "😂", "😎", "🤖", "🐱", "🦊", "🐶", "🦄", "🐼", "🦉"]
const hovered = ref<string | null>(null)

// Provider & models
const providers = ref<Api.Provider.ProviderInfo[]>([])
const providerModels = ref<Api.Provider.ModelInfo[]>([])
const modelsLoading = ref(false)
const providerLoading = ref(false)

const visible = computed({
  get: () => props.show,
  set: (value) => emit("update:show", value),
})

const form = reactive<Api.Dataset.FormModel>({
  id: undefined,
  icon: "🤖",
  description: "",
  name: "",
  search_type: "hybrid",
  embedding_model: "",
  provider_id: undefined,
})

function resetForm() {
  Object.assign(form, {
    id: undefined,
    icon: "🤖",
    description: "",
    name: "",
    search_type: "hybrid",
    embedding_model: "",
    provider_id: undefined,
  })
  providerModels.value = []
}

async function syncForm() {
  Object.assign(form, {
    id: props.model.id,
    icon: props.model.icon || "🤖",
    description: props.model.description || "",
    name: props.model.name || "",
    search_type: props.model.search_type || "hybrid",
    embedding_model: props.model.embedding_model || "",
    provider_id: props.model.provider_id || undefined,
  })
  // Fetch providers list
  fetchProviders()
  // If editing with an existing provider, fetch its models
  if (props.model.provider_id) {
    fetchModels(props.model.provider_id)
  }
}

watch(
  () => props.show,
  (show) => {
    if (show) {
      syncForm()
    } else {
      resetForm()
    }
  },
)

watch(
  () => props.model,
  () => {
    if (props.show) {
      syncForm()
    }
  },
  { deep: true },
)

function selectEmoji(emoji: string) {
  form.icon = emoji
}

function handlePositiveClick() {
  emit("submit", { ...form })
}

// Provider & models
async function fetchProviders() {
  providerLoading.value = true
  try {
    const { response: res } = await getProviderList()
    if (res?.data?.code === 0) {
      providers.value = res.data.data?.providers ?? []
    }
  } finally {
    providerLoading.value = false
  }
}

async function fetchModels(providerId: number) {
  if (!providerId) {
    providerModels.value = []
    return
  }
  modelsLoading.value = true
  try {
    const { response: res } = await getProviderModels(providerId)
    if (res?.data?.code === 0) {
      providerModels.value = res.data.data?.models ?? []
    }
  } finally {
    modelsLoading.value = false
  }
}

function onProviderChange(providerId: number) {
  form.provider_id = providerId
  form.embedding_model = "" // reset model when provider changes
  fetchModels(providerId)
}
</script>

<template>
  <NModal
    v-model:show="visible"
    :mask-closable="false"
    preset="dialog"
    :show-icon="false"
    :style="bodyStyle"
    positive-text="确认"
    negative-text="取消"
    title="知识库设置"
    @positive-click="handlePositiveClick"
  >
    <div class="flex flex-col gap-4">
      <!-- 知识库名称 -->
      <div class="flex items-start gap-4">
        <span class="w-80px shrink-0 pt-6px text-14px text-gray-600">知识库名称</span>
        <div class="flex flex-1 items-center gap-2">
          <NPopover trigger="click" placement="bottom-start">
            <template #trigger>
              <NAvatar
                :style="{ color: 'black', backgroundColor: '#FFEAD5', cursor: 'pointer' }"
                :size="36"
              >
                {{ form.icon }}
              </NAvatar>
            </template>
            <div class="flex gap-2 p-2">
              <span
                v-for="emoji in emojiList"
                :key="emoji"
                class="cursor-pointer rounded p-1 text-20px transition-colors"
                :class="hovered === emoji ? 'bg-gray-200' : ''"
                @click="selectEmoji(emoji)"
                @mouseenter="hovered = emoji"
                @mouseleave="hovered = null"
              >
                {{ emoji }}
              </span>
            </div>
          </NPopover>
          <NInput
            v-model:value="form.name"
            class="flex-1"
            size="small"
            placeholder="请输入知识库名称"
          />
        </div>
      </div>

      <!-- 描述 -->
      <div class="flex items-start gap-4">
        <span class="w-80px shrink-0 pt-6px text-14px text-gray-600">描述</span>
        <NInput
          v-model:value="form.description"
          type="textarea"
          size="small"
          class="flex-1 resize-y"
          :autosize="{ minRows: 2, maxRows: 10 }"
          placeholder="描述该数据集的内容。详细描述可以让AI更快地访问数据集的内容。"
        />
      </div>

      <!-- 检索方式 -->
      <div class="flex items-center gap-4">
        <span class="w-80px shrink-0 text-14px text-gray-600">检索方式</span>
        <NSelect
          v-model:value="form.search_type"
          :options="searchTypeOptions"
          size="small"
          class="flex-1"
        />
      </div>

      <!-- Provider & Embedding 模型 -->
      <div class="flex items-center gap-4">
        <span class="w-80px shrink-0 text-14px text-gray-600">模型提供商</span>
        <NSelect
          v-model:value="form.provider_id"
          :options="providers.map((p) => ({ label: p.name, value: p.id }))"
          :loading="providerLoading"
          size="small"
          class="flex-1"
          placeholder="选择模型提供商"
          @update:value="onProviderChange"
        />
      </div>
      <div class="flex items-center gap-4">
        <span class="w-80px shrink-0 text-14px text-gray-600">Embedding 模型</span>
        <NSelect
          v-model:value="form.embedding_model"
          :options="providerModels.map((m) => ({ label: m.id, value: m.id }))"
          :loading="modelsLoading"
          size="small"
          class="flex-1"
          placeholder="选择 Embedding 模型"
          filterable
        />
      </div>
    </div>
  </NModal>
</template>
