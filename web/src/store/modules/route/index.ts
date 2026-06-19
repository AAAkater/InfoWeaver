import { computed, nextTick, ref, shallowRef } from "vue"
import type { RouteRecordRaw } from "vue-router"
import { defineStore } from "pinia"
import { useBoolean } from "@sa/hooks"
import type { RouteKey } from "@/typings/router"
import { router } from "@/router"
import { fetchIsRouteExist } from "@/service/api"
import { SetupStoreId } from "@/enum"
import { createStaticRoutes, getAuthVueRoutes } from "@/router/routes"
import { getRouteName, nestFlatRoutes, flattenNestedRoutes } from "@/router/helpers"
import { useAuthStore } from "../auth"
import { useTabStore } from "../tab"
import {
  filterAuthRoutesByRoles,
  getBreadcrumbsByRoute,
  getCacheRouteNames,
  getGlobalMenusByAuthRoutes,
  getSelectedMenuKeyPathByKey,
  isRouteExistByRouteName,
  transformMenuToSearchMenus,
  updateLocaleOfGlobalMenus,
} from "./shared"

export const useRouteStore = defineStore(SetupStoreId.Route, () => {
  const authStore = useAuthStore()
  const tabStore = useTabStore()
  const { bool: isInitConstantRoute, setBool: setIsInitConstantRoute } = useBoolean()
  const { bool: isInitAuthRoute, setBool: setIsInitAuthRoute } = useBoolean()

  /** Home route key */
  const routeHome = ref(import.meta.env.VITE_ROUTE_HOME)

  /** constant routes */
  const constantRoutes = shallowRef<RouteRecordRaw[]>([])

  function addConstantRoutes(routes: RouteRecordRaw[]) {
    const constantRoutesMap = new Map<string, RouteRecordRaw>([])

    routes.forEach((route) => {
      constantRoutesMap.set(route.name as string, route)
    })

    constantRoutes.value = Array.from(constantRoutesMap.values())
  }

  /** auth routes */
  const authRoutes = shallowRef<RouteRecordRaw[]>([])

  function addAuthRoutes(routes: RouteRecordRaw[]) {
    const authRoutesMap = new Map<string, RouteRecordRaw>([])

    routes.forEach((route) => {
      authRoutesMap.set(route.name as string, route)
    })

    authRoutes.value = Array.from(authRoutesMap.values())
  }

  const removeRouteFns: (() => void)[] = []

  /** Global menus */
  const menus = ref<App.Global.Menu[]>([])
  const searchMenus = computed(() => transformMenuToSearchMenus(menus.value))

  /** Get global menus */
  function getGlobalMenus(routes: RouteRecordRaw[]) {
    menus.value = getGlobalMenusByAuthRoutes(routes)
  }

  /** Update global menus by locale */
  function updateGlobalMenusByLocale() {
    menus.value = updateLocaleOfGlobalMenus(menus.value)
  }

  /** Cache routes */
  const cacheRoutes = ref<RouteKey[]>([])

  /**
   * Exclude cache routes
   *
   * for reset route cache
   */
  const excludeCacheRoutes = ref<RouteKey[]>([])

  /**
   * Get cache routes
   *
   * @param routes Vue routes
   */
  function getCacheRoutes(routes: RouteRecordRaw[]) {
    cacheRoutes.value = getCacheRouteNames(routes)
  }

  /**
   * Reset route cache
   *
   * @default router.currentRoute.value.name current route name
   * @param routeKey
   */
  async function resetRouteCache(routeKey?: RouteKey) {
    const routeName = routeKey || (router.currentRoute.value.name as RouteKey)

    excludeCacheRoutes.value.push(routeName)

    await nextTick()

    excludeCacheRoutes.value = []
  }

  /** Global breadcrumbs */
  const breadcrumbs = computed(() => getBreadcrumbsByRoute(router.currentRoute.value, menus.value))

  /** Reset store */
  async function resetStore() {
    const routeStore = useRouteStore()

    routeStore.$reset()

    resetVueRoutes()

    // after reset store, need to re-init constant route
    await initConstantRoute()
  }

  /** Reset vue routes */
  function resetVueRoutes() {
    removeRouteFns.forEach((fn) => fn())
    removeRouteFns.length = 0
  }

  /** init constant route */
  async function initConstantRoute() {
    if (isInitConstantRoute.value) return

    const staticRoute = createStaticRoutes()
    addConstantRoutes(nestFlatRoutes(staticRoute.constantRoutes))

    handleConstantAndAuthRoutes()

    setIsInitConstantRoute(true)

    tabStore.initHomeTab()
  }

  /** Init auth route */
  async function initAuthRoute() {
    if (!authStore.userInfo.id) {
      await authStore.initUserInfo()
    }

    initStaticAuthRoute()

    tabStore.initHomeTab()
  }

  /** Init static auth route */
  function initStaticAuthRoute() {
    const { authRoutes: staticAuthRoutes } = createStaticRoutes()

    if (authStore.isStaticSuper) {
      addAuthRoutes(nestFlatRoutes(staticAuthRoutes))
    } else {
      const nestedRoutes = nestFlatRoutes(staticAuthRoutes)
      const filteredAuthRoutes = filterAuthRoutesByRoles(nestedRoutes, authStore.userInfo.roles)

      addAuthRoutes(filteredAuthRoutes)
    }

    handleConstantAndAuthRoutes()

    setIsInitAuthRoute(true)
  }

  /** handle constant and auth routes */
  function handleConstantAndAuthRoutes() {
    const allRoutes = [...constantRoutes.value, ...authRoutes.value]

    const flatRoutes = flattenNestedRoutes(allRoutes)

    const vueRoutes = getAuthVueRoutes(flatRoutes)

    resetVueRoutes()

    addRoutesToVueRouter(vueRoutes)

    getGlobalMenus(flatRoutes)

    getCacheRoutes(vueRoutes)
  }

  /**
   * Add routes to vue router
   *
   * @param routes Vue routes
   */
  function addRoutesToVueRouter(routes: RouteRecordRaw[]) {
    routes.forEach((route) => {
      const removeFn = router.addRoute(route)
      addRemoveRouteFn(removeFn)
    })
  }

  /**
   * Add remove route fn
   *
   * @param fn
   */
  function addRemoveRouteFn(fn: () => void) {
    removeRouteFns.push(fn)
  }

  /**
   * Get is auth route exist
   *
   * @param routePath Route path
   */
  async function getIsAuthRouteExist(routePath: string) {
    const routeName = getRouteName(routePath)

    if (!routeName) {
      return false
    }

    if (authRouteMode.value === "static") {
      const { authRoutes: staticAuthRoutes } = createStaticRoutes()
      return isRouteExistByRouteName(routeName, nestFlatRoutes(staticAuthRoutes))
    }

    const { data } = await fetchIsRouteExist(routeName)

    return data
  }

  /**
   * Get selected menu key path
   *
   * @param selectedKey Selected menu key
   */
  function getSelectedMenuKeyPath(selectedKey: string) {
    return getSelectedMenuKeyPathByKey(selectedKey, menus.value)
  }

  async function onRouteSwitchWhenLoggedIn() {
    // some global init logic when logged in and switch route
  }

  async function onRouteSwitchWhenNotLoggedIn() {
    // some global init logic if it does not need to be logged in
  }

  return {
    resetStore,
    routeHome,
    menus,
    searchMenus,
    updateGlobalMenusByLocale,
    cacheRoutes,
    excludeCacheRoutes,
    resetRouteCache,
    breadcrumbs,
    initConstantRoute,
    isInitConstantRoute,
    initAuthRoute,
    isInitAuthRoute,
    setIsInitAuthRoute,
    getIsAuthRouteExist,
    getSelectedMenuKeyPath,
    onRouteSwitchWhenLoggedIn,
    onRouteSwitchWhenNotLoggedIn,
  }
})
