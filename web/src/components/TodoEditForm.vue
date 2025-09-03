<template>
  <div class="edit-form">
    <h4>Edit Todo #{{ todo.id }}</h4>
    <div class="form-group">
      <label>Task:</label>
      <select v-model.number="formData.taskId" required :disabled="loading">
        <option v-if="loading" value="">Loading tasks...</option>
        <option v-else value="">Select a task</option>
        <option
          v-for="task in tasks"
          :key="task.id"
          :value="task.id"
        >
          {{ task.description }} (ID: {{ task.id }})
        </option>
      </select>
    </div>
    <div class="form-group">
      <label>Status:</label>
      <select v-model="formData.status" required>
        <option value="pending">Pending</option>
        <option value="doing">In Progress</option>
        <option value="done">Completed</option>
        <option value="cancelled">Cancelled</option>
      </select>
    </div>
    <div class="form-group">
      <label>Planned Start:</label>
      <input v-model="formData.plannedStart" type="datetime-local" />
    </div>
    <div class="form-group">
      <label>Planned End:</label>
      <input v-model="formData.plannedEnd" type="datetime-local" />
    </div>
    <div class="form-group">
      <label>Actual Start:</label>
      <input v-model="formData.actualStart" type="datetime-local" />
    </div>
    <div class="form-group">
      <label>Actual End:</label>
      <input v-model="formData.actualEnd" type="datetime-local" />
    </div>
    <div class="form-group">
      <label>Completed Time:</label>
      <input v-model="formData.completedTime" type="datetime-local" />
    </div>
    <div class="edit-actions">
      <button @click="handleSave" :disabled="saving">Save</button>
      <button @click="$emit('cancel')" class="cancel-btn">Cancel</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import type { TodoResponse } from '../api/types'
import { getAllTasks } from '../api/task'
import type { TaskResponse } from '../api/types'

const props = defineProps<{
  todo: TodoResponse
}>()

const emit = defineEmits<{
  (e: 'save', data: any): void
  (e: 'cancel'): void
}>()

const saving = ref(false)
const tasks = ref<TaskResponse[]>([])
const loading = ref(false)

const formData = ref({
  taskId: props.todo.taskId,
  status: props.todo.status,
  plannedStart: props.todo.plannedTime.start || '',
  plannedEnd: props.todo.plannedTime.end || '',
  actualStart: props.todo.actualTime.start || '',
  actualEnd: props.todo.actualTime.end || '',
  completedTime: props.todo.completedTime || ''
})

// Fetch all tasks when component is mounted
onMounted(async () => {
  loading.value = true
  try {
    tasks.value = await getAllTasks()
  } catch (error) {
    console.error('Failed to fetch tasks:', error)
  } finally {
    loading.value = false
  }
})

// Update form data when todo prop changes
watch(() => props.todo, (newTodo) => {
  formData.value = {
    taskId: newTodo.taskId,
    status: newTodo.status,
    plannedStart: newTodo.plannedTime.start || '',
    plannedEnd: newTodo.plannedTime.end || '',
    actualStart: newTodo.actualTime.start || '',
    actualEnd: newTodo.actualTime.end || '',
    completedTime: newTodo.completedTime || ''
  }
})

const handleSave = () => {
  if (!formData.value.taskId || !formData.value.status) {
    // You might want to handle validation error here
    return
  }

  // Convert datetime strings to ISO format
  const toISODate = (dateString: string) => {
    if (!dateString) return undefined
    const date = new Date(dateString)
    return isNaN(date.getTime()) ? undefined : date.toISOString()
  }

  const updateData = {
    taskId: formData.value.taskId,
    status: formData.value.status,
    plannedStart: toISODate(formData.value.plannedStart),
    plannedEnd: toISODate(formData.value.plannedEnd),
    actualStart: toISODate(formData.value.actualStart),
    actualEnd: toISODate(formData.value.actualEnd),
    completedTime: toISODate(formData.value.completedTime)
  }

  saving.value = true
  emit('save', updateData)
  saving.value = false
}
</script>

<style scoped>
.edit-form {
  flex: 1;
  background: #fff8e1;
  padding: 15px;
  border-radius: 8px;
  margin-right: 10px;
}

.edit-form h4 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #2c3e50;
}

.edit-form .form-group {
  margin-bottom: 10px;
}

.edit-form .form-group label {
  display: block;
  margin-bottom: 3px;
  font-size: 12px;
  font-weight: 500;
  color: #555;
}

.edit-form .form-group input,
.edit-form .form-group select {
  width: 100%;
  padding: 6px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 12px;
}

.edit-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
}

.edit-actions button {
  padding: 8px 16px;
  font-size: 12px;
}

.cancel-btn {
  background: #95a5a6;
}

.cancel-btn:hover {
  background: #7f8c8d;
}
</style>