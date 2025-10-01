import { createSignal } from 'solid-js'
import { createTask } from '../api/task'
import type { TaskCreateRequest } from '../api/types'
import styles from './AddTask.module.css'

interface Props {
  onTaskAdded: () => void
}

/**
 * AddTask组件 - 负责添加新任务
 * 
 * SolidJS概念解释:
 * - 使用createSignal管理本地表单状态
 * - 通过props接收回调函数，在任务添加成功后通知父组件
 * - 表单提交使用事件处理函数
 */
export default function AddTask(props: Props) {
  const [formData, setFormData] = createSignal({
    eventId: '',
    parentTaskId: '',
    description: '',
    allowedStart: '',
    allowedEnd: '',
    plannedStart: '',
    plannedEnd: '',
    status: 'pending'
  })

  const [loading, setLoading] = createSignal(false)
  const [error, setError] = createSignal('')

  // 更新表单字段的辅助函数
  type FormFields = 'eventId' | 'parentTaskId' | 'description' | 'allowedStart' | 'allowedEnd' | 'plannedStart' | 'plannedEnd' | 'status'
  const updateFormField = (field: FormFields, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  // 处理表单提交
  const handleSubmit = async (e: Event) => {
    e.preventDefault()
    
    const currentFormData = formData()
    if (!currentFormData.eventId || !currentFormData.description) {
      setError('Event ID and Description are required')
      return
    }

    try {
      setLoading(true)
      setError('')

      const taskData: TaskCreateRequest = {
        eventId: parseInt(currentFormData.eventId),
        parentTaskId: currentFormData.parentTaskId ? parseInt(currentFormData.parentTaskId) : undefined,
        preTaskIds: [], // 暂时为空，可以根据需要扩展
        description: currentFormData.description,
        allowedStart: currentFormData.allowedStart || undefined,
        allowedEnd: currentFormData.allowedEnd || undefined,
        plannedStart: currentFormData.plannedStart || undefined,
        plannedEnd: currentFormData.plannedEnd || undefined,
        status: currentFormData.status
      }

      await createTask(taskData)
      
      // 重置表单
      setFormData({
        eventId: '',
        parentTaskId: '',
        description: '',
        allowedStart: '',
        allowedEnd: '',
        plannedStart: '',
        plannedEnd: '',
        status: 'pending'
      })

      // 通知父组件任务已添加
      props.onTaskAdded()
      
    } catch (err) {
      setError('Failed to add task')
      console.error('Error adding task:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div class={styles.addTaskForm}>
      <h3>Add New Task</h3>
      <form onSubmit={handleSubmit}>
        <div class={styles.formGroup}>
          <label for="eventId">Event ID:</label>
          <input
            id="eventId"
            type="number"
            value={formData().eventId}
            onInput={(e) => updateFormField('eventId', e.currentTarget.value)}
            required
            placeholder="Enter event ID"
          />
        </div>

        <div class={styles.formGroup}>
          <label for="parentTaskId">Parent Task ID (optional):</label>
          <input
            id="parentTaskId"
            type="number"
            value={formData().parentTaskId}
            onInput={(e) => updateFormField('parentTaskId', e.currentTarget.value)}
            placeholder="Enter parent task ID"
          />
        </div>

        <div class={styles.formGroup}>
          <label for="description">Description:</label>
          <textarea
            id="description"
            value={formData().description}
            onInput={(e) => updateFormField('description', e.currentTarget.value)}
            required
            placeholder="Enter task description"
            rows={3}
          />
        </div>

        <div class={styles.formGroup}>
          <label for="status">Status:</label>
          <select
            id="status"
            value={formData().status}
            onChange={(e) => updateFormField('status', e.target.value)}
          >
            <option value="pending">Pending</option>
            <option value="in_progress">In Progress</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>

        <div class={styles.formRow}>
          <div class={styles.formGroup}>
            <label for="allowedStart">Allowed Start (optional):</label>
            <input
              id="allowedStart"
              type="datetime-local"
              value={formData().allowedStart}
              onChange={(e) => updateFormField('allowedStart', e.target.value)}
            />
          </div>

          <div class={styles.formGroup}>
            <label for="allowedEnd">Allowed End (optional):</label>
            <input
              id="allowedEnd"
              type="datetime-local"
              value={formData().allowedEnd}
              onChange={(e) => updateFormField('allowedEnd', e.target.value)}
            />
          </div>
        </div>

        <div class={styles.formRow}>
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
        </div>

        {error() && <div class={styles.errorMessage}>{error()}</div>}

        <button type="submit" disabled={loading()} class={styles.button}>
          {loading() ? 'Adding...' : 'Add Task'}
        </button>
      </form>
    </div>
  )
}