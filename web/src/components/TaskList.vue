<template>
  <div class="task-list">
    <h2>Tasks</h2>
    <div class="task-form">
      <input v-model="newTask.eventId" placeholder="Event ID" />
      <input v-model="newTask.subTaskId" placeholder="Sub Task ID (optional)" />
      <input v-model="newTask.description" placeholder="Description" />
      <input v-model="newTask.startTime" type="datetime-local" />
      <input v-model="newTask.endTime" type="datetime-local" />
      <input v-model="newTask.estimateTime" type="number" placeholder="Estimate Time (minutes)" />
      <button @click="createTask">Create Task</button>
    </div>
    <ul class="task-items">
      <li v-for="task in tasks" :key="task.id" class="task-item">
        <div class="task-info">
          <h3>{{ task.description }}</h3>
          <p>Event ID: {{ task.eventId }}</p>
          <p v-if="task.subTaskId">Sub Task ID: {{ task.subTaskId }}</p>
          <p>Start: {{ formatDate(task.startTime) }}</p>
          <p>End: {{ formatDate(task.endTime) }}</p>
          <p>Estimate: {{ task.estimateTime }} minutes</p>
        </div>
        <div class="task-actions">
          <button @click="editTask(task)">Edit</button>
          <button @click="deleteTask(task.id)">Delete</button>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Task {
  id: string
  eventId: string
  subTaskId: string
  description: string
  startTime: string
  endTime: string
  estimateTime: number
}

const tasks = ref<Task[]>([])
const newTask = ref({
  eventId: '',
  subTaskId: '',
  description: '',
  startTime: '',
  endTime: '',
  estimateTime: 0
})

const API_BASE = 'http://localhost:8080/api'

const fetchTasks = async () => {
  try {
    const response = await fetch(`${API_BASE}/tasks`)
    if (response.ok) {
      tasks.value = await response.json()
    }
  } catch (error) {
    console.error('Failed to fetch tasks:', error)
  }
}

const createTask = async () => {
  try {
    const response = await fetch(`${API_BASE}/tasks`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        ...newTask.value,
        estimateTime: Number(newTask.value.estimateTime)
      })
    })
    if (response.ok) {
      newTask.value = { 
        eventId: '', 
        subTaskId: '', 
        description: '', 
        startTime: '', 
        endTime: '', 
        estimateTime: 0 
      }
      fetchTasks()
    }
  } catch (error) {
    console.error('Failed to create task:', error)
  }
}

const editTask = (task: Task) => {
  console.log('Edit task:', task)
}

const deleteTask = async (id: string) => {
  try {
    const response = await fetch(`${API_BASE}/tasks/${id}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      fetchTasks()
    }
  } catch (error) {
    console.error('Failed to delete task:', error)
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

onMounted(() => {
  fetchTasks()
})
</script>

<style scoped>
.task-list {
  margin: 20px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
}

.task-form {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.task-form input {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.task-form button {
  padding: 8px 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.task-items {
  list-style: none;
  padding: 0;
}

.task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.task-info h3 {
  margin: 0 0 5px 0;
}

.task-actions {
  display: flex;
  gap: 10px;
}

.task-actions button {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.task-actions button:first-child {
  background-color: #2196F3;
  color: white;
}

.task-actions button:last-child {
  background-color: #f44336;
  color: white;
}
</style>