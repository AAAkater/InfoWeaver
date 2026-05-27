import type { RouteRecordRaw } from "vue-router"
import type { RawRouteComponent, RouteLayoutKey } from "@/typings/router"
import { routes as generatedRoutes } from "./definitions"

/** Layout components */
const layoutMap: Record<RouteLayoutKey, RawRouteComponent> = {
  base: () => import("@/layouts/base-layout/index.vue"),
  blank: () => import("@/layouts/blank-layout/index.vue"),
}

/**
 * Group flat routes by meta.layout and wrap each group in a layout component.
 * Routes without a known layout (e.g. catch-all NotFound) are kept standalone.
 */
export function getAuthVueRoutes(routes: RouteRecordRaw[]): RouteRecordRaw[] {
  const redirects: RouteRecordRaw[] = []
  const grouped = new Map<RouteLayoutKey, RouteRecordRaw[]>()
  const standalone: RouteRecordRaw[] = []

  for (const route of routes) {
    if ("redirect" in route && route.redirect) {
      redirects.push(route)
      continue
    }

    const meta = route.meta as Record<string, unknown> | undefined
    const layout = meta?.layout as RouteLayoutKey | undefined
    if (!layout || !layoutMap[layout]) {
      standalone.push(route)
      continue
    }

    const cleanMeta = meta ? { ...meta } : {}
    delete cleanMeta.layout
    const items = grouped.get(layout) || []
    items.push({ ...route, meta: cleanMeta })
    grouped.set(layout, items)
  }

  const result: RouteRecordRaw[] = [...redirects]

  for (const item of standalone) {
    result.push(item)
  }

  for (const [layout, items] of grouped) {
    result.push({
      path: `/${layout}-layout`,
      component: layoutMap[layout],
      children: items,
    })
  }

  return result
}

/** Split routes into constant (always registered) and auth (requires login/permission) */
export function createStaticRoutes() {
  const constantRoutes: RouteRecordRaw[] = []
  const authRoutes: RouteRecordRaw[] = []

  generatedRoutes.forEach((item) => {
    const meta = item.meta as Record<string, unknown> | undefined
    if (meta?.constant) {
      constantRoutes.push(item)
    } else {
      authRoutes.push(item)
    }
  })

  return { constantRoutes, authRoutes }
}
