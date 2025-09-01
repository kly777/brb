
import { get, post, put, del } from './api';
import type { TodoCreateRequest, TodoUpdateRequest, TodoResponse } from './types';

/**
 * Todo API 函数封装
 */

/**
 * 创建新的todo
 * @param data - Todo创建请求数据
 * @returns 创建的todo响应
 */
export function createTodo(data: TodoCreateRequest): Promise<TodoResponse> {
  return post<TodoResponse>('/todos', data);
}

/**
 * 获取所有todo
 * @returns todo响应数组
 */
export function getAllTodos(): Promise<TodoResponse[]> {
  return get<TodoResponse[]>('/todos');
}

/**
 * 根据ID获取单个todo
 * @param id - todo的ID
 * @returns todo响应数据
 */
export function getTodo(id: number): Promise<TodoResponse> {
  return get<TodoResponse>(`/todos/${id}`);
}

/**
 * 更新todo
 * @param id - todo的ID
 * @param data - Todo更新请求数据
 * @returns 空响应
 */
export function updateTodo(id: number, data: TodoUpdateRequest): Promise<void> {
  return put<void>(`/todos/${id}`, data);
}

/**
 * 删除todo
 * @param id - todo的ID
 * @returns 空响应
 */
export function deleteTodo(id: number): Promise<void> {
  return del<void>(`/todos/${id}`);
}

export default {
  createTodo,
  getAllTodos,
  getTodo,
  updateTodo,
  deleteTodo,
};
