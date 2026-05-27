import process from "node:process"
import path from "node:path"
import { URL, fileURLToPath } from "node:url"
import { defineConfig, loadEnv } from "vite"
import vue from "@vitejs/plugin-vue"
import vueJsx from "@vitejs/plugin-vue-jsx"
import progress from "vite-plugin-progress"
import VueDevtools from "vite-plugin-vue-devtools"
import unocss from "@unocss/vite"
import presetIcons from "@unocss/preset-icons"
import { FileSystemIconLoader } from "@iconify/utils/lib/loader/node-loaders"
import Icons from "unplugin-icons/vite"
import IconsResolver from "unplugin-icons/resolver"
import Components from "unplugin-vue-components/vite"
import { NaiveUiResolver } from "unplugin-vue-components/resolvers"
import { ProNaiveUIResolver } from "pro-naive-ui-resolver"
import { createSvgIconsPlugin } from "vite-plugin-svg-icons"
import dayjs from "dayjs"
import utc from "dayjs/plugin/utc"
import timezone from "dayjs/plugin/timezone"

export function getBuildTime() {
  dayjs.extend(utc)
  dayjs.extend(timezone)
  return dayjs.tz(Date.now(), "Asia/Shanghai").format("YYYY-MM-DD HH:mm:ss")
}
export default defineConfig((configEnv) => {
  const viteEnv = loadEnv(configEnv.mode, process.cwd()) as unknown as Env.ImportMeta

  const buildTime = getBuildTime()

  const { VITE_ICON_PREFIX, VITE_ICON_LOCAL_PREFIX } = viteEnv
  const localIconPath = path.join(process.cwd(), "src/assets/svg-icon")
  const collectionName = VITE_ICON_LOCAL_PREFIX.replace(`${VITE_ICON_PREFIX}-`, "")
  return {
    base: viteEnv.VITE_BASE_URL,
    resolve: {
      alias: {
        "~": fileURLToPath(new URL("./", import.meta.url)),
        "@": fileURLToPath(new URL("./src", import.meta.url)),
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          api: "modern-compiler",
          additionalData: `@use "@/styles/scss/global.scss" as *;`,
        },
      },
    },
    plugins: [
      vue(),
      vueJsx(),
      VueDevtools(),
      unocss({
        presets: [
          presetIcons({
            prefix: `${VITE_ICON_PREFIX}-`,
            scale: 1,
            extraProperties: {
              display: "inline-block",
            },
            collections: {
              [collectionName]: FileSystemIconLoader(localIconPath, (svg) =>
                svg.replace(/^<svg\s/, '<svg width="1em" height="1em" '),
              ),
            },
            warn: true,
          }),
        ],
      }),
      Icons({
        compiler: "vue3",
        customCollections: {
          [collectionName]: FileSystemIconLoader(localIconPath, (svg) =>
            svg.replace(/^<svg\s/, '<svg width="1em" height="1em" '),
          ),
        },
        scale: 1,
        defaultClass: "inline-block",
      }),
      Components({
        dts: "src/typings/components.d.ts",
        types: [{ from: "vue-router", names: ["RouterLink", "RouterView"] }],
        resolvers: [
          NaiveUiResolver(),
          ProNaiveUIResolver(),
          IconsResolver({ customCollections: [collectionName], componentPrefix: VITE_ICON_PREFIX }),
        ],
      }),
      createSvgIconsPlugin({
        iconDirs: [localIconPath],
        symbolId: `${VITE_ICON_LOCAL_PREFIX}-[dir]-[name]`,
        inject: "body-last",
        customDomId: "__SVG_ICON_LOCAL__",
      }),
      progress(),
    ],
    define: {
      BUILD_TIME: JSON.stringify(buildTime),
    },
    server: {
      host: "0.0.0.0",
      port: 9527,
      open: true,
      proxy: {
        "/api/v1": {
          target: viteEnv.VITE_SERVICE_BASE_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api\/v1/, ""),
        },
      },
    },
    preview: {
      port: 9725,
    },
    build: {
      reportCompressedSize: false,
      sourcemap: viteEnv.VITE_SOURCE_MAP === "Y",
      commonjsOptions: {
        ignoreTryCatch: false,
      },
    },
  }
})
