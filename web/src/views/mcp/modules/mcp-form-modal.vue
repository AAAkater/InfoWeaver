<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue"
import { NIcon } from "naive-ui"
import {
  ArrowRouting24Regular as RouteIcon,
  Braces24Regular as BracesIcon,
  PlugDisconnected24Regular as PlugIcon,
  Server20Regular as ServerIcon,
  TextHeader124Regular as HeaderIcon,
  Wrench24Regular as WrenchIcon,
} from "@vicons/fluent"

const props = defineProps<{
  isEdit?: boolean
  model: Api.Mcp.McpUpdateReq
  show: boolean
}>()

const emit = defineEmits<{
  (event: "submit", model: Api.Mcp.McpUpdateReq): void
  (event: "update:show", value: boolean): void
}>()

const transportOptions = [
  { label: "stdio", value: "stdio" },
  { label: "SSE", value: "sse" },
  { label: "Streamable HTTP", value: "streamable_http" },
]

const bodyStyle = { width: "600px" }

const visible = computed({
  get: () => props.show,
  set: (value) => emit("update:show", value),
})

const form = reactive<Api.Mcp.McpUpdateReq>({
  id: 0,
  name: "",
  transport: "stdio",
  command: "",
  args: "",
  url: "",
  enabled: true,
  env_vars: {},
  headers: {},
})

// env_vars editing
const envVarsText = ref("")

// headers editing
const headersText = ref("")

function resetForm() {
  Object.assign(form, {
    id: 0,
    name: "",
    transport: "stdio" as const,
    command: "",
    args: "",
    url: "",
    enabled: true,
    env_vars: {},
    headers: {},
  })
  envVarsText.value = ""
  headersText.value = ""
}

function syncForm() {
  Object.assign(form, {
    id: props.model.id ?? 0,
    name: props.model.name || "",
    transport: props.model.transport || "stdio",
    command: props.model.command || "",
    args: props.model.args || "",
    url: props.model.url || "",
    enabled: props.model.enabled ?? true,
    env_vars: props.model.env_vars ?? {},
    headers: props.model.headers ?? {},
  })
  envVarsText.value = JSON.stringify(props.model.env_vars ?? {}, null, 2)
  headersText.value = JSON.stringify(props.model.headers ?? {}, null, 2)
}

watch(
  () => props.show,
  (show) => {
    if (show) syncForm()
    else resetForm()
  },
)

watch(
  () => props.model,
  () => {
    if (props.show) syncForm()
  },
  { deep: true },
)

function parseJsonField(text: string, fallback: Record<string, any>) {
  if (!text.trim()) return {}
  try {
    return JSON.parse(text)
  } catch {
    return fallback
  }
}

function handlePositiveClick() {
  form.env_vars = parseJsonField(envVarsText.value, form.env_vars ?? {})
  form.headers = parseJsonField(headersText.value, form.headers ?? {})
  emit("submit", { ...form })
}
</script>

<template>
  <NModal
    v-model:show="visible"
    :mask-closable="false"
    preset="card"
    :show-icon="false"
    :style="bodyStyle"
    :title="props.isEdit ? '编辑 MCP 服务器' : '创建 MCP 服务器'"
    class="mcp-form-modal"
  >
    <template #header>
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <NIcon :component="PlugIcon" size="20" color="#7c3aed" />
          <span class="font-medium">{{
            props.isEdit ? "编辑 MCP 服务器" : "创建 MCP 服务器"
          }}</span>
        </div>
        <NSwitch v-model:value="form.enabled" size="small">
          <template #checked>启用</template>
          <template #unchecked>禁用</template>
        </NSwitch>
      </div>
    </template>

    <div class="flex flex-col gap-4">
      <!-- 基本信息区 -->
      <div class="rounded-xl border border-gray-100 bg-gray-50/50 p-4">
        <div class="mb-3 flex items-center gap-1.5">
          <NIcon :component="ServerIcon" size="15" color="#6366f1" />
          <span class="text-xs font-medium text-gray-500">基本信息</span>
        </div>
        <div class="flex gap-3">
          <div class="flex-1">
            <div class="mb-1 text-xs text-gray-400">名称</div>
            <NInput
              v-model:value="form.name"
              size="small"
              placeholder="例如: time-server, filesystem"
            />
          </div>
          <div style="width: 180px">
            <div class="mb-1 text-xs text-gray-400">传输方式</div>
            <NSelect v-model:value="form.transport" :options="transportOptions" size="small" />
          </div>
        </div>
      </div>

      <!-- 连接配置区 -->
      <div class="rounded-xl border border-gray-100 bg-gray-50/50 p-4">
        <div class="mb-3 flex items-center gap-1.5">
          <NIcon :component="RouteIcon" size="15" color="#6366f1" />
          <span class="text-xs font-medium text-gray-500">连接配置</span>
        </div>

        <!-- stdio 模式 -->
        <template v-if="form.transport === 'stdio'">
          <div class="flex gap-3">
            <div class="flex-1">
              <div class="mb-1 text-xs text-gray-400">命令</div>
              <NInput v-model:value="form.command" size="small" placeholder="npx, python, node" />
            </div>
            <div class="flex-1">
              <div class="mb-1 text-xs text-gray-400">参数</div>
              <NInput
                v-model:value="form.args"
                size="small"
                placeholder="-m @modelcontextprotocol/server-time"
              />
            </div>
          </div>
        </template>

        <!-- HTTP 模式 -->
        <template v-else>
          <div>
            <div class="mb-1 text-xs text-gray-400">服务器 URL</div>
            <NInput
              v-model:value="form.url"
              size="small"
              placeholder="https://mcp.example.com/sse"
            />
          </div>
        </template>
      </div>

      <!-- 高级配置区 -->
      <div class="rounded-xl border border-gray-100 bg-gray-50/50 p-4">
        <div class="mb-3 flex items-center gap-1.5">
          <NIcon :component="WrenchIcon" size="15" color="#6366f1" />
          <span class="text-xs font-medium text-gray-500">高级配置</span>
        </div>
        <div class="flex flex-col gap-3">
          <!-- 环境变量 -->
          <div>
            <div class="mb-1 flex items-center gap-1">
              <NIcon :component="BracesIcon" size="13" color="#a78bfa" />
              <span class="text-xs text-gray-400">环境变量 (JSON)</span>
            </div>
            <NInput
              v-model:value="envVarsText"
              type="textarea"
              size="small"
              :autosize="{ minRows: 2, maxRows: 4 }"
              placeholder='{"NODE_ENV": "production"}'
            />
          </div>
          <!-- Headers -->
          <div>
            <div class="mb-1 flex items-center gap-1">
              <NIcon :component="HeaderIcon" size="13" color="#a78bfa" />
              <span class="text-xs text-gray-400">Headers (JSON)</span>
            </div>
            <NInput
              v-model:value="headersText"
              type="textarea"
              size="small"
              :autosize="{ minRows: 2, maxRows: 4 }"
              placeholder='{"Authorization": "Bearer sk-xxx"}'
            />
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
        <NButton size="small" @click="visible = false">取消</NButton>
        <NButton size="small" type="primary" @click="handlePositiveClick">确认</NButton>
      </div>
    </template>
  </NModal>
</template>
