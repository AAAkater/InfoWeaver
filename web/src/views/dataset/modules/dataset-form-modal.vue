<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';

const props = defineProps<{
  isEdit?: boolean;
  model: Api.Dataset.FormModel;
  show: boolean;
}>();

const emit = defineEmits<{
  (event: 'submit', model: Api.Dataset.FormModel): void;
  (event: 'update:show', value: boolean): void;
}>();

const bodyStyle = {
  width: '500px'
};

const searchTypeOptions = [
  { label: '关键词检索', value: 'sparse' },
  { label: '语义检索', value: 'dense' },
  { label: '混合检索', value: 'hybrid' }
];

const emojiList = ['😀', '😂', '😎', '🤖', '🐱', '🦊', '🐶', '🦄', '🐼', '🦉'];
const hovered = ref<string | null>(null);

const visible = computed({
  get: () => props.show,
  set: value => emit('update:show', value)
});

const form = reactive<Api.Dataset.FormModel>({
  id: undefined,
  icon: '🤖',
  description: '',
  name: '',
  search_type: 'hybrid',
  embedding_model: '',
  provider_id: 0
});

function resetForm() {
  Object.assign(form, {
    id: undefined,
    icon: '🤖',
    description: '',
    name: '',
    search_type: 'hybrid',
    embedding_model: '',
    provider_id: 0
  });
}

function syncForm() {
  Object.assign(form, {
    id: props.model.id,
    icon: props.model.icon || '🤖',
    description: props.model.description || '',
    name: props.model.name || '',
    search_type: props.model.search_type || 'hybrid',
    embedding_model: props.model.embedding_model || '',
    provider_id: props.model.provider_id || 0
  });
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

function selectEmoji(emoji: string) {
  form.icon = emoji;
}

function handlePositiveClick() {
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
    title="知识库设置"
    @positive-click="handlePositiveClick"
  >
    <NSpace :size="10" vertical>
      <NCard title="知识库名称" :bordered="false" size="small" content-style="display:flex;gap: 8px;">
        <NPopover trigger="click" placement="bottom-start">
          <template #trigger>
            <NAvatar
              :style="{
                color: 'black',
                backgroundColor: '#FFEAD5',
                cursor: 'pointer'
              }"
            >
              {{ form.icon }}
            </NAvatar>
          </template>
          <div style="display: flex; gap: 8px; padding: 5px; flex-wrap: wrap">
            <span
              v-for="emoji in emojiList"
              :key="emoji"
              style="font-size: 20px; cursor: pointer; padding: 4px; border-radius: 4px"
              :style="{
                backgroundColor: hovered === emoji ? '#e0e0e0' : 'transparent'
              }"
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
          type="text"
          style="background-color: #f1f3f6"
          size="tiny"
          placeholder="请输入知识库名称"
        />
      </NCard>
      <NCard title="描述" :bordered="false" size="small">
        <NInput
          v-model:value="form.description"
          type="textarea"
          size="tiny"
          style="background-color: #f1f3f6"
          placeholder="描述该数据集的内容。详细描述可以让AI更快地访问数据集的内容。如果为空，将使用默认的命中策略。"
        />
      </NCard>
      <NCard title="检索方式" :bordered="false" size="small">
        <NSelect
          v-model:value="form.search_type"
          :options="searchTypeOptions"
          size="tiny"
          style="background-color: #f1f3f6"
        />
      </NCard>
      <NCard title="Embedding模型" :bordered="false" size="small">
        <NInput
          v-model:value="form.embedding_model"
          type="text"
          size="tiny"
          style="background-color: #f1f3f6"
          placeholder="请输入Embedding模型名称"
        />
      </NCard>
      <NCard title="Provider ID" :bordered="false" size="small">
        <NInputNumber
          v-model:value="form.provider_id"
          size="tiny"
          style="background-color: #f1f3f6"
          placeholder="请输入Provider ID"
        />
      </NCard>
    </NSpace>
    <NSpace />
  </NModal>
</template>
