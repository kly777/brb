import { get, post, put, del } from './api';
import type { EventCreateRequest, EventUpdateRequest, EventResponse } from './types';

/**
 * Event API 函数封装
 */

/**
 * 创建新的event
 * @param data - Event创建请求数据
 * @returns 创建的event响应
 */
export function createEvent(data: EventCreateRequest): Promise<EventResponse> {
  return post<EventResponse>('/events', data);
}

/**
 * 获取所有event
 * @returns event响应数组
 */
export function getAllEvents(): Promise<EventResponse[]> {
  return get<EventResponse[]>('/events');
}

/**
 * 根据ID获取单个event
 * @param id - event的ID
 * @returns event响应数据
 */
export function getEvent(id: number): Promise<EventResponse> {
  return get<EventResponse>(`/events/${id}`);
}

/**
 * 更新event
 * @param id - event的ID
 * @param data - Event更新请求数据
 * @returns 空响应
 */
export function updateEvent(id: number, data: EventUpdateRequest): Promise<void> {
  return put<void>(`/events/${id}`, data);
}

/**
 * 删除event
 * @param id - event的ID
 * @returns 空响应
 */
export function deleteEvent(id: number): Promise<void> {
  return del<void>(`/events/${id}`);
}

export default {
  createEvent,
  getAllEvents,
  getEvent,
  updateEvent,
  deleteEvent,
};