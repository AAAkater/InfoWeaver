import type { Ref } from "vue"
import type {
  AxiosError,
  CustomAxiosRequestConfig,
  MappedType,
  RequestInstanceCommon,
  ResponseType,
} from "@sa/axios"

export type HookRequestInstanceResponseSuccessData<ApiData> = {
  data: Ref<ApiData>
  error: Ref<null>
}

export type HookRequestInstanceResponseFailData<ResponseData> = {
  data: Ref<null>
  error: Ref<AxiosError<ResponseData>>
}

export type HookRequestInstanceResponseData<ResponseData, ApiData> = {
  loading: Ref<boolean>
} & (
  | HookRequestInstanceResponseSuccessData<ApiData>
  | HookRequestInstanceResponseFailData<ResponseData>
)

export interface HookRequestInstance<
  ResponseData,
  ApiData,
  State extends Record<string, unknown>,
> extends RequestInstanceCommon<State> {
  <T extends ApiData = ApiData, R extends ResponseType = "json">(
    config: CustomAxiosRequestConfig,
  ): HookRequestInstanceResponseData<ResponseData, MappedType<R, T>>
}
