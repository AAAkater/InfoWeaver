import type { RouteMeta } from 'vue-router';
import ElegantVueRouter from '@elegant-router/vue/vite';
import type { RouteKey } from '@elegant-router/types';

export function setupElegantRouter() {
  return ElegantVueRouter({
    layouts: {
      base: 'src/layouts/base-layout/index.vue',
      blank: 'src/layouts/blank-layout/index.vue'
    },
    customRoutes: {
      names: [
        'exception_403',
        'exception_404',
        'exception_500',
        'document_project',
        'document_project-link',
        'document_video',
        'document_vue',
        'document_vite',
        'document_unocss',
        'document_naive',
        'document_pro-naive',
        'document_antd',
        'document_alova'
      ]
    },
    routePathTransformer(routeName, routePath) {
      const key = routeName as RouteKey;

      if (key === 'login') {
        const modules: UnionKey.LoginModule[] = ['pwd-login', 'code-login', 'register', 'reset-pwd', 'bind-wechat'];

        const moduleReg = modules.join('|');

        return `/login/:module(${moduleReg})?`;
      }

      return routePath;
    },
    onRouteMetaGen(routeName) {
      const key = routeName as RouteKey;

      const constantRoutes: RouteKey[] = ['login', '403', '404', '500'];

      const hideInMenuRoutes: Partial<Record<RouteKey, string>> = {
        'dataset-info': 'dataset'
      };

      const orderMap: Partial<Record<RouteKey, number>> = {
        home: 1,
        chat: 2,
        dataset: 3,
        provider: 4,
        mcp: 5
      };

      const iconMap: Partial<Record<RouteKey, string>> = {
        home: 'mdi:monitor-dashboard',
        chat: 'mdi:chat-plus',
        dataset: 'material-symbols-light:book-4-spark-rounded',
        provider: 'mdi:cloud-braces',
        mcp: 'mdi:server'
      };

      const meta: Partial<RouteMeta> = {
        title: key,
        i18nKey: `route.${key}` as App.I18n.I18nKey
      };

      if (key in orderMap) {
        meta.order = orderMap[key];
      }

      if (key in iconMap) {
        meta.icon = iconMap[key];
      }

      if (constantRoutes.includes(key)) {
        meta.constant = true;
      }

      if (key in hideInMenuRoutes) {
        meta.hideInMenu = true;
        meta.activeMenu = hideInMenuRoutes[key];
      }

      return meta;
    }
  });
}
