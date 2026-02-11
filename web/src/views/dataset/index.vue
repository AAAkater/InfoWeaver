<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue';
import type { Component } from 'vue';
import { NIcon, useMessage } from 'naive-ui';
import {
  Add12Filled,
  Delete16Regular as DeleteIcon,
  Edit24Regular as EditIcon,
  MoreHorizontal28Regular,
  Search48Filled
} from '@vicons/fluent';
import { createDataset, deleteDataset, getDatasets } from '@/service/api/dataset';
const showModal = ref(false);

const bodyStyle = {
  width: '500px'
};
const datasets = ref<Api.Dataset.DatasetItem[]>();
const model: Api.Dataset.FormModel = reactive({
  icon: '🤖',
  description: '',
  name: ''
});
const message = useMessage();
function renderIcon(icon: Component) {
  return () => {
    return h(NIcon, null, {
      default: () => h(icon)
    });
  };
}
async function fetchDatasets() {
  const { response: res } = await getDatasets();
  if (res?.data?.code === 0) {
    const datasetList = res.data.data?.datasets ?? [];
    datasets.value = datasetList;
  } else {
    message.error('获取数据失败');
    datasets.value = [];
  }
}
async function handleCreateDataset(md: Api.Dataset.FormModel) {
  const { response: res } = await createDataset(md);
  if (res.data.code === 0) {
    message.success('创建成功');
  } else {
    const errorMsg = res.data.msg;
    window.$message?.error(errorMsg);
  }
  fetchDatasets();
}
async function handleDeleteDataset(id: number) {
  const { response: res } = await deleteDataset(id);
  if (res.data.code === 0) {
    message.success('删除成功');
  } else {
    const errorMsg = res.data.msg;
    window.$message?.error(errorMsg);
  }
  fetchDatasets();
}
const emojiList = ['😀', '😂', '😎', '🤖', '🐱', '🦊', '🐶', '🦄', '🐼', '🦉'];

const hovered = ref<string | null>(null);
const dropdownOptions = [
  {
    label: '编辑',
    key: 'edit',
    icon: renderIcon(EditIcon)
  },
  {
    label: '删除',
    key: 'delete',
    icon: renderIcon(DeleteIcon)
  }
];
function handleSelect(key: string, id: number) {
  switch (key) {
    case 'edit':
      showModal.value = true;
      break;
    case 'delete':
      handleDeleteDataset(id);
      break;
    default:
  }
}
function selectEmoji(emoji: string) {
  model.icon = emoji;
}

onMounted(() => {
  fetchDatasets();
});
</script>

<template>
  <NSpace vertical :size="16">
    <NSpace>
      <NInput round placeholder="请输入知识库id" clearable>
        <template #prefix>
          <NIcon :component="Search48Filled" />
        </template>
      </NInput>
    </NSpace>

    <NGrid cols="3" x-gap="16" y-gap="17" responsive="screen" item-responsive>
      <NGridItem class="grid-item-equal-height">
        <NCard size="huge" class="bg-[#E8ECEF]" content-style="padding-left: 5px;padding-right: 5px; height: 100%;">
          <NButton quaternary class="w-full justify-start" @click="showModal = true">
            <NIcon :component="Add12Filled" class="mr-2" />
            创建知识库
          </NButton>
          <NModal
            v-model:show="showModal"
            :mask-closable="false"
            preset="dialog"
            :show-icon="false"
            :style="bodyStyle"
            positive-text="确认"
            negative-text="取消"
            title="知识库设置"
            @positive-click="handleCreateDataset(model)"
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
                      {{ model.icon }}
                    </NAvatar>
                  </template>
                  <!-- Popover 内容：Emoji 选择器 -->
                  <div style="display: flex; gap: 8px; padding: 5px; flex-wrap: wrap">
                    <span
                      v-for="emoji in emojiList"
                      :key="emoji"
                      style="font-size: 20px; cursor: pointer; padding: 4px; border-radius: 4px"
                      :style="{ backgroundColor: hovered === emoji ? '#e0e0e0' : 'transparent' }"
                      @click="selectEmoji(emoji)"
                      @mouseenter="hovered = emoji"
                      @mouseleave="hovered = null"
                    >
                      {{ emoji }}
                    </span>
                  </div>
                </NPopover>

                <NInput
                  v-model:value="model.name"
                  type="text"
                  style="background-color: #f1f3f6"
                  size="tiny"
                  placeholder="请输入知识库名称"
                />
              </NCard>
              <NCard title="描述" :bordered="false" size="small">
                <NInput
                  v-model:value="model.description"
                  type="textarea"
                  size="tiny"
                  style="background-color: #f1f3f6"
                  placeholder="描述该数据集的内容。详细描述可以让AI更快地访问数据集的内容。如果为空，将使用默认的命中策略。"
                />
              </NCard>
            </NSpace>
            <NSpace></NSpace>
          </NModal>
        </NCard>
      </NGridItem>
      <NGridItem v-for="dataset in datasets" :key="dataset.id">
        <NCard hoverable size="huge" style="cursor: pointer">
          <NSpace vertical>
            <div style="display: flex; align-items: flex-start; gap: 8px; width: 100%">
              <NAvatar
                size="large"
                :style="{
                  color: 'black',
                  backgroundColor: '#E0F2FE',
                  cursor: 'pointer'
                }"
              >
                {{ dataset.icon || '🤖' }}
              </NAvatar>

              <div style="flex: 1">
                <div :style="{ fontWeight: 'bold' }">{{ dataset.name }}</div>
                <div style="color: #949494; font-size: 10px; line-height: 1.2">id · {{ dataset.id }}</div>
              </div>

              <NDropdown
                :options="dropdownOptions"
                trigger="click"
                size="small"
                @select="key => handleSelect(key, dataset.id)"
              >
                <NButton size="small" secondary>
                  <NIcon>
                    <MoreHorizontal28Regular />
                  </NIcon>
                </NButton>
              </NDropdown>
            </div>
            <div style="color: #666; font-size: 10px; line-height: 1.4">{{ dataset.description }}</div>
            <div
              style="
                color: #949494;
                font-size: 10px;
                font-size: 10px;
                transform: scale(0.85);
                transform-origin: left top;
              "
            >
              更新于 · {{ dataset.updated_at }}
            </div>
          </NSpace>
        </NCard>
      </NGridItem>
    </NGrid>
  </NSpace>
</template>

<style scoped>
.grid-item-equal-height {
  display: flex;
}
.grid-item-equal-height .n-card {
  flex: 1;
  min-height: 0;
}
</style>
