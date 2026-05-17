<script setup lang="ts">
import { ref } from 'vue';
import { createReusableTemplate } from '@vueuse/core';
import { useThemeStore } from '@/store/modules/theme';

defineOptions({
  name: 'CardData'
});

interface StatCard {
  key: string;
  title: string;
  value: string;
  subtitle: string;
  color: {
    start: string;
    end: string;
  };
  icon: string;
}

const statCards = ref<StatCard[]>([
  {
    key: 'todayTokens',
    title: '今日 Token',
    value: '12,580',
    subtitle: '较昨日 +18%',
    color: { start: '#5da8ff', end: '#8e9dff' },
    icon: 'mdi:lightning-bolt'
  },
  {
    key: 'apiCalls',
    title: 'API 调用',
    value: '347',
    subtitle: '今日调用次数',
    color: { start: '#26deca', end: '#56cdf3' },
    icon: 'mdi:api'
  },
  {
    key: 'activeModels',
    title: '活跃模型',
    value: '5',
    subtitle: '共 8 个模型',
    color: { start: '#fcbc25', end: '#f68057' },
    icon: 'mdi:brain-circuit'
  },
  {
    key: 'totalDocs',
    title: '文档数',
    value: '128',
    subtitle: '已索引文档',
    color: { start: '#865ec0', end: '#b955a4' },
    icon: 'mdi:file-document-multiple'
  }
]);

interface GradientBgProps {
  gradientColor: string;
}

const [DefineGradientBg, GradientBg] = createReusableTemplate<GradientBgProps>();

const themeStore = useThemeStore();

function getGradientColor(color: StatCard['color']) {
  return `linear-gradient(to bottom right, ${color.start}, ${color.end})`;
}
</script>

<template>
  <NCard :bordered="false" size="small" class="card-wrapper">
    <!-- define component start: GradientBg -->
    <DefineGradientBg v-slot="{ $slots, gradientColor }">
      <div
        class="px-16px pb-4px pt-8px text-white"
        :style="{ backgroundImage: gradientColor, borderRadius: themeStore.themeRadius + 'px' }"
      >
        <component :is="$slots.default" />
      </div>
    </DefineGradientBg>
    <!-- define component end: GradientBg -->

    <NGrid cols="s:1 m:2 l:4" responsive="screen" :x-gap="16" :y-gap="16">
      <NGi v-for="item in statCards" :key="item.key">
        <GradientBg :gradient-color="getGradientColor(item.color)" class="flex-1">
          <h3 class="text-14px opacity-85">{{ item.title }}</h3>
          <div class="flex items-end justify-between pt-12px">
            <div>
              <div class="text-28px font-700">{{ item.value }}</div>
              <div class="mt-2px text-12px opacity-80">{{ item.subtitle }}</div>
            </div>
            <SvgIcon :icon="item.icon" class="text-36px opacity-30" />
          </div>
        </GradientBg>
      </NGi>
    </NGrid>
  </NCard>
</template>

<style scoped></style>
