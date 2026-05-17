<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';

const props = defineProps<{
  isEdit?: boolean;
  model: Api.Mcp.McpUpdateReq;
  show: boolean;
}>();

const emit = defineEmits<{
  (event: 'submit', model: Api.Mcp.McpUpdateReq): void;
  (event: 'update:show', value: boolean): void;
}>();

const transportOptions = [
  { label: 'stdio', value: 'stdio' },
  { label: 'SSE', value: 'sse' },
  { label: 'Streamable HTTP', value: 'streamable_http' }
];

const bodyStyle = { width: '600px' };

const visible = computed({
  get: () => props.show,
  set: value => emit('update:show', value)
});

const form = reactive<Api.Mcp.McpUpdateReq>({
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

// env_vars editing
const envVarsText = ref('');

// headers editing
const headersText = ref('');

function resetForm() {
  Object.assign(form, {
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
  envVarsText.value = '';
  headersText.value = '';
}

function syncForm() {
  Object.assign(form, {
    id: props.model.id ?? 0,
    name: props.model.name || '',
    transport: props.model.transport || 'stdio',
    command: props.model.command || '',
    args: props.model.args || '',
    url: props.model.url || '',
    enabled: props.model.enabled ?? true,
    env_vars: props.model.env_vars ?? {},
    headers: props.model.headers ?? {}
  });
  envVarsText.value = JSON.stringify(props.model.env_vars ?? {}, null, 2);
  headersText.value = JSON.stringify(props.model.headers ?? {}, null, 2);
}

watch(
  () => props.show,
  show => {
    if (show) syncForm();
    else resetForm();
  }
);

watch(
  () => props.model,
  () => {
    if (props.show) syncForm();
  },
  { deep: true }
);

function parseJsonField(text: string, fallback: Record<string, any>) {
  if (!text.trim()) return {};
  try {
    return JSON.parse(text);
  } catch {
    return fallback;
  }
}

function handlePositiveClick() {
  form.env_vars = parseJsonField(envVarsText.value, form.env_vars ?? {});
  form.headers = parseJsonField(headersText.value, form.headers ?? {});
  emit('submit', { ...form });
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
    :title="props.isEdit ? '编辑 MCP 服务器' : '创建 MCP 服务器'"
    @positive-click="handlePositiveClick"
  >
    <NSpace :size="10" vertical>
      <NCard title="名称" :bordered="false" size="small">
        <NInput
          v-model:value="form.name"
          size="tiny"
          style="background-color: #f1f3f6"
          placeholder="请输入 MCP 服务器名称"
        />
      </NCard>

      <NCard title="传输方式" :bordered="false" size="small">
        <NSelect
          v-model:value="form.transport"
          :options="transportOptions"
          size="tiny"
          style="background-color: #f1f3f6"
        />
      </NCard>

      <NCard v-if="form.transport === 'stdio'" title="命令" :bordered="false" size="small">
        <NInput
          v-model:value="form.command"
          size="tiny"
          style="background-color: #f1f3f6"
          placeholder="例如: npx, python, node"
        />
      </NCard>

      <NCard v-if="form.transport === 'stdio'" title="参数" :bordered="false" size="small">
        <NInput
          v-model:value="form.args"
          size="tiny"
          style="background-color: #f1f3f6"
          placeholder="例如: -m mcp-server-time"
        />
      </NCard>

      <NCard v-if="form.transport !== 'stdio'" title="URL" :bordered="false" size="small">
        <NInput
          v-model:value="form.url"
          size="tiny"
          style="background-color: #f1f3f6"
          placeholder="请输入 MCP 服务器 URL"
        />
      </NCard>

      <NCard title="环境变量 (JSON)" :bordered="false" size="small">
        <NInput
          v-model:value="envVarsText"
          type="textarea"
          size="tiny"
          :autosize="{ minRows: 2, maxRows: 6 }"
          style="background-color: #f1f3f6"
          placeholder='{"KEY": "value"}'
        />
      </NCard>

      <NCard title="Headers (JSON)" :bordered="false" size="small">
        <NInput
          v-model:value="headersText"
          type="textarea"
          size="tiny"
          :autosize="{ minRows: 2, maxRows: 6 }"
          style="background-color: #f1f3f6"
          placeholder='{"Authorization": "Bearer xxx"}'
        />
      </NCard>

      <NCard title="启用" :bordered="false" size="small">
        <NSwitch v-model:value="form.enabled" />
      </NCard>
    </NSpace>
  </NModal>
</template>
