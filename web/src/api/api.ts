/**
 * 基础API请求封装
 */

const BASE_URL = '/api';

/**
 * 请求选项接口
 */
interface RequestOptions extends RequestInit {
  params?: Record<string, string | number>;
}

/**
 * 处理HTTP请求
 * @param endpoint - API端点
 * @param options - 请求选项
 * @returns 响应数据
 */
async function request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const { params, ...restOptions } = options;

  // 构建URL，处理查询参数
  let url = `${BASE_URL}${endpoint}`;
  if (params && Object.keys(params).length > 0) {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        searchParams.append(key, String(value));
      }
    });
    url += `?${searchParams.toString()}`;
  }

  // 设置默认请求头
  const headers = {
    'Content-Type': 'application/json',
    ...restOptions.headers,
  };

  try {
    const response = await fetch(url, {
      ...restOptions,
      headers,
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    // 处理204 No Content响应
    if (response.status === 204) {
      return null as T;
    }

    const data = await response.json();
    return data;
  } catch (error) {

    console.error('API请求失败:', error);
    throw error;
  }
}

/**
 * GET请求
 * @param endpoint - API端点
 * @param params - 查询参数
 * @returns 响应数据
 */
export function get<T>(endpoint: string, params?: Record<string, string | number>): Promise<T> {
  return request<T>(endpoint, { method: 'GET', params });
}

/**
 * POST请求
 * @param endpoint - API端点
 * @param data - 请求体数据
 * @returns 响应数据
 */
export function post<T>(endpoint: string, data?: any): Promise<T> {
  return request<T>(endpoint, {
    method: 'POST',
    body: data ? JSON.stringify(data) : undefined,
  });
}

/**
 * PUT请求
 * @param endpoint - API端点
 * @param data - 请求体数据
 * @returns 响应数据
 */
export function put<T>(endpoint: string, data?: any): Promise<T> {
  return request<T>(endpoint, {
    method: 'PUT',
    body: data ? JSON.stringify(data) : undefined,
  });
}

/**
 * DELETE请求
 * @param endpoint - API端点
 * @returns 响应数据
 */
export function del<T>(endpoint: string): Promise<T> {
  return request<T>(endpoint, { method: 'DELETE' });
}

export default {
  get,
  post,
  put,
  delete: del,
};