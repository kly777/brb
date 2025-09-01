<template>
  <div class="todo-calendar">
    <div class="calendar-header">
      <h2>Todo Calendar</h2>
      <div class="controls">
        <button @click="prevMonth">← Previous</button>
        <span class="current-month">{{ currentMonth }}</span>
        <button @click="nextMonth">Next →</button>
      </div>
    </div>

    <FullCalendar
      ref="calendarRef"
      :options="calendarOptions"
      class="calendar-container"
    />

    <div v-if="loading" class="loading">Loading todos...</div>
    <div v-if="error" class="error">{{ error }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import timeGridPlugin from '@fullcalendar/timegrid'
import interactionPlugin from '@fullcalendar/interaction'
import type { EventInput } from '@fullcalendar/core'
import { getAllTodos } from '../api/todo'
import type { TodoResponse } from '../api/types'

const calendarRef = ref()
const todos = ref<TodoResponse[]>([])
const loading = ref(false)
const error = ref('')
const currentDate = ref(new Date())

// 获取当前月份显示文本
const currentMonth = computed(() => {
  return currentDate.value.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long'
  })
})

// 转换todo到日历事件
const todoEvents = computed<EventInput[]>(() => {
  return todos.value.map(todo => ({
    id: todo.id.toString(),
    title: `Todo #${todo.id}`,
    start: todo.plannedTime.start || new Date().toISOString(),
    end: todo.plannedTime.end || undefined,
    extendedProps: {
      status: todo.status,
      taskId: todo.taskId,
      actualStart: todo.actualTime.start,
      actualEnd: todo.actualTime.end,
      completedTime: todo.completedTime
    },
    backgroundColor: getStatusColor(todo.status),
    borderColor: getStatusColor(todo.status),
    textColor: '#fff'
  }))
})

// 根据状态获取颜色
const getStatusColor = (status: string): string => {
  const colors: { [key: string]: string } = {
    pending: '#ff6b6b',
    in_progress: '#4ecdc4',
    completed: '#1dd1a1',
    cancelled: '#8395a7'
  }
  return colors[status] || '#3498db'
}

// 日历配置
const calendarOptions = {
  plugins: [dayGridPlugin, timeGridPlugin, interactionPlugin],
  initialView: 'dayGridMonth',
  headerToolbar: false as const,
  events: todoEvents.value,
  eventTimeFormat: {
    hour: '2-digit' as const,
    minute: '2-digit' as const,
    hour12: false
  },
  eventClick: (info: any) => {
    const event = info.event
    console.log('Todo clicked:', event.extendedProps)
    // 这里可以添加点击事件的处理逻辑，比如显示详情弹窗
  }
}

// 获取todos数据
const fetchTodos = async () => {
  try {
    loading.value = true
    error.value = ''
    const response = await getAllTodos()
    todos.value = response
  } catch (err) {
    error.value = 'Failed to load todos'
    console.error('Error fetching todos:', err)
  } finally {
    loading.value = false
  }
}

// 切换到上个月
const prevMonth = () => {
  const calendarApi = calendarRef.value?.getApi()
  if (calendarApi) {
    calendarApi.prev()
    currentDate.value = calendarApi.getDate()
  }
}

// 切换到下个月
const nextMonth = () => {
  const calendarApi = calendarRef.value?.getApi()
  if (calendarApi) {
    calendarApi.next()
    currentDate.value = calendarApi.getDate()
  }
}

onMounted(() => {
  fetchTodos()
})
</script>

<style scoped>
.todo-calendar {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
}

.calendar-header h2 {
  margin: 0;
  color: #2c3e50;
}

.controls {
  display: flex;
  align-items: center;
  gap: 15px;
}

.controls button {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.controls button:hover {
  background: #f0f0f0;
}

.current-month {
  font-size: 1.2em;
  font-weight: 500;
  color: #2c3e50;
}

.calendar-container {
  height: 600px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.loading, .error {
  text-align: center;
  padding: 20px;
  margin-top: 20px;
  border-radius: 8px;
}

.loading {
  background: #e3f2fd;
  color: #1976d2;
}

.error {
  background: #ffebee;
  color: #d32f2f;
}

:deep(.fc-event) {
  border-radius: 4px;
  font-size: 12px;
  padding: 2px 4px;
}

:deep(.fc-daygrid-event) {
  margin-bottom: 2px;
}
</style>