<script setup lang="ts">
import { useAuthStore } from '@/store/modules/auth';
const authStore = useAuthStore();
</script>

<template>
  <NCard>
    <div>
      <NH3>我的账户</NH3>
      <div class="flex">
        <SoybeanAvatar class="size-65px!" />
        <div>
          <p class="text-15px font-550">{{ authStore.userInfo.username }}</p>
          <p class="text-9px text-gray font-400">{{ authStore.userInfo.email }}</p>
        </div>
      </div>
    </div>
  </NCard>
  <div>
    <NH4>用户名</NH4>
    <NH4>邮箱</NH4>
    <NH4>密码</NH4>
  </div>
</template>

<style scoped></style>
