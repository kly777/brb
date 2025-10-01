import { lazy } from 'solid-js'
import './App.css'

/**
 * SolidJS主应用组件 - 从Vue版本转换而来
 * 
 * SolidJS概念解释:
 * - 使用lazy进行组件懒加载，提高应用性能
 * - 组件按需加载，减少初始包大小
 * - 使用CSS Grid布局创建响应式界面
 */

// 使用lazy进行组件懒加载
const TodoList = lazy(() => import('./components/TodoList'))
const TodoCalendar = lazy(() => import('./components/TodoCalendar'))
const EventTaskView = lazy(() => import('./components/EventTaskView'))

function App() {
  return (
    <div class="app-container">
      <header class="app-header">
        <h1>BRB Task Management System</h1>
        <p>Manage your events, tasks, and todos in one place</p>
      </header>

      <main class="app-main">
        <div class="app-content">
          {/* 
            SolidJS中的懒加载组件使用方式:
            - 使用lazy包装导入的组件
            - 这些组件会在需要时自动加载
            - 提高应用初始加载速度
          */}
          <TodoList />
          <TodoCalendar />
          <EventTaskView />
        </div>
      </main>
    </div>
  )
}

export default App
