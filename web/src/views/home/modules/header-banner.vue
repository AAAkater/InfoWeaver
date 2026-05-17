<script setup lang="ts">
import { computed, ref } from 'vue';
import { useAuthStore } from '@/store/modules/auth';
import defaultAvatar from '@/assets/imgs/soybean.jpg';

defineOptions({
  name: 'HeaderBanner'
});

const authStore = useAuthStore();

// Mock user stats
const totalTokens = ref('1,245,680');
const monthlyTokens = ref('326,400');
const totalCalls = ref('4,832');

const avatarUrl = computed(() => {
  return authStore.userInfo.avatar_url || '';
});
</script>

<template>
  <NCard :bordered="false" class="card-wrapper">
    <div class="flex flex-col gap-16px sm:flex-row sm:items-center sm:justify-between">
      <div class="flex-y-center">
        <div class="size-72px shrink-0 overflow-hidden rd-1/2">
          <img :src="avatarUrl || defaultAvatar" class="size-full" />
        </div>
        <div class="pl-12px">
          <h3 class="text-18px font-semibold">你好，{{ authStore.userInfo.username }}</h3>
          <p class="text-12px text-gray-400">欢迎回到 InfoWeaver 工作台</p>
        </div>
      </div>
      <div class="flex gap-24px">
        <div class="text-center">
          <div class="text-24px text-[#5da8ff] font-700">{{ totalTokens }}</div>
          <div class="text-12px text-gray-400">累计 Token</div>
        </div>
        <div class="text-center">
          <div class="text-24px text-[#26deca] font-700">{{ monthlyTokens }}</div>
          <div class="text-12px text-gray-400">本月 Token</div>
        </div>
        <div class="text-center">
          <div class="text-24px text-[#8e9dff] font-700">{{ totalCalls }}</div>
          <div class="text-12px text-gray-400">API 调用</div>
        </div>
      </div>
    </div>
  </NCard>
</template>

<style scoped></style>
