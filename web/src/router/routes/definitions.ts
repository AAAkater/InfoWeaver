import type { RouteRecordRaw } from "vue-router"

export const routes: RouteRecordRaw[] = [
  // ── Builtin ──────────────────────────────────────────────
  {
    name: "Root",
    path: "/",
    redirect: "/home",
    meta: { title: "Root", hideInMenu: true },
  },
  {
    name: "NotFound",
    path: "/:pathMatch(.*)*",
    component: () => import("@/views/system/404/index.vue"),
    meta: { title: "NotFound", hideInMenu: true },
  },

  // ── App pages ────────────────────────────────────────────
  {
    name: "Home",
    path: "/home",
    component: () => import("@/views/home/index.vue"),
    meta: { layout: "base", title: "Home", i18nKey: "route.home", icon: "mdi:home-outline" },
  },
  {
    name: "Chat",
    path: "/chat",
    component: () => import("@/views/chat/index.vue"),
    meta: { layout: "base", title: "Chat", i18nKey: "route.chat", icon: "mdi:message-outline" },
  },
  {
    name: "Dataset",
    path: "/dataset",
    component: () => import("@/views/dataset/index.vue"),
    meta: {
      layout: "base",
      title: "Dataset",
      i18nKey: "route.dataset",
      icon: "mdi:database-outline",
    },
  },
  {
    name: "DatasetInfoId",
    path: "/dataset-info/:id",
    component: () => import("@/views/dataset-info/[id].vue"),
    props: true,
    meta: {
      layout: "base",
      title: "DatasetInfoId",
      i18nKey: "route.dataset-info",
      hideInMenu: true,
    },
  },
  {
    name: "Mcp",
    path: "/mcp",
    component: () => import("@/views/mcp/index.vue"),
    meta: { layout: "base", title: "Mcp", i18nKey: "route.mcp", icon: "mdi:server" },
  },
  {
    name: "Provider",
    path: "/provider",
    component: () => import("@/views/provider/index.vue"),
    meta: {
      layout: "base",
      title: "Provider",
      i18nKey: "route.provider",
      icon: "mdi:store-outline",
    },
  },
  {
    name: "UserCenter",
    path: "/user-center",
    component: () => import("@/views/user-center/index.vue"),
    meta: { layout: "base", title: "UserCenter", i18nKey: "route.user-center", hideInMenu: true },
  },

  // ── System pages ──────────────────────────────────────────
  {
    name: "system403",
    path: "/system/403",
    component: () => import("@/views/system/403/index.vue"),
    meta: {
      layout: "base",
      title: "system403",
      i18nKey: "route.403",
      constant: true,
      hideInMenu: true,
    },
  },
  {
    name: "system404",
    path: "/system/404",
    component: () => import("@/views/system/404/index.vue"),
    meta: {
      layout: "base",
      title: "system404",
      i18nKey: "route.404",
      constant: true,
      hideInMenu: true,
    },
  },
  {
    name: "system500",
    path: "/system/500",
    component: () => import("@/views/system/500/index.vue"),
    meta: {
      layout: "base",
      title: "system500",
      i18nKey: "route.500",
      constant: true,
      hideInMenu: true,
    },
  },
  {
    name: "systemIframe",
    path: "/system/iframe/:url",
    component: () => import("@/views/system/iframe-page/[url].vue"),
    meta: { layout: "base", title: "systemIframe", i18nKey: "route.iframe-page", hideInMenu: true },
  },
  {
    name: "systemLogin",
    path: "/system/login",
    component: () => import("@/views/system/login/index.vue"),
    meta: {
      layout: "blank",
      title: "systemLogin",
      i18nKey: "route.login",
      constant: true,
      hideInMenu: true,
    },
  },
]
