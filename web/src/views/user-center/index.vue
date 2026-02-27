<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useMessage } from 'naive-ui';
import { updateUserPassword, updateUserProfile } from '@/service/api/auth';
import { useAuthStore } from '@/store/modules/auth';

const authStore = useAuthStore();
const message = useMessage();

// 用户信息编辑
const isEditingProfile = ref(false);
const profileForm = reactive({
  username: authStore.userInfo.username,
  email: authStore.userInfo.email
});
const isLoadingProfile = ref(false);

// 密码修改
const isEditingPassword = ref(false);
const passwordForm = reactive({
  first_password: '',
  second_password: ''
});
const isLoadingPassword = ref(false);

async function handleUpdateProfile() {
  if (!profileForm.username.trim()) {
    message.error('用户名不能为空');
    return;
  }
  if (!profileForm.email.trim()) {
    message.error('邮箱不能为空');
    return;
  }

  isLoadingProfile.value = true;

  const { response: res } = await updateUserProfile({
    username: profileForm.username,
    email: profileForm.email
  });

  if (res.data.code === 0) {
    message.success('用户信息更新成功');
    authStore.userInfo.username = profileForm.username;
    authStore.userInfo.email = profileForm.email;
    isEditingProfile.value = false;
  } else {
    message.error(res.data.msg || '更新失败');
  }

  isLoadingProfile.value = false;
}

async function handleUpdatePassword() {
  if (!passwordForm.first_password.trim()) {
    message.error('新密码不能为空');
    return;
  }
  if (passwordForm.first_password.length < 6) {
    message.error('新密码长度不能少于6位');
    return;
  }
  if (passwordForm.first_password !== passwordForm.second_password) {
    message.error('两次输入的密码不一致');
    return;
  }

  isLoadingPassword.value = true;

  const { response: res } = await updateUserPassword({
    first_password: passwordForm.first_password,
    second_password: passwordForm.second_password
  });

  if (res.data.code === 0) {
    message.success('密码修改成功');
    passwordForm.first_password = '';
    passwordForm.second_password = '';
    isEditingPassword.value = false;
  } else {
    message.error(res.data.msg || '密码修改失败');
  }

  isLoadingPassword.value = false;
}

function cancelEditProfile() {
  isEditingProfile.value = false;
  profileForm.username = authStore.userInfo.username;
  profileForm.email = authStore.userInfo.email;
}

function cancelEditPassword() {
  isEditingPassword.value = false;
  passwordForm.first_password = '';
  passwordForm.second_password = '';
}
</script>

<template>
  <div class="flex flex-col gap-32px">
    <!-- 账户信息卡片 -->
    <NCard class="profile-card" :bordered="false" :segmented="false">
      <div class="user-header">
        <div class="flex items-center gap-16px">
          <SoybeanAvatar class="size-80px!" />
          <div class="flex-1">
            <h2 class="mb-4px text-18px font-600">{{ authStore.userInfo.username }}</h2>
            <p class="text-13px text-gray-500">{{ authStore.userInfo.email }}</p>
          </div>
        </div>
      </div>
    </NCard>

    <!-- 账户安全设置 -->
    <div>
      <div class="section-header mb-16px">
        <h3 class="text-16px font-600">账户安全</h3>
        <div class="section-divider"></div>
      </div>

      <!-- 用户名修改 -->
      <NCard class="setting-card" :bordered="false" content-style="padding: 0;">
        <div class="setting-item">
          <div class="setting-label">
            <span class="label-title">用户名</span>
            <span class="label-desc">用于登录和识别身份</span>
          </div>
          <div v-if="!isEditingProfile" class="setting-value">
            <span class="value-text">{{ authStore.userInfo.username }}</span>
            <NButton text type="primary" class="edit-btn" @click="isEditingProfile = true">修改</NButton>
          </div>
          <div v-else class="setting-form">
            <NInput v-model:value="profileForm.username" placeholder="请输入新用户名" clearable class="form-input" />
            <div class="form-actions">
              <NButton type="primary" :loading="isLoadingProfile" size="small" @click="handleUpdateProfile">
                保存
              </NButton>
              <NButton size="small" @click="cancelEditProfile">取消</NButton>
            </div>
          </div>
        </div>
      </NCard>

      <!-- 邮箱修改 -->
      <NCard class="setting-card" :bordered="false" content-style="padding: 0;">
        <div class="setting-item">
          <div class="setting-label">
            <span class="label-title">邮箱地址</span>
            <span class="label-desc">用于账户恢复和通知</span>
          </div>
          <div v-if="!isEditingProfile" class="setting-value">
            <span class="value-text">{{ authStore.userInfo.email }}</span>
          </div>
          <div v-else class="setting-form">
            <NInput v-model:value="profileForm.email" type="text" placeholder="请输入新邮箱" clearable />
            <span class="form-tip">邮箱修改会在下次登录时需要验证</span>
          </div>
        </div>
      </NCard>

      <!-- 密码修改 -->
      <NCard class="setting-card" :bordered="false" content-style="padding: 0;">
        <div class="setting-item">
          <div class="setting-label">
            <span class="label-title">密码</span>
            <span class="label-desc">定期修改密码保护账户安全</span>
          </div>
          <div v-if="!isEditingPassword" class="setting-value">
            <span class="value-text">••••••••</span>
            <NButton text type="primary" class="edit-btn" @click="isEditingPassword = true">修改</NButton>
          </div>
          <div v-else class="setting-form">
            <NInput
              v-model:value="passwordForm.first_password"
              type="password"
              placeholder="请输入新密码（至少6位）"
              clearable
              show-password-on="click"
              class="form-input"
            />
            <NInput
              v-model:value="passwordForm.second_password"
              type="password"
              placeholder="请确认新密码"
              clearable
              show-password-on="click"
              class="form-input"
            />
            <div class="form-actions">
              <NButton type="primary" :loading="isLoadingPassword" size="small" @click="handleUpdatePassword">
                保存
              </NButton>
              <NButton size="small" @click="cancelEditPassword">取消</NButton>
            </div>
          </div>
        </div>
      </NCard>
    </div>
  </div>
</template>

<style scoped>
.profile-card {
  background: linear-gradient(135deg, rgba(79, 172, 254, 0.1) 0%, rgba(180, 198, 255, 0.1) 100%);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.profile-card:hover {
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.15);
}

.user-header {
  padding: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 4px;
}

.section-header h3 {
  color: #1f2937;
  white-space: nowrap;
}

.section-divider {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, #e5e7eb 0%, transparent 100%);
}

.setting-card {
  margin-bottom: 12px;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  transition: all 0.3s ease;
}

.setting-card:hover {
  border-color: #d0d0d0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  gap: 24px;
}

.setting-label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 150px;
}

.label-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.label-desc {
  font-size: 12px;
  color: #9ca3af;
}

.setting-value {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.value-text {
  font-size: 14px;
  color: #4b5563;
  min-width: 200px;
  word-break: break-all;
}

.edit-btn {
  font-size: 13px;
}

.setting-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex: 1;
  min-width: 300px;
}

.form-input {
  width: 100%;
}

.form-tip {
  font-size: 12px;
  color: #9ca3af;
  padding: 0 4px;
}

.form-actions {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

@media (max-width: 768px) {
  .setting-item {
    flex-direction: column;
    align-items: flex-start;
    padding: 16px;
  }

  .value-text {
    min-width: unset;
  }

  .setting-form {
    width: 100%;
    min-width: unset;
  }

  .user-header {
    padding: 16px;
  }

  .label-desc {
    display: none;
  }
}
</style>
