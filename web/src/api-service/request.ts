import axios, { type AxiosRequestConfig } from "axios";
import { ElLoading, ElMessage } from "element-plus";
// 携带cookie
axios.defaults.withCredentials = true;
/**
 * 请求返回格式
 */
export type RequestResult<T = any> = {
  /**
   * 状态码
   */
  status: "failure" | "success";
  /**
   * 错误码
   */
  code: string;
  /**
   * 错误信息
   */
  msg: string;
  data?: T;
  [key: string]: any;
};
export type DeleteResult = {
  data: number;
} & RequestResult;

type Request = {
  <T>(config: AxiosRequestConfig): Promise<T>;
  get: <T>(url: string, data?: any, config?: AxiosRequestConfig) => Promise<T>;
  post: <T>(url: string, data?: any, config?: AxiosRequestConfig) => Promise<T>;
  put: <T>(url: string, data?: any, config?: AxiosRequestConfig) => Promise<T>;
  delete: <T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ) => Promise<T>;
};

const request: Request = <T>(config: AxiosRequestConfig): Promise<T> => {
  let url = config.url;


  return axios
    .request<T>({
      ...config,
      headers: {
        ...config.headers,
        n: "n",
      },
      url,
    })
    .then((response) => {
      return new Promise((resolve, reject) => {
        // @ts-ignore
        if (response.data?.status !== "failure") {
          // @ts-ignore
          resolve(replaceNullWithUndefined(response.data) || {});
        } else {
          const errorMsg = (response.data as any).msg;
          if (errorMsg) {
            ElMessage.error(errorMsg);
          }
          reject(response);
        }
      });
    });
};

request.get = <T>(
  url: string,
  data?: any,
  config?: AxiosRequestConfig,
): Promise<T> => {
  return request({
    ...config,
    method: "GET",
    url: url,
    params: data ? data : {},
  });
};

request.post = <T>(
  url: string,
  data?: any,
  config?: AxiosRequestConfig,
): Promise<T> => {
  return request({
    ...config,
    method: "POST",
    url: url,
    data: data ? data : {},
  });
};

request.put = <T>(
  url: string,
  data?: any,
  config?: AxiosRequestConfig,
): Promise<T> => {
  return request({
    ...config,
    method: "PUT",
    url: url,
    data: data ? data : {},
  });
};

request.delete = <T>(
  url: string,
  data?: any,
  config?: AxiosRequestConfig,
): Promise<T> => {
  return request({
    ...config,
    method: "DELETE",
    url: url,
    data: data ? data : {},
  });
};

function replaceNullWithUndefined(obj: any): any {
  if (obj === null || obj === undefined) {
    return undefined;
  }

  if (Array.isArray(obj)) {
    // 如果是数组，遍历数组的元素并递归处理
    return obj.map((item) => replaceNullWithUndefined(item));
  }

  if (typeof obj === "object") {
    // 如果是对象，遍历对象的属性并递归处理
    for (const key in obj) {
      obj[key] = replaceNullWithUndefined(obj[key]);
    }
  }

  return obj;
}

const requestWithLoading = <T>(
  f: () => Promise<T>,
  message?: string,
): Promise<T> => {
  const loadingInstance = ElLoading.service({ fullscreen: true });
  return f().finally(() => {
    loadingInstance.close();
  });
};

export { requestWithLoading };
export default request;
