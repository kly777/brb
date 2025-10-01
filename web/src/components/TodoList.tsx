import { createSignal, createEffect, For, Show } from 'solid-js'
import { createTodo, getAllTodos, updateTodo, deleteTodo as apiDeleteTodo } from '../api/todo'
import { getAllTasks } from '../api/task'
import { getAllEvents } from '../api/event'
import type { TodoResponse, TodoCreateRequest, TodoUpdateRequest, TaskResponse, EventResponse } from '../api/types'
import AddTodoForm from './AddTodoForm'
import TodoItem from './TodoItem'
import styles from './TodoList.module.css'

/**
 * SolidJS TodoList组件 - 重构版本，使用拆分后的子组件
 * 
 * SolidJS核心概念解释:
 * - 使用createSignal管理组件状态
 * - 使用createEffect处理副作用（相当于Vue的onMounted + watch）
 * - 使用For和Show组件进行列表渲染和条件渲染
 * - 将大型组件拆分为多个小组件，提高可维护性
 */
export default function TodoList() {
  // 使用createSignal创建响应式状态
  const [todos, setTodos] = createSignal<TodoResponse[]>([])
  const [tasks, setTasks] = createSignal<TaskResponse[]>([])
  const [events, setEvents] = createSignal<EventResponse[]>([])
  const [loading, setLoading] = createSignal(false)
  const [error, setError] = createSignal('')
  const [adding, setAdding] = createSignal(false)
  const [updating, setUpdating] = createSignal(false)
  const [editingTodoId, setEditingTodoId] = createSignal<number | null>(null)

  // createEffect用于副作用，相当于Vue的onMounted + watch
  createEffect(() => {
    fetchData()
  })

  // 获取所有数据
  const fetchData = async () => {
    try {
      setLoading(true)
      setError('')
      
      // 使用Promise.all并发请求
      const [todosResponse, tasksResponse, eventsResponse] = await Promise.all([
        getAllTodos(),
        getAllTasks(),
        getAllEvents()
      ])
      
      setTodos(todosResponse)
      setTasks(tasksResponse)
      setEvents(eventsResponse)
    } catch (err) {
      setError('Failed to load data')
      console.error('Error fetching data:', err)
    } finally {
      setLoading(false)
    }
  }

  // 只获取todos数据（用于刷新）
  const fetchTodos = async () => {
    try {
      setLoading(true)
      setError('')
      const response = await getAllTodos()
      setTodos(response)
    } catch (err) {
      setError('Failed to load todos')
      console.error('Error fetching todos:', err)
    } finally {
      setLoading(false)
    }
  }

  // 添加新todo
  const handleAddTodo = async (todoData: TodoCreateRequest) => {
    try {
      setAdding(true)
      setError('')
      await createTodo(todoData)
      await fetchTodos() // 刷新列表
    } catch (err) {
      setError('Failed to add todo')
      console.error('Error adding todo:', err)
    } finally {
      setAdding(false)
    }
  }

  // 删除todo
  const handleDeleteTodo = async (id: number) => {
    if (!confirm('Are you sure you want to delete this todo?')) {
      return
    }

    try {
      setError('')
      await apiDeleteTodo(id)
      await fetchTodos() // 刷新列表
    } catch (err) {
      setError('Failed to delete todo')
      console.error('Error deleting todo:', err)
    }
  }

  // 保存编辑
  const handleSaveEdit = async (id: number, data: TodoUpdateRequest) => {
    try {
      setUpdating(true)
      setError('')
      await updateTodo(id, data)
      setEditingTodoId(null)
      await fetchTodos() // 刷新列表
    } catch (err) {
      setError('Failed to update todo')
      console.error('Error updating todo:', err)
    } finally {
      setUpdating(false)
    }
  }

  // 开始编辑
  const handleStartEdit = (id: number) => {
    setEditingTodoId(id)
  }

  // 取消编辑
  const handleCancelEdit = () => {
    setEditingTodoId(null)
  }

  return (
    <div class={styles.todoList}>
      <h2>Todo Management</h2>

      {/* 使用AddTodoForm组件 */}
      <AddTodoForm
        tasks={tasks()}
        events={events()}
        onAddTodo={handleAddTodo}
        loading={adding()}
      />

      {/* Todo列表 */}
      <div class={styles.todosContainer}>
        <h3>Existing Todos</h3>
        
        {/* Show组件用于条件渲染 */}
        <Show when={loading()}>
          <div class={styles.loading}>Loading todos...</div>
        </Show>
        
        <Show when={error()}>
          <div class={styles.error}>{error()}</div>
        </Show>
        
        <Show when={!loading() && todos().length === 0}>
          <div class={styles.empty}>No todos found.</div>
        </Show>
        
        <Show when={!loading() && todos().length > 0}>
          <ul class={styles.todos}>
            <For each={todos()}>
              {(todo) => (
                <TodoItem
                  todo={todo}
                  tasks={tasks()}
                  events={events()}
                  onEdit={handleSaveEdit}
                  onDelete={handleDeleteTodo}
                  isEditing={editingTodoId() === todo.id}
                  onStartEdit={handleStartEdit}
                  onCancelEdit={handleCancelEdit}
                />
              )}
            </For>
          </ul>
        </Show>
      </div>
    </div>
  )
}
