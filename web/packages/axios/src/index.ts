import axios, { AxiosError } from "axios"
import type {
  AxiosInstance,
  AxiosResponse,
  CreateAxiosDefaults,
  InternalAxiosRequestConfig,
} from "axios"
import axiosRetry from "axios-retry"
import { stringify } from "qs"
import { nanoid } from "@sa/utils"
import { isHttpSuccess, transformResponse } from "./shared"
import { BACKEND_ERROR_CODE, REQUEST_ID_KEY } from "./constant"
import type {
  CustomAxiosRequestConfig,
  FlatRequestInstance,
  MappedType,
  RequestOption,
  ResponseType,
} from "./type"

export class CommonRequest<
  ResponseData,
  ApiData = ResponseData,
  State extends Record<string, unknown> = Record<string, unknown>,
> {
  private abortControllerMap = new Map<string, AbortController>()

  constructor(
    axiosConfig?: CreateAxiosDefaults,
    options?: Partial<RequestOption<ResponseData, ApiData, State>>,
  ) {
    const opts: RequestOption<ResponseData, ApiData, State> = Object.assign(
      {
        defaultState: {} as State,
        transform: async (response) => response.data as unknown as ApiData,
        onRequest: async (config) => config,
        isBackendSuccess: () => true,
        onBackendFail: async () => {},
        onError: async () => {},
      },
      options,
    )
    const instance = this.#createAxiosInstance(axiosConfig)

    this.#setupRequestInterceptor(instance, opts)
    this.#setupResponseInterceptor(instance, opts)

    const flatRequest = this.#createCallable(instance, opts)
    flatRequest.cancelAllRequest = this.#cancelAllRequest
    flatRequest.state = { ...opts.defaultState } as State

    return flatRequest as FlatRequestInstance<ResponseData, ApiData, State> & this
  }

  #createAxiosInstance(axiosConfig?: CreateAxiosDefaults) {
    const TEN_SECONDS = 10 * 1000
    const defaultConfig: CreateAxiosDefaults = {
      timeout: TEN_SECONDS,
      headers: { "Content-Type": "application/json" },
      validateStatus: isHttpSuccess,
      paramsSerializer: (params: any) => stringify(params),
    }
    Object.assign(defaultConfig, axiosConfig)

    const instance = axios.create(defaultConfig)
    axiosRetry(instance, Object.assign({ retries: 0 }, defaultConfig))
    return instance
  }

  #setupRequestInterceptor(
    instance: AxiosInstance,
    opts: RequestOption<ResponseData, ApiData, State>,
  ) {
    instance.interceptors.request.use((conf) => {
      const config: InternalAxiosRequestConfig = { ...conf }

      const requestId = nanoid()
      config.headers.set(REQUEST_ID_KEY, requestId)

      if (!config.signal) {
        const ac = new AbortController()
        config.signal = ac.signal
        this.abortControllerMap.set(requestId, ac)
      }

      return opts.onRequest?.(config) || config
    })
  }

  #setupResponseInterceptor(
    instance: AxiosInstance,
    opts: RequestOption<ResponseData, ApiData, State>,
  ) {
    instance.interceptors.response.use(
      async (response) => {
        const responseType: ResponseType = (response.config?.responseType as ResponseType) || "json"

        await transformResponse(response)

        if (responseType !== "json" || opts.isBackendSuccess(response)) {
          return response
        }

        const fail = await opts.onBackendFail(response, instance)
        if (fail) return fail

        const backendError = new AxiosError<ResponseData>(
          "the backend request error",
          BACKEND_ERROR_CODE,
          response.config,
          response.request,
          response,
        )
        await opts.onError(backendError)
        return Promise.reject(backendError)
      },
      async (error: AxiosError<ResponseData>) => {
        await opts.onError(error)
        return Promise.reject(error)
      },
    )
  }

  #createCallable(
    instance: AxiosInstance,
    opts: RequestOption<ResponseData, ApiData, State>,
  ): FlatRequestInstance<ResponseData, ApiData, State> {
    return (async <T extends ApiData = ApiData, R extends ResponseType = "json">(
      config: CustomAxiosRequestConfig,
    ) => {
      try {
        const response: AxiosResponse<ResponseData> = await instance(config)

        if ((response.config?.responseType || "json") === "json") {
          return { data: await opts.transform(response), error: null, response }
        }
        return { data: response.data as MappedType<R, T>, error: null, response }
      } catch (error) {
        return { data: null, error, response: (error as AxiosError<ResponseData>).response }
      }
    }) as FlatRequestInstance<ResponseData, ApiData, State>
  }

  #cancelAllRequest = () => {
    this.abortControllerMap.forEach((ac) => ac.abort())
    this.abortControllerMap.clear()
  }
}

export { BACKEND_ERROR_CODE, REQUEST_ID_KEY }
export type * from "./type"
export type { CreateAxiosDefaults, AxiosError }
