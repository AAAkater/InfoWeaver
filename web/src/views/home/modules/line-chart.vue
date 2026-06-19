<script setup lang="ts">
import { useEcharts } from "@/hooks/common/echarts"

defineOptions({
  name: "LineChart",
})

const { domRef, updateOptions } = useEcharts(() => ({
  title: {
    text: "Token 用量趋势",
    textStyle: { fontSize: 14, fontWeight: 600 },
    left: "0",
  },
  tooltip: {
    trigger: "axis",
    axisPointer: { type: "cross" },
  },
  legend: {
    data: ["输入 Token", "输出 Token"],
    top: "0",
    right: "0",
  },
  grid: {
    left: "3%",
    right: "4%",
    bottom: "3%",
    top: "50px",
  },
  xAxis: {
    type: "category",
    boundaryGap: false,
    data: [] as string[],
  },
  yAxis: {
    type: "value",
    axisLabel: { formatter: (v: number) => (v >= 1000 ? `${(v / 1000).toFixed(0)}k` : String(v)) },
  },
  series: [
    {
      color: "#5da8ff",
      name: "输入 Token",
      type: "line",
      smooth: true,
      areaStyle: { color: "rgba(93,168,255,0.12)" },
      data: [] as number[],
    },
    {
      color: "#26deca",
      name: "输出 Token",
      type: "line",
      smooth: true,
      areaStyle: { color: "rgba(38,222,202,0.12)" },
      data: [] as number[],
    },
  ],
}))

const days = ["5/11", "5/12", "5/13", "5/14", "5/15", "5/16", "5/17"]
const inputTokens = [8200, 9300, 7500, 11200, 9800, 13200, 10500]
const outputTokens = [2100, 2800, 1900, 3500, 2600, 4200, 3200]

setTimeout(() => {
  updateOptions((opts) => {
    opts.xAxis.data = days
    opts.series[0].data = inputTokens
    opts.series[1].data = outputTokens
    return opts
  })
}, 400)
</script>

<template>
  <div ref="domRef" class="h-300px" />
</template>

<style scoped></style>
