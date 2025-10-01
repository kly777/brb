import { Show, createSignal } from 'solid-js'
import type { TodoResponse, TaskResponse, EventResponse, TodoUpdateRequest } from '../api/types'
import styles from './TodoItem.module.css'

interface Props {
  todo: TodoResponse
  tasks: TaskResponse[]
  events: EventResponse[]
  onEdit: (id: number, data: TodoUpdateRequest) => Promise<void>
  onDelete: (id: number) => Promise<void>
  isEditing: boolean
  onStartEdit: (id: number) => void
  onCancelEdit: () => void
}

/**
 * TodoItem组件 - 负责显示单个Todo项
 * 
 * SolidJS概念解释:
 * - 组件接收props并渲染单个todo项
 * - 使用Show组件进行条件渲染
 * - 事件处理函数通过props从父组件传递
 */
export default function TodoItem(props: Props) {
  // 格式化日期显示
  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return ''
    return new Date(dateString).toLocaleString('zh-CN')
  }

  // 根据ID获取任务描述
  const getTaskDescription = (taskId: number): string => {
    const task = props.tasks.find(t => t.id === taskId)
    return task ? task.description : `Task #${taskId}`
  }

  // 根据ID获取事件标题
  const getEventTitle = (eventId: number): string => {
    const event = props.events.find(e => e.id === eventId)
    return event ? event.title : `Event #${eventId}`
  }

  return (
    <li class={styles.todoItem}>
      <Show when={!props.isEditing}>
        <div class={styles.todoInfo}>
          <span class={styles.todoId}>#{props.todo.id}</span>
          <span class={styles.todoTask}>Task: {getTaskDescription(props.todo.taskId)}</span>
          <Show when={props.todo.eventId}>
            <span class={styles.todoEvent}>Event: {getEventTitle(props.todo.eventId!)}</span>
          </Show>
          <span class={`${styles.todoStatus} ${styles[props.todo.status]}`}>{props.todo.status}</span>
          
          <Show when={props.todo.plannedTime.start}>
            <div class={styles.todoTime}>
              Planned: {formatDate(props.todo.plannedTime.start)}
              <Show when={props.todo.plannedTime.end}>
                - {formatDate(props.todo.plannedTime.end)}
              </Show>
            </div>
          </Show>
          
          <Show when={props.todo.actualTime.start}>
            <div class={styles.todoTime}>
              Actual: {formatDate(props.todo.actualTime.start)}
              <Show when={props.todo.actualTime.end}>
                - {formatDate(props.todo.actualTime.end)}
              </Show>
            </div>
          </Show>
          
          <Show when={props.todo.completedTime}>
            <div class={styles.todoTime}>
              Completed: {formatDate(props.todo.completedTime)}
            </div>
          </Show>
        </div>
      </Show>

      {/* 编辑表单 - 暂时使用简单版本 */}
      <Show when={props.isEditing}>
        <div class={styles.editForm}>
          <h4>Edit Todo #{props.todo.id}</h4>
          <p>编辑功能待实现...</p>
          <div class={styles.editActions}>
            <button onClick={() => props.onCancelEdit()} class={styles.cancelBtn}>
              Cancel
            </button>
          </div>
        </div>
      </Show>

      <Show when={!props.isEditing}>
        <div class={styles.todoActions}>
          <button onClick={() => props.onStartEdit(props.todo.id)} class={`${styles.todoButton} ${styles.editButton}`}>
            Edit
          </button>
          <button onClick={() => props.onDelete(props.todo.id)} class={`${styles.todoButton} ${styles.deleteButton}`}>
            Delete
          </button>
        </div>
      </Show>
    </li>
  )
}