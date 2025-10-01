import { createSignal, createEffect, createMemo, Show, For } from 'solid-js'
import { getAllTodos } from '../api/todo'
import type { TodoResponse } from '../api/types'
import styles from './TodoCalendar.module.css'

/**
 * TodoCalendar组件 - 简化版本，不使用FullCalendar
 * 
 * SolidJS概念解释:
 * - createMemo: 创建派生状态，当依赖项变化时自动重新计算
 * - 手动实现简单的日历视图
 * - 使用CSS Grid布局创建日历网格
 */
export default function TodoCalendar() {
  const [todos, setTodos] = createSignal<TodoResponse[]>([])
  const [loading, setLoading] = createSignal(false)
  const [error, setError] = createSignal('')
  const [currentDate, setCurrentDate] = createSignal(new Date())

  // 获取当前月份显示文本 - 使用createMemo进行性能优化
  const currentMonth = createMemo(() => {
    return currentDate().toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long'
    })
  })

  // 生成当前月份的日期数组
  const monthDays = createMemo(() => {
    const date = new Date(currentDate())
    const year = date.getFullYear()
    const month = date.getMonth()
    
    // 获取月份的第一天和最后一天
    const firstDay = new Date(year, month, 1)
    const lastDay = new Date(year, month + 1, 0)
    
    // 获取第一天是星期几（0-6，0是星期日）
    const firstDayOfWeek = firstDay.getDay()
    
    // 获取月份的总天数
    const daysInMonth = lastDay.getDate()
    
    // 生成日期数组
    const days = []
    
    // 添加上个月的最后几天（填充日历开头）
    const prevMonthLastDay = new Date(year, month, 0).getDate()
    for (let i = firstDayOfWeek - 1; i >= 0; i--) {
      days.push({
        date: new Date(year, month - 1, prevMonthLastDay - i),
        isCurrentMonth: false,
        todos: []
      })
    }
    
    // 添加当前月的日期
    for (let i = 1; i <= daysInMonth; i++) {
      const dayDate = new Date(year, month, i)
      const dayTodos = todos().filter(todo => {
        const todoDate = todo.plannedTime.start ? new Date(todo.plannedTime.start) : null
        return todoDate && 
               todoDate.getDate() === dayDate.getDate() &&
               todoDate.getMonth() === dayDate.getMonth() &&
               todoDate.getFullYear() === dayDate.getFullYear()
      })
      
      days.push({
        date: dayDate,
        isCurrentMonth: true,
        todos: dayTodos
      })
    }
    
    // 添加下个月的前几天（填充日历结尾）
    const totalCells = 42 // 6行 * 7列
    const remainingCells = totalCells - days.length
    for (let i = 1; i <= remainingCells; i++) {
      days.push({
        date: new Date(year, month + 1, i),
        isCurrentMonth: false,
        todos: []
      })
    }
    
    return days
  })

  // 根据状态获取颜色
  const getStatusColor = (status: string): string => {
    const colors: { [key: string]: string } = {
      pending: '#ff6b6b',
      in_progress: '#4ecdc4',
      completed: '#1dd1a1',
      cancelled: '#8395a7'
    }
    return colors[status] || '#3498db'
  }

  // 获取todos数据
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

  // 切换到上个月
  const prevMonth = () => {
    const newDate = new Date(currentDate())
    newDate.setMonth(newDate.getMonth() - 1)
    setCurrentDate(newDate)
  }

  // 切换到下个月
  const nextMonth = () => {
    const newDate = new Date(currentDate())
    newDate.setMonth(newDate.getMonth() + 1)
    setCurrentDate(newDate)
  }

  // 格式化日期显示
  const formatDate = (date: Date) => {
    return date.getDate()
  }

  // 组件挂载时获取数据
  createEffect(() => {
    fetchTodos()
  })

  return (
    <div class={styles.calendar}>
      <div class={styles.calendarHeader}>
        <h2>Todo Calendar</h2>
        <div class={styles.controls}>
          <button onClick={prevMonth} class={styles.navButton}>← Previous</button>
          <span class={styles.currentMonth}>{currentMonth()}</span>
          <button onClick={nextMonth} class={styles.navButton}>Next →</button>
        </div>
      </div>

      <Show when={loading()}>
        <div class={styles.loading}>Loading todos...</div>
      </Show>
      
      <Show when={error()}>
        <div class={styles.error}>{error()}</div>
      </Show>

      <Show when={!loading()}>
        <div class={styles.calendarGrid}>
          {/* 星期标题 */}
          <div class={styles.calendarWeekdays}>
            <div class={styles.weekday}>Sun</div>
            <div class={styles.weekday}>Mon</div>
            <div class={styles.weekday}>Tue</div>
            <div class={styles.weekday}>Wed</div>
            <div class={styles.weekday}>Thu</div>
            <div class={styles.weekday}>Fri</div>
            <div class={styles.weekday}>Sat</div>
          </div>

          {/* 日期网格 */}
          <div class={styles.calendarDays}>
            <For each={monthDays()}>
              {(day) => (
                <div
                  class={`${styles.calendarDay} ${day.isCurrentMonth ? styles.currentMonth : styles.otherMonth}`}
                >
                  <div class={styles.dayNumber}>{formatDate(day.date)}</div>
                  <div class={styles.dayTodos}>
                    <For each={day.todos}>
                      {(todo) => (
                        <div
                          class={styles.todoMarker}
                          style={{ background: getStatusColor(todo.status) }}
                          title={`Todo #${todo.id} - ${todo.status}`}
                        />
                      )}
                    </For>
                  </div>
                </div>
              )}
            </For>
          </div>
        </div>
      </Show>
    </div>
  )
}