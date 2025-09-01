<template>
  <div class="todo-list">
    <h2>Todo Management</h2>
    
    <!-- Add Todo Form -->
    <div class="add-todo-form">
      <h3>Add New Todo</h3>
      <form @submit.prevent="addTodo">
        <div class="form-group">
          <label for="taskId">Task ID:</label>
          <input
            id="taskId"
            v-model.number="newTodo.taskId"
            type="number"
            required
            placeholder="Enter task ID"
          />
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
          <input
            id="plannedStart"
            v-model="newTodo.plannedStart"
            type="datetime-local"
          />
        </div>
        
        <div class="form-group">
          <label for="plannedEnd">Planned End (optional):</label>
          <input
            id="plannedEnd"
            v-model="newTodo.plannedEnd"
            type="datetime-local"
          />
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
          <div class="todo-info">
            <span class="todo-id">#{{ todo.id }}</span>
            <span class="todo-task">Task: {{ todo.taskId }}</span>
            <span :class="['todo-status', todo.status]">{{ todo.status }}</span>
            <div v-if="todo.plannedTime.start" class="todo-time">
              Planned: {{ formatDate(todo.plannedTime.start) }}
              <span v-if="todo.plannedTime.end">- {{ formatDate(todo.plannedTime.end) }}</span>
            </div>
          </div>
          <button @click="deleteTodo(todo.id)" class="delete-btn">Delete</button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAllTodos, createTodo, deleteTodo as _deleteTodo } from '../api/todo'
import type { TodoResponse, TodoCreateRequest } from '../api/types'

const todos = ref<TodoResponse[]>([])
const loading = ref(false)
const error = ref('')
const adding = ref(false)

const newTodo = ref({
  taskId: null as number | null,
  status: '',
  plannedStart: '',
  plannedEnd: ''
})

// Fetch todos on mount
onMounted(() => {
  fetchTodos()
})

// Fetch all todos
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
      status: newTodo.value.status,
      plannedStart: newTodo.value.plannedStart || undefined,
      plannedEnd: newTodo.value.plannedEnd || undefined
    }

    await createTodo(todoData)
    await fetchTodos() // Refresh the list
    
    // Reset form
    newTodo.value = {
      taskId: null,
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

.delete-btn {
  background: #e74c3c;
  padding: 6px 12px;
  font-size: 12px;
}

.delete-btn:hover {
  background: #c0392b;
}

.loading, .error, .empty {
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