<template>
  <div class="add-task-form">
    <h3>Add New Task</h3>
    <form @submit.prevent="addTask">
      <div class="form-group">
        <label for="eventId">Event:</label>
        <select id="eventId" v-model.number="newTask.eventId" required :disabled="loading">
          <option v-if="loading" value="">Loading events...</option>
          <option v-else value="">Select an event</option>
          <option
            v-for="event in events"
            :key="event.id"
            :value="event.id"
          >
            {{ event.title }} (ID: {{ event.id }})
          </option>
        </select>
      </div>

      <div class="form-group">
        <label for="description">Description:</label>
        <textarea id="description" v-model="newTask.description" required placeholder="Enter task description"></textarea>
      </div>

      <div class="form-group">
        <label for="status">Status:</label>
        <select id="status" v-model="newTask.status" required>
          <option value="">Select status</option>
          <option value="pending">Pending</option>
          <option value="in_progress">In Progress</option>
          <option value="completed">Completed</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>

      <div class="form-group">
        <label for="allowedStart">Allowed Start (optional):</label>
        <input id="allowedStart" v-model="newTask.allowedStart" type="datetime-local" />
      </div>

      <div class="form-group">
        <label for="allowedEnd">Allowed End (optional):</label>
        <input id="allowedEnd" v-model="newTask.allowedEnd" type="datetime-local" />
      </div>

      <div class="form-group">
        <label for="plannedStart">Planned Start (optional):</label>
        <input id="plannedStart" v-model="newTask.plannedStart" type="datetime-local" />
      </div>

      <div class="form-group">
        <label for="plannedEnd">Planned End (optional):</label>
        <input id="plannedEnd" v-model="newTask.plannedEnd" type="datetime-local" />
      </div>

      <button type="submit" :disabled="adding">Add Task</button>
      <div v-if="error" class="error">{{ error }}</div>
      <div v-if="success" class="success">Task added successfully!</div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { createTask } from '../api/task'
import { getAllEvents } from '../api/event'
import type { TaskCreateRequest, EventResponse } from '../api/types'

const emit = defineEmits<{
  (e: 'task-added'): void
}>()

const newTask = ref({
  eventId: null as number | null,
  description: '',
  status: '',
  allowedStart: '',
  allowedEnd: '',
  plannedStart: '',
  plannedEnd: ''
})

const adding = ref(false)
const error = ref('')
const success = ref(false)
const events = ref<EventResponse[]>([])
const loading = ref(false)

// Fetch all events when component is mounted
onMounted(async () => {
  loading.value = true
  try {
    events.value = await getAllEvents()
  } catch (err) {
    console.error('Failed to fetch events:', err)
    error.value = 'Failed to load events'
  } finally {
    loading.value = false
  }
})

const addTask = async () => {
  if (!newTask.value.eventId || !newTask.value.description || !newTask.value.status) {
    error.value = 'Event ID, Description, and Status are required'
    return
  }

  try {
    adding.value = true
    error.value = ''
    success.value = false

    const taskData: TaskCreateRequest = {
      eventId: newTask.value.eventId,
      description: newTask.value.description,
      status: newTask.value.status,
      allowedStart: newTask.value.allowedStart || undefined,
      allowedEnd: newTask.value.allowedEnd || undefined,
      plannedStart: newTask.value.plannedStart || undefined,
      plannedEnd: newTask.value.plannedEnd || undefined
    }

    await createTask(taskData)
    success.value = true
    emit('task-added')

    // Reset form
    newTask.value = {
      eventId: null,
      description: '',
      status: '',
      allowedStart: '',
      allowedEnd: '',
      plannedStart: '',
      plannedEnd: ''
    }
  } catch (err) {
    error.value = 'Failed to add task'
    console.error('Error adding task:', err)
  } finally {
    adding.value = false
  }
}
</script>

<style scoped>
.add-task-form {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}

.add-task-form h3 {
  margin-top: 0;
  color: #2c3e50;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group textarea {
  min-height: 80px;
  resize: vertical;
}

button {
  padding: 10px 20px;
  background: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

button:hover {
  background: #2980b9;
}

button:disabled {
  background: #bdc3c7;
  cursor: not-allowed;
}

.error {
  color: #e74c3c;
  margin-top: 10px;
}

.success {
  color: #27ae60;
  margin-top: 10px;
}
</style>