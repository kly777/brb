<script setup lang="ts">
import TodoCalendar from './components/TodoCalendar.vue'
import TodoList from './components/TodoList.vue'
import EventTaskView from './components/EventTaskView.vue'
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <h1>BRB Task Management System</h1>
      <p>Manage your events, tasks, and todos in one place</p>
    </header>
    
    <main class="app-main">
      <div class="app-content">
        <TodoList />
        <TodoCalendar />
        <EventTaskView />
      </div>
    </main>
  </div>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.app-header {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  padding: 2rem;
  text-align: center;
  color: white;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
}

.app-header h1 {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.app-header p {
  font-size: 1.1rem;
  opacity: 0.9;
  margin: 0;
}

.app-main {
  padding: 2rem;
}

.app-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .app-header {
    padding: 1.5rem;
  }
  
  .app-header h1 {
    font-size: 2rem;
  }
  
  .app-main {
    padding: 1rem;
  }
  
  .app-content {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
}

@media (max-width: 480px) {
  .app-header {
    padding: 1rem;
  }
  
  .app-header h1 {
    font-size: 1.75rem;
  }
  
  .app-main {
    padding: 0.5rem;
  }
}
</style>
