/**
 * API类型定义 - 基于后端DTO结构
 */

// Sign相关类型
export interface SignCreateRequest {
  signifier: string;
  signified: string;
}

export interface SignUpdateRequest {
  signifier: string;
  signified: string;
}

export interface SignResponse {
  id: number;
  signifier: string;
  signified: string;
}

// Todo相关类型
export interface TodoCreateRequest {
  eventId?: number;
  taskId: number;
  status: string;
  plannedStart?: string;
  plannedEnd?: string;
  actualStart?: string;
  actualEnd?: string;
}

export interface TodoUpdateRequest {
  eventId?: number;
  taskId: number;
  status: string;
  plannedStart?: string;
  plannedEnd?: string;
  actualStart?: string;
  actualEnd?: string;
}

export interface TimeSpan {
  start?: string;
  end?: string;
}

export interface TodoResponse {
  id: number;
  eventId?: number;
  taskId: number;
  status: string;
  plannedTime: TimeSpan;
  actualTime: TimeSpan;
  completedTime?: string;
}

// Event相关类型
export interface EventCreateRequest {
  isTemplate: boolean;
  title: string;
  description: string;
  location: string;
  priority: number;
  category: string;
}

export interface EventUpdateRequest {
  isTemplate: boolean;
  title: string;
  description: string;
  location: string;
  priority: number;
  category: string;
}

export interface EventResponse {
  id: number;
  isTemplate: boolean;
  title: string;
  description: string;
  location: string;
  priority: number;
  category: string;
}

// Task相关类型
export interface TaskCreateRequest {
  eventId: number;
  parentTaskId?: number;
  preTaskIds?: number[];
  description: string;
  allowedStart?: string;
  allowedEnd?: string;
  plannedStart?: string;
  plannedEnd?: string;
  status: string;
}

export interface TaskUpdateRequest {
  eventId: number;
  parentTaskId?: number;
  preTaskIds?: number[];
  description: string;
  allowedStart?: string;
  allowedEnd?: string;
  plannedStart?: string;
  plannedEnd?: string;
  status: string;
}

export interface TaskResponse {
  id: number;
  eventId: number;
  parentTaskId?: number;
  preTaskIds: number[];
  description: string;
  allowedTime: TimeSpan;
  plannedTime: TimeSpan;
  status: string;
  createdAt: string;
}