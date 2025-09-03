<template>
  <div class="todo-list">
    <h2>Todo Management</h2>

    <!-- Add Todo Form -->
    <div class="add-todo-form">
      <h3>Add New Todo</h3>
      <form @submit.prevent="addTodo">
        <div class="form-group">
          <label for="taskId">Task:</label>
          <select id="taskId" v-model.number="newTodo.taskId" required>
            <option value="">Select a task</option>
            <option v-for="task in tasks" :key="task.id" :value="task.id">
              {{ task.description }} (ID: {{ task.id }})
            </option>
          </select>
        </div>

        <div class="form-group">
          <label for="eventId">Event (optional):</label>
          <select id="eventId" v-model.number="newTodo.eventId">
            <option value="">Select an event (optional)</option>
            <option v-for="event in events" :key="event.id" :value="event.id">
              {{ event.title }} (ID: {{ event.id }})
            </option>
          </select>
        </div>

        <div class="form-group">
          <label for="status">Status:</label>
          <select id="status" v-model="newTodo.status" required>
            <option value="">Select status</option>
            <option value="pending">Pending</option>
            <option value="in_progress">In Progress</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>

        <div class="form-group">
          <label for="plannedStart">Planned Start (optional):</label>
          <input id="plannedStart" v-model="newTodo.plannedStart" type="datetime-local" />
        </div>

        <div class="form-group">
          <label for="plannedEnd">Planned End (optional):</label>
          <input id="plannedEnd" v-model="newTodo.plannedEnd" type="datetime-local" />
        </div>

        <button type="submit" :disabled="adding">Add Todo</button>
      </form>
    </div>

    <!-- Todo List -->
    <div class="todos-container">
      <h3>Existing Todos</h3>
      <div v-if="loading" class="loading">Loading todos...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else-if="todos.length === 0" class="empty">No todos found.</div>
      <ul v-else class="todos">
        <li v-for="todo in todos" :key="todo.id" class="todo-item">
          <div v-if="editingTodoId !== todo.id" class="todo-info">
            <span class="todo-id">#{{ todo.id }}</span>
            <span class="todo-task">Task: {{ getTaskDescription(todo.taskId) }}</span>
            <span v-if="todo.eventId" class="todo-event">Event: {{ getEventTitle(todo.eventId) }}</span>
            <span :class="['todo-status', todo.status]">{{ todo.status }}</span>
            <div v-if="todo.plannedTime.start" class="todo-time">
              Planned: {{ formatDate(todo.plannedTime.start) }}
              <span v-if="todo.plannedTime.end">- {{ formatDate(todo.plannedTime.end) }}</span>
            </div>
            <div v-if="todo.actualTime.start" class="todo-time">
              Actual: {{ formatDate(todo.actualTime.start) }}
              <span v-if="todo.actualTime.end">- {{ formatDate(todo.actualTime.end) }}</span>
            </div>
            <div v-if="todo.completedTime" class="todo-time">
              Completed: {{ formatDate(todo.completedTime) }}
            </div>
          </div>

          <TodoEditForm v-if="editingTodoId === todo.id" :todo="todo"
            @save="(data: TodoUpdateRequest) => saveEdit(todo.id, data)"
            @cancel="cancelEdit" />

          <div class="todo-actions" v-if="editingTodoId !== todo.id">
            <button @click="editingTodoId = todo.id" class="edit-btn">Edit</button>
            <button @click="deleteTodo(todo.id)" class="delete-btn">Delete</button>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAllTodos, createTodo, updateTodo, deleteTodo as _deleteTodo } from '../api/todo'
import { getAllTasks } from '../api/task'
import { getAllEvents } from '../api/event'
import type { TodoResponse, TodoCreateRequest, TodoUpdateRequest, TaskResponse, EventResponse } from '../api/types'
import TodoEditForm from './TodoEditForm.vue'

const todos = ref<TodoResponse[]>([])
const tasks = ref<TaskResponse[]>([])
const events = ref<EventResponse[]>([])
const loading = ref(false)
const error = ref('')
const adding = ref(false)
const updating = ref(false)
const editingTodoId = ref<number | null>(null)

const newTodo = ref({
  taskId: null as number | null,
  eventId: null as number | null,
  status: '',
  plannedStart: '',
  plannedEnd: ''
})


// Fetch todos on mount
onMounted(() => {
  fetchData()
})

// Fetch all data (todos, tasks, events)
const fetchData = async () => {
  try {
    loading.value = true
    error.value = ''
    const [todosResponse, tasksResponse, eventsResponse] = await Promise.all([
      getAllTodos(),
      getAllTasks(),
      getAllEvents()
    ])
    todos.value = todosResponse
    tasks.value = tasksResponse
    events.value = eventsResponse
  } catch (err) {
    error.value = 'Failed to load data'
    console.error('Error fetching data:', err)
  } finally {
    loading.value = false
  }
}

// Fetch only todos (for refreshing after operations)
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

// Add new todo
const addTodo = async () => {
  if (!newTodo.value.taskId || !newTodo.value.status) {
    error.value = 'Task ID and Status are required'
    return
  }

  try {
    adding.value = true
    error.value = ''

    const todoData: TodoCreateRequest = {
      taskId: newTodo.value.taskId,
      eventId: newTodo.value.eventId || undefined,
      status: newTodo.value.status,
      plannedStart: newTodo.value.plannedStart || undefined,
      plannedEnd: newTodo.value.plannedEnd || undefined
    }

    await createTodo(todoData)
    await fetchTodos() // Refresh the list

    // Reset form
    newTodo.value = {
      taskId: null,
      eventId: null,
      status: '',
      plannedStart: '',
      plannedEnd: ''
    }
  } catch (err) {
    error.value = 'Failed to add todo'
    console.error('Error adding todo:', err)
  } finally {
    adding.value = false
  }
}

// Delete todo
const deleteTodo = async (id: number) => {
  if (!confirm('Are you sure you want to delete this todo?')) {
    return
  }

  try {
    error.value = ''
    await _deleteTodo(id)
    await fetchTodos() // Refresh the list
  } catch (err) {
    error.value = 'Failed to delete todo'
    console.error('Error deleting todo:', err)
  }
}

// Format date for display
const formatDate = (dateString: string | undefined) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleString('zh-CN')
}


// Save edited todo
const saveEdit = async (id: number, data: TodoUpdateRequest) => {
  try {
    updating.value = true
    error.value = ''
    await updateTodo(id, data)
    editingTodoId.value = null
    await fetchTodos() // Refresh the list
  } catch (err) {
    error.value = 'Failed to update todo'
    console.error('Error updating todo:', err)
  } finally {
    updating.value = false
  }
}

// Cancel editing
const cancelEdit = () => {
  editingTodoId.value = null
}

// Helper function to get task description by ID
const getTaskDescription = (taskId: number): string => {
  const task = tasks.value.find(t => t.id === taskId)
  return task ? task.description : `Task #${taskId}`
}

// Helper function to get event title by ID
const getEventTitle = (eventId: number): string => {
  const event = events.value.find(e => e.id === eventId)
  return event ? event.title : `Event #${eventId}`
}
</script>

<style scoped>
.todo-list {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.add-todo-form {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}

.add-todo-form h3 {
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
.form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
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

.todos-container h3 {
  color: #2c3e50;
  margin-bottom: 15px;
}

.todos {
  list-style: none;
  padding: 0;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  margin-bottom: 10px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.todo-info {
  flex: 1;
}

.todo-id {
  font-weight: bold;
  color: #7f8c8d;
  margin-right: 10px;
}

.todo-task {
  font-weight: 500;
  margin-right: 10px;
}

.todo-status {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  color: white;
}

.todo-status.pending {
  background: #ff6b6b;
}

.todo-status.in_progress {
  background: #4ecdc4;
}

.todo-status.completed {
  background: #1dd1a1;
}

.todo-status.cancelled {
  background: #8395a7;
}

.todo-time {
  font-size: 12px;
  color: #7f8c8d;
  margin-top: 5px;
}

.todo-event {
  font-style: italic;
  color: #95a5a6;
  margin-right: 10px;
}

.delete-btn {
  background: #e74c3c;
  padding: 6px 12px;
  font-size: 12px;
}

.delete-btn:hover {
  background: #c0392b;
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

/* Edit form styles */
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

/* Todo actions styles */
.todo-actions {
  display: flex;
  gap: 8px;
}

.edit-btn {
  background: #f39c12;
  padding: 6px 12px;
  font-size: 12px;
}

.edit-btn:hover {
  background: #e67e22;
}

.delete-btn {
  background: #e74c3c;
  padding: 6px 12px;
  font-size: 12px;
}

.delete-btn:hover {
  background: #c0392b;
}
</style>