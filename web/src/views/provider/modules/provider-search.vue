<script setup lang="ts">
import { Add12Filled, Search48Filled } from '@vicons/fluent';

const props = defineProps<{
  loading?: boolean;
  modelValue: string;
}>();

const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void;
  (event: 'create'): void;
}>();
</script>

<template>
  <NSpace justify="space-between">
    <NInput
      :value="props.modelValue"
      round
      placeholder="搜索 Provider 名称或 URL"
      clearable
      :loading="props.loading"
      style="width: 300px"
      @update:value="value => emit('update:modelValue', value)"
    >
      <template #prefix>
        <NIcon :component="Search48Filled" />
      </template>
    </NInput>

    <NButton type="primary" @click="emit('create')">
      <NIcon :component="Add12Filled" class="mr-2" />
      创建 Provider
    </NButton>
  </NSpace>
</template>
