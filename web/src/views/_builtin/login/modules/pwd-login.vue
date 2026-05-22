<script setup lang="ts">
import { computed, reactive } from "vue"
import { loginModuleRecord } from "@/constants/app"
import { useAuthStore } from "@/store/modules/auth"
import { useRouterPush } from "@/hooks/common/router"
import { useFormRules, useNaiveForm } from "@/hooks/common/form"
import { $t } from "@/locales"

defineOptions({
  name: "PwdLogin",
})

const authStore = useAuthStore()
const { toggleLoginModule } = useRouterPush()
const { validate } = useNaiveForm()

interface FormModel {
  email: string
  password: string
}

const model: FormModel = reactive({
  email: "",
  password: "",
})

const rules = computed<Record<keyof FormModel, App.Global.FormRule[]>>(() => {
  // inside computed to make locale reactive, if not apply i18n, you can define it without computed
  const { formRules } = useFormRules()

  return {
    email: formRules.email,
    password: formRules.pwd,
  }
})

async function handleSubmit() {
  await validate()
  await authStore.login(model.email, model.password)
}
</script>

<template>
  <NForm :model="model" :rules="rules" size="large" :show-label="false" @keyup.enter="handleSubmit">
    <NFormItem path="userName">
      <NInput v-model:value="model.email" :placeholder="$t('page.login.common.emailPlaceholder')" />
    </NFormItem>
    <NFormItem path="password">
      <NInput
        v-model:value="model.password"
        type="password"
        show-password-on="click"
        :placeholder="$t('page.login.common.passwordPlaceholder')"
      />
    </NFormItem>
    <NSpace vertical :size="24">
      <div class="flex-y-center justify-between">
        <NCheckbox>{{ $t("page.login.pwdLogin.rememberMe") }}</NCheckbox>
        <NButton quaternary @click="toggleLoginModule('reset-pwd')">
          {{ $t("page.login.pwdLogin.forgetPassword") }}
        </NButton>
      </div>
      <NButton
        type="primary"
        size="large"
        round
        block
        :loading="authStore.loginLoading"
        @click="handleSubmit"
      >
        {{ $t("common.confirm") }}
      </NButton>
      <div class="flex-y-center justify-between gap-12px">
        <NButton class="flex-1" block @click="toggleLoginModule('register')">
          {{ $t(loginModuleRecord.register) }}
        </NButton>
      </div>
    </NSpace>
  </NForm>
</template>

<style scoped></style>
