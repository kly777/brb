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
  max-width: 900px;
  margin: 0 auto;
  padding: 24px;
}

.add-todo-form {
  background: rgba(255, 255, 255, 0.95);
  padding: 28px;
  border-radius: 16px;
  margin-bottom: 36px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.add-todo-form h3 {
  margin-top: 0;
  color: #2c3e50;
  font-size: 1.5rem;
  font-weight: 600;
  border-bottom: 3px solid #667eea;
  padding-bottom: 16px;
  margin-bottom: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #495057;
  font-size: 14px;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 14px 18px;
  border: 2px solid #e9ecef;
  border-radius: 10px;
  font-size: 15px;
  transition: all 0.3s ease;
  background: white;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

button {
  padding: 14px 28px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-size: 15px;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
}

button:disabled {
  background: #bdc3c7;
  box-shadow: none;
  transform: none;
  cursor: not-allowed;
}

.todos-container h3 {
  color: #2c3e50;
  margin-bottom: 20px;
  font-size: 1.4rem;
  font-weight: 600;
  border-bottom: 3px solid #667eea;
  padding-bottom: 12px;
}

.todos {
  list-style: none;
  padding: 0;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  margin-bottom: 16px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 14px;
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.todo-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 28px rgba(0, 0, 0, 0.12);
}

.todo-info {
  flex: 1;
}

.todo-id {
  font-weight: 700;
  color: #667eea;
  margin-right: 12px;
  font-size: 14px;
  background: rgba(102, 126, 234, 0.1);
  padding: 4px 10px;
  border-radius: 12px;
}

.todo-task {
  font-weight: 600;
  margin-right: 12px;
  color: #2c3e50;
}

.todo-event {
  font-style: italic;
  color: #6c757d;
  margin-right: 12px;
  background: rgba(108, 117, 125, 0.1);
  padding: 4px 10px;
  border-radius: 12px;
}

.todo-status {
  padding: 6px 14px;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 600;
  color: white;
  display: inline-block;
  margin-right: 8px;
}

.todo-status.pending {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a52 100%);
}

.todo-status.in_progress {
  background: linear-gradient(135deg, #4ecdc4 0%, #44a08d 100%);
}

.todo-status.completed {
  background: linear-gradient(135deg, #1dd1a1 0%, #00b894 100%);
}

.todo-status.cancelled {
  background: linear-gradient(135deg, #8395a7 0%, #576574 100%);
}

.todo-time {
  font-size: 13px;
  color: #6c757d;
  margin-top: 8px;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 8px;
  border-left: 4px solid #667eea;
}

.loading,
.error,
.empty {
  text-align: center;
  padding: 40px;
  border-radius: 14px;
  font-size: 16px;
  margin: 20px 0;
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

/* Edit form styles */
.edit-form {
  flex: 1;
  background: rgba(255, 248, 225, 0.95);
  padding: 20px;
  border-radius: 12px;
  margin-right: 15px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.edit-form h4 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #2c3e50;
  font-size: 1.2rem;
  font-weight: 600;
}

.edit-form .form-group {
  margin-bottom: 16px;
}

.edit-form .form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 13px;
  font-weight: 600;
  color: #495057;
}

.edit-form .form-group input,
.edit-form .form-group select {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 13px;
  background: white;
}

.edit-form .form-group input:focus,
.edit-form .form-group select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.edit-actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}

.edit-actions button {
  padding: 10px 20px;
  font-size: 13px;
  border-radius: 8px;
}

.cancel-btn {
  background: linear-gradient(135deg, #95a5a6 0%, #7f8c8d 100%);
}

.cancel-btn:hover {
  background: linear-gradient(135deg, #7f8c8d 0%, #6c7b80 100%);
}

/* Todo actions styles */
.todo-actions {
  display: flex;
  gap: 10px;
}

.edit-btn {
  background: linear-gradient(135deg, #f39c12 0%, #e67e22 100%);
  padding: 8px 16px;
  font-size: 13px;
  border-radius: 8px;
  box-shadow: 0 3px 8px rgba(243, 156, 18, 0.3);
}

.edit-btn:hover {
  background: linear-gradient(135deg, #e67e22 0%, #d35400 100%);
  transform: translateY(-2px);
  box-shadow: 0 5px 12px rgba(243, 156, 18, 0.4);
}

.delete-btn {
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
  padding: 8px 16px;
  font-size: 13px;
  border-radius: 8px;
  box-shadow: 0 3px 8px rgba(231, 76, 60, 0.3);
}

.delete-btn:hover {
  background: linear-gradient(135deg, #c0392b 0%, #a93226 100%);
  transform: translateY(-2px);
  box-shadow: 0 5px 12px rgba(231, 76, 60, 0.4);
}

/* Responsive design */
@media (max-width: 768px) {
  .todo-list {
    padding: 16px;
  }
  
  .add-todo-form {
    padding: 20px;
    margin-bottom: 28px;
  }
  
  .todo-item {
    padding: 20px;
    flex-direction: column;
    gap: 16px;
  }
  
  .todo-info {
    width: 100%;
  }
  
  .todo-actions {
    width: 100%;
    justify-content: flex-end;
  }
  
  .edit-form {
    margin-right: 0;
    margin-bottom: 15px;
  }
}

@media (max-width: 480px) {
  .todo-list {
    padding: 12px;
  }
  
  .add-todo-form {
    padding: 16px;
  }
  
  .form-group input,
  .form-group select {
    padding: 12px 16px;
  }
  
  .todo-item {
    padding: 16px;
  }
  
  .todo-status {
    font-size: 12px;
    padding: 5px 12px;
  }
  
  .edit-actions,
  .todo-actions {
    flex-direction: column;
    gap: 8px;
  }
  
  .edit-actions button,
  .todo-actions button {
    width: 100%;
    text-align: center;
  }
}
</style>