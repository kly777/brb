<template>
  <div class="todo-list">
    <h2>Todos</h2>
    <div class="todo-form">
      <input v-model="newTodo.taskId" placeholder="Task ID" />
      <select v-model="newTodo.status">
        <option value="pending">Pending</option>
        <option value="in_progress">In Progress</option>
        <option value="completed">Completed</option>
      </select>
      <input v-model="newTodo.priority" type="number" placeholder="Priority" />
      <input v-model="newTodo.startTime" type="datetime-local" />
      <input v-model="newTodo.endTime" type="datetime-local" />
      <button @click="createTodo">Create Todo</button>
    </div>
    <ul class="todo-items">
      <li v-for="todo in todos" :key="todo.id" class="todo-item">
        <div class="todo-info">
          <h3>Task ID: {{ todo.taskId }}</h3>
          <p>Status: {{ todo.status }}</p>
          <p>Priority: {{ todo.priority }}</p>
          <p v-if="todo.completedTime">Completed: {{ formatDate(todo.completedTime) }}</p>
          <p>Start: {{ formatDate(todo.startTime) }}</p>
          <p>End: {{ formatDate(todo.endTime) }}</p>
        </div>
        <div class="todo-actions">
          <button @click="completeTodo(todo)" v-if="todo.status !== 'completed'">Complete</button>
          <button @click="editTodo(todo)">Edit</button>
          <button @click="deleteTodo(todo.id)">Delete</button>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Todo {
  id: string
  taskId: string
  status: string
  priority: number
  completedTime: string | null
  startTime: string
  endTime: string
}

const todos = ref<Todo[]>([])
const newTodo = ref({
  taskId: '',
  status: 'pending',
  priority: 0,
  startTime: '',
  endTime: ''
})

const API_BASE = 'http://localhost:8080/api'

const fetchTodos = async () => {
  try {
    const response = await fetch(`${API_BASE}/todos`)
    if (response.ok) {
      todos.value = await response.json()
    }
  } catch (error) {
    console.error('Failed to fetch todos:', error)
  }
}

const createTodo = async () => {
  try {
    const response = await fetch(`${API_BASE}/todos`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        ...newTodo.value,
        priority: Number(newTodo.value.priority)
      })
    })
    if (response.ok) {
      newTodo.value = { 
        taskId: '', 
        status: 'pending', 
        priority: 0, 
        startTime: '', 
        endTime: '' 
      }
      fetchTodos()
    }
  } catch (error) {
    console.error('Failed to create todo:', error)
  }
}

const completeTodo = async (todo: Todo) => {
  try {
    const updatedTodo = {
      ...todo,
      status: 'completed',
      completedTime: new Date().toISOString()
    }
    
    const response = await fetch(`${API_BASE}/todos/${todo.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(updatedTodo)
    })
    if (response.ok) {
      fetchTodos()
    }
  } catch (error) {
    console.error('Failed to complete todo:', error)
  }
}

const editTodo = (todo: Todo) => {
  console.log('Edit todo:', todo)
}

const deleteTodo = async (id: string) => {
  try {
    const response = await fetch(`${API_BASE}/todos/${id}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      fetchTodos()
    }
  } catch (error) {
    console.error('Failed to delete todo:', error)
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

onMounted(() => {
  fetchTodos()
})
</script>

<style scoped>
.todo-list {
  margin: 20px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
}

.todo-form {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.todo-form input,
.todo-form select {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.todo-form button {
  padding: 8px 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.todo-items {
  list-style: none;
  padding: 0;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.todo-info h3 {
  margin: 0 0 5px 0;
}

.todo-actions {
  display: flex;
  gap: 10px;
}

.todo-actions button {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.todo-actions button:first-child {
  background-color: #FF9800;
  color: white;
}

.todo-actions button:nth-child(2) {
  background-color: #2196F3;
  color: white;
}

.todo-actions button:last-child {
  background-color: #f44336;
  color: white;
}
</style>