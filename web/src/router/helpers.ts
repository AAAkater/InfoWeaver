import type { RouteRecordRaw } from "vue-router"
import type { RouteKey, RoutePath } from "@/typings/router"

/** Mirror of the routePathMap, kept in sync manually. */
const routePathMap: Record<string, string> = {
  Root: "/",
  NotFound: "/:pathMatch(.*)*",
  system403: "/system/403",
  system404: "/system/404",
  system500: "/system/500",
  systemIframe: "/system/iframe/:url",
  systemLogin: "/system/login",
  Chat: "/chat",
  Dataset: "/dataset",
  DatasetInfoId: "/dataset-info/:id",
  Home: "/home",
  Mcp: "/mcp",
  Provider: "/provider",
  UserCenter: "/user-center",
}

export function getRoutePath(key: RouteKey): string | undefined {
  return routePathMap[key]
}

/**
 * get route name by route path
 * @param path route path
 */
export function getRouteName(path: RoutePath) {
  const entries = Object.entries(routePathMap) as [RouteKey, RoutePath][]
  return entries.find(([, routePath]) => routePath === path)?.[0] ?? null
}

/**
 * Convert flat RouteRecordRaw[] to nested (with children).
 * Parent-child is inferred from path prefix: /a/b is a child of /a.
 */
export function nestFlatRoutes(routes: RouteRecordRaw[]): RouteRecordRaw[] {
  const sorted = [...routes].sort((a, b) => a.path.split("/").length - b.path.split("/").length)
  const result: Map<string, RouteRecordRaw> = new Map()

  for (const route of sorted) {
    const parentPath = getParentPath(route.path)

    if (parentPath) {
      const parent = findRouteByPathInMap(result, parentPath)
      if (parent) {
        parent.children = parent.children || []
        parent.children.push({ ...route })
        continue
      }
    }

    result.set(route.path, { ...route })
  }

  return [...result.values()]
}

function getParentPath(path: string): string | null {
  if (path === "/") return null
  const segments = path.split("/")
  segments.pop()
  const parentPath = segments.join("/")
  return parentPath || "/"
}

function findRouteByPathInMap(
  map: Map<string, RouteRecordRaw>,
  path: string,
): RouteRecordRaw | undefined {
  const direct = map.get(path)
  if (direct) return direct

  for (const [, route] of map) {
    if (route.children?.length) {
      const found = findRouteByPath(route, path)
      if (found) return found
    }
  }

  return undefined
}

function findRouteByPath(route: RouteRecordRaw, path: string): RouteRecordRaw | undefined {
  if (route.path === path) return route
  if (route.children?.length) {
    for (const child of route.children) {
      const found = findRouteByPath(child, path)
      if (found) return found
    }
  }
  return undefined
}

/** Flatten nested routes back to flat format */
export function flattenNestedRoutes(routes: RouteRecordRaw[]): RouteRecordRaw[] {
  const result: RouteRecordRaw[] = []

  for (const route of routes) {
    const { children, ...rest } = route
    result.push(rest)
    if (children?.length) {
      result.push(...flattenNestedRoutes(children))
    }
  }

  return result
}

export { routePathMap }
export type { RouteKey, RoutePath }
