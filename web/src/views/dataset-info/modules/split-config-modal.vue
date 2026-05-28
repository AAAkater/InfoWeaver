<script setup lang="ts">
import { NButton, NInputNumber, NModal } from "naive-ui"

defineOptions({ name: "SplitConfigModal" })

interface Props {
  show: boolean
  loading: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  confirm: [chunkSize: number, chunkOverlap: number]
  cancel: []
}>()

const chunkSize = defineModel<number>("chunkSize", { default: 512 })
const chunkOverlap = defineModel<number>("chunkOverlap", { default: 50 })
</script>

<template>
  <NModal
    :show="show"
    title="分块配置"
    preset="card"
    style="width: 420px"
    @update:show="!$event && emit('cancel')"
  >
    <div class="flex flex-col gap-16px">
      <div>
        <div class="mb-4px text-13px text-gray-500">Chunk Size（分块大小）</div>
        <NInputNumber v-model:value="chunkSize" :min="64" :max="4096" style="width: 100%" />
      </div>
      <div>
        <div class="mb-4px text-13px text-gray-500">Chunk Overlap（重叠长度）</div>
        <NInputNumber v-model:value="chunkOverlap" :min="0" :max="2048" style="width: 100%" />
      </div>
      <div class="flex justify-end gap-8px">
        <NButton @click="emit('cancel')">取消</NButton>
        <NButton
          type="primary"
          :loading="loading"
          @click="emit('confirm', chunkSize, chunkOverlap)"
        >
          确认拆分
        </NButton>
      </div>
    </div>
  </NModal>
</template>
