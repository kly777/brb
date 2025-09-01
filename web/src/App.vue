<script setup lang="ts">
import TodoCalendar from './components/TodoCalendar.vue'
import TodoList from './components/TodoList.vue'
</script>

<template>
  <div class="app-content">
    <TodoList />
    <TodoCalendar />
  </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}

.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}

.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}

.app-content {
  display: flex;
  gap: 20px;
  margin-top: 20px;
}

.app-content>* {
  flex: 1;
  min-width: 0;
}
</style>
