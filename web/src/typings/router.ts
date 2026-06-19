import type { RouteRecordSingleView, RouteRecordRedirect, RouteComponent } from "vue-router"

/**
 * Route path map — maps route key to route path.
 * Keep in sync with the actual route definitions.
 */
export interface RoutePathMap {
  Root: "/"
  NotFound: "/:pathMatch(.*)*"
  system403: "/system/403"
  system404: "/system/404"
  system500: "/system/500"
  systemIframe: "/system/iframe/:url"
  systemLogin: "/system/login"
  Chat: "/chat"
  Dataset: "/dataset"
  DatasetInfoId: "/dataset-info/:id"
  Home: "/home"
  Mcp: "/mcp"
  Provider: "/provider"
  UserCenter: "/user-center"
}

/** Route key — a route name */
export type RouteKey = keyof RoutePathMap

/** Route path — the URL path for a route */
export type RoutePath = RoutePathMap[RouteKey]

/** Root route key */
export type RootRouteKey = "Root"

/** Not-found route key */
export type NotFoundRouteKey = "NotFound"

/** Builtin route keys (always registered) */
export type BuiltinRouteKey = RootRouteKey | NotFoundRouteKey

/** Reuse route key (none) */
export type ReuseRouteKey = never

/** Route file key — routes that have their own view file */
export type RouteFileKey = Exclude<RouteKey, BuiltinRouteKey | ReuseRouteKey>

/** Route layout key */
export type RouteLayoutKey = "base" | "blank"

/** Raw route component (sync or async) */
export type RawRouteComponent = RouteComponent | (() => Promise<RouteComponent>)

type MappedNamePath = {
  [K in RouteKey]: { name: K; path: RoutePathMap[K] }
}[RouteKey]

/** AutoRouter single-view route definition */
export type AutoRouterSingleView = Omit<RouteRecordSingleView, "component" | "name" | "path"> & {
  component: RawRouteComponent
  layout: RouteLayoutKey
} & MappedNamePath

/** AutoRouter redirect route definition */
export type AutoRouterRedirect = Omit<RouteRecordRedirect, "children" | "name" | "path"> &
  MappedNamePath

/** AutoRouter route (single-view or redirect) */
export type AutoRouterRoute = AutoRouterSingleView | AutoRouterRedirect

/** API response route type (from backend) */
export interface ElegantConstRoute {
  name: string
  path: string
  component?: string
  layout?: string
  redirect?: string
  props?: Record<string, unknown>
  meta?: Record<string, unknown>
  children?: ElegantConstRoute[]
}
