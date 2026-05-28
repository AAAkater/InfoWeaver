<script setup lang="ts">
import { computed, reactive, watch } from "vue"

type ProviderFormModel = {
  id?: number
  api_key?: string
  base_url: string
  mode: Api.Provider.ProviderMode
  name: string
}

const props = defineProps<{
  isEdit?: boolean
  model: ProviderFormModel
  show: boolean
}>()

const emit = defineEmits<{
  (event: "submit", model: ProviderFormModel): void
  (event: "update:show", value: boolean): void
}>()

const bodyStyle = {
  width: "500px",
}

const modeOptions = [
  { label: "OpenAI", value: "openai" },
  { label: "OpenAI Response", value: "openai_response" },
  { label: "Gemini", value: "gemini" },
  { label: "Anthropic", value: "anthropic" },
  { label: "Ollama", value: "ollama" },
]

const visible = computed({
  get: () => props.show,
  set: (value) => emit("update:show", value),
})

const form = reactive<ProviderFormModel>({
  id: undefined,
  api_key: "",
  base_url: "",
  mode: "openai",
  name: "",
})

function resetForm() {
  Object.assign(form, {
    id: undefined,
    api_key: "",
    base_url: "",
    mode: "openai",
    name: "",
  })
}

function syncForm() {
  Object.assign(form, {
    id: props.model.id,
    api_key: "",
    base_url: props.model.base_url || "",
    mode: props.model.mode || "openai",
    name: props.model.name || "",
  })
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

function handlePositiveClick() {
  emit("submit", { ...form })
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
    :title="props.isEdit ? '编辑 Provider' : '创建 Provider'"
    @positive-click="handlePositiveClick"
  >
    <div class="flex flex-col gap-4">
      <!-- 名称 -->
      <div class="flex items-center gap-4">
        <span class="w-80px shrink-0 text-14px text-gray-600">名称</span>
        <NInput
          v-model:value="form.name"
          size="small"
          class="flex-1"
          placeholder="请输入 Provider 名称"
        />
      </div>

      <!-- Base URL -->
      <div class="flex items-center gap-4">
        <span class="w-80px shrink-0 text-14px text-gray-600">Base URL</span>
        <NInput
          v-model:value="form.base_url"
          size="small"
          class="flex-1"
          placeholder="请输入 API Base URL"
        />
      </div>

      <!-- API Key（仅创建时显示） -->
      <div v-if="!props.isEdit" class="flex items-center gap-4">
        <span class="w-80px shrink-0 text-14px text-gray-600">API Key</span>
        <NInput
          v-model:value="form.api_key"
          type="password"
          size="small"
          class="flex-1"
          placeholder="请输入 API Key"
          show-password-on="click"
        />
      </div>

      <!-- 模式 -->
      <div class="flex items-center gap-4">
        <span class="w-80px shrink-0 text-14px text-gray-600">模式</span>
        <NSelect v-model:value="form.mode" :options="modeOptions" size="small" class="flex-1" />
      </div>
    </div>
  </NModal>
</template>
