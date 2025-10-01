import { createSignal, createEffect, For, Show } from 'solid-js'
import { getAllEvents } from '../api/event'
import { getAllTasks } from '../api/task'
import { createTodo } from '../api/todo'
import type { EventResponse, TaskResponse, TodoCreateRequest } from '../api/types'
import AddTask from './AddTask'
import styles from './EventTaskView.module.css'

/**
 * EventTaskView组件 - 显示事件和任务管理界面
 * 
 * SolidJS概念解释:
 * - 使用createSignal管理多个状态
 * - 使用createEffect在组件挂载时获取数据
 * - 将事件和任务分别管理，提高代码清晰度
 */
export default function EventTaskView() {
  // 事件相关状态
  const [events, setEvents] = createSignal<EventResponse[]>([])
  const [eventsLoading, setEventsLoading] = createSignal(false)
  const [eventsError, setEventsError] = createSignal('')

  // 任务相关状态
  const [tasks, setTasks] = createSignal<TaskResponse[]>([])
  const [tasksLoading, setTasksLoading] = createSignal(false)
  const [tasksError, setTasksError] = createSignal('')

  // 快速添加Todo表单状态
  const [activeTodoForm, setActiveTodoForm] = createSignal<number | null>(null)
  const [addingTodo, setAddingTodo] = createSignal(false)
  const [quickTodo, setQuickTodo] = createSignal({
    status: '',
    plannedStart: '',
    plannedEnd: ''
  })

  // 组件挂载时获取数据
  createEffect(() => {
    fetchEvents()
    fetchTasks()
  })

  // 获取所有事件
  const fetchEvents = async () => {
    try {
      setEventsLoading(true)
      setEventsError('')
      const response = await getAllEvents()
      setEvents(response)
    } catch (err) {
      setEventsError('Failed to load events')
      console.error('Error fetching events:', err)
    } finally {
      setEventsLoading(false)
    }
  }

  // 获取所有任务
  const fetchTasks = async () => {
    try {
      setTasksLoading(true)
      setTasksError('')
      const response = await getAllTasks()
      setTasks(response)
    } catch (err) {
      setTasksError('Failed to load tasks')
      console.error('Error fetching tasks:', err)
    } finally {
      setTasksLoading(false)
    }
  }

  // 刷新任务列表（用于AddTask组件回调）
  const handleTaskAdded = () => {
    fetchTasks()
  }

  // 格式化日期显示
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('zh-CN')
  }

  // 格式化时间段显示
  const formatTimeSpan = (timeSpan: { start?: string; end?: string }) => {
    if (!timeSpan.start && !timeSpan.end) return 'No time set'
    
    const start = timeSpan.start ? new Date(timeSpan.start).toLocaleString('zh-CN') : 'Not set'
    const end = timeSpan.end ? new Date(timeSpan.end).toLocaleString('zh-CN') : 'Not set'
    
    return `${start} - ${end}`
  }

  // 显示添加Todo表单
  const showAddTodoForm = (taskId: number) => {
    setActiveTodoForm(taskId)
    setQuickTodo({
      status: '',
      plannedStart: '',
      plannedEnd: ''
    })
  }

  // 取消添加Todo
  const cancelAddTodo = () => {
    setActiveTodoForm(null)
  }

  // 快速添加Todo
  const addQuickTodo = async (taskId: number) => {
    const currentQuickTodo = quickTodo()
    if (!currentQuickTodo.status) {
      return
    }

    try {
      setAddingTodo(true)
      const todoData: TodoCreateRequest = {
        taskId: taskId,
        status: currentQuickTodo.status,
        plannedStart: currentQuickTodo.plannedStart || undefined,
        plannedEnd: currentQuickTodo.plannedEnd || undefined
      }

      await createTodo(todoData)
      setActiveTodoForm(null)
      setQuickTodo({
        status: '',
        plannedStart: '',
        plannedEnd: ''
      })
      
      alert('Todo added successfully!')
    } catch (err) {
      console.error('Error adding todo:', err)
      alert('Failed to add todo')
    } finally {
      setAddingTodo(false)
    }
  }

  // 更新快速Todo表单字段
  const updateQuickTodoField = (field: 'status' | 'plannedStart' | 'plannedEnd', value: string) => {
    setQuickTodo(prev => ({ ...prev, [field]: value }))
  }

  return (
    <div class={styles.eventTaskView}>
      <h2>Events and Tasks Management</h2>

      {/* 添加任务表单 */}
      <AddTask onTaskAdded={handleTaskAdded} />

      {/* 事件列表 */}
      <div class={styles.section}>
        <h3>Events</h3>
        <Show when={eventsLoading()}>
          <div class={styles.loading}>Loading events...</div>
        </Show>
        
        <Show when={eventsError()}>
          <div class={styles.error}>{eventsError()}</div>
        </Show>
        
        <Show when={!eventsLoading() && events().length === 0}>
          <div class={styles.empty}>No events found.</div>
        </Show>
        
        <Show when={!eventsLoading() && events().length > 0}>
          <div class={styles.eventsList}>
            <For each={events()}>
              {(event) => (
                <div class={styles.eventItem}>
                  <div class={styles.eventInfo}>
                    <h4>{event.title}</h4>
                    <p>{event.description}</p>
                    <div class={styles.eventDetails}>
                      <span class={styles.id}>ID: {event.id}</span>
                      <span class={styles.location}>Location: {event.location}</span>
                      <span class={styles.priority}>Priority: {event.priority}</span>
                      <span class={styles.category}>Category: {event.category}</span>
                      <span class={styles.template}>Template: {event.isTemplate ? 'Yes' : 'No'}</span>
                    </div>
                  </div>
                </div>
              )}
            </For>
          </div>
        </Show>
      </div>

      {/* 任务列表 */}
      <div class={styles.section}>
        <h3>Tasks</h3>
        <Show when={tasksLoading()}>
          <div class={styles.loading}>Loading tasks...</div>
        </Show>
        
        <Show when={tasksError()}>
          <div class={styles.error}>{tasksError()}</div>
        </Show>
        
        <Show when={!tasksLoading() && tasks().length === 0}>
          <div class={styles.empty}>No tasks found.</div>
        </Show>
        
        <Show when={!tasksLoading() && tasks().length > 0}>
          <div class={styles.tasksList}>
            <For each={tasks()}>
              {(task) => (
                <div class={styles.taskItem}>
                  <div class={styles.taskInfo}>
                    <h4>Task #{task.id}</h4>
                    <p>{task.description}</p>
                    <div class={styles.taskDetails}>
                      <span class={styles.eventId}>Event ID: {task.eventId}</span>
                      <span class={styles.status}>Status: {task.status}</span>
                      <span class={styles.createdAt}>Created: {formatDate(task.createdAt)}</span>
                      
                      <Show when={task.allowedTime.start || task.allowedTime.end}>
                        <div class={styles.timeRange}>
                          <span>Allowed: {formatTimeSpan(task.allowedTime)}</span>
                        </div>
                      </Show>
                      
                      <Show when={task.plannedTime.start || task.plannedTime.end}>
                        <div class={styles.timeRange}>
                          <span>Planned: {formatTimeSpan(task.plannedTime)}</span>
                        </div>
                      </Show>
                    </div>
                  </div>
                  
                  <div class={styles.taskActions}>
                    <button onClick={() => showAddTodoForm(task.id)} class={styles.addTodoBtn}>
                      Add Todo
                    </button>
                    
                    <Show when={activeTodoForm() === task.id}>
                      <div class={styles.quickTodoForm}>
                        <h5>Add Todo for Task #{task.id}</h5>
                        <form onSubmit={(e) => { e.preventDefault(); addQuickTodo(task.id) }}>
                          <div class={styles.formGroup}>
                            <label>Status:</label>
                            <select
                              value={quickTodo().status}
                              onChange={(e) => updateQuickTodoField('status', e.target.value)}
                              required
                            >
                              <option value="">Select status</option>
                              <option value="pending">Pending</option>
                              <option value="in_progress">In Progress</option>
                              <option value="completed">Completed</option>
                              <option value="cancelled">Cancelled</option>
                            </select>
                          </div>
                          
                          <div class={styles.formGroup}>
                            <label>Planned Start:</label>
                            <input
                              type="datetime-local"
                              value={quickTodo().plannedStart}
                              onChange={(e) => updateQuickTodoField('plannedStart', e.target.value)}
                            />
                          </div>
                          
                          <div class={styles.formGroup}>
                            <label>Planned End:</label>
                            <input
                              type="datetime-local"
                              value={quickTodo().plannedEnd}
                              onChange={(e) => updateQuickTodoField('plannedEnd', e.target.value)}
                            />
                          </div>
                          
                          <button type="submit" disabled={addingTodo()}>
                            {addingTodo() ? 'Adding...' : 'Add Todo'}
                          </button>
                          <button type="button" onClick={cancelAddTodo} class={styles.cancelBtn}>
                            Cancel
                          </button>
                        </form>
                      </div>
                    </Show>
                  </div>
                </div>
              )}
            </For>
          </div>
        </Show>
      </div>
    </div>
  )
}