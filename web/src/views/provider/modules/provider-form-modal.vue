<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';

type ProviderFormModel = {
  id?: number;
  api_key: string;
  base_url: string;
  mode: Api.Provider.ProviderMode;
  name: string;
};

const props = defineProps<{
  isEdit?: boolean;
  model: ProviderFormModel;
  show: boolean;
}>();

const emit = defineEmits<{
  (event: 'submit', model: ProviderFormModel): void;
  (event: 'update:show', value: boolean): void;
}>();

const bodyStyle = {
  width: '500px'
};

const modeOptions = [
  { label: 'OpenAI', value: 'openai' },
  { label: 'OpenAI Response', value: 'openai_response' },
  { label: 'Gemini', value: 'gemini' },
  { label: 'Anthropic', value: 'anthropic' },
  { label: 'Ollama', value: 'ollama' }
];

const visible = computed({
  get: () => props.show,
  set: value => emit('update:show', value)
});

const form = reactive<ProviderFormModel>({
  id: undefined,
  api_key: '',
  base_url: '',
  mode: 'openai',
  name: ''
});

const apiKeyDisplay = ref('');

function resetForm() {
  Object.assign(form, {
    id: undefined,
    api_key: '',
    base_url: '',
    mode: 'openai',
    name: ''
  });
  apiKeyDisplay.value = '';
}

function syncForm() {
  Object.assign(form, {
    id: props.model.id,
    base_url: props.model.base_url || '',
    mode: props.model.mode || 'openai',
    name: props.model.name || ''
  });
  apiKeyDisplay.value = props.isEdit ? '********' : '';
  form.api_key = props.isEdit ? '' : props.model.api_key || '';
}

watch(
  () => props.show,
  show => {
    if (show) {
      syncForm();
    } else {
      resetForm();
    }
  }
);

watch(
  () => props.model,
  () => {
    if (props.show) {
      syncForm();
    }
  },
  { deep: true }
);

function handlePositiveClick() {
  emit('submit', {
    ...form,
    api_key: props.isEdit ? (apiKeyDisplay.value === '********' ? '' : apiKeyDisplay.value) : apiKeyDisplay.value
  });
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
    <NSpace :size="16" vertical>
      <NFormItem label="名称" required>
        <NInput v-model:value="form.name" placeholder="请输入 Provider 名称" />
      </NFormItem>
      <NFormItem label="Base URL" required>
        <NInput v-model:value="form.base_url" placeholder="请输入 API Base URL" />
      </NFormItem>
      <NFormItem label="API Key" required>
        <NInput
          v-model:value="apiKeyDisplay"
          type="password"
          :placeholder="props.isEdit ? '默认已隐藏，修改时直接输入新值' : '请输入 API Key'"
          show-password-on="click"
        />
      </NFormItem>
      <NFormItem label="模式" required>
        <NSelect v-model:value="form.mode" :options="modeOptions" />
      </NFormItem>
    </NSpace>
  </NModal>
</template>
