<template>
  <div class="event-task-view">
    <h2>Events and Tasks Management</h2>

    <!-- Add Task Form -->
    <AddTask @task-added="fetchTasks" />

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
              <span class="id">ID: {{ event.id }}</span>
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
          <div class="task-actions">
            <button @click="showAddTodoForm(task.id)" class="add-todo-btn">
              Add Todo
            </button>
            <div v-if="activeTodoForm === task.id" class="quick-todo-form">
              <h5>Add Todo for Task #{{ task.id }}</h5>
              <form @submit.prevent="addQuickTodo(task.id)">
                <div class="form-group">
                  <label>Status:</label>
                  <select v-model="quickTodo.status" required>
                    <option value="">Select status</option>
                    <option value="pending">Pending</option>
                    <option value="in_progress">In Progress</option>
                    <option value="completed">Completed</option>
                    <option value="cancelled">Cancelled</option>
                  </select>
                </div>
                <div class="form-group">
                  <label>Planned Start:</label>
                  <input v-model="quickTodo.plannedStart" type="datetime-local" />
                </div>
                <div class="form-group">
                  <label>Planned End:</label>
                  <input v-model="quickTodo.plannedEnd" type="datetime-local" />
                </div>
                <button type="submit" :disabled="addingTodo">Add Todo</button>
                <button type="button" @click="cancelAddTodo" class="cancel-btn">Cancel</button>
              </form>
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
import { createTodo } from '../api/todo'
import type { EventResponse, TaskResponse, TodoCreateRequest } from '../api/types'
import AddTask from './AddTask.vue'

const events = ref<EventResponse[]>([])
const tasks = ref<TaskResponse[]>([])
const eventsLoading = ref(false)
const tasksLoading = ref(false)
const eventsError = ref('')
const tasksError = ref('')

// Quick todo form state
const activeTodoForm = ref<number | null>(null)
const addingTodo = ref(false)
const quickTodo = ref({
  status: '',
  plannedStart: '',
  plannedEnd: ''
})

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

// Show add todo form for a specific task
const showAddTodoForm = (taskId: number) => {
  activeTodoForm.value = taskId
  quickTodo.value = {
    status: '',
    plannedStart: '',
    plannedEnd: ''
  }
}

// Cancel adding todo
const cancelAddTodo = () => {
  activeTodoForm.value = null
}

// Add quick todo for a specific task
const addQuickTodo = async (taskId: number) => {
  if (!quickTodo.value.status) {
    return
  }

  try {
    addingTodo.value = true
    const todoData: TodoCreateRequest = {
      taskId: taskId,
      status: quickTodo.value.status,
      plannedStart: quickTodo.value.plannedStart || undefined,
      plannedEnd: quickTodo.value.plannedEnd || undefined
    }

    await createTodo(todoData)
    activeTodoForm.value = null
    quickTodo.value = {
      status: '',
      plannedStart: '',
      plannedEnd: ''
    }
    
    // Show success message or refresh todos if needed
    alert('Todo added successfully!')
  } catch (err) {
    console.error('Error adding todo:', err)
    alert('Failed to add todo')
  } finally {
    addingTodo.value = false
  }
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
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.section h3 {
  color: #2c3e50;
  border-bottom: 3px solid #667eea;
  padding-bottom: 12px;
  margin-bottom: 24px;
  font-size: 1.5rem;
  font-weight: 600;
}

.events-list,
.tasks-list {
  display: grid;
  gap: 20px;
}

.event-item,
.task-item {
  background: white;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.event-item:hover,
.task-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.event-item h4,
.task-item h4 {
  margin-top: 0;
  color: #2c3e50;
  font-size: 1.2rem;
  font-weight: 600;
  margin-bottom: 12px;
}

.event-info p,
.task-info p {
  color: #666;
  line-height: 1.6;
  margin-bottom: 16px;
}

.event-details,
.task-details {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 16px;
  font-size: 14px;
}

.event-details span,
.task-details span {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.time-range {
  width: 100%;
  margin-top: 12px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  border-left: 4px solid #667eea;
}

.time-range span {
  background: none;
  color: #2c3e50;
  padding: 0;
  font-size: 14px;
}

.loading,
.error,
.empty {
  text-align: center;
  padding: 40px;
  border-radius: 12px;
  font-size: 16px;
}

.loading {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  color: #1976d2;
}

.error {
  background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
  color: #d32f2f;
}

.empty {
  background: linear-gradient(135deg, #f5f5f5 0%, #eeeeee 100%);
  color: #757575;
}

.task-actions {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 2px solid #eee;
}

.add-todo-btn {
  padding: 12px 24px;
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.3);
}

.add-todo-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(39, 174, 96, 0.4);
}

.quick-todo-form {
  margin-top: 20px;
  padding: 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
  border: 1px solid #dee2e6;
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.quick-todo-form h5 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #2c3e50;
  font-size: 1.1rem;
  font-weight: 600;
}

.quick-todo-form .form-group {
  margin-bottom: 16px;
}

.quick-todo-form .form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  font-weight: 600;
  color: #495057;
}

.quick-todo-form .form-group input,
.quick-todo-form .form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s ease;
}

.quick-todo-form .form-group input:focus,
.quick-todo-form .form-group select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.quick-todo-form button {
  padding: 12px 24px;
  margin-right: 12px;
  font-size: 14px;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.quick-todo-form button[type="submit"] {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.quick-todo-form button[type="submit"]:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
}

.cancel-btn {
  background: #95a5a6;
  color: white;
  border: none;
}

.cancel-btn:hover {
  background: #7f8c8d;
  transform: translateY(-2px);
}

/* Responsive design */
@media (max-width: 768px) {
  .event-task-view {
    padding: 16px;
  }
  
  .section {
    padding: 20px;
    margin-bottom: 32px;
  }
  
  .event-item,
  .task-item {
    padding: 20px;
  }
  
  .event-details,
  .task-details {
    gap: 8px;
  }
  
  .event-details span,
  .task-details span {
    padding: 4px 10px;
    font-size: 11px;
  }
  
  .quick-todo-form {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .section {
    padding: 16px;
  }
  
  .event-item,
  .task-item {
    padding: 16px;
  }
  
  .quick-todo-form .form-group input,
  .quick-todo-form .form-group select {
    padding: 10px 14px;
  }
  
  .quick-todo-form button {
    padding: 10px 20px;
    margin-right: 8px;
  }
}
</style>