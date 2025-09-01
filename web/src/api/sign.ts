import { get, post, put, del } from './api';
import type { SignCreateRequest, SignUpdateRequest, SignResponse } from './types';

/**
 * Sign API 函数封装
 */

/**
 * 创建新的sign
 * @param data - Sign创建请求数据
 * @returns 创建的sign响应
 */
export function createSign(data: SignCreateRequest): Promise<SignResponse> {
  return post<SignResponse>('/signs', data);
}

/**
 * 根据ID获取单个sign
 * @param id - sign的ID
 * @returns sign响应数据
 */
export function getSign(id: number): Promise<SignResponse> {
  return get<SignResponse>(`/signs/${id}`);
}

/**
 * 更新sign
 * @param id - sign的ID
 * @param data - Sign更新请求数据
 * @returns 空响应
 */
export function updateSign(id: number, data: SignUpdateRequest): Promise<void> {
  return put<void>(`/signs/${id}`, data);
}

/**
 * 删除sign
 * @param id - sign的ID
 * @returns 空响应
 */
export function deleteSign(id: number): Promise<void> {
  return del<void>(`/signs/${id}`);
}

export default {
  createSign,
  getSign,
  updateSign,
  deleteSign,
};