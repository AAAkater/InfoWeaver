import type { AxiosResponse } from "axios"
import { BACKEND_ERROR_CODE, createFlatRequest, createRequest } from "@sa/axios"
import { useAuthStore } from "@/store/modules/auth"
import { getAuthorization, showErrorMsg } from "./shared"
import type { RequestInstanceState } from "./type"

export const request = createFlatRequest(
  {
    baseURL: "/api/v1",
    headers: {
      apifoxToken: "XL299LiMEDZ0H5h3A29PxwQXdMJqWyY2",
    },
  },
  {
    defaultState: {
      errMsgStack: [],
      refreshTokenPromise: null,
    } as RequestInstanceState,
    transform(response: AxiosResponse<App.Service.Response<any>>) {
      return response.data.data
    },
    async onRequest(config) {
      const Authorization = getAuthorization()
      Object.assign(config.headers, { Authorization })

      return config
    },
    isBackendSuccess(response) {
      // when the backend response code is "0000"(default), it means the request is success
      // to change this logic by yourself, you can modify the `VITE_SERVICE_SUCCESS_CODE` in `.env` file
      return String(response.data.code) === import.meta.env.VITE_SERVICE_SUCCESS_CODE
    },
    async onBackendFail(response, instance) {
      const authStore = useAuthStore()
      const responseCode = String(response.data.code)

      function handleLogout() {
        authStore.resetStore()
      }

      // when the backend response code is in `logoutCodes`, it means the user will be logged out and redirected to login page
      const logoutCodes = import.meta.env.VITE_SERVICE_LOGOUT_CODES?.split(",") || []
      if (logoutCodes.includes(responseCode)) {
        handleLogout()
        return null
      }

      return null
    },
    onError(error) {
      // when the request is fail, you can show error message

      let message = error.message
      let backendErrorCode = ""

      // get backend error message and code
      if (error.code === BACKEND_ERROR_CODE) {
        message = error.response?.data?.msg || message
        backendErrorCode = String(error.response?.data?.code || "")
      }

      showErrorMsg(request.state, message)
    },
  },
)
