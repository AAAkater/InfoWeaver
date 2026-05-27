import type { App } from "vue"
import { createRouter, createWebHistory } from "vue-router"
import { createBuiltinVueRoutes } from "./routes/builtin"
import { createRouterGuard } from "./guard"

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: createBuiltinVueRoutes(),
})

/** Setup Vue Router */
export async function setupRouter(app: App) {
  app.use(router)
  createRouterGuard(router)
  await router.isReady()
}
