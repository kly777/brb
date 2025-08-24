<template>
  <div class="app">
    <h1>BRB Management System</h1>
    <div class="dashboard">
      <EventList class="dashboard-item" />
      <TaskList class="dashboard-item" />
      <TodoList class="dashboard-item" />
    </div>
  </div>
</template>

<script setup lang="ts">
import EventList from './components/EventList.vue'
import TaskList from './components/TaskList.vue'
import TodoList from './components/TodoList.vue'
</script>

<style scoped>
.app {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.app h1 {
  text-align: center;
  color: #333;
  margin-bottom: 30px;
}

.dashboard {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
}

.dashboard-item {
  background: #f9f9f9;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>
