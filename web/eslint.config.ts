import {
  defineConfigWithVueTs,
  configureVueProject,
  vueTsConfigs,
} from "@vue/eslint-config-typescript"
import skipFormatting from "@vue/eslint-config-prettier/skip-formatting"
import { globalIgnores } from "eslint/config"
import pluginVue from "eslint-plugin-vue"
import pluginVitest from "@vitest/eslint-plugin"

configureVueProject({ scriptLangs: ["ts", "tsx"] })

export default defineConfigWithVueTs(
  {
    name: "app/files-to-lint",
    files: ["**/*.{ts,mts,tsx,vue}"],
  },

  globalIgnores(["**/dist/**", "**/dist-ssr/**", "**/coverage/**"]),

  pluginVue.configs["flat/essential"],
  vueTsConfigs.recommended,

  {
    ...pluginVitest.configs.recommended,
    files: ["src/**/__tests__/*"],
  },
  {
    rules: {
      "vue/multi-word-component-names": [
        "warn",
        {
          ignores: ["index", "App", "Register", "[id]", "[url]"],
        },
      ],
      "vue/component-name-in-template-casing": [
        "warn",
        "PascalCase",
        {
          registeredComponentsOnly: false,
          ignores: ["/^icon-/"],
        },
      ],
      semi: "off",
      "@typescript-eslint/semi": "off",
    },
  },
  skipFormatting,
)
