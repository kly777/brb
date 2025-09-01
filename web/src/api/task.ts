import { get, post, put, del } from './api';
import type { TaskCreateRequest, TaskUpdateRequest, TaskResponse } from './types';

/**
 * Task API 函数封装
 */

/**
 * 创建新的task
 * @param data - Task创建请求数据
 * @returns 创建的task响应
 */
export function createTask(data: TaskCreateRequest): Promise<TaskResponse> {
  return post<TaskResponse>('/tasks', data);
}

/**
 * 获取所有task
 * @returns task响应数组
 */
export function getAllTasks(): Promise<TaskResponse[]> {
  return get<TaskResponse[]>('/tasks');
}

/**
 * 根据ID获取单个task
 * @param id - task的ID
 * @returns task响应数据
 */
export function getTask(id: number): Promise<TaskResponse> {
  return get<TaskResponse>(`/tasks/${id}`);
}

/**
 * 更新task
 * @param id - task的ID
 * @param data - Task更新请求数据
 * @returns 空响应
 */
export function updateTask(id: number, data: TaskUpdateRequest): Promise<void> {
  return put<void>(`/tasks/${id}`, data);
}

/**
 * 删除task
 * @param id - task的ID
 * @returns 空响应
 */
export function deleteTask(id: number): Promise<void> {
  return del<void>(`/tasks/${id}`);
}

export default {
  createTask,
  getAllTasks,
  getTask,
  updateTask,
  deleteTask,
};