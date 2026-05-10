<script setup lang="ts">
import { computed } from 'vue';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';

dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

const props = defineProps<{
  dataset: Api.Dataset.DatasetItem | null;
  show: boolean;
}>();

const emit = defineEmits<{
  (event: 'update:show', value: boolean): void;
}>();

const visible = computed({
  get: () => props.show,
  set: value => emit('update:show', value)
});

const searchTypeOptions = [
  { label: '关键词检索', value: 'sparse' },
  { label: '语义检索', value: 'dense' },
  { label: '混合检索', value: 'hybrid' }
];

function formatTime(isoString: string) {
  return dayjs(isoString).fromNow();
}
</script>

<template>
  <NModal
    v-model:show="visible"
    :mask-closable="true"
    :show-icon="false"
    style="width: 700px; max-width: 90vw"
  >
    <NCard :bordered="false" size="large" style="border-radius: 12px">
      <NSpace vertical :size="20">
        <div
          style="
            display: flex;
            align-items: center;
            gap: 16px;
            padding-bottom: 16px;
            border-bottom: 1px solid #e8e8e8;
          "
        >
          <NAvatar
            :size="64"
            :style="{
              color: '#333',
              backgroundColor: '#f5f5f5',
              fontSize: '32px'
            }"
          >
            {{ props.dataset?.icon || '🤖' }}
          </NAvatar>
          <div style="flex: 1">
            <h2 style="margin: 0; color: #333; font-size: 24px; font-weight: 600">
              {{ props.dataset?.name }}
            </h2>
            <div style="color: #999; font-size: 12px; margin-top: 4px">
              创建于 {{ props.dataset ? dayjs(props.dataset.created_at).format('YYYY-MM-DD HH:mm') : '-' }}
            </div>
          </div>
          <NButton text style="color: #666" @click="emit('update:show', false)">
            关闭
          </NButton>
        </div>

        <div v-if="props.dataset?.description" style="background: #f8f9fa; padding: 16px; border-radius: 8px">
          <div style="color: #666; font-size: 12px; margin-bottom: 8px; font-weight: 500">描述</div>
          <div style="color: #333; font-size: 14px; line-height: 1.6">{{ props.dataset.description }}</div>
        </div>

        <NGrid cols="2" x-gap="12" y-gap="12" responsive="screen">
          <NGridItem>
            <div class="info-card">
              <div class="info-label">检索方式</div>
              <div class="info-value">
                <NTag
                  :type="
                    props.dataset?.search_type === 'hybrid'
                      ? 'success'
                      : props.dataset?.search_type === 'dense'
                        ? 'info'
                        : 'warning'
                  "
                  size="small"
                >
                  {{
                    searchTypeOptions.find(o => o.value === props.dataset?.search_type)?.label ||
                    props.dataset?.search_type
                  }}
                </NTag>
              </div>
            </div>
          </NGridItem>
          <NGridItem>
            <div class="info-card">
              <div class="info-label">Embedding 模型</div>
              <div class="info-value">{{ props.dataset?.embedding_model || '未设置' }}</div>
            </div>
          </NGridItem>
          <NGridItem>
            <div class="info-card">
              <div class="info-label">Provider ID</div>
              <div class="info-value">{{ props.dataset?.provider_id || '-' }}</div>
            </div>
          </NGridItem>
          <NGridItem>
            <div class="info-card">
              <div class="info-label">Owner ID</div>
              <div class="info-value">{{ props.dataset?.owner_id || '-' }}</div>
            </div>
          </NGridItem>
        </NGrid>

        <div
          style="
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding-top: 16px;
            border-top: 1px solid #e8e8e8;
          "
        >
          <div style="color: #999; font-size: 12px">ID: {{ props.dataset?.id }}</div>
          <div style="color: #999; font-size: 12px">
            更新于 {{ props.dataset ? formatTime(props.dataset.updated_at) : '-' }}
          </div>
        </div>
      </NSpace>
    </NCard>
  </NModal>
</template>

<style scoped>
.info-card {
  background: #f8f9fa;
  padding: 12px 16px;
  border-radius: 8px;
  border: 1px solid #e8e8e8;
}

.info-label {
  color: #666;
  font-size: 12px;
  margin-bottom: 6px;
  font-weight: 500;
}

.info-value {
  color: #333;
  font-size: 14px;
  font-weight: 500;
}
</style>
