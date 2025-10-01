import { createSignal, For } from 'solid-js'
import type { TaskResponse, EventResponse, TodoCreateRequest } from '../api/types'
import styles from './AddTodoForm.module.css'

interface Props {
  tasks: TaskResponse[]
  events: EventResponse[]
  onAddTodo: (data: TodoCreateRequest) => Promise<void>
  loading: boolean
}

/**
 * AddTodoForm组件 - 负责添加新Todo的表单
 * 
 * SolidJS概念解释:
 * - 使用props接收父组件传递的数据和方法
 * - 在Solid中，props是响应式的，可以直接在JSX中使用
 * - 使用createSignal管理本地表单状态
 */
export default function AddTodoForm(props: Props) {
  // 本地表单状态管理
  const [formData, setFormData] = createSignal({
    taskId: null as number | null,
    eventId: null as number | null,
    status: '',
    plannedStart: '',
    plannedEnd: ''
  })

  // 更新表单字段的辅助函数
  type FormFields = 'taskId' | 'eventId' | 'status' | 'plannedStart' | 'plannedEnd'
  const updateFormField = (field: FormFields, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  // 处理表单提交
  const handleSubmit = async (e: Event) => {
    e.preventDefault()
    
    const currentFormData = formData()
    if (!currentFormData.taskId || !currentFormData.status) {
      alert('Task ID and Status are required')
      return
    }

    const todoData: TodoCreateRequest = {
      taskId: currentFormData.taskId,
      eventId: currentFormData.eventId || undefined,
      status: currentFormData.status,
      plannedStart: currentFormData.plannedStart || undefined,
      plannedEnd: currentFormData.plannedEnd || undefined
    }

    await props.onAddTodo(todoData)

    // 重置表单
    setFormData({
      taskId: null,
      eventId: null,
      status: '',
      plannedStart: '',
      plannedEnd: ''
    })
  }

  return (
    <div class={styles.addTodoForm}>
      <h3>Add New Todo</h3>
      <form onSubmit={handleSubmit}>
        <div class={styles.formGroup}>
          <label for="taskId">Task:</label>
          <select
            id="taskId"
            value={formData().taskId || ''}
            onChange={(e) => updateFormField('taskId', e.target.value ? parseInt(e.target.value) : null)}
            required
          >
            <option value="">Select a task</option>
            <For each={props.tasks}>
              {(task) => (
                <option value={task.id}>
                  {task.description} (ID: {task.id})
                </option>
              )}
            </For>
          </select>
        </div>

        <div class={styles.formGroup}>
          <label for="eventId">Event (optional):</label>
          <select
            id="eventId"
            value={formData().eventId || ''}
            onChange={(e) => updateFormField('eventId', e.target.value ? parseInt(e.target.value) : null)}
          >
            <option value="">Select an event (optional)</option>
            <For each={props.events}>
              {(event) => (
                <option value={event.id}>
                  {event.title} (ID: {event.id})
                </option>
              )}
            </For>
          </select>
        </div>

        <div class={styles.formGroup}>
          <label for="status">Status:</label>
          <select
            id="status"
            value={formData().status}
            onChange={(e) => updateFormField('status', e.target.value)}
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
          <label for="plannedStart">Planned Start (optional):</label>
          <input
            id="plannedStart"
            type="datetime-local"
            value={formData().plannedStart}
            onChange={(e) => updateFormField('plannedStart', e.target.value)}
          />
        </div>

        <div class={styles.formGroup}>
          <label for="plannedEnd">Planned End (optional):</label>
          <input
            id="plannedEnd"
            type="datetime-local"
            value={formData().plannedEnd}
            onChange={(e) => updateFormField('plannedEnd', e.target.value)}
          />
        </div>

        <button type="submit" disabled={props.loading} class={styles.button}>
          {props.loading ? 'Adding...' : 'Add Todo'}
        </button>
      </form>
    </div>
  )
}