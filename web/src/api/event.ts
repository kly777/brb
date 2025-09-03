import { del, get, put } from './api';
import type { EventResponse, EventUpdateRequest } from './types';

/**
 * 获取所有事件
 * @returns 事件响应数组
 */
export function getAllEvents(): Promise<EventResponse[]> {
  return get<EventResponse[]>('/events');
}

/**
 * 根据ID获取单个事件
 * @param id - 事件的ID
 * @returns 事件响应数据
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
  getAllEvents,
  getEvent
};