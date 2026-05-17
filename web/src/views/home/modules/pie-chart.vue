<script setup lang="ts">
import { useEcharts } from '@/hooks/common/echarts';

defineOptions({
  name: 'PieChart'
});

const { domRef, updateOptions } = useEcharts(() => ({
  title: {
    text: '模型调用分布',
    textStyle: { fontSize: 14, fontWeight: 600 },
    left: '0'
  },
  tooltip: {
    trigger: 'item',
    formatter: '{b}: {c} 次 ({d}%)'
  },
  legend: {
    bottom: '0',
    left: 'center',
    itemStyle: { borderWidth: 0 },
    textStyle: { fontSize: 11 }
  },
  series: [
    {
      color: ['#5da8ff', '#26deca', '#8e9dff', '#fedc69', '#ec4786', '#fcbc25'],
      name: '模型调用分布',
      type: 'pie',
      radius: ['50%', '75%'],
      center: ['50%', '45%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 6,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: { show: false },
      emphasis: {
        label: { show: true, fontSize: 12, fontWeight: 'bold' }
      },
      data: [] as { name: string; value: number }[]
    }
  ]
}));

setTimeout(() => {
  updateOptions(opts => {
    opts.series[0].data = [
      { name: 'text-embedding-3-small', value: 1256 },
      { name: 'gpt-4o-mini', value: 987 },
      { name: 'claude-3-haiku', value: 654 },
      { name: 'text-embedding-3-large', value: 543 },
      { name: 'gpt-4o', value: 321 },
      { name: 'claude-3-sonnet', value: 198 }
    ];
    return opts;
  });
}, 600);
</script>

<template>
  <div ref="domRef" class="h-300px" />
</template>
