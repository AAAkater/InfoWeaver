<script setup lang="ts">
import { computed, ref, watch } from 'vue';

const props = defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{
  (event: 'submit', modelName: string): void;
  (event: 'update:show', value: boolean): void;
}>();

const bodyStyle = {
  width: '400px'
};

const visible = computed({
  get: () => props.show,
  set: value => emit('update:show', value)
});

const modelName = ref('');

watch(
  () => props.show,
  show => {
    if (show) {
      modelName.value = '';
    }
  }
);

function handlePositiveClick() {
  emit('submit', modelName.value.trim());
}
</script>

<template>
  <NModal
    v-model:show="visible"
    :mask-closable="false"
    preset="dialog"
    :show-icon="false"
    :style="bodyStyle"
    positive-text="添加"
    negative-text="取消"
    title="添加模型"
    @positive-click="handlePositiveClick"
  >
    <NSpace :size="16" vertical>
      <NFormItem label="模型名称" required>
        <NInput v-model:value="modelName" placeholder="请输入模型名称" />
      </NFormItem>
    </NSpace>
  </NModal>
</template>
