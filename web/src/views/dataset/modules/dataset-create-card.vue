<script setup lang="ts">
import { Add12Filled } from "@vicons/fluent"

const emit = defineEmits<{
  (event: "create"): void
}>()
</script>

<template>
  <NCard
    size="huge"
    class="bg-[#E8ECEF]"
    content-style="padding-left: 5px;padding-right: 5px; height: 100%;"
  >
    <NButton quaternary class="w-full justify-start" @click="emit('create')">
      <NIcon :component="Add12Filled" class="mr-2" />
      创建知识库
    </NButton>
  </NCard>
</template>
