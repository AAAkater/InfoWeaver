import type { RouteRecordRaw } from "vue-router"
import { routes as generatedRoutes } from "./definitions"
import { getAuthVueRoutes } from "./index"

/** builtin routes, must be constant and setup in vue-router */
const builtinRouteNames = new Set(["Root", "NotFound", "systemLogin"])

const builtinRoutes: RouteRecordRaw[] = generatedRoutes.filter((route) =>
  builtinRouteNames.has(route.name as string),
)

/** create builtin vue routes */
export function createBuiltinVueRoutes() {
  return getAuthVueRoutes(builtinRoutes)
}
