<script setup lang="ts">
import { Search48Filled } from '@vicons/fluent';

const props = defineProps<{
  loading?: boolean;
  modelValue: string;
}>();

const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void;
  (event: 'search'): void;
}>();

function handleSearch() {
  emit('search');
}
</script>

<template>
  <NSpace>
    <NInput
      :value="props.modelValue"
      round
      placeholder="请输入关键字"
      clearable
      :loading="props.loading"
      @update:value="value => emit('update:modelValue', value)"
      @blur="handleSearch"
      @keyup.enter="handleSearch"
    >
      <template #prefix>
        <NIcon :component="Search48Filled" />
      </template>
    </NInput>
  </NSpace>
</template>
