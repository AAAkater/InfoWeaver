import { request } from "../request"

/**
 * whether the route is exist
 *
 * @param routeName route name
 */
export function fetchIsRouteExist(routeName: string) {
  return request<boolean>({ url: "/route/isRouteExist", params: { routeName } })
}
