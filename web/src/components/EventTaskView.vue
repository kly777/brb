<template>
  <div class="event-task-view">
    <h2>Events and Tasks Management</h2>

    <!-- Events Section -->
    <div class="section">
      <h3>Events</h3>
      <div v-if="eventsLoading" class="loading">Loading events...</div>
      <div v-else-if="eventsError" class="error">{{ eventsError }}</div>
      <div v-else-if="events.length === 0" class="empty">No events found.</div>
      <div v-else class="events-list">
        <div v-for="event in events" :key="event.id" class="event-item">
          <div class="event-info">
            <h4>{{ event.title }}</h4>
            <p>{{ event.description }}</p>
            <div class="event-details">
              <span class="location">Location: {{ event.location }}</span>
              <span class="priority">Priority: {{ event.priority }}</span>
              <span class="category">Category: {{ event.category }}</span>
              <span class="template">Template: {{ event.isTemplate ? 'Yes' : 'No' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Tasks Section -->
    <div class="section">
      <h3>Tasks</h3>
      <div v-if="tasksLoading" class="loading">Loading tasks...</div>
      <div v-else-if="tasksError" class="error">{{ tasksError }}</div>
      <div v-else-if="tasks.length === 0" class="empty">No tasks found.</div>
      <div v-else class="tasks-list">
        <div v-for="task in tasks" :key="task.id" class="task-item">
          <div class="task-info">
            <h4>Task #{{ task.id }}</h4>
            <p>{{ task.description }}</p>
            <div class="task-details">
              <span class="event-id">Event ID: {{ task.eventId }}</span>
              <span class="status">Status: {{ task.status }}</span>
              <span class="created-at">Created: {{ formatDate(task.createdAt) }}</span>
              <div v-if="task.allowedTime.start || task.allowedTime.end" class="time-range">
                <span>Allowed: {{ formatTimeSpan(task.allowedTime) }}</span>
              </div>
              <div v-if="task.plannedTime.start || task.plannedTime.end" class="time-range">
                <span>Planned: {{ formatTimeSpan(task.plannedTime) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAllEvents } from '../api/event'
import { getAllTasks } from '../api/task'
import type { EventResponse, TaskResponse } from '../api/types'

const events = ref<EventResponse[]>([])
const tasks = ref<TaskResponse[]>([])
const eventsLoading = ref(false)
const tasksLoading = ref(false)
const eventsError = ref('')
const tasksError = ref('')

// Fetch events and tasks on mount
onMounted(() => {
  fetchEvents()
  fetchTasks()
})

// Fetch all events
const fetchEvents = async () => {
  try {
    eventsLoading.value = true
    eventsError.value = ''
    const response = await getAllEvents()
    events.value = response
  } catch (err) {
    eventsError.value = 'Failed to load events'
    console.error('Error fetching events:', err)
  } finally {
    eventsLoading.value = false
  }
}

// Fetch all tasks
const fetchTasks = async () => {
  try {
    tasksLoading.value = true
    tasksError.value = ''
    const response = await getAllTasks()
    tasks.value = response
  } catch (err) {
    tasksError.value = 'Failed to load tasks'
    console.error('Error fetching tasks:', err)
  } finally {
    tasksLoading.value = false
  }
}

// Format date for display
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// Format time span for display
const formatTimeSpan = (timeSpan: { start?: string; end?: string }) => {
  if (!timeSpan.start && !timeSpan.end) return 'No time set'
  
  const start = timeSpan.start ? new Date(timeSpan.start).toLocaleString('zh-CN') : 'Not set'
  const end = timeSpan.end ? new Date(timeSpan.end).toLocaleString('zh-CN') : 'Not set'
  
  return `${start} - ${end}`
}
</script>

<style scoped>
.event-task-view {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.section {
  margin-bottom: 40px;
}

.section h3 {
  color: #2c3e50;
  border-bottom: 2px solid #3498db;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.events-list,
.tasks-list {
  display: grid;
  gap: 15px;
}

.event-item,
.task-item {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.event-item h4,
.task-item h4 {
  margin-top: 0;
  color: #2c3e50;
}

.event-details,
.task-details {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 10px;
  font-size: 14px;
}

.event-details span,
.task-details span {
  background: #f8f9fa;
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid #e9ecef;
}

.time-range {
  width: 100%;
  margin-top: 5px;
}

.loading,
.error,
.empty {
  text-align: center;
  padding: 20px;
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

.empty {
  background: #f5f5f5;
  color: #757575;
}
</style>